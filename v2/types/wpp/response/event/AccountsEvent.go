package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type AccountUpdateEvent struct {
	*Components
}

func (*AccountUpdateEvent) GetType() wpp.EventType { return wpp.EventTypeAccountUpdate }

func (m *AccountUpdateEvent) String() string { return response.Val(m) }

type AccountAlertsEvent struct {
	*Components
}

func (*AccountAlertsEvent) GetType() wpp.EventType { return wpp.EventTypeAccountAlerts }

func (m *AccountAlertsEvent) String() string { return response.Val(m) }

type AccountReviewUpdateEvent struct {
	*Components
}

func (*AccountReviewUpdateEvent) GetType() wpp.EventType { return wpp.EventTypeAccountAlerts }

func (m *AccountReviewUpdateEvent) String() string { return response.Val(m) }
