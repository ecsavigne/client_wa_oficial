package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type EventErrorSocketConnect struct {
	response.Error `json:"error"`
}

func (*EventErrorSocketConnect) GetType() types.EventType { return types.EventTypeErrorSocketConnect }
