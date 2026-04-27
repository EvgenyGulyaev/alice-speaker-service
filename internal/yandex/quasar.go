package yandex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	quasarBaseURL            = "https://iot.quasar.yandex.ru"
	quasarAuthURL            = "https://mobileproxy.passport.yandex.net/1/bundle/auth/x_token/"
	quasarPageURL            = "https://yandex.ru/quasar"
	quasarScenarioNamePrefix = "Codex "
	quasarScenarioTTL        = 5 * time.Minute
)

var quasarCSRFRegexp = regexp.MustCompile(`"csrfToken2":"(.+?)"`)
var quasarScenarioDeletionScheduler = newScenarioDeletionScheduler()

type quasarSession struct {
	httpClient *http.Client
	csrfToken  string
	baseURL    string
	authURL    string
	pageURL    string
}

type quasarStatusResponse struct {
	Status       string           `json:"status"`
	Message      string           `json:"message"`
	RequestID    string           `json:"request_id"`
	ScenarioID   string           `json:"scenario_id"`
	TrackID      string           `json:"track_id"`
	PassportHost string           `json:"passport_host"`
	Scenarios    []quasarScenario `json:"scenarios"`
}

type quasarScenario struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) RunCloudTTS(xToken, deviceID, text, voice string) (ScenarioActionResult, error) {
	if strings.TrimSpace(xToken) == "" {
		return ScenarioActionResult{}, errors.New("unofficial x_token is required")
	}
	if strings.TrimSpace(deviceID) == "" {
		return ScenarioActionResult{}, errors.New("device id is required")
	}
	if strings.TrimSpace(text) == "" {
		return ScenarioActionResult{}, errors.New("text is required")
	}

	session, err := newQuasarSession(c.httpClient.Timeout)
	if err != nil {
		return ScenarioActionResult{}, err
	}
	if err := session.loginWithXToken(strings.TrimSpace(xToken)); err != nil {
		return ScenarioActionResult{}, err
	}

	scenarioID, staleScenarioIDs, err := session.ensureScenario(deviceID)
	if err != nil {
		return ScenarioActionResult{}, err
	}
	voiceUsed, voiceFallback, err := session.updateScenarioTTS(scenarioID, deviceID, text, voice)
	if err != nil {
		return ScenarioActionResult{}, err
	}
	result, err := session.runScenarioAction(scenarioID)
	if err != nil {
		return ScenarioActionResult{}, err
	}
	session.deleteScenarios(staleScenarioIDs)
	quasarScenarioDeletionScheduler.schedule(strings.TrimSpace(xToken), strings.TrimSpace(deviceID), scenarioID, func() {
		cleanupSession, cleanupErr := newQuasarSession(c.httpClient.Timeout)
		if cleanupErr != nil {
			return
		}
		if err := cleanupSession.loginWithXToken(strings.TrimSpace(xToken)); err != nil {
			return
		}
		if err := cleanupSession.deleteScenario(scenarioID); err != nil {
			return
		}
	})
	result.VoiceUsed = voiceUsed
	result.VoiceFallback = voiceFallback
	return result, nil
}

func newQuasarSession(timeout time.Duration) (*quasarSession, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &quasarSession{
		httpClient: &http.Client{
			Timeout: timeout,
			Jar:     jar,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		baseURL: quasarBaseURL,
		authURL: quasarAuthURL,
		pageURL: quasarPageURL,
	}, nil
}

func (s *quasarSession) loginWithXToken(xToken string) error {
	form := url.Values{
		"type":    {"x-token"},
		"retpath": {"https://www.yandex.ru"},
	}

	request, err := http.NewRequest(http.MethodPost, s.authURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Ya-Consumer-Authorization", "OAuth "+xToken)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var payload quasarStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return err
	}
	if payload.Status != "ok" || payload.TrackID == "" || payload.PassportHost == "" {
		if payload.Message != "" {
			return errors.New(payload.Message)
		}
		return errors.New("failed to start Yandex x-token session")
	}

	authURL := strings.TrimRight(payload.PassportHost, "/") + "/auth/session/?track_id=" + url.QueryEscape(payload.TrackID)
	followRequest, err := http.NewRequest(http.MethodGet, authURL, nil)
	if err != nil {
		return err
	}
	followResponse, err := s.httpClient.Do(followRequest)
	if err != nil {
		return err
	}
	defer followResponse.Body.Close()

	if followResponse.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("yandex auth session failed with status %d", followResponse.StatusCode)
	}

	return s.loadCSRFToken()
}

