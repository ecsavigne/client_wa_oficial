package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/internal"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/message"
)

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number,omitempty"`
	PhoneNumberID      string `json:"phone_number_id,omitempty"`
}

func (self Metadata) GetDisplayPhoneNumber() string { return self.DisplayPhoneNumber }
func (self Metadata) GetPhoneNumberID() string      { return self.PhoneNumberID }

type Profile struct {
	Name string `json:"name,omitempty"` //  The customer's name.
}

type Contact struct {
	WaId     string `json:"wa_id,omitempty"`
	UserId   string `json:"user_id,omitempty"` // Additional unique, alphanumeric identifier for a WhatsApp user.
	*Profile `json:"profile,omitempty"`
}

type ErrorData struct {
	MessagingProduct string `json:"messaging_product,omitempty"` // "whatsapp",
	Details          string `json:"details,omitempty"`
}

type Error struct {
	Message      string     `json:"message,omitempty"` // Combination of the error code and title. ej: (#130429) Rate limit hit.
	Type         string     `json:"type,omitempty"`
	Code         int64      `json:"code,omitempty"`
	ErrorData    *ErrorData `json:"error_data,omitempty"`
	Title        string     `json:"title,omitempty"`
	ErrorSubcode int64      `json:"error_subcode,omitempty"`
	FbtraceID    string     `json:"fbtrace_id,omitempty"`
	Href         string     `json:"href,omitempty"`
}

type Text struct {
	Body string `json:"body,omitempty"`
}

type Reaction struct {
	MessageID string `json:"message_id,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
}

type InfoMedia struct {
	MIMEType string `json:"mime_type,omitempty"`
	ID       string `json:"id,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Sha256   string `json:"sha256,omitempty"`
}

type Image struct {
	*InfoMedia `json:",omitempty"`
}

type Sticker struct {
	*InfoMedia `json:",omitempty"`
	Animated   bool `json:"animated,omitempty"`
}

