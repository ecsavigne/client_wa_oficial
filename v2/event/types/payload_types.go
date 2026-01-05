package types

// type_notification_webhook
type TYPE_NOTIFICATION_WEBHOOK int

const (
	WEBHOOK_NOTIFICATION_UNKNOWN                  TYPE_NOTIFICATION_WEBHOOK = iota
	WEBHOOK_NOTIFICATION_MESSAGE                                            // message
	WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_CATEGORY                           // template_category_update
	WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_STATUS                             // message_template_status_update
)

var type_NOTIFICATION_WEBHOOK = map[string]TYPE_NOTIFICATION_WEBHOOK{
	"unknown":                        WEBHOOK_NOTIFICATION_UNKNOWN,
	"message":                        WEBHOOK_NOTIFICATION_MESSAGE,
	"template_category_update":       WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_CATEGORY,
	"message_template_status_update": WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_STATUS,
}

func (t TYPE_NOTIFICATION_WEBHOOK) Enum() string {
	return []string{"unknown", "message", "template_category_update", "message_template_status_update"}[t]
}

func (t TYPE_NOTIFICATION_WEBHOOK) String() string {
	return t.Enum()
}

// ParseTypeNotificationWebhook parses a given string and returns a TYPE_NOTIFICATION_WEBHOOK.
// If the string does not match any of the predefined types, it returns WEBHOOK_NOTIFICATION_UNKNOWN.
func ParseTypeNotificationWebhook(str string) TYPE_NOTIFICATION_WEBHOOK {
	if v, ok := type_NOTIFICATION_WEBHOOK[str]; ok {
		return v
	}

	return WEBHOOK_NOTIFICATION_UNKNOWN
}
