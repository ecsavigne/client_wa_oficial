package types

type Context struct {
	MessageId string `json:"message_id" valid:"required"`
}

type Header struct {
	MessagingProduct      string `json:"messaging_product,omitempty" validate:"required"`
	RecipientType         string `json:"recipient_type,omitempty" validate:"required"` // "individual"
	To                    string `json:"to,omitempty" validate:"required"`
	Type                  string `json:"type" validate:"required"` // "text" | "image" | "audio" | "document" | "location" | "video" | "button" | "interactive" | "template" | "sticker" | "contacts" | "reaction"
	Status                string `json:"status,omitempty"`
	BizOpaqueCallbackData string `json:"biz_opaque_callback_data,omitempty"`
	*Context              `json:"context,omitempty"`
}
