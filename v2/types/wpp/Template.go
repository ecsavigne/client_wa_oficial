package wpp

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	str "strings"

	"google.golang.org/protobuf/types/known/structpb"
)

// type param template pre-approval
type BODY_PARAM_TPL_LIBRARY_TYPE = string

const (
	PTLT_ADDRESS      BODY_PARAM_TPL_LIBRARY_TYPE = "ADDRESS"
	PTLT_TEXT         BODY_PARAM_TPL_LIBRARY_TYPE = "TEXT"
	PTLT_AMOUNT       BODY_PARAM_TPL_LIBRARY_TYPE = "AMOUNT"
	PTLT_DATE         BODY_PARAM_TPL_LIBRARY_TYPE = "DATE"
	PTLT_PHONE_NUMBER BODY_PARAM_TPL_LIBRARY_TYPE = "PHONE NUMBER"
	PTLT_EMAIL        BODY_PARAM_TPL_LIBRARY_TYPE = "EMAIL"
	PTLT_NUMBER       BODY_PARAM_TPL_LIBRARY_TYPE = "NUMBER"
)

// Industry
type INDUSTRY_TYPE = string

const (
	IT_E_COMMERCE         INDUSTRY_TYPE = "E_COMMERCE"
	IT_FINANCIAL_SERVICES INDUSTRY_TYPE = "FINANCIAL_SERVICES"
	IT_TELECOMMUNICATION  INDUSTRY_TYPE = "TELECOMMUNICATION"
)

// Topic
type TOPIC_TYPE = string

const (
	TT_ACOUNT_UPDATE         TOPIC_TYPE = "ACCOUNT_UPDATE"
	TT_ACCOUNT_UPDATES       TOPIC_TYPE = "ACCOUNT_UPDATES"
	TT_CALL_PERMISSIONS      TOPIC_TYPE = "CALL_PERMISSIONS"
	TT_EVENT_REMINDER        TOPIC_TYPE = "EVENT_REMINDER"
	TT_GROUP_INVITE_LINK     TOPIC_TYPE = "GROUP_INVITE_LINK"
	TT_IDENTITY_VERIFICATION TOPIC_TYPE = "IDENTITY_VERIFICATION"
	TT_CUSTOMER_FEEDBACK     TOPIC_TYPE = "CUSTOMER_FEEDBACK"
	TT_ORDER_MANAGEMENT      TOPIC_TYPE = "ORDER_MANAGEMENT"
	TT_PAYMENTS              TOPIC_TYPE = "PAYMENTS"
)

// Use Case
type USECASE_TYPE = string

