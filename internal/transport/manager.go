package transport

import (
	"aliceSpeakerService/internal/model"
	"fmt"
)

type Manager struct {
	transports map[string]SpeakerTransport
}

func NewManager(items ...SpeakerTransport) *Manager {
	transports := make(map[string]SpeakerTransport, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		transports[item.Name()] = item
	}
	return &Manager{transports: transports}
}

func (m *Manager) Announce(account model.Account, request Request) (Result, error) {
	name := model.NormalizeTransport(account.Transport)
	transport, ok := m.transports[name]
	if !ok {
		return Result{}, fmt.Errorf("transport %q is not available", name)
	}
	return transport.Announce(account, request)
}
