package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type SecurityEvent struct {
	*Components
}

func (*SecurityEvent) GetType() wpp.EventType {
	return wpp.EventTypeSegurity
}

func (m *SecurityEvent) String() string { return response.Val(m) }
