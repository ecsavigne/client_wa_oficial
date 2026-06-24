package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type BusinessCapabilityUpdateEvent struct {
	*Components
}

func (*BusinessCapabilityUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypeBusinessCapabilityUpdate
}

func (m *BusinessCapabilityUpdateEvent) String() string { return response.Val(m) }