const (
	UCT_ACCOUNT_CREATION_CONFIRMATION  USECASE_TYPE = "ACCOUNT_CREATION_CONFIRMATION"
	UCT_PAYMENT_DUE_REMINDER           USECASE_TYPE = "PAYMENT_DUE_REMINDER"
	UCT_FEEDBACK_SURVEY                USECASE_TYPE = "FEEDBACK_SURVEY"
	UCT_PAYMENT_ACTION_REQUIRED        USECASE_TYPE = "PAYMENT_ACTION_REQUIRED"
	UCT_SHIPMENT_CONFIRMATION          USECASE_TYPE = "SHIPMENT_CONFIRMATION"
	UCT_PAYMENT_OVERDUE                USECASE_TYPE = "PAYMENT_OVERDUE"
	UCT_DELIVERY_UPDATE                USECASE_TYPE = "DELIVERY_UPDATE"
	UCT_PAYMENT_CONFIRMATION           USECASE_TYPE = "PAYMENT_CONFIRMATION"
	UCT_ORDER_DELAY                    USECASE_TYPE = "ORDER_DELAY"
	UCT_FRAUD_ALERT                    USECASE_TYPE = "FRAUD_ALERT"
	UCT_DELIVERY_FAILED                USECASE_TYPE = "DELIVERY_FAILED"
	UCT_AUTO_PAY_REMINDER              USECASE_TYPE = "AUTO_PAY_REMINDER"
	UCT_DELIVERY_CONFIRMATION          USECASE_TYPE = "DELIVERY_CONFIRMATION"
	UCT_PAYMENT_SCHEDULED              USECASE_TYPE = "PAYMENT_SCHEDULED"
	UCT_ORDER_PICK_UP                  USECASE_TYPE = "ORDER_PICK_UP"
	UCT_PAYMENT_REJECT_FAIL            USECASE_TYPE = "PAYMENT_REJECT_FAIL"
	UCT_ORDER_ACTION_NEEDED            USECASE_TYPE = "ORDER_ACTION_NEEDED"
	UCT_STATEMENT_AVAILABLE            USECASE_TYPE = "STATEMENT_AVAILABLE"
	UCT_ORDER_CONFIRMATION             USECASE_TYPE = "ORDER_CONFIRMATION"
	UCT_LOW_BALANCE_WARNING            USECASE_TYPE = "LOW_BALANCE_WARNING"
	UCT_ORDER_OR_TRANSACTION_CANCEL    USECASE_TYPE = "ORDER_OR_TRANSACTION_CANCEL"
	UCT_RECEIPT_ATTACHMENT             USECASE_TYPE = "RECEIPT_ATTACHMENT"
	UCT_RETURN_CONFIRMATION            USECASE_TYPE = "RETURN_CONFIRMATION"
	UCT_STATEMENT_ATTACHMENT           USECASE_TYPE = "STATEMENT_ATTACHMENT"
	UCT_TRANSACTION_ALERT              USECASE_TYPE = "TRANSACTION_ALERT"
	UCT_APPOINTMENT_MISSED_CALLS       USECASE_TYPE = "APPOINTMENT_MISSED_CALLS"
	UCT_APPOINTMENT_REMINDER           USECASE_TYPE = "APPOINTMENT_REMINDER"
	UCT_APPOINTMENT_SCHEDULING_REQUEST USECASE_TYPE = "APPOINTMENT_SCHEDULING_REQUEST"
	UCT_EVENT_DETAILS_REMINDER         USECASE_TYPE = "EVENT_DETAILS_REMINDER"
	UCT_EVENT_RSVP_CONFIRMATON         USECASE_TYPE = "EVENT_RSVP_CONFIRMATON"
	UCT_EVENT_RSVP_REMINDER            USECASE_TYPE = "EVENT_RSVP_REMINDER"
	UCT_FOLLOW_UP_MISSED_CALLS         USECASE_TYPE = "FOLLOW_UP_MISSED_CALLS"
	UCT_GROUP_INVITE_UPON_REQUEST      USECASE_TYPE = "GROUP_INVITE_UPON_REQUEST"
	UCT_IN_PERSON_VERIFICATION         USECASE_TYPE = "IN_PERSON_VERIFICATION"
	UCT_NETWORK_TROUBLESHOOTING        USECASE_TYPE = "NETWORK_TROUBLESHOOTING"
	UCT_ORDER_REFUND_REMINDER          USECASE_TYPE = "ORDER_REFUND_REMINDER"
	UCT_PAYMENT_NOTICE                 USECASE_TYPE = "PAYMENT_NOTICE"
	UCT_PAYMENT_SUCCESSFUL             USECASE_TYPE = "PAYMENT_SUCCESSFUL"
	UCT_RECHARGE_CONFIRMATION          USECASE_TYPE = "RECHARGE_CONFIRMATION"
	UCT_RECHARGE_REJECT_FAIL           USECASE_TYPE = "RECHARGE_REJECT_FAIL"
	UCT_RENEWAL_CONFIRMATION           USECASE_TYPE = "RENEWAL_CONFIRMATION"
	UCT_RENEWAL_REMINDER               USECASE_TYPE = "RENEWAL_REMINDER"
	UCT_RESCHEDULING_REQUEST           USECASE_TYPE = "RESCHEDULING_REQUEST"
	UCT_ROAMING_REMINDER               USECASE_TYPE = "ROAMING_REMINDER"
	UCT_SERVICE_DISRUPTION             USECASE_TYPE = "SERVICE_DISRUPTION"
	UCT_TECHNICIAN_ARRIVAL             USECASE_TYPE = "TECHNICIAN_ARRIVAL"
	UCT_TICKET_ACKNOWLEDGEMENT         USECASE_TYPE = "TICKET_ACKNOWLEDGEMENT"
	UCT_UPGRADE_CONFIRMATION           USECASE_TYPE = "UPGRADE_CONFIRMATION"
	UCT_VERIFY_TRANSACTION             USECASE_TYPE = "VERIFY_TRANSACTION"
	UCT_VERIFY_USER                    USECASE_TYPE = "VERIFY_USER"
	UCT_CALL_PERMISSION_REQUEST        USECASE_TYPE = "CALL_PERMISSION_REQUEST"
	UCT_APPOINTMENT_MISSED             USECASE_TYPE = "APPOINTMENT_MISSED"
	UCT_APPOINTMENT_SCHEDULING         USECASE_TYPE = "APPOINTMENT_SCHEDULING"
	UCT_ORDER_REFUND                   USECASE_TYPE = "ORDER_REFUND"
)

