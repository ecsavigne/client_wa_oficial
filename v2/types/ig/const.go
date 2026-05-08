package ig

// types attachment
type IG_ATTACHMENT_TYPE string

const (
	IG_ATTACHMENT_TYPE_AUDIO       IG_ATTACHMENT_TYPE = "audio"
	IG_ATTACHMENT_TYPE_IMAGE       IG_ATTACHMENT_TYPE = "image"
	IG_ATTACHMENT_TYPE_VIDEO       IG_ATTACHMENT_TYPE = "video"
	IG_ATTACHMENT_TYPE_FILE        IG_ATTACHMENT_TYPE = "file"
	IG_ATTACHMENT_TYPE_LIKE_HEART  IG_ATTACHMENT_TYPE = "like_heart"
	IG_ATTACHMENT_TYPE_MEDIA_SHARE IG_ATTACHMENT_TYPE = "media_share"
	IG_ATTACHMENT_TYPE_TEMPLATE    IG_ATTACHMENT_TYPE = "template"
	// IG_REACTION_MESSAGE_TYPE = "reaction"
	// IG_BUTTON_MESSAGE_TYPE    = "button"
	// IG_CONTACT_MESSAGE_TYPE   = "contact"
	// IG_ORDER_MESSAGE_TYPE     = "order"
	// IG_SYSTEM_MESSAGE_TYPE    = "system"
	// IG_UNKNOWN_MESSAGE_TYPE   = "unknown"
	// IG_LOCATION_MESSAGE_TYPE  = "location"
)

func (self IG_ATTACHMENT_TYPE) String() string {
	return string(self)
}

// type get info
type IG_GET_INFO_TYPE string

const (
	IG_GET_INFO_ACCOUNT_BUSINESS      IG_GET_INFO_TYPE = "account_business"
	IG_GET_INFO_LINK                  IG_GET_INFO_TYPE = "ig_link"
	IG_GET_INFO_PERSISTENT_MENU       IG_GET_INFO_TYPE = "persistent_menu"
	IG_GET_INFO_ICE_BREAKERS          IG_GET_INFO_TYPE = "ice_breakers"
	IG_GET_INFO_WELCOME_MESSAGE_FLOWS IG_GET_INFO_TYPE = "welcome_message_flows"
	IG_GET_COMMENT                    IG_GET_INFO_TYPE = "comment"
	IG_GET_REPLIES_COMMENTS           IG_GET_INFO_TYPE = "replies_comments"
	IG_GET_SUBSCRIBE_WEBHOOK_FIELD    IG_GET_INFO_TYPE = "subscribe_webhook_field"
	IG_GET_METRICS_MEDIA              IG_GET_INFO_TYPE = "metrics_media"
	IG_GET_METRICS_MEDIA_INSIGHT      IG_GET_INFO_TYPE = "metrics_media_insight"
	IG_GET_METRICS_USER_INSIGHT       IG_GET_INFO_TYPE = "metrics_user_insight"
	IG_GET_LIST_CONVERSATION          IG_GET_INFO_TYPE = "list_conversation"
	IG_GET_MESSAGES_CONVERSATION      IG_GET_INFO_TYPE = "messages_conversation"
	IG_GET_USER_CONVERSATION          IG_GET_INFO_TYPE = "user_conversation"
	IG_GET_INFO_MESSAGE               IG_GET_INFO_TYPE = "info_message"
)

// type delete
type IG_DELETE_TYPE string

const (
	IG_DELETE_MESSAGE               IG_DELETE_TYPE = "message"
	IG_DELETE_PERSISTENT_MENU       IG_DELETE_TYPE = "persistent_menu"
	IG_DELETE_ICE_BREAKERS          IG_DELETE_TYPE = "ice_breakers"
	IG_DELETE_WELCOME_MESSAGE_FLOWS IG_DELETE_TYPE = "welcome_message_flows"
	IG_DELETE_COMMENT               IG_DELETE_TYPE = "comment"
)

// type ig template
type IG_TEMPLATE_TYPE string

const (
	IG_TEMPLATE_BUTTON  IG_TEMPLATE_TYPE = "button"
	IG_TEMPLATE_GENERIC IG_TEMPLATE_TYPE = "generic"
)

func (self IG_TEMPLATE_TYPE) String() string {
	return string(self)
}

// type update
const (
	IG_UPDATE_WELCOME_MESSAGE_FLOWS = "welcome_message_flows"
)

// type create
type IG_CREATE_TYPE string

const (
	IG_CREATE_POST           IG_CREATE_TYPE = "post"
	IG_CREATE_STORY          IG_CREATE_TYPE = "story"
	IG_CREATE_COMMENT        IG_CREATE_TYPE = "comment"
	IG_CREATE_REPLY_COMMENT  IG_CREATE_TYPE = "reply_comment"
	IG_CREATE_HIDE_COMMENT   IG_CREATE_TYPE = "hide_comment"
	IG_CREATE_ENABLE_COMMENT IG_CREATE_TYPE = "enable_comment"
)

// MediaType Post can be "STORIES", "REELS", "CAROUSEL"
type IG_MEDIA_TYPE string

const (
	IG_MEDIA_TYPE_STORIES  IG_MEDIA_TYPE = "STORIES"
	IG_MEDIA_TYPE_REELS    IG_MEDIA_TYPE = "REELS"
	IG_MEDIA_TYPE_CAROUSEL IG_MEDIA_TYPE = "CAROUSEL"
	IG_MEDIA_TYPE_POST     IG_MEDIA_TYPE = ""
)

func (self IG_MEDIA_TYPE) String() string {
	return string(self)
}

type IG_TYPE_BUTTON string

const (
	IG_TYPE_BUTTON_POSTBACK IG_TYPE_BUTTON = "postback"
	IG_TYPE_BUTTON_WEB_URL  IG_TYPE_BUTTON = "web_url"
)

func (self IG_TYPE_BUTTON) String() string {
	return string(self)
}
