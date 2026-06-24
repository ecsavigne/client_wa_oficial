package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type HistoryEvent struct {
	*Components
}

func (*HistoryEvent) GetType() wpp.EventType {
	return wpp.EventTypeHistory
}

func (m *HistoryEvent) String() string { return response.Val(m) }
