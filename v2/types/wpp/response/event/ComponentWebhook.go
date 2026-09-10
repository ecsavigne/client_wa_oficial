package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response/event/types"
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

func (self Profile) GetName() string { return self.Name }

type Contact struct {
	WaId     string `json:"wa_id,omitempty"`
	UserId   string `json:"user_id,omitempty"` // Additional unique, alphanumeric identifier for a WhatsApp user.
	*Profile `json:"profile,omitempty"`
}

func (self Contact) GetWaID() string      { return self.WaId }
func (self Contact) GetUserId() string    { return self.UserId }
func (self Contact) GetProfile() *Profile { return self.Profile }

type ErrorData struct {
	MessagingProduct string `json:"messaging_product,omitempty"` // "whatsapp",
	Details          string `json:"details,omitempty"`
}

func (self ErrorData) GetMessagingProduct() string { return self.MessagingProduct }
func (self ErrorData) GetDetails() string          { return self.Details }

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

func (self Error) GetMessage() string       { return self.Message }
func (self Error) GetType() string          { return self.Type }
func (self Error) GetCode() int64           { return self.Code }
func (self Error) GetErrorData() *ErrorData { return self.ErrorData }
func (self Error) GetTitle() string         { return self.Title }
func (self Error) GetErrorSubcode() int64   { return self.ErrorSubcode }
func (self Error) GetFbtraceID() string     { return self.FbtraceID }
func (self Error) GetHref() string          { return self.Href }

type Text struct {
	Body string `json:"body,omitempty"`
}

func (self Text) GetText() string { return self.Body }

type Reaction struct {
	MessageID string `json:"message_id,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
}

func (self Reaction) GetMessageID() string { return self.MessageID }
func (self Reaction) GetEmoji() string     { return self.Emoji }

type InfoMedia struct {
	MIMEType string `json:"mime_type,omitempty"`
	ID       string `json:"id,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Sha256   string `json:"sha256,omitempty"`
	Url      string `json:"url,omitempty"`
}

func (self InfoMedia) GetMimeType() string { return self.MIMEType }
func (self InfoMedia) GetID() string       { return self.ID }
func (self InfoMedia) GetCaption() string  { return self.Caption }
func (self InfoMedia) GetSha256() string   { return self.Sha256 }
func (self InfoMedia) GetURL() string      { return self.Url }

type Image struct {
	*InfoMedia `json:",omitempty"`
}

func (v Image) GetCaption() string  { return v.InfoMedia.Caption }
func (v Image) GetID() string       { return v.InfoMedia.ID }
func (v Image) GetSha256() string   { return v.InfoMedia.Sha256 }
func (v Image) GetURL() string      { return v.InfoMedia.Url }
func (v Image) GetMimeType() string { return v.InfoMedia.MIMEType }

type Sticker struct {
	*InfoMedia `json:",omitempty"`
	Animated   bool `json:"animated,omitempty"`
}

func (v Sticker) GetCaption() string  { return v.InfoMedia.Caption }
func (v Sticker) GetID() string       { return v.InfoMedia.ID }
func (v Sticker) GetSha256() string   { return v.InfoMedia.Sha256 }
func (v Sticker) GetURL() string      { return v.InfoMedia.Url }
func (v Sticker) GetMimeType() string { return v.InfoMedia.MIMEType }