type FORMAT_TYPE = string

const (
	FT_TEXT       FORMAT_TYPE = "TEXT"
	FT_IMAGE      FORMAT_TYPE = "IMAGE"
	FT_DOCUMENT   FORMAT_TYPE = "DOCUMENT"
	FT_VIDEO      FORMAT_TYPE = "VIDEO"
	FT_LOCATION   FORMAT_TYPE = "LOCATION"
	FT_GIF        FORMAT_TYPE = "GIF"
	FT_COLLECTION FORMAT_TYPE = "COLLECTION"
	// FT_PRODUCT    FORMAT_TYPE = "PRODUCT"
)

type REJECTED_REASON = string

const (
	RR_ABUSIVE_CONTENT      REJECTED_REASON = "ABUSIVE_CONTENT"
	RR_INVALID_FORMAT       REJECTED_REASON = "INVALID_FORMAT"
	RR_NONE                 REJECTED_REASON = "NONE"
	RR_PROMOTIONAL          REJECTED_REASON = "PROMOTIONAL"
	RR_TAG_CONTENT_MISMATCH REJECTED_REASON = "TAG_CONTENT_MISMATCH"
	RR_SCAM                 REJECTED_REASON = "SCAM"
)

type CATEGORY = string

const (
	C_AUTHENTICATION          CATEGORY = "AUTHENTICATION"
	C_MARKETING               CATEGORY = "MARKETING"
	C_UTILITY                 CATEGORY = "UTILITY"
	C_ACCOUNT_UPDATE          CATEGORY = "ACCOUNT_UPDATE"
	C_PAYMENT_UPDATE          CATEGORY = "PAYMENT_UPDATE"
	C_PERSONAL_FINANCE_UPDATE CATEGORY = "PERSONAL_FINANCE_UPDATE"
	C_SHIPPING_UPDATE         CATEGORY = "SHIPPING_UPDATE"
	C_RESERVATION_UPDATE      CATEGORY = "RESERVATION_UPDATE"
	C_ISSUE_RESOLUTION        CATEGORY = "ISSUE_RESOLUTION"
	C_APPOINTMENT_UPDATE      CATEGORY = "APPOINTMENT_UPDATE"
	C_TRANSPORTATION_UPDATE   CATEGORY = "TRANSPORTATION_UPDATE"
	C_TICKET_UPDATE           CATEGORY = "TICKET_UPDATE"
	C_ALERT_UPDATE            CATEGORY = "ALERT_UPDATE"
	C_AUTO_REPLY              CATEGORY = "AUTO_REPLY"
	C_TRANSACTIONAL           CATEGORY = "TRANSACTIONAL"
	C_OTP                     CATEGORY = "OTP"
)

type SUB_CATEGORY = string

const (
	SC_ORDER_DETAILS SUB_CATEGORY = "ORDER_DETAILS"
	SC_ORDER_STATUS  SUB_CATEGORY = "ORDER_STATUS"
)

type PARAMETER_FORMAT = string

const (
	PF_NAMED      PARAMETER_FORMAT = "NAMED"
	PF_POSITIONAL PARAMETER_FORMAT = "POSITIONAL"
)

type STATUS = string

