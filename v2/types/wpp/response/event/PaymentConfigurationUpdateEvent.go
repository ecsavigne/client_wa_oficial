package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type PaymentConfigurationUpdateEvent struct {
	*Components
}

func (*PaymentConfigurationUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypePaymentConfigurationUpdate
}

func (m *PaymentConfigurationUpdateEvent) String() string { return response.Val(m) }
