package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type ErrorSocketConnectEvent struct {
	*response.Error `json:"error"`
}

func (*ErrorSocketConnectEvent) GetType() wpp.EventType { return wpp.EventTypeErrorSocketConnect }
func (m *ErrorSocketConnectEvent) String() string       { return response.Val(m) }
