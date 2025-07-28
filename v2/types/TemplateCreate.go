package types

type REJECTED_REASON = string

const (
	ABUSIVE_CONTENT      REJECTED_REASON = "ABUSIVE_CONTENT"
	INVALID_FORMAT       REJECTED_REASON = "INVALID_FORMAT"
	NONE                 REJECTED_REASON = "NONE"
	PROMOTIONAL          REJECTED_REASON = "PROMOTIONAL"
	TAG_CONTENT_MISMATCH REJECTED_REASON = "TAG_CONTENT_MISMATCH"
	SCAM                 REJECTED_REASON = "SCAM"
)

type CATEGORY = string

const (
	AUTHENTICATION CATEGORY = "AUTHENTICATION"
	MARKETING      CATEGORY = "MARKETING"
	UTILITY        CATEGORY = "UTILITY"
)

type SUB_CATEGORY = string

const (
	ORDER_DETAILS SUB_CATEGORY = "ORDER_DETAILS"
	ORDER_STATUS  SUB_CATEGORY = "ORDER_STATUS"
)

const (
	HEADER CATEGORY = "header"
	BODY   CATEGORY = "body"
	FOOTER CATEGORY = "footer"
)

type PARAMETER_FORMAT = string

const (
	NAMED      PARAMETER_FORMAT = "NAMED"
	POSITIONAL PARAMETER_FORMAT = "POSITIONAL"
)

type STATUS = string

const (
	ACTIVE   STATUS = "ACTIVE"
	INACTIVE STATUS = "INACTIVE"
)

type OTP_TYPE = string

const (
	COPY_CODE  OTP_TYPE = "COPY_CODE"
	ONE_TAP    OTP_TYPE = "ONE_TAP"
	ZERO_TAP   OTP_TYPE = "ZERO_TAP"
	NO_BUTTONS OTP_TYPE = "NO_BUTTONS"
)

type ArrayButton = []Button

type PositionalParams = []string
type NamedParam struct {
	ParamName string `json:"param_name" validate:"required"`
	Example   string `json:"example" validate:"required"`
}

type HeaderText = []PositionalParams
type HeaderHandle = []string
type HeaderTextNamedParam = []NamedParam
type BodyTextNamedParam = []NamedParam
type BodyText = []PositionalParams

type Example struct {
	*HeaderHandle         `json:"header_handle,omitempty"`
	*HeaderText           `json:"header_text,omitempty"`
	*BodyText             `json:"body_text,omitempty"`
	*BodyTextNamedParam   `json:"body_text_named_params,omitempty"`
	*HeaderTextNamedParam `json:"header_text_named_params,omitempty"`
}

type Header struct {
	// type of assets of media content. Configurable to "IMAGE", "VIDEO" o "DOCUMENT"
	Format string `json:"format" validate:"required"`
}

type MockupComponent struct {
	Type string `json:"type" validate:"required"`
	// accept parameters for programming ex:
	// 		Posicional: se incluye una matriz de parámetros posicionales numerados que corresponden a posiciones numéricas en el texto del cuerpo con ejemplos.
	// 		Por ejemplo: “Hello {{1}}, your account balance is {{2}}” | [ “John”, “$1,000” ]
	// 		Nominal: se incluyen objetos JSON que contengan un parámetro con nombre y ejemplos.
	// Footer not accept parameters
	// Por ejemplo: { "param_name": "order_id", "example": "335628"}
	// less than 60 characters in header and footer, 1024 characters in body,
	Text         string `json:"text,omitempty"`
	*Example     `json:"example,omitempty"`
	*Header      `json:",omitempty"`
	*ArrayButton `json:"buttons,omitempty"`
}

// Optional data during creation of a template from a library template. These are optional fields for the body component.
type LibraryTemplateBodyInputs struct {
	AddContactNumber          bool  `json:"add_contact_number,omitempty"`
	AddLearnMoreLink          bool  `json:"add_learn_more_link,omitempty"`
	AddSecurityRecommendation bool  `json:"add_security_recommendation,omitempty"`
	AddTrackPackageLink       bool  `json:"add_track_package_link,omitempty"`
	CodeExpirationMinutes     int64 `json:"code_expiration_minutes,omitempty"`
}

type URL_ struct {
	BaseUrl          string `json:"base_url" validate:"required"`
	UrlSuffixExample string `json:"url_suffix_example,omitempty"`
}

type SupportedApp struct {
	PackageName   string `json:"package_name" validate:"required"`
	SignatureHash string `json:"signature_hash" validate:"required"`
}

// Optional data during creation of a template from a library template. These are optional fields for the button component.
type LibraryTemplateButtonInputs struct {
	Type                 TYPE_BUTTON `json:"type" validate:"required"`
	PhoneNumber          string      `json:"phone_number,omitempty"`
	*URL_                `json:"url,omitempty"`
	OtpType              OTP_TYPE `json:"otp_type,omitempty"`
	ZeroTapTermsAccepted bool     `json:"zero_tap_terms_accepted,omitempty"`
	*SupportedApp        `json:"supported_apps,omitempty"`
}

// All template have limit of one body component
type MockupTemplate struct {
	// less than 512 characters
	ID                           string           `json:"id,omitempty"`
	Name                         string           `json:"name" validate:"required"`
	Category                     CATEGORY         `json:"category" validate:"required"`
	CorrectCategory              CATEGORY         `json:"correct_category,omitempty"`
	PreviousCategory             CATEGORY         `json:"previous_category,omitempty"`
	SubCategory                  SUB_CATEGORY     `json:"sub_category,omitempty"`
	CtaUrlLinkTrackingOptedOut   bool             `json:"cta_url_link_tracking_opted_out,omitempty"`
	ParameterFormat              PARAMETER_FORMAT `json:"parameter_format,omitempty"` // The parameter format, can be Named or Positional
	RejectedReason               REJECTED_REASON  `json:"rejected_reason,omitempty"`
	Status                       STATUS           `json:"status,omitempty"`
	Language                     string           `json:"language" validate:"required"`
	AllowCategoryChange          bool             `json:"allow_category_change,omitempty"`
	*LibraryTemplateBodyInputs   `json:"library_template_body_inputs,omitempty"`
	*LibraryTemplateButtonInputs `json:"library_template_button_inputs,omitempty"`
	LibraryTemplateName          string            `json:"library_template_name,omitempty"`
	MessageSendTtlSeconds        int64             `json:"message_send_ttl_seconds,omitempty"`
	Components                   []MockupComponent `json:"components,omitempty"`
}
