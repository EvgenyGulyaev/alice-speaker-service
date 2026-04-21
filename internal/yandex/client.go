package yandex

import (
	"aliceSpeakerService/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const smartHomeBaseURL = "https://api.iot.yandex.net/v1.0"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type ScenarioActionResult struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    smartHomeBaseURL,
	}
}

func (c *Client) LoadResources(token string) (model.Resources, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseURL+"/user/info", nil)
	if err != nil {
		return model.Resources{}, err
	}
	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return model.Resources{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return model.Resources{}, fmt.Errorf("yandex user info failed with status %d", response.StatusCode)
	}

	var payload userInfoResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return model.Resources{}, err
	}

	result := model.Resources{
		Rooms:     make([]model.Room, 0, len(payload.Rooms)),
		Devices:   make([]model.Device, 0, len(payload.Devices)),
		Scenarios: make([]model.Scenario, 0, len(payload.Scenarios)),
	}

	for _, room := range payload.Rooms {
		result.Rooms = append(result.Rooms, model.Room{
			ID:          room.ID,
			Name:        room.Name,
			HouseholdID: room.HouseholdID,
		})
	}
	for _, device := range payload.Devices {
		if !isSupportedSpeakerType(device.Type) {
			continue
		}
		result.Devices = append(result.Devices, model.Device{
			ID:          device.ID,
			Name:        device.Name,
			Type:        device.Type,
			RoomID:      device.Room,
			HouseholdID: device.HouseholdID,
		})
	}
	for _, scenario := range payload.Scenarios {
		result.Scenarios = append(result.Scenarios, model.Scenario{
			ID:       scenario.ID,
			Name:     scenario.Name,
			IsActive: scenario.IsActive,
		})
	}

	return result, nil
}

func (c *Client) RunScenario(token, scenarioID string) (ScenarioActionResult, error) {
	if strings.TrimSpace(scenarioID) == "" {
		return ScenarioActionResult{}, errors.New("scenario id is required")
	}

	request, err := http.NewRequest(http.MethodPost, c.baseURL+"/scenarios/"+scenarioID+"/actions", bytes.NewReader([]byte("{}")))
	if err != nil {
		return ScenarioActionResult{}, err
	}
	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return ScenarioActionResult{}, err
	}
	defer response.Body.Close()

	var payload scenarioActionResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return ScenarioActionResult{}, err
	}
	if response.StatusCode >= http.StatusBadRequest {
		if payload.Message != "" {
			return ScenarioActionResult{}, errors.New(payload.Message)
		}
		return ScenarioActionResult{}, fmt.Errorf("scenario run failed with status %d", response.StatusCode)
	}

	return ScenarioActionResult{
		Status:    payload.Status,
		RequestID: payload.RequestID,
	}, nil
}

func isSupportedSpeakerType(deviceType string) bool {
	return strings.HasPrefix(deviceType, "devices.types.smart_speaker.") || deviceType == "devices.types.media_device.tv_box"
}