const (
	S_ACTIVE           STATUS = "ACTIVE"
	S_INACTIVE         STATUS = "INACTIVE"
	S_APPROVED         STATUS = "APPROVED"
	S_IN_APPEAL        STATUS = "IN_APPEAL"
	S_PENDING          STATUS = "PENDING"
	S_REJECTED         STATUS = "REJECTED"
	S_PENDING_DELETION STATUS = "PENDING_DELETION"
	S_DELETED          STATUS = "DELETED"
	S_DISABLED         STATUS = "DISABLED"
	S_PAUSED           STATUS = "PAUSED"
	S_LIMIT_EXCEEDED   STATUS = "LIMIT_EXCEEDED"
	S_ARCHIVED         STATUS = "ARCHIVED"
)

type ENUM_OTP_TYPE string

const (
	OT_COPY_CODE  ENUM_OTP_TYPE = "COPY_CODE"
	OT_ONE_TAP    ENUM_OTP_TYPE = "ONE_TAP"
	OT_ZERO_TAP   ENUM_OTP_TYPE = "ZERO_TAP"
	OT_NO_BUTTONS ENUM_OTP_TYPE = "NO_BUTTONS"
)

type TYPE_COMPONENT = string

const (
	TC_FOOTER                   TYPE_COMPONENT = "FOOTER"
	TC_HEADER                   TYPE_COMPONENT = "HEADER"
	TC_BODY                     TYPE_COMPONENT = "BODY"
	TC_GREETING                 TYPE_COMPONENT = "GREETING"
	TC_BUTTONS                  TYPE_COMPONENT = "BUTTONS"
	TC_CAROUSEL                 TYPE_COMPONENT = "CAROUSEL"
	TC_LIMITED_TIME_OFFER       TYPE_COMPONENT = "LIMITED_TIME_OFFER"
	TC_CALL_PERMISSION_REQUEST  TYPE_COMPONENT = "CALL_PERMISSION_REQUEST"
	TC_TAP_TARGET_CONFIGURATION TYPE_COMPONENT = "TAP_TARGET_CONFIGURATION"
)

type QUALITYSCORE = string

const (
	QS_GREEN   QUALITYSCORE = "GREEN"
	QS_YELLOW  QUALITYSCORE = "YELLOW"
	QS_RED     QUALITYSCORE = "RED"
	QS_UNKNOWN QUALITYSCORE = "UNKNOWN"
)

type ArrayButton = []Button

type PositionalParams = any // []string or string
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
	HeaderHandle          []string           `json:"header_handle,omitempty"`
	HeaderText            []PositionalParams `json:"header_text,omitempty"`
	BodyText              []PositionalParams `json:"body_text,omitempty"`
	BodyTextNamedParams   []NamedParam       `json:"body_text_named_params,omitempty"`
	HeaderTextNamedParams []NamedParam       `json:"header_text_named_params,omitempty"`
}

/*
	func ListValueToSlice(lv *structpb.ListValue) []any {
		res := make([]any, 0)

		res = lv.AsSlice()

		return res
	}
*/
func (e *Example) GetPositionalParams(lv *structpb.ListValue) []PositionalParams {
	res := make([]any, 0)

	res = lv.AsSlice()

	return res
}

type LimitedTimeOffer struct {
	Text string `json:"text" validate:"required"`
	// Set in true for that offer expiration details appear in the message sent
	// Offer details text. Maximum of 16 characters
	HasExpiration bool `json:"has_expiration,omitempty"`
}

type Card = []MockupComponent

type MockupComponent struct {
	Type TYPE_COMPONENT `json:"type" validate:"required"`
	// type of assets of media content. Configurable to "IMAGE",  "VIDEO" o "DOCUMENT"
	Format FORMAT_TYPE `json:"format,omitempty"`
	// accept parameters for programming ex:
	// 		Positional: se incluye una matriz de parámetros posicionales numerados que corresponden a posiciones numéricas en el texto del cuerpo con ejemplos.
	// 		Por ejemplo: “Hello {{1}}, your account balance is {{2}}” | [ “John”, “$1,000” ]
	// 		Nominal: se incluyen objetos JSON que contengan un parámetro con nombre y ejemplos.
	// Footer not accept parameters
	// Por ejemplo: { "param_name": "order_id", "example": "335628"}
	// less than 60 characters in header and footer, 1024 characters in body,
	// Required for components with type HEADER,BODY
	Text              string `json:"text,omitempty"`
	*ArrayButton      `json:"buttons,omitempty"`
	*LimitedTimeOffer `json:"limited_time_offer,omitempty"`
	*Card             `json:"cards,omitempty"`
	*Example          `json:"example,omitempty"`
}