type ErrorMessage struct {
	Code    int64  `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
	Title   string `json:"title,omitempty"`
}

func (self ErrorMessage) GetCode() int64     { return self.Code }
func (self ErrorMessage) GetDetails() string { return self.Details }
func (self ErrorMessage) GetTitle() string   { return self.Title }

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

func (self Location) GetLatitude() float64  { return self.Latitude }
func (self Location) GetLongitude() float64 { return self.Longitude }
func (self Location) GetName() string       { return self.Name }
func (self Location) GetAddress() string    { return self.Address }

type ReferredProduct struct {
	CatalogId         string `json:"catalog_id"`          // Unique identifier of the Meta catalog linked to the WhatsApp Business Account.
	ProductRetailerId string `json:"product_retailer_id"` // Unique identifier of the product in a catalog.
}

func (self ReferredProduct) GetCatalogId() string         { return self.CatalogId }
func (self ReferredProduct) GetProductRetailerId() string { return self.ProductRetailerId }

type Revoke struct {
	OriginalMessageId string `json:"original_message_id"`
}

func (self Revoke) GetOriginalMessageId() string { return self.OriginalMessageId }

type Edit struct {
	OriginalMessageId string `json:"original_message_id"`
	Message           *struct {
		Context *struct {
			ID string `json:"id,omitempty"`
		} `json:"context,omitempty"`
	} `json:"message,omitempty"`
	Type         string `json:"type,omitempty"`
	*TypeMessage `json:",omitempty"`
}

func (self Edit) GetOriginalMessageId() string { return self.OriginalMessageId }
func (self Edit) GetType() string              { return self.Type }
func (self Edit) GetTypeMessage() *TypeMessage { return self.TypeMessage }
func (self Edit) GetMessage() *struct {
	Context *struct {
		ID string `json:"id,omitempty"`
	} `json:"context,omitempty"`
} {
	return self.Message
}

type TypeMessageArray struct {
	Contacts []message.Contact `json:"contacts,omitempty"`
	Errors   []ErrorMessage    `json:"errors,omitempty"`
}

func (self TypeMessageArray) GetContacts() []message.Contact { return self.Contacts }
func (self TypeMessageArray) GetErrors() []ErrorMessage      { return self.Errors }

type TypeMessage struct {
	*Revoke           `json:"revoke,omitempty"`
	*Edit             `json:"edit,omitempty"`
	*Audio            `json:"audio,omitempty"`
	*Button           `json:"button,omitempty"`
	*TypeMessageArray `json:",omitempty"`
	*Context          `json:"context,omitempty"`
	*Document         `json:"document,omitempty"`
	*Identity         `json:"identity,omitempty"`
	*Image            `json:"image,omitempty"`
	*Interactive      `json:"interactive,omitempty"`
	*Location         `json:"location,omitempty"`
	*Reaction         `json:"reaction,omitempty"`
	*Order            `json:"order,omitempty"`
	*Referral         `json:"referral,omitempty"`
	*Sticker          `json:"sticker,omitempty"`
	*System           `json:"system,omitempty"`
	*Text             `json:"text,omitempty"`
	*Video            `json:"video,omitempty"`
}

func (self TypeMessage) GetRevoke() *Revoke                     { return self.Revoke }
func (self TypeMessage) GetEdit() *Edit                         { return self.Edit }
func (self TypeMessage) GetAudio() *Audio                       { return self.Audio }
func (self TypeMessage) GetButton() *Button                     { return self.Button }
func (self TypeMessage) GetTypeMessageArray() *TypeMessageArray { return self.TypeMessageArray }
func (self TypeMessage) GetContext() *Context                   { return self.Context }
func (self TypeMessage) GetDocument() *Document                 { return self.Document }
func (self TypeMessage) GetIdentity() *Identity                 { return self.Identity }
func (self TypeMessage) GetImage() *Image                       { return self.Image }
func (self TypeMessage) GetInteractive() *Interactive           { return self.Interactive }
func (self TypeMessage) GetLocation() *Location                 { return self.Location }
func (self TypeMessage) GetReaction() *Reaction                 { return self.Reaction }
func (self TypeMessage) GetOrder() *Order                       { return self.Order }
func (self TypeMessage) GetReferral() *Referral                 { return self.Referral }
func (self TypeMessage) GetSticker() *Sticker                   { return self.Sticker }
func (self TypeMessage) GetSystem() *System                     { return self.System }
func (self TypeMessage) GetText() *Text                         { return self.Text }
func (self TypeMessage) GetVideo() *Video                       { return self.Video }

type Context struct {
	Forwarded           bool            `json:"forwarded"`            //Set to true if the message received by the business has been forwarded.
	FrequentlyForwarded bool            `json:"frequently_forwarded"` //Set to true if the message received by the business has been forwarded more than 5 times.
	From                string          `json:"from"`
	ID                  string          `json:"id"`
	ReferredProduct     ReferredProduct `json:"referred_product"`
}

func (ctx Context) GetForwarded() bool                  { return ctx.Forwarded }
func (ctx Context) GetFrequentlyForwarded() bool        { return ctx.FrequentlyForwarded }
func (ctx Context) GetFrom() string                     { return ctx.From }
func (ctx Context) GetID() string                       { return ctx.ID }
func (ctx Context) GetReferredProduct() ReferredProduct { return ctx.ReferredProduct }
func (ctx Context) String() string                      { return internal.ConvertStr(&ctx) }

type Button struct {
	Text    string `json:"text,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type InteractiveCommon struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type ListReply struct {
	*InteractiveCommon `json:",omitempty"`
	Description        string `json:"description,omitempty"`
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

func (r Referral) GetSourceURL() string    { return r.SourceURL }
func (r Referral) GetSourceId() string     { return r.SourceID }
func (r Referral) GetSourceType() string   { return r.SourceType }
func (r Referral) GetHeadline() string     { return r.Headline }
func (r Referral) GetBody() string         { return r.Body }
func (r Referral) GetMediaType() string    { return r.MediaType }
func (r Referral) GetImageURL() string     { return r.ImageURL }
func (r Referral) GetVideoURL() string     { return r.VideoURL }
func (r Referral) GetThumbnailURL() string { return r.ThumbnailURL }
func (r Referral) GetCtwaClid() string     { return r.CtwaClid }

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

func (v Audio) GetCaption() string  { return v.InfoMedia.Caption }
func (v Audio) GetID() string       { return v.InfoMedia.ID }
func (v Audio) GetSha256() string   { return v.InfoMedia.Sha256 }
func (v Audio) GetURL() string      { return v.InfoMedia.Url }
func (v Audio) GetMimeType() string { return v.InfoMedia.MIMEType }

type Video struct {
	*InfoMedia `json:",omitempty"`
}

func (v Video) GetCaption() string  { return v.InfoMedia.Caption }
func (v Video) GetID() string       { return v.InfoMedia.ID }
func (v Video) GetSha256() string   { return v.InfoMedia.Sha256 }
func (v Video) GetURL() string      { return v.InfoMedia.Url }
func (v Video) GetMimeType() string { return v.InfoMedia.MIMEType }

type Document struct {
	*InfoMedia `json:",omitempty"`
	Filename   string `json:"filename,omitempty"`
}

func (v Document) GetCaption() string  { return v.InfoMedia.Caption }
func (v Document) GetID() string       { return v.InfoMedia.ID }
func (v Document) GetSha256() string   { return v.InfoMedia.Sha256 }
func (v Document) GetURL() string      { return v.InfoMedia.Url }
func (v Document) GetMimeType() string { return v.InfoMedia.MIMEType }

func (v Document) GetFilename() string {
	return v.Filename
}

type Message struct {
	From         string `json:"from,omitempty"`
	To           string `json:"to,omitempty"`
	ID           string `json:"id,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	Type         string `json:"type,omitempty"`
	*TypeMessage `json:",omitempty"`
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

func (self WabaInfo) GetWabaID() string                       { return self.WabaID }
func (self WabaInfo) GetOwnerBusinessID() string              { return self.OwnerBusinessID }
func (self WabaInfo) GetSolutionID() string                   { return self.SolutionID }
func (self WabaInfo) GetSolutionPartnerBusinessIDs() []string { return self.SolutionPartnerBusinessIDs }
func (self WabaInfo) GetAdAccountLinked() string              { return self.AdAccountLinked }
func (self WabaInfo) GetPartnerAppID() string                 { return self.PartnerAppID }

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

type MessageEchoes struct {
	From         string `json:"from,omitempty"`
	To           string `json:"to,omitempty"`
	ID           string `json:"id,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	Type         string `json:"type,omitempty"`
	*TypeMessage `json:",omitempty"`
}

type Value struct {
	AuthInternationalRateEligibility *AuthInternationalRateEligibility `json:"auth_international_rate_eligibility,omitempty"`
	BanInfo                          *BanInfo                          `json:"ban_info,omitempty"`
	Country                          string                            `json:"country,omitempty"`
	Contacts                         []Contact                         `json:"contacts,omitempty"`
	DisconnectionInfo                *DisconnectionInfo                `json:"disconnection_info,omitempty"`
	Errors                           []Error                           `json:"errors,omitempty"`
	MessagingProduct                 string                            `json:"messaging_product,omitempty"` //Product used to send the message. Value is always whatsapp.
	Messages                         []Message                         `json:"messages,omitempty"`
	MessageEchoes                    []MessageEchoes                   `json:"message_echoes,omitempty"`
	*MessageTemplateWebhook          `json:",omitempty"`
	Metadata                         *Metadata                       `json:"metadata,omitempty"`
	WabaInfo                         *WabaInfo                       `json:"waba_info,omitempty"`
	RestrictionInfo                  []RestrictionInfo               `json:"restriction_info,omitempty"`
	PartnerClientCertificationInfo   *PartnerClientCertificationInfo `json:"partner_client_certification_info,omitempty"`
	Statuses                         []Statuse                       `json:"statuses,omitempty"`
	ViolationInfo                    *ViolationInfo                  `json:"violation_info,omitempty"`
	VolumeTierInfo                   *VolumeTierInfo                 `json:"volume_tier_info,omitempty"`
}

func (self *Value) GetMessagingProduct() string        { return self.MessagingProduct }
func (self *Value) GetMetadata() *Metadata             { return self.Metadata }
func (self *Value) GetContacts() []Contact             { return self.Contacts }
func (self *Value) GetMessages() []Message             { return self.Messages }
func (self *Value) GetMessagesEchoes() []MessageEchoes { return self.MessageEchoes }

// GetMessages() or GetMessagesEchoes()
func (self *Value) GetMessage() any {
	if len(self.GetMessages()) > 0 {
		return self.GetMessages()
	}

	if len(self.GetMessagesEchoes()) > 0 {
		return self.GetMessagesEchoes()
	}

	return nil
}

func (self *Value) GetTypeMessage() string {
	switch v := self.GetMessage().(type) {
	case []Message:
		return v[0].Type
	case []MessageEchoes:
		return v[0].Type
	default:
		return ""
	}
}

func (self *Value) GetStatuses() []Statuse { return self.Statuses }

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

func (self *Entry) GetChange() []Change {
	return self.Changes
}

func (self *Entry) GetID() string {
	return self.ID
}

type Components struct {
	Object string  `json:"object,omitempty"` // "whatsapp_business_account"
	Entry  []Entry `json:"entry,omitempty"`
}

func (self *Components) GetObject() string { return self.Object }
func (self *Components) GetEntry() []Entry { return self.Entry }

func (self *Components) GetSatusMessage() (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	return self.GetEntry()[0].GetChange()[0].GetValue().GetStatuses()[0].Status
}

func (self *Components) GetTypeMessage() (typ string) {
	defer func() {
		if r := recover(); r != nil {
			typ = ""
		}
	}()

	v := self.Entry[0].Changes[0].Value
	return v.GetTypeMessage()
}

func (self *Components) GetNotificationType() (vNotification string) {
	defer func() {
		if r := recover(); r != nil {
			vNotification = ""
		}
	}()

	return self.GetEntry()[0].GetChange()[0].Field
}

// Only in event WEBHOOK_NOTIFICATION_ACCOUNT_UPDATE
func (self *Components) GetEvent() (vEvt string) {
	defer func() {
		if r := recover(); r != nil {
			vEvt = ""
		}
	}()

	return self.GetEntry()[0].GetChange()[0].GetValue().Event
}

func (self *Components) GetChange() (vCh []Change) {
	defer func() {
		if r := recover(); r != nil {
			vCh = []Change{}
		}
	}()

	return self.GetEntry()[0].GetChange()
}

func (self *Components) GetValue() (vVal *Value) {
	defer func() {
		if r := recover(); r != nil {
			vVal = nil
		}
	}()

	return self.GetChange()[0].GetValue()
}

/*
Return []Message, []MessageEchoes or nil
*/
func (self *Components) GetMessages() (vMsg any) {
	defer func() {
		if r := recover(); r != nil {
			vMsg = []Message{}
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		return self.GetValue().GetMessage()
	} else {
		return nil
	}
}

func (self *Components) GetMessageID() (vMessageID string) {
	defer func() {
		if r := recover(); r != nil {
			vMessageID = ""
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].ID
		case []MessageEchoes:
			return v[0].ID
		default:
			return ""
		}
	} else {
		panic("this is not a message")
	}
}

// Get Form
func (self *Components) GetSenderMessage() (vMessageID string) {
	defer func() {
		if r := recover(); r != nil {
			vMessageID = ""
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].From
		case []MessageEchoes:
			return v[0].From
		default:
			return ""
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetMessageTimestamp() (timestamp string) {
	defer func() {
		if r := recover(); r != nil {
			timestamp = ""
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Timestamp
		case []MessageEchoes:
			return v[0].Timestamp
		default:
			return ""
		}
	} else {
		panic("this is not a message")
	}
}

// Phone number of whatsapp business
func (self *Components) GetDisplayPhoneNumber() (phoneNumber string) {
	return self.GetValue().GetMetadata().GetDisplayPhoneNumber()
}
func (self *Components) GetContactPhoneNumber() (phoneNumber string) {
	return self.GetValue().GetContacts()[0].GetWaID()
}
func (self *Components) GetContactUserID() (phoneNumber string) {
	return self.GetValue().GetContacts()[0].GetUserId()
}
func (self *Components) GetAudioMessage() (vAudio *Audio) {
	defer func() {
		if r := recover(); r != nil {
			vAudio = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Audio
		case []MessageEchoes:
			return v[0].Audio
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetVideoMessage() (v *Video) {
	defer func() {
		if r := recover(); r != nil {
			v = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Video
		case []MessageEchoes:
			return v[0].Video
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetImageMessage() (vImage *Image) {
	defer func() {
		if r := recover(); r != nil {
			vImage = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Image
		case []MessageEchoes:
			return v[0].Image
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetDocumentMessage() (vDocument *Document) {
	defer func() {
		if r := recover(); r != nil {
			vDocument = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Document
		case []MessageEchoes:
			return v[0].Document
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetStickerMessage() (vSticker *Sticker) {
	defer func() {
		if r := recover(); r != nil {
			vSticker = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Sticker
		case []MessageEchoes:
			return v[0].Sticker
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetLocationMessage() (vLocation *Location) {
	defer func() {
		if r := recover(); r != nil {
			vLocation = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Location
		case []MessageEchoes:
			return v[0].Location
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetContactMessage() (vContact []message.Contact) {
	defer func() {
		if r := recover(); r != nil {
			vContact = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Contacts
		case []MessageEchoes:
			return v[0].Contacts
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetTextMessage() (vText *Text) {
	defer func() {
		if r := recover(); r != nil {
			vText = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Text
		case []MessageEchoes:
			return v[0].Text
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetOrderMessage() (vOrder *Order) {
	defer func() {
		if r := recover(); r != nil {
			vOrder = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Order
		case []MessageEchoes:
			return v[0].Order
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetSystemMessage() (vSystem *System) {
	defer func() {
		if r := recover(); r != nil {
			vSystem = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].System
		case []MessageEchoes:
			return v[0].System
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetReferralMessage() (vReferral *Referral) {
	defer func() {
		if r := recover(); r != nil {
			vReferral = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Referral
		case []MessageEchoes:
			return v[0].Referral
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetErrors() (vErrors []ErrorMessage) {
	defer func() {
		if r := recover(); r != nil {
			vErrors = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Errors
		case []MessageEchoes:
			return v[0].Errors
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetButton() (vButtons *Button) {
	defer func() {
		if r := recover(); r != nil {
			vButtons = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Button
		case []MessageEchoes:
			return v[0].Button
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetInteractive() (vInteractive *Interactive) {
	defer func() {
		if r := recover(); r != nil {
			vInteractive = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Interactive
		case []MessageEchoes:
			return v[0].Interactive
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetIdentity() (vIdentity *Identity) {
	defer func() {
		if r := recover(); r != nil {
			vIdentity = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Identity
		case []MessageEchoes:
			return v[0].Identity
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetContext() (vContext *Context) {
	defer func() {
		if r := recover(); r != nil {
			vContext = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Context
		case []MessageEchoes:
			return v[0].Context
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}

func (self *Components) GetReactionMessage() (vReaction *Reaction) {
	defer func() {
		if r := recover(); r != nil {
			vReaction = nil
		}
	}()

	if self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_MESSAGE.Enum() || self.GetNotificationType() != types.WEBHOOK_NOTIFICATION_SMB_MESSAGE_ECHOES.Enum() {
		switch v := self.GetMessages().(type) {
		case []Message:
			return v[0].Reaction
		case []MessageEchoes:
			return v[0].Reaction
		default:
			return nil
		}
	} else {
		panic("this is not a message")
	}
}
