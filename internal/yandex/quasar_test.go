package yandex

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScenarioSpeakerTTSIncludesVoiceWhenProvided(t *testing.T) {
	payload := scenarioSpeakerTTS("name", "trigger", "device", "hello", "oksana")

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if !strings.Contains(string(raw), `"voice":"oksana"`) {
		t.Fatalf("expected payload to include voice, got %s", raw)
	}
}

func TestUpdateScenarioTTSFallsBackWithoutVoice(t *testing.T) {
	var bodies []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/m/v4/user/scenarios/scenario-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		bodies = append(bodies, string(body))

		if len(bodies) == 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"status":"error","message":"voice is not supported"}`))
			return
		}

		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	session := &quasarSession{
		httpClient: server.Client(),
		csrfToken:  "csrf",
		baseURL:    server.URL,
	}

	if err := session.updateScenarioTTS("scenario-1", "device-1", "hello", "oksana"); err != nil {
		t.Fatalf("expected fallback to succeed, got %v", err)
	}

	if len(bodies) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(bodies))
	}
	if !strings.Contains(bodies[0], `"voice":"oksana"`) {
		t.Fatalf("expected first request to include voice, got %s", bodies[0])
	}
	if strings.Contains(bodies[1], `"voice":"oksana"`) {
		t.Fatalf("expected fallback request without voice, got %s", bodies[1])
	}
}
