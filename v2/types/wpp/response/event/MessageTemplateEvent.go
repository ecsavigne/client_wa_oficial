package event

import (
	evt_types "github.com/ecsavigne/client_wa_oficial/v2/types/general/response/event/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type MessageTemplateEvent struct {
	*Components
}

func (*MessageTemplateEvent) GetType() wpp.EventType { return wpp.EventTypeTemplateMessage }
func (t *MessageTemplateEvent) String() string       { return response.Val(t) }
func (t *MessageTemplateEvent) GetTypeNotifications() evt_types.TYPE_NOTIFICATION_WEBHOOK {
	return evt_types.ParseTypeNotificationWebhook(t.Entry[0].Changes[0].Field)
}

func (t *MessageTemplateEvent) GetStatus() (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	status = t.Entry[0].Changes[0].Value.Event

	return status
}

func (t *MessageTemplateEvent) GetCategory() (cat string) {
	defer func() {
		if r := recover(); r != nil {
			cat = ""
		}
	}()

	cat = t.Entry[0].Changes[0].Value.NewCategory
	return cat
}

func (t *MessageTemplateEvent) GetName() (name string) {
	defer func() {
		if r := recover(); r != nil {
			name = ""
		}
	}()

	name = t.Entry[0].Changes[0].Value.MessageTemplateName
	return name
}
