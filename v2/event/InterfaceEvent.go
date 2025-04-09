package event

import "github.com/ecsavigne/client_wa_oficial/v2/types"

type EventInterface interface {
	GetType() types.EventType
}
