package event

import "github.com/ecsavigne/client_wa_oficial/types"

type EventErrorSocketConnect struct {
	types.Error `json:"error"`
}

func (*EventErrorSocketConnect) GetType() types.EventType { return types.EventTypeErrorSocketConnect }
