package transport

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/yandex"
	"errors"
)

type UnofficialTransport struct {
	client *yandex.Client
}

func NewUnofficialTransport(client *yandex.Client) *UnofficialTransport {
	return &UnofficialTransport{client: client}
}

func (t *UnofficialTransport) Name() string {
	return model.TransportUnofficial
}

func (t *UnofficialTransport) Announce(account model.Account, request Request) (Result, error) {
	if request.DeviceID == "" {
		return Result{}, errors.New("device id is required for unofficial Alice transport")
	}
	if request.Text == "" {
		return Result{}, errors.New("text is required for unofficial Alice transport")
	}

	response, err := t.client.RunCloudTTS(account.UnofficialXToken, request.DeviceID, request.Text, request.Voice)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Status:    response.Status,
		RequestID: response.RequestID,
	}, nil
}
