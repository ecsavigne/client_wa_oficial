package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type TemplateEventer interface {
	GetName() string
}

// type MessageTemplateEvent struct {
// 	*Components
// 	TemplateEventer
// }

// func (*MessageTemplateEvent) GetType() wpp.EventType { return wpp.EventTypeTemplateMessage }
// func (t *MessageTemplateEvent) String() string       { return response.Val(t) }
// func (t *MessageTemplateEvent) GetName() (name string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			name = ""
// 		}
// 	}()

// 	return t.GetEntry()[0].GetChange()[0].GetValue().MessageTemplateName
// }

type MessageTemplateComponentsUpdateEvent struct {
	*Components
}

var _ TemplateEventer = &MessageTemplateComponentsUpdateEvent{}

func (*MessageTemplateComponentsUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypeMessageTemplateComponentsUpdate
}
func (t *MessageTemplateComponentsUpdateEvent) String() string { return response.Val(t) }
func (t *MessageTemplateComponentsUpdateEvent) GetName() (name string) {
	defer func() {
		if r := recover(); r != nil {
			name = ""
		}
	}()

	return t.GetEntry()[0].GetChange()[0].GetValue().MessageTemplateName
}

type MessageTemplateQualityUpdateEvent struct {
	*Components
}

var _ TemplateEventer = &MessageTemplateQualityUpdateEvent{}

func (*MessageTemplateQualityUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypeMessageTemplateQualityUpdate
}
func (t *MessageTemplateQualityUpdateEvent) String() string { return response.Val(t) }
func (t *MessageTemplateQualityUpdateEvent) GetQualityScore() (quality string) {
	defer func() {
		if r := recover(); r != nil {
			quality = ""
		}
	}()

	return t.GetEntry()[0].GetChange()[0].GetValue().NewQualityScore
}
func (t *MessageTemplateQualityUpdateEvent) GetName() (name string) {
	defer func() {
		if r := recover(); r != nil {
			name = ""
		}
	}()

	return t.GetEntry()[0].GetChange()[0].GetValue().MessageTemplateName
}

type MessageTemplateStatusUpdateEvent struct {
	*Components
}

var _ TemplateEventer = &MessageTemplateStatusUpdateEvent{}

func (*MessageTemplateStatusUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypeMessageTemplateStatusUpdate
}
func (t *MessageTemplateStatusUpdateEvent) String() string { return response.Val(t) }
func (t *MessageTemplateStatusUpdateEvent) GetStatus() (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	return t.GetEntry()[0].GetChange()[0].GetValue().Event
}

func (t *MessageTemplateStatusUpdateEvent) GetName() (name string) {
	defer func() {
		if r := recover(); r != nil {
			name = ""
		}
	}()

	return t.GetEntry()[0].GetChange()[0].GetValue().MessageTemplateName
}

// func (t *MessageTemplateEvent) GetTypeNotifications() evt_types.TYPE_NOTIFICATION_WEBHOOK {
// 	return evt_types.ParseTypeNotificationWebhook(t.GetEntry()[0].GetChange()[0].Field)
// }

type TemplateCategoryUpdateEvent struct {
	*Components
}

var _ TemplateEventer = &TemplateCategoryUpdateEvent{}

func (*TemplateCategoryUpdateEvent) GetType() wpp.EventType {
	return wpp.EventTypeTemplateCategoryUpdate
}
func (t *TemplateCategoryUpdateEvent) String() string { return response.Val(t) }

func (t *TemplateCategoryUpdateEvent) GetCategory() (cat string) {
	defer func() {
		if r := recover(); r != nil {
			cat = ""
		}
	}()

	// return t.Entry[0].Changes[0].Value.NewCategory
	return t.GetEntry()[0].GetChange()[0].GetValue().NewCategory
}

func (t *TemplateCategoryUpdateEvent) GetName() (name string) {
	defer func() {
		if r := recover(); r != nil {
			name = ""
		}
	}()

	// return t.Entry[0].Changes[0].Value.MessageTemplateName
	return t.GetEntry()[0].GetChange()[0].GetValue().MessageTemplateName
}
