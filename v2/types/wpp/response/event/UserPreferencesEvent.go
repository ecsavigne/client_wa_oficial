package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type UserPreferencesEvent struct {
	*Components
}

func (*UserPreferencesEvent) GetType() wpp.EventType {
	return wpp.EventTypeUserPreferences
}

func (m *UserPreferencesEvent) String() string { return response.Val(m) }
