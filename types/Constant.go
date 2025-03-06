package types

type (
	MessageType     = string
	interactiveType = string
)

const (
	MessageTypeText        MessageType = "text"
	MessageTypeTemplate    MessageType = "template"
	MessageTypeReaction    MessageType = "reaction"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeAudio       MessageType = "audio"
	MessageTypeImage       MessageType = "image"
	MessageTypeVideo       MessageType = "video"
	MessageTypeDocument    MessageType = "document"
	MessageTypeLocation    MessageType = "location"
	MessageTypeContact     MessageType = "contacts"
	MessageTypeSticker     MessageType = "sticker"
)

const (
	InteractiveTypeProduct        interactiveType = "product"
	InteractiveTypeMultiProduct   interactiveType = "product_list"
	InteractiveTypeCatalog        interactiveType = "catalog_message"
	InteractiveTypeProcess        interactiveType = "flow"
	InteractiveTypeButtonUrl      interactiveType = "cta_url"
	InteractiveTypeButtonResponse interactiveType = "button"
	InteractiveTypeList           interactiveType = "list"
)
