package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type ErrorSocketConnectEvent struct {
	*response.Error `json:"error"`
}

func (*ErrorSocketConnectEvent) GetType() types.EventType { return types.EventTypeErrorSocketConnect }
