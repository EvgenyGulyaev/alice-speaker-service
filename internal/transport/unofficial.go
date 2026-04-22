package transport

import (
	"aliceSpeakerService/internal/model"
	"errors"
)

type UnofficialTransportStub struct{}

func NewUnofficialTransportStub() *UnofficialTransportStub {
	return &UnofficialTransportStub{}
}

func (t *UnofficialTransportStub) Name() string {
	return model.TransportUnofficial
}

func (t *UnofficialTransportStub) Announce(account model.Account, request Request) (Result, error) {
	return Result{}, errors.New("unofficial alice transport is not configured yet")
}
