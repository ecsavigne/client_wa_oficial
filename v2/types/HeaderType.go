package types

type Header struct {
	MessagingProduct string `json:"messaging_product" validate:"required"`
	RecipientType    string `json:"recipient_type" validate:"required"` // "individual"
	To               string `json:"to" validate:"required"`
	Type             string `json:"type" validate:"required"` // "text" | "image" | "audio" | "document" | "location" | "video" | "button" | "interactive" | "template" | "sticker" | "contacts" | "reaction"
}