func (s *quasarSession) loadCSRFToken() error {
	request, err := http.NewRequest(http.MethodGet, s.pageURL, nil)
	if err != nil {
		return err
	}
	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	matches := quasarCSRFRegexp.FindSubmatch(body)
	if len(matches) < 2 {
		return errors.New("could not extract Yandex csrf token")
	}
	s.csrfToken = string(matches[1])
	return nil
}

func (s *quasarSession) ensureScenario(deviceID string) (string, []string, error) {
	scenarios, err := s.listScenarios()
	if err != nil {
		return "", nil, err
	}

	name := quasarScenarioName(deviceID)
	staleScenarioIDs := make([]string, 0)
	matchingScenarioID := ""
	for _, scenario := range scenarios {
		if strings.HasPrefix(scenario.Name, quasarScenarioNamePrefix) && scenario.Name != name {
			staleScenarioIDs = append(staleScenarioIDs, scenario.ID)
		}
		if scenario.Name == name {
			matchingScenarioID = scenario.ID
		}
	}
	if matchingScenarioID != "" {
		return matchingScenarioID, staleScenarioIDs, nil
	}

	scenarioID, err := s.createScenario(deviceID)
	return scenarioID, staleScenarioIDs, err
}

func (s *quasarSession) listScenarios() ([]quasarScenario, error) {
	var payload quasarStatusResponse
	if err := s.doJSON(http.MethodGet, s.baseURL+"/m/user/scenarios", nil, &payload); err != nil {
		return nil, err
	}
	if payload.Status != "ok" {
		if payload.Message != "" {
			return nil, errors.New(payload.Message)
		}
		return nil, errors.New("failed to load Yandex scenarios")
	}
	return payload.Scenarios, nil
}

func (s *quasarSession) createScenario(deviceID string) (string, error) {
	body := scenarioSpeakerTTS(quasarScenarioName(deviceID), encodeScenarioTrigger(deviceID), deviceID, "пустышка", "")

	var payload quasarStatusResponse
	if err := s.doJSON(http.MethodPost, s.baseURL+"/m/v4/user/scenarios", body, &payload); err != nil {
		return "", err
	}
	if payload.Status != "ok" || payload.ScenarioID == "" {
		if payload.Message != "" {
			return "", errors.New(payload.Message)
		}
		return "", errors.New("failed to create Yandex TTS scenario")
	}
	return payload.ScenarioID, nil
}

func (s *quasarSession) updateScenarioTTS(scenarioID, deviceID, text, voice string) (string, bool, error) {
	err := s.updateScenarioTTSRequest(scenarioID, deviceID, text, voice)
	if err == nil {
		return strings.TrimSpace(voice), false, nil
	}
	if strings.TrimSpace(voice) == "" {
		return "", false, err
	}

	fallbackErr := s.updateScenarioTTSRequest(scenarioID, deviceID, text, "")
	if fallbackErr != nil {
		return "", false, fmt.Errorf("failed to update Yandex TTS scenario with voice: %w; fallback without voice failed: %v", err, fallbackErr)
	}
	return "", true, nil
}

func (s *quasarSession) updateScenarioTTSRequest(scenarioID, deviceID, text, voice string) error {
	body := scenarioSpeakerTTS(quasarScenarioName(deviceID), encodeScenarioTrigger(deviceID), deviceID, text, voice)
	var payload quasarStatusResponse
	if err := s.doJSON(http.MethodPut, s.baseURL+"/m/v4/user/scenarios/"+scenarioID, body, &payload); err != nil {
		return err
	}
	if payload.Status != "ok" {
		if payload.Message != "" {
			return errors.New(payload.Message)
		}
		return errors.New("failed to update Yandex TTS scenario")
	}
	return nil
}

