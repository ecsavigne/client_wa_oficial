package types

// types objects of webhook payload of the events of wpp and ig
const (
	TYPE_OBJECT_WEBHOOK_AD_ACCOUNT                = "ad account"
	TYPE_OBJECT_WEBHOOK_APPLICATION               = "application"
	TYPE_OBJECT_WEBHOOK_CATALOG                   = "catalog"
	TYPE_OBJECT_WEBHOOK_INSTAGRAM                 = "instagram"
	TYPE_OBJECT_WEBHOOK_MANAGED_META_ACCOUNT      = "managed meta account"
	TYPE_OBJECT_WEBHOOK_PAGE                      = "page"
	TYPE_OBJECT_WEBHOOK_PERMISSIONS               = "permissions"
	TYPE_OBJECT_WEBHOOK_USER                      = "user"
	TYPE_OBJECT_WEBHOOK_WHATSAPP_BUSINESS_ACCOUNT = "whatsapp business account"
)

// type_notification_webhook (is value of the field in payload of the webhook)
type TYPE_NOTIFICATION_WEBHOOK int

const (
	WEBHOOK_NOTIFICATION_UNKNOWN                  TYPE_NOTIFICATION_WEBHOOK = iota
	WEBHOOK_NOTIFICATION_MESSAGE                                            // message
	WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_CATEGORY                           // template_category_update
	WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_STATUS
	// objeto ig
	WEBHOOK_NOTIFICATION_COMMENTS
	WEBHOOK_NOTIFICATION_LIVE_COMMENTS
	WEBHOOK_NOTIFICATION_MENTIONS
	WEBHOOK_NOTIFICATION_MESSAGE_ECHOES
	WEBHOOK_NOTIFICATION_MESSAGE_REACTIONS
	WEBHOOK_NOTIFICATION_MESSAGING_HANDOVER
	WEBHOOK_NOTIFICATION_MESSAGING_OPTINS
	WEBHOOK_NOTIFICATION_MESSAGING_POLICY_ENFORCEMENT
	WEBHOOK_NOTIFICATION_MESSAGING_POSTBACKS
	WEBHOOK_NOTIFICATION_MESSAGING_REFERRAL
	WEBHOOK_NOTIFICATION_MESSAGING_SEEN
	WEBHOOK_NOTIFICATION_RESPONSE_FEEDBACK
	WEBHOOK_NOTIFICATION_STANDBY
	WEBHOOK_NOTIFICATION_STORY_INSIGHTS
)

var type_NOTIFICATION_WEBHOOK = map[string]TYPE_NOTIFICATION_WEBHOOK{
	// fields Wpp
	"unknown":                        WEBHOOK_NOTIFICATION_UNKNOWN,
	"messages":                       WEBHOOK_NOTIFICATION_MESSAGE,
	"template_category_update":       WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_CATEGORY,
	"message_template_status_update": WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_STATUS,
	//field ig
	"comments":                     WEBHOOK_NOTIFICATION_COMMENTS,
	"live_comments":                WEBHOOK_NOTIFICATION_LIVE_COMMENTS,
	"mentions":                     WEBHOOK_NOTIFICATION_MENTIONS,
	"message_echoes":               WEBHOOK_NOTIFICATION_MESSAGE_ECHOES,
	"message_reactions":            WEBHOOK_NOTIFICATION_MESSAGE_REACTIONS,
	"messaging_handover":           WEBHOOK_NOTIFICATION_MESSAGING_HANDOVER,
	"messaging_optins":             WEBHOOK_NOTIFICATION_MESSAGING_OPTINS,
	"messaging_policy_enforcement": WEBHOOK_NOTIFICATION_MESSAGING_POLICY_ENFORCEMENT,
	"messaging_postbacks":          WEBHOOK_NOTIFICATION_MESSAGING_POSTBACKS,
	"messaging_referral":           WEBHOOK_NOTIFICATION_MESSAGING_REFERRAL,
	"messaging_seen":               WEBHOOK_NOTIFICATION_MESSAGING_SEEN,
	"response_feedback":            WEBHOOK_NOTIFICATION_RESPONSE_FEEDBACK,
	"standby":                      WEBHOOK_NOTIFICATION_STANDBY,
	"story_insights":               WEBHOOK_NOTIFICATION_STORY_INSIGHTS,
}

func (t TYPE_NOTIFICATION_WEBHOOK) Enum() string {
	return []string{"unknown", "messages", "template_category_update", "message_template_status_update", "comments", "live_comments", "mentions", "message_echoes", "message_reactions", "messaging_handover", "messaging_optins", "messaging_policy_enforcement", "messaging_postbacks", "messaging_referral", "messaging_seen", "response_feedback", "standby", "story_insights"}[t]
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