type ErrorMessage struct {
	Code    int64  `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
	Title   string `json:"title,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

type ReferredProduct struct {
	CatalogId         string `json:"catalog_id"`          // Unique identifier of the Meta catalog linked to the WhatsApp Business Account.
	ProductRetailerId string `json:"product_retailer_id"` // Unique identifier of the product in a catalog.
}

type Context struct {
	Forwarded           bool            `json:"forwarded"`            //Set to true if the message received by the business has been forwarded.
	FrequentlyForwarded bool            `json:"frequently_forwarded"` //Set to true if the message received by the business has been forwarded more than 5 times.
	From                string          `json:"from"`
	ID                  string          `json:"id"`
	ReferredProduct     ReferredProduct `json:"referred_product"`
}

func (ctx Context) GetForwarded() bool {
	return ctx.Forwarded
}

func (ctx Context) GetFrequentlyForwarded() bool {
	return ctx.FrequentlyForwarded
}

func (ctx Context) GetFrom() string {
	return ctx.From
}

func (ctx Context) GetID() string {
	return ctx.ID
}

func (ctx Context) GetReferredProduct() ReferredProduct {
	return ctx.ReferredProduct
}

func (ctx Context) String() string {
	return internal.ConvertStr(&ctx)
}

type Button struct {
	Text    string `json:"text,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type InteractiveCommon struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type ListReply struct {
	*InteractiveCommon
	Description string `json:"description,omitempty"`
}

type ButtonReply struct {
	*InteractiveCommon
}

type Interactive struct {
	ListReply   *ListReply   `json:"list_reply,omitempty"`
	ButtonReply *ButtonReply `json:"button_reply,omitempty"`
	Type        string       `json:"type,omitempty"`
}

type Referral struct {
	SourceURL    string `json:"source_url,omitempty"`
	SourceID     string `json:"source_id,omitempty"`
	SourceType   string `json:"source_type,omitempty"`
	Headline     string `json:"headline,omitempty"`
	Body         string `json:"body,omitempty"`
	MediaType    string `json:"media_type,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	VideoURL     string `json:"video_url,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	CtwaClid     string `json:"ctwa_clid,omitempty"`
}

type ProductItem struct {
	ProductRetailerID string  `json:"product_retailer_id,omitempty"`
	Quantity          int     `json:"quantity,omitempty"`
	ItemPrice         float32 `json:"item_price,omitempty"`
	Currency          string  `json:"currency,omitempty"`
}

type Order struct {
	CatalogID    string        `json:"catalog_id,omitempty"`
	ProductItems []ProductItem `json:"product_items,omitempty"`
	Text         string        `json:"text,omitempty"`
}

type System struct {
	Body     string `json:"body,omitempty"`
	WaId     string `json:"wa_id,omitempty"`
	Identity string `json:"identity,omitempty"`
	NewWaID  int64  `json:"new_wa_id,omitempty"`
	Type     string `json:"type,omitempty"`
	Customer string `json:"customer,omitempty"`
}

type Identity struct {
	Acknowledged     string `json:"acknowledged,omitempty"`
	CreatedTimestamp string `json:"created_timestamp,omitempty"`
	Hash             string `json:"hash,omitempty"`
}

type Audio struct {
	*InfoMedia `json:",omitempty"`
}
type Video struct {
	*InfoMedia `json:",omitempty"`
}

type Document struct {
	*InfoMedia `json:",omitempty"`
	Filename   string `json:"filename,omitempty"`
}

type Message struct {
	From        string            `json:"from,omitempty"`
	ID          string            `json:"id,omitempty"`
	Type        string            `json:"type,omitempty"`
	Timestamp   string            `json:"timestamp,omitempty"`
	Text        *Text             `json:"text,omitempty"`
	Reaction    *Reaction         `json:"reaction,omitempty"`
	Image       *Image            `json:"image,omitempty"`
	Audio       *Audio            `json:"audio,omitempty"`
	Video       *Video            `json:"video,omitempty"`
	Document    *Document         `json:"document,omitempty"`
	Sticker     *Sticker          `json:"sticker,omitempty"`
	Location    *Location         `json:"location,omitempty"`
	Context     *Context          `json:"context,omitempty"`
	Button      *Button           `json:"button,omitempty"`
	Order       *Order            `json:"order,omitempty"`
	System      *System           `json:"system,omitempty"`
	Referral    *Referral         `json:"referral,omitempty"`
	Identity    *Identity         `json:"identity,omitempty"`
	Interactive *Interactive      `json:"interactive,omitempty"`
	Contacts    []message.Contact `json:"contacts,omitempty"`
	Errors      []ErrorMessage    `json:"errors,omitempty"`
}

type Pricing struct {
	Billable     bool   `json:"billable,omitempty"`
	PricingModel string `json:"pricing_model,omitempty"`
	Category     string `json:"category,omitempty"`
}

type Origin struct {
	Type string `json:"type,omitempty"`
}

type Conversation struct {
	ID                  string  `json:"id,omitempty"`
	ExpirationTimestamp string  `json:"expiration_timestamp,omitempty"`
	Origin              *Origin `json:"origin,omitempty"`
}

type Statuse struct {
	ID                    string        `json:"id,omitempty"`
	BizOpaqueCallbackData string        `json:"biz_opaque_callback_data,omitempty"`
	Status                string        `json:"status,omitempty"` //delivered, read, sent, failed
	Timestamp             string        `json:"timestamp,omitempty"`
	RecipientID           string        `json:"recipient_id,omitempty"`
	Conversation          *Conversation `json:"conversation,omitempty"`
	Pricing               *Pricing      `json:"pricing,omitempty"`
	Errors                []Error       `json:"errors,omitempty"`
}

type MessageTemplateButton struct {
	MessageTemplateButtonType        string `json:"message_template_button_type,omitempty"`
	MessageTemplateButtonText        string `json:"message_template_button_text,omitempty"`
	MessageTemplateButtonUrl         string `json:"message_template_button_url,omitempty"`
	MessageTemplateButtonPhoneNumber string `json:"message_template_button_phone_number,omitempty"`
}

type MessageTemplateComponentsUpdate struct {
	MessageTemplateId      string                  `json:"message_template_id,omitempty"`
	MessageTemplateName    string                  `json:"message_template_name,omitempty"`
	MessageTemplateBody    string                  `json:"message_template_element,omitempty"`
	MessageTemplateTitle   string                  `json:"message_template_title,omitempty"`
	MessageTemplateFooter  string                  `json:"message_template_footer,omitempty"`
	MessageTemplateButtons []MessageTemplateButton `json:"message_template_buttons,omitempty"`
}

type MessageTemplateQualityUpdate struct {
	PreviousQualityScore    string `json:"previous_quality_score,omitempty"`
	NewQualityScore         string `json:"new_quality_score,omitempty"`
	MessageTemplateLanguage string `json:"message_template_language,omitempty"`
}

type DisableInfo struct {
	DisableDate uint64 `json:"disable_date,omitempty"`
}

type OtherInfo struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type RejectionInfo struct {
	Reason         string `json:"reason,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

type MessageTemplateStatusUpdate struct {
	Event                   string `json:"event,omitempty"`
	Reason                  string `json:"reason,omitempty"`
	MessageTemplateCategory string `json:"message_template_category,omitempty"`
	*DisableInfo            `json:"disable_info,omitempty"`
	*OtherInfo              `json:"other_info,omitempty"`
	*RejectionInfo          `json:"rejection_info,omitempty"`
}

type TemplateCategoryUpdate struct {
	CorrectCategory         string `json:"correct_category,omitempty"`
	PreviousCategory        string `json:"previous_category,omitempty"`
	NewCategory             string `json:"new_category,omitempty"`
	CategoryUpdateTimestamp uint64 `json:"category_update_timestamp,omitempty"`
}

type MessageTemplateWebhook struct {
	*MessageTemplateComponentsUpdate `json:",omitempty"`
	*MessageTemplateQualityUpdate    `json:",omitempty"`
	*MessageTemplateStatusUpdate     `json:",omitempty"`
	*TemplateCategoryUpdate          `json:",omitempty"`
}

type WabaInfo struct {
	WabaID                     string   `json:"waba_id,omitempty"`
	AdAccountLinked            string   `json:"ad_account_linked,omitempty"`
	OwnerBusinessID            string   `json:"owner_business_id,omitempty"`
	PartnerAppID               string   `json:"partner_app_id,omitempty"`
	SolutionID                 string   `json:"solution_id,omitempty"`
	SolutionPartnerBusinessIDs []string `json:"solution_partner_business_ids,omitempty"`
}

type ViolationInfo struct {
	ViolationType string `json:"violation_type,omitempty"`
}

type ExceptionCountries struct {
	CountryCode string `json:"country_code,omitempty"`
	StartTime   uint64 `json:"start_time,omitempty"`
}
type AuthInternationalRateEligibility struct {
	ExceptionCountries []ExceptionCountries `json:"exception_countries"`
	StartTime          uint64               `json:"start_time,omitempty"`
}

type BanInfo struct {
	WabaBanState string `json:"waba_ban_state,omitempty"`
	WabaBanDate  string `json:"waba_ban_date,omitempty"`
}

type VolumeTierInfo struct {
	TierUpdateTime  uint64 `json:"tier_update_time,omitempty"`
	PricingCategory string `json:"pricing_category,omitempty"`
	Tier            string `json:"tier,omitempty"`
	EffectiveMonth  string `json:"effective_month,omitempty"`
	Region          string `json:"region,omitempty"`
}

type DisconnectionInfo struct {
	Reason      string `json:"reason,omitempty"`
	InitiatedBy string `json:"initiated_by,omitempty"`
}

type PartnerClientCertificationInfo struct {
	ClientBusinessID string   `json:"client_business_id,omitempty"`
	Status           string   `json:"status,omitempty"`
	RejectionReasons []string `json:"rejection_reasons,omitempty"`
}

type RestrictionInfo struct {
	RestrictionType string `json:"restriction_type,omitempty"`
	Expiration      uint64 `json:"expiration,omitempty"`
	Remediation     string `json:"remediation,omitempty"`
}

type Value struct {
	MessagingProduct                 string                            `json:"messaging_product,omitempty"` //Product used to send the message. Value is always whatsapp.
	Metadata                         *Metadata                         `json:"metadata,omitempty"`
	WabaInfo                         *WabaInfo                         `json:"waba_info,omitempty"`
	ViolationInfo                    *ViolationInfo                    `json:"violation_info,omitempty"`
	AuthInternationalRateEligibility *AuthInternationalRateEligibility `json:"auth_international_rate_eligibility,omitempty"`
	BanInfo                          *BanInfo                          `json:"ban_info,omitempty"`
	VolumeTierInfo                   *VolumeTierInfo                   `json:"volume_tier_info,omitempty"`
	DisconnectionInfo                *DisconnectionInfo                `json:"disconnection_info,omitempty"`
	PartnerClientCertificationInfo   *PartnerClientCertificationInfo   `json:"partner_client_certification_info,omitempty"`
	RestrictionInfo                  []RestrictionInfo                 `json:"restriction_info,omitempty"`
	Country                          string                            `json:"country,omitempty"`
	Contacts                         []Contact                         `json:"contacts,omitempty"`
	Messages                         []Message                         `json:"messages,omitempty"`
	Statuses                         []Statuse                         `json:"statuses,omitempty"`
	Errors                           []Error                           `json:"errors,omitempty"`
	*MessageTemplateWebhook          `json:",omitempty"`
}

func (self Value) GetMessagingProduct() string {
	return self.MessagingProduct
}

func (self Value) GetMetadata() *Metadata {
	return self.Metadata
}

func (self Value) GetContacts() []Contact {
	return self.Contacts
}

func (self Value) GetMessages() []Message {
	return self.Messages
}

func (self Value) GetStatuses() []Statuse {
	return self.Statuses
}

type Change struct {
	Value *Value `json:"value,omitempty"`
	Field string `json:"field,omitempty"` // Notification type. value will be "messages"
}

func (self Change) GetValue() (val *Value) {
	defer func() {
		if r := recover(); r != nil {
			val = (*Value)(nil)
		}
	}()

	return self.Value
}

func (self Change) GetField() (val string) {
	defer func() {
		if r := recover(); r != nil {
			val = ""
		}
	}()

	return self.Field
}

type Entry struct {
	ID      string   `json:"id,omitempty"`
	Time    uint64   `json:"time,omitempty"`
	Changes []Change `json:"changes,omitempty"`
}

func (self Entry) GetChange() []Change {
	return self.Changes
}

func (self Entry) GetID() string {
	return self.ID
}

type Components struct {
	Object string  `json:"object,omitempty"` // "whatsapp_business_account"
	Entry  []Entry `json:"entry,omitempty"`
}

func (self Components) GetEntry() []Entry {
	return self.Entry
}

func (self Components) GetSatusMessage() (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	return self.GetEntry()[0].GetChange()[0].GetValue().GetStatuses()[0].Status
}

func (self Components) GetTypeMessage() (typ string) {
	defer func() {
		if r := recover(); r != nil {
			typ = ""
		}
	}()

	return self.Entry[0].Changes[0].Value.Messages[0].Type
}