func (s *quasarSession) runScenarioAction(scenarioID string) (ScenarioActionResult, error) {
	var payload quasarStatusResponse
	if err := s.doJSON(http.MethodPost, s.baseURL+"/m/user/scenarios/"+scenarioID+"/actions", map[string]any{}, &payload); err != nil {
		return ScenarioActionResult{}, err
	}
	if payload.Status != "ok" {
		if payload.Message != "" {
			return ScenarioActionResult{}, errors.New(payload.Message)
		}
		return ScenarioActionResult{}, errors.New("failed to run Yandex TTS scenario")
	}
	return ScenarioActionResult{
		Status:    payload.Status,
		RequestID: coalesceNonEmpty(payload.RequestID, scenarioID),
	}, nil
}

func (s *quasarSession) deleteScenario(scenarioID string) error {
	scenarioID = strings.TrimSpace(scenarioID)
	if scenarioID == "" {
		return nil
	}

	return s.doJSON(http.MethodDelete, s.baseURL+"/m/v4/user/scenarios/"+scenarioID, nil, nil)
}

func (s *quasarSession) deleteScenarios(scenarioIDs []string) {
	for _, scenarioID := range scenarioIDs {
		if err := s.deleteScenario(scenarioID); err != nil {
			continue
		}
	}
}

func (s *quasarSession) doJSON(method, requestURL string, body any, target any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}

	request, err := http.NewRequest(method, requestURL, reader)
	if err != nil {
		return err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if method != http.MethodGet {
		if s.csrfToken == "" {
			return errors.New("csrf token is not initialized")
		}
		request.Header.Set("x-csrf-token", s.csrfToken)
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if target == nil {
		if response.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("quasar request failed with status %d", response.StatusCode)
		}
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return err
	}
	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("quasar request failed with status %d", response.StatusCode)
	}
	return nil
}

func scenarioSpeakerTTS(name, trigger, deviceID, text, voice string) map[string]any {
	ttsValue := map[string]any{
		"text": text,
	}
	if strings.TrimSpace(voice) != "" {
		ttsValue["voice"] = strings.TrimSpace(voice)
	}

	return map[string]any{
		"name": name,
		"icon": "home",
		"triggers": []map[string]any{
			{
				"trigger": map[string]any{
					"type":  "scenario.trigger.voice",
					"value": trigger,
				},
			},
		},
		"steps": []map[string]any{
			{
				"type": "scenarios.steps.actions.v2",
				"parameters": map[string]any{
					"items": []map[string]any{
						{
							"id":   deviceID,
							"type": "step.action.item.device",
							"value": map[string]any{
								"id":        deviceID,
								"item_type": "device",
								"capabilities": []map[string]any{
									{
										"type": "devices.capabilities.quasar",
										"state": map[string]any{
											"instance": "tts",
											"value":    ttsValue,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func quasarScenarioName(deviceID string) string {
	return quasarScenarioNamePrefix + strings.TrimSpace(deviceID)
}

func encodeScenarioTrigger(uid string) string {
	const maskEN = "0123456789abcdef-"
	const maskRU = "оеаинтсрвлкмдпуяы"

	var builder strings.Builder
	builder.Grow(len(uid))
	for _, char := range strings.ToLower(uid) {
		index := strings.IndexRune(maskEN, char)
		if index == -1 {
			continue
		}
		builder.WriteByte(maskRU[index])
	}
	return builder.String()
}

func coalesceNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

type scenarioDeletionScheduler struct {
	mu     sync.Mutex
	timers map[string]*time.Timer
}

func newScenarioDeletionScheduler() *scenarioDeletionScheduler {
	return &scenarioDeletionScheduler{
		timers: make(map[string]*time.Timer),
	}
}

func (s *scenarioDeletionScheduler) schedule(xToken, deviceID, scenarioID string, cleanup func()) {
	key := strings.TrimSpace(xToken) + "|" + strings.TrimSpace(deviceID)
	if key == "|" || strings.TrimSpace(scenarioID) == "" {
		return
	}

	s.mu.Lock()
	if existing := s.timers[key]; existing != nil {
		existing.Stop()
	}
	s.timers[key] = time.AfterFunc(quasarScenarioTTL, func() {
		cleanup()
		s.mu.Lock()
		delete(s.timers, key)
		s.mu.Unlock()
	})
	s.mu.Unlock()
}
