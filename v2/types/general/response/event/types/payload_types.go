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

// fields =  (WEBHOOK_NOTIFICATION_[FIELD_NAME])
const (
	// obj whatsapp
	WEBHOOK_NOTIFICATION_UNKNOWN TYPE_NOTIFICATION_WEBHOOK = iota
	/*
		- Se rechaza el aumento del límite de mensajes de todos los números de teléfono de un portfolio comercial, se pospone la decisión sobre el aumento o se necesita más información antes de tomar una decisión.
		- El estado de la cuenta de empresa oficial de un número de teléfono del negocio se aprueba o se rechaza.
		- Se elimina la foto del perfil de un número de teléfono del negocio.
	*/
	WEBHOOK_NOTIFICATION_ACCOUNT_ALERTS
	/*
		- WhatsApp Business account is approved.
		- WhatsApp Business account is rejected.
		- Business account approval is deferred pending further review or information.
	*/
	WEBHOOK_NOTIFICATION_ACCOUNT_REVIEW_UPDATE
	/*
		- Business verification request is approved, rejected, or dismissed.
		- WhatsApp Business account is deleted.
		- Business account is shared with or removed from a partner.
		- Business account violates Meta policies or terms.
		- Business account becomes eligible for international authentication rates.
		- Primary business location is configured.
		- Business account grants partner access to ad accounts.
		- Business account is restricted due to policy violations.
		- Business account accepts WhatsApp API Terms of Service.
		- Customer grants or revokes app permissions for the business account.
		- Business account volume-based pricing tier is updated.
		- Business account is deregistered after device or phone registration change.
		- Business account reconnects after device or phone registration change.
	*/
	WEBHOOK_NOTIFICATION_ACCOUNT_UPDATE
	/*
		- WhatsApp Business account is created.
		- Business account or portfolio feature limits are increased or decreased.
	*/
	WEBHOOK_NOTIFICATION_BUSINESS_CAPABILITY_UPDATE
	/*
		- Solution provider syncs chat history after customer approves sharing.
		- Chat history sync is denied because the customer rejects sharing.
	*/
	WEBHOOK_NOTIFICATION_HISTORY
	WEBHOOK_NOTIFICATION_MESSAGE // message
	/*
		- Template is edited.
	*/
	WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_COMPONENTS_UPDATE
	/*
		- Template quality rating changes.
	*/
	WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_QUALITY_UPDATE
	/*
		- Template is approved.
		- Template is rejected.
		- Template is disabled.
		- Template is archived.
		- Template is unarchived.
	*/
	WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_STATUS_UPDATE
	/*
		- Multi-partner solution is saved as a draft.
		- Multi-partner solution request is sent to a partner.
		- Partner accepts a multi-partner solution request.
		- Partner rejects a multi-partner solution request.
		- Partner requests deactivation of a multi-partner solution.
		- Multi-partner solution is deactivated.
	*/
	WEBHOOK_NOTIFICATION_PARTNER_SOLUTIONS
	/*
		- WhatsApp Business payment settings linked to a payment gateway account.
		- WhatsApp Business payment settings unlinked from a payment gateway account.
		- WhatsApp Business payment settings become active.
	*/
	WEBHOOK_NOTIFICATION_PAYMENT_CONFIGURATION_UPDATE
	/*
		- New business phone display name enters review.
		- Approved business phone display name is edited and reviewed.
	*/
	WEBHOOK_NOTIFICATION_PHONE_NUMBER_NAME_UPDATE
	/*
		- Business phone number messaging volume tier changes.
	*/
	WEBHOOK_NOTIFICATION_PHONE_NUMBER_QUALITY_UPDATE
	/*
		- Meta Business Suite user requests disabling WhatsApp two-step verification.
		- Meta Business Suite user disables two-step verification via email reset process.
		- Meta Business Suite user enables or changes the business phone PIN.
	*/
	WEBHOOK_NOTIFICATION_SECURITY
	/*
		- Solution provider syncs a business customer's WhatsApp Business app contacts.
		- Business App user adds a new contact.
		- Business App user deletes an existing contact.
		- Business App user updates an existing contact.
	*/
	WEBHOOK_NOTIFICATION_SMB_APP_STATE_SYNC
	/*
		- Business App user sends a message to a WhatsApp user or another business.
		- Business App user deletes a previously sent message.
		- Business App user edits a previously sent message.
	*/
	WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES

	/*
		- WhatsApp template category will be changed automatically.
		- WhatsApp template category changed manually or through an automated proces
	*/
	WEBHOOK_NOTIFICATION_TEMPLATE_CATEGORY_UPDATE // template_category_update
	/*
		- WhatsApp user stops marketing messages.
		- WhatsApp user resumes marketing messages.
		- Webhook only triggers for stop/resume actions, not “Interested” or “Not interested” preferences.
	*/
	WEBHOOK_NOTIFICATION_USER_PREFERENCES

	// objeto ig
	WEBHOOK_NOTIFICATION_COMMENTS
	WEBHOOK_NOTIFICATION_LIVE_COMMENTS
	WEBHOOK_NOTIFICATION_MENTIONS
	WEBHOOK_NOTIFICATION_MESSAGE_ECHOES
	WEBHOOK_NOTIFICATION_MESSAGE_REACTIONS
	WEBHOOK_NOTIFICATION_MESSAGE_EDIT
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
	"unknown":                            WEBHOOK_NOTIFICATION_UNKNOWN,
	"account_alerts":                     WEBHOOK_NOTIFICATION_ACCOUNT_ALERTS,
	"account_update":                     WEBHOOK_NOTIFICATION_ACCOUNT_UPDATE,
	"account_review_update":              WEBHOOK_NOTIFICATION_ACCOUNT_REVIEW_UPDATE,
	"business_capability_update":         WEBHOOK_NOTIFICATION_BUSINESS_CAPABILITY_UPDATE,
	"history":                            WEBHOOK_NOTIFICATION_HISTORY,
	"messages":                           WEBHOOK_NOTIFICATION_MESSAGE,
	"message_template_components_update": WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_COMPONENTS_UPDATE,
	"message_template_quality_update":    WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_QUALITY_UPDATE,
	"message_template_status_update":     WEBHOOK_NOTIFICATION_MESSAGE_TEMPLATE_STATUS_UPDATE,
	"partner_solutions":                  WEBHOOK_NOTIFICATION_PARTNER_SOLUTIONS,
	"payment_configuration_update":       WEBHOOK_NOTIFICATION_PAYMENT_CONFIGURATION_UPDATE,
	"phone_number_name_update":           WEBHOOK_NOTIFICATION_PHONE_NUMBER_NAME_UPDATE,
	"phone_number_quality_update":        WEBHOOK_NOTIFICATION_PHONE_NUMBER_QUALITY_UPDATE,
	"security":                           WEBHOOK_NOTIFICATION_SECURITY,
	"smb_app_state_sync":                 WEBHOOK_NOTIFICATION_SMB_APP_STATE_SYNC,
	"smb_message_echoes":                 WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES,
	"template_category_update":           WEBHOOK_NOTIFICATION_TEMPLATE_CATEGORY_UPDATE,
	"user_preferences":                   WEBHOOK_NOTIFICATION_USER_PREFERENCES,
	//field ig and type notification
	"comments":                     WEBHOOK_NOTIFICATION_COMMENTS,
	"live_comments":                WEBHOOK_NOTIFICATION_LIVE_COMMENTS,
	"mentions":                     WEBHOOK_NOTIFICATION_MENTIONS,
	"message_echoes":               WEBHOOK_NOTIFICATION_MESSAGE_ECHOES,
	"message_reactions":            WEBHOOK_NOTIFICATION_MESSAGE_REACTIONS,
	"message_edit":                 WEBHOOK_NOTIFICATION_MESSAGE_EDIT,
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
	return []string{"unknown", "messages", "template_category_update", "message_template_status_update", "comments", "live_comments", "mentions", "message_echoes", "message_reactions", "messaging_handover", "messaging_optins", "messaging_policy_enforcement", "messaging_postbacks", "messaging_referral", "messaging_seen", "response_feedback", "standby", "story_insights", "message_edit"}[t]
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