// Optional data during creation of a template from a library template. These are optional fields for the body component.
type LibraryTemplateBodyInput struct {
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
type LibraryTemplateButtonInput struct {
	*URL_   `json:"url,omitempty"`
	OtpType ENUM_OTP_TYPE `json:"otp_type,omitempty"`
	*Button `json:",omitempty"`
	SApp    []SupportedApp `json:"supported_apps,omitempty"`
}

// Payload for MockupTemplate.LibraryTemplateButtonInput
type LibraryTemplateButtonInputPayload = []LibraryTemplateButtonInput

// Payload for MockupTemplate.LibraryTemplateBodyInput
type LibraryTemplateBodyInputPayload = LibraryTemplateBodyInput

type TemplateLibrary struct {
	*ArrayButton               `json:"buttons,omitempty"`
	Body                       string                        `json:"body,omitempty"`
	BodyParam                  []string                      `json:"body_params,omitempty"`
	BodyParamType              []BODY_PARAM_TPL_LIBRARY_TYPE `json:"body_param_types,omitempty"`
	Header                     string                        `json:"header,omitempty"`
	Topic                      TOPIC_TYPE                    `json:"topic,omitempty"`
	Industry                   []INDUSTRY_TYPE               `json:"industry,omitempty"`
	Usecase                    USECASE_TYPE                  `json:"usecase,omitempty"`
	LibraryTemplateBodyInput   *json.RawMessage              `json:"library_template_body_inputs,omitempty"`
	LibraryTemplateButtonInput *json.RawMessage              `json:"library_template_button_inputs,omitempty"`
}

// All template have limit of one body component
type MockupTemplate struct {
	ID                         string           `json:"id,omitempty"`
	Name                       string           `json:"name" validate:"required"` // Maximum of 512 characters
	Category                   CATEGORY         `json:"category" validate:"required"`
	CorrectCategory            CATEGORY         `json:"correct_category,omitempty"`
	Content                    string           `json:"content,omitempty"`
	PreviousCategory           CATEGORY         `json:"previous_category,omitempty"`
	SubCategory                SUB_CATEGORY     `json:"sub_category,omitempty"`
	CtaUrlLinkTrackingOptedOut bool             `json:"cta_url_link_tracking_opted_out,omitempty"`
	ParameterFormat            PARAMETER_FORMAT `json:"parameter_format,omitempty"` // The parameter format, can be Named or Positional
	RejectedReason             REJECTED_REASON  `json:"rejected_reason,omitempty"`
	Status                     STATUS           `json:"status,omitempty"`
	Language                   string           `json:"language" validate:"required"`
	AllowCategoryChange        bool             `json:"allow_category_change,omitempty"`
	QualityScore               QUALITYSCORE     `json:"quality_score,omitempty"`
	// The name exact of the library template
	LibraryTemplateName string `json:"library_template_name,omitempty"`
	// Time to live for message template sent. If users are offline for more than TTL
	// duration after message template is sent, we will retry the delivery for a period
	// of time known as a time-to-live, TTL, or the message validity period.
	// TTL can be configured for certain message types. See Time-To-Live.
	// The TTL can be customized in 1-second increments.
	// Valid values for the message_send_ttl_seconds property:
	// Authentication templates: 30 to 900 seconds (30 seconds to 15 minutes)
	// Utility templates: 30 to 43,200 seconds (30 seconds to 12 hours)
	// Marketing templates: 43,200 to 2,592,000 seconds (12 hours to 30 days)
	// For authentication and utility templates, you can set the message_send_ttl_seconds property to -1, which will apply a custom TTL of 30 days.
	MessageSendTtlSeconds int64             `json:"message_send_ttl_seconds,omitempty"`
	Components            []MockupComponent `json:"components,omitempty"`
	*TemplateLibrary      `json:",omitempty"`
}

type TYPE_TEMPLATE = string

const (
	TTPL_TEXT           TYPE_TEMPLATE = "TEXT"
	TTPL_MEDIA          TYPE_TEMPLATE = "MEDIA"
	TTPL_INTERACTIVE    TYPE_TEMPLATE = "INTERACTIVE"
	TTPL_LOCATION       TYPE_TEMPLATE = "LOCATION"
	TTPL_AUTH           TYPE_TEMPLATE = "AUTHENTICATION"
	TTPL_MULTI_PRODUCTS TYPE_TEMPLATE = "MULTI_PRODUCTS"
)

func (self MockupTemplate) GetHeader() *MockupComponent {
	if len(self.Components) == 0 {
		return nil
	}

	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_HEADER {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) GetBody() *MockupComponent {
	if len(self.Components) == 0 {
		return nil
	}

	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_BODY {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) GetFooter() *MockupComponent {
	if len(self.Components) == 0 {
		return nil
	}

	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_FOOTER {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) getButtons() *MockupComponent {
	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_BUTTONS {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) getCarusel() *MockupComponent {
	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_CAROUSEL {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) getLimitedTimeOffer() *MockupComponent {
	if len(self.Components) == 0 {
		return nil
	}

	for _, component := range self.Components {
		if strings.ToUpper(component.Type) == TC_LIMITED_TIME_OFFER {
			c := component
			return &c
		}
	}
	return nil
}

func (self MockupTemplate) getProductsButtons() bool {
	if len(self.Components) == 0 {
		return false
	}

	btns := self.getButtons()
	if btns == nil {
		return false
	}

	productBtn := []string{TB_CATALOG, TB_MPM, TB_SPM}
	for _, btn := range *btns.ArrayButton {
		if slices.Contains(productBtn, str.ToUpper(btn.Type)) {
			return true
		}
	}

	return false
}

func (self MockupTemplate) hasComponents() bool {
	return len(self.Components) != 0

}

// isText.
func (self MockupTemplate) isText() bool {
	if self.isAuth() {
		return false
	}

	h := self.GetHeader()
	b := self.GetBody()
	f := self.GetFooter()

	if !self.hasComponents() {
		if self.TemplateLibrary != nil && self.TemplateLibrary.ArrayButton == nil {
			return true
		}

		return false
	}

	if len(self.Components) == 3 {
		if h != nil && h.Format != FT_TEXT {
			return false
		}

		if f == nil {
			return false
		}
	} else {
		if len(self.Components) == 2 {
			if h != nil {
				if h.Format != FT_TEXT {
					return false
				}
			} else {
				if f == nil {
					return false
				}
			}
		}
	}

	return b != nil
}

// isMedia
func (self MockupTemplate) isMedia() bool {
	typeMedia := []string{FT_DOCUMENT, FT_IMAGE, FT_VIDEO, FT_GIF}
	return self.GetHeader() != nil && slices.Contains(typeMedia, strings.ToUpper(self.GetHeader().Format))
}

// isButton
func (self MockupTemplate) isButton() bool {
	btns := self.getButtons()
	if btns == nil {
		return false
	}

	return true
}

// isInteractive
func (self MockupTemplate) isInteractive() bool {
	if !self.hasComponents() {
		if self.TemplateLibrary != nil && self.TemplateLibrary.ArrayButton != nil {
			return true
		}

		return false
	}

	return self.isButton()
}

// isLocation
func (self MockupTemplate) isLocation() bool {
	return self.GetHeader() != nil && strings.ToUpper(self.GetHeader().Format) == FT_LOCATION
}

// isAuth
func (self MockupTemplate) isAuth() bool {
	return strings.ToUpper(self.Category) == C_AUTHENTICATION
}

// isProducts
func (self MockupTemplate) isProducts() bool {
	if self.getCarusel() != nil {
		return true
	}

	return self.getProductsButtons()
}

// type MockupTemplate are (text, media, interactive msg, location, authentication and products various message )
func (self MockupTemplate) GetTypeTpl() string {
	switch {
	case self.isText():
		return TTPL_TEXT
	case self.isMedia():
		return TTPL_MEDIA
	case self.isLocation():
		return TTPL_LOCATION
	case self.isAuth():
		return TTPL_AUTH
	case self.isProducts():
		return TTPL_MULTI_PRODUCTS
	case self.isInteractive():
		return TTPL_INTERACTIVE
	default:
		return "UNKNOWN"
	}
}

func (self MockupTemplate) GetContentHeaderHandle() (content string) {
	h := self.GetHeader()
	if h == nil {
		return content
	}

	if len(h.HeaderHandle) > 0 {
		return h.HeaderHandle[0]
	}

	return content
}

const (
	LOCATE_LIBERARY_TEMPLATE = "LIBRARY_TEMPLATE"
	LOCATE_WABA_TEMPLATE     = "WABA_TEMPLATE"
)

func (self MockupTemplate) GetLocatedTpl() string {
	if self.Components != nil {
		return LOCATE_WABA_TEMPLATE
	}

	return LOCATE_LIBERARY_TEMPLATE
}

var getParamPositional = func(pp PositionalParams, arrParam *[]map[string]any) {
	for _, p := range pp.([]string) {
		*arrParam = append(*arrParam, map[string]any{
			"type": "text",
			"text": p,
		})
	}
}

var getParamName = func(np NamedParam, arrParam *[]map[string]any) {
	*arrParam = append(*arrParam, map[string]any{
		"type":           "text",
		"parameter_name": np.ParamName,
		"text":           np.Example,
	})
}

var getParamButtons = func(cmp MockupComponent, arrBtnWithParam *[]map[string]any) {
	getParamType := func(b Button) (parmIndex int, paramName string) {
		if b.Type != TB_URL {
			return 0, ""
		}

		re := regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)
		if re.MatchString(b.Url) {
			matches := re.FindStringSubmatch(b.Url)
			if len(matches) <= 1 {
				return 0, ""
			}

			matchStr := matches[1]
			isNumber := regexp.MustCompile(`[0-9]`)
			if len(matchStr) == 1 && isNumber.MatchString(matchStr) {
				val, _ := strconv.Atoi(matchStr)
				return val, ""
			} else {
				return 0, matches[1]
			}
		}

		return 0, ""
	}

	for index, b := range *cmp.ArrayButton {
		// search if has params {{name}} or {{#}}
		pos, name := getParamType(b)
		if pos == 0 && name == "" {
			continue
		}

		var tempParam, temp map[string]any = make(map[string]any, 0), make(map[string]any, 0)

		if pos != 0 {
			temp = map[string]any{
				"type": "text",
				"text": pos,
			}
		}

		if name != "" {
			temp = map[string]any{
				"type":           "text",
				"text":           "",
				"parameter_name": name,
			}
		}

		tempParam[fmt.Sprintf("btn%d", index)] = temp

		*arrBtnWithParam = append(*arrBtnWithParam, tempParam)
	}
}

func (cmp MockupComponent) GetParams() (param []map[string]any) {
	param = make([]map[string]any, 0)

	switch cmp.Type {
	case TC_HEADER:
		fallthrough
	case TC_BODY:
		switch {
		case cmp.Example != nil && (cmp.BodyText != nil || cmp.HeaderText != nil):
			if cmp.BodyText != nil {
				for _, v := range cmp.BodyText {
					getParamPositional(v, &param)
				}
			} else {
				for _, v := range cmp.HeaderText {
					getParamPositional(v, &param)
				}
			}

		case cmp.Example != nil && (cmp.BodyTextNamedParams != nil || cmp.HeaderTextNamedParams != nil):
			if cmp.BodyTextNamedParams != nil {
				for _, v := range cmp.BodyTextNamedParams {
					getParamName(v, &param)
				}
			} else {
				for _, v := range cmp.HeaderTextNamedParams {
					getParamName(v, &param)
				}
			}
		}
	case TC_BUTTONS:
		getParamButtons(cmp, &param)
	}

	return param
}
