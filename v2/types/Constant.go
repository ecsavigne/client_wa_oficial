package types

type (
	MessageType     = string
	interactiveType = string
)

const (
	MessageTypeAudio       MessageType = "audio"
	MessageTypeContact     MessageType = "contacts"
	MessageTypeDocument    MessageType = "document"
	MessageTypeImage       MessageType = "image"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeLocation    MessageType = "location"
	MessageTypeReaction    MessageType = "reaction"
	MessageTypeResponse    MessageType = "response"
	MessageTypeSticker     MessageType = "sticker"
	MessageTypeTemplate    MessageType = "template"
	MessageTypeText        MessageType = "text"
	MessageTypeVideo       MessageType = "video"
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

// Event types
type EventType = string

const (
	EventTypeErrorSocketConnect EventType = "ErrorSocketConnectEvent"
	EventTypeMessageAudio       EventType = "EventTypeMessageAudio"
	EventTypeMessageButton      EventType = "EventTypeMessageButton"
	EventTypeMessageDocument    EventType = "EventTypeMessageDocument"
	EventTypeMessageContact     EventType = "EventTypeMessageContact"
	EventTypeMessageText        EventType = "EventTypeMessageText"
	EventTypeMessageImage       EventType = "EventTypeMessageImage"
	EventTypeMessageInteractive EventType = "EventTypeMessageInteractive"
	EventTypeMessageOrder       EventType = "EventTypeMessageOrder"
	EventTypeMessageSticker     EventType = "EventTypeMessageSticker"
	EventTypeMessageSystem      EventType = "EventTypeMessageSystem"
	EventTypeMessageVideo       EventType = "EventTypeMessageVideo"
	EventTypeMessageUnknown     EventType = "EventTypeMessageUnknown"
	EventTypeMessageReaction    EventType = "EventTypeMessageReaction"
	EventTypeMessageLocation    EventType = "EventTypeMessageLocation"
	EventTypeMessage            EventType = "EventTypeMessage"
	EventTypeStatusMessage      EventType = "EventTypeStatusMessage"
)
