package event

import "github.com/ecsavigne/client_wa_oficial/types"

type EventInterface interface {
	GetType() types.EventType
}
