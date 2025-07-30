package types

type TYPE_BUTTON = string

const (
	TB_QUICK_REPLY          TYPE_BUTTON = "QUICK_REPLY"
	TB_URL                  TYPE_BUTTON = "URL"
	TB_PHONE_NUMBER         TYPE_BUTTON = "PHONE_NUMBER"
	TB_OTP                  TYPE_BUTTON = "OTP"
	TB_MPM                  TYPE_BUTTON = "MPM"
	TB_SPM                  TYPE_BUTTON = "SPM"
	TB_CATALOG              TYPE_BUTTON = "CATALOG"
	TB_FLOW                 TYPE_BUTTON = "FLOW"
	TB_VOICE_CALL           TYPE_BUTTON = "VOICE_CALL"
	TB_APP                  TYPE_BUTTON = "APP"
	TB_POSTBACK             TYPE_BUTTON = "POSTBACK"
	TB_BOOKING_CONFIRMATION TYPE_BUTTON = "BOOKING_CONFIRMATION"
)

type Url_Button struct {
	// Admite one var max, ej: "https://www.luckyshrub.com/shop/", ej2: https://www.luckyshrub.com/shop?promo={{1}}. max size = 2.000 chars
	Url string `json:"url" validate:"required"`
	// Required if <URL> contains a variable
	Example []string `json:"example,omitempty"`
}

type Phone_Number_Button struct {
	// less than 20 chars. ej: "+13057652345"
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type Copy_Code_Button struct {
	// 15 chars maximum. String that will be copied in Clipboard
	Example string `json:"example,omitempty"`
}

type ICON_TYPE = string

const (
	DOCUMENT_I  ICON_TYPE = "DOCUMENT"
	PROMOTION_I ICON_TYPE = "PROMOTION"
	REVIEW_I    ICON_TYPE = "REVIEW"
)

type Flow_Button struct {
	FlowId uint64 `json:"flow_id" validate:"required"`
	// name of proccess only allowed on Cloud Api
	FlowName string `json:"flow_name,omitempty" validate:"required"`
	//  string in format json ex: "{\"version\": \"3.1\", \"screens\": [...]}"
	FlowJson string `json:"flow_json" validate:"required"`
	// Allowed values: "navigate", "data_exchange". Use "navigate" for define first screen with part of template message.
	FlowAction string `json:"flow_action,omitempty"`
	// Optional only if <flow_action> is "navigate". The id of the first screen of process
	NavigateScreen string `json:"navigate_screen,omitempty"`
	//  Allowed values: "DOCUMENT", "PROMOTION", "REVIEW". Default "PROMOTION"
	Icon ICON_TYPE `json:"icon,omitempty"`
}

type Button struct {
	Type TYPE_BUTTON `json:"type" validate:"required"`
	// 25 characters maximum
	Text                 string `json:"text" validate:"required"`
	*Url_Button          `json:",omitempty"`
	*Phone_Number_Button `json:",omitempty"`
	*Copy_Code_Button    `json:",omitempty"`
	*Flow_Button         `json:",omitempty"`
	EndpointUri          string `json:"endpoint_uri,omitempty"`
	ZeroTapTermsAccepted bool   `json:"zero_tap_terms_accepted,omitempty"`
	TtlMinutes           int64  `json:"ttl_minutes,omitempty"`
	*SupportedApp        `json:"supported_apps,omitempty"`
}
