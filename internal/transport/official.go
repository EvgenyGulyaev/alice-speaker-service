package transport

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/yandex"
)

type OfficialScenarioTransport struct {
	client *yandex.Client
}

func NewOfficialScenarioTransport(client *yandex.Client) *OfficialScenarioTransport {
	return &OfficialScenarioTransport{client: client}
}

func (t *OfficialScenarioTransport) Name() string {
	return model.TransportOfficial
}

func (t *OfficialScenarioTransport) Announce(account model.Account, request Request) (Result, error) {
	response, err := t.client.RunScenario(account.OAuthToken, request.ScenarioID)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Status:    response.Status,
		RequestID: response.RequestID,
	}, nil
}
