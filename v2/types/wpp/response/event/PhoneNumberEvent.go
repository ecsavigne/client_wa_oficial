package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type PhoneNumberNameUpdateEvent struct {
	*Components
}

func (*PhoneNumberNameUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypePhoneNumberNameUpdate
}

func (m *PhoneNumberNameUpdateEvent) String() string { return response.Val(m) }

type PhoneNumberQualityUpdateEvent struct {
	*Components
}

func (*PhoneNumberQualityUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypePhoneNumberQualityUpdate
}

func (m *PhoneNumberQualityUpdateEvent) String() string { return response.Val(m) }
