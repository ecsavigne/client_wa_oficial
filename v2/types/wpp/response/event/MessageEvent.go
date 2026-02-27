package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type MessageEvent struct {
	*Components
}

func (*MessageEvent) GetType() wpp.EventType { return wpp.EventTypeMessage }

func (m *MessageEvent) String() string { return response.Val(m) }

type StatusMessageEvent struct {
	*Components
}

func (*StatusMessageEvent) GetType() wpp.EventType { return wpp.EventTypeStatusMessage }
func (m *StatusMessageEvent) String() string       { return response.Val(m) }

type MessageAudioEvent struct {
	*Components
}

func (*MessageAudioEvent) GetType() wpp.EventType { return wpp.EventTypeMessageAudio }
func (m *MessageAudioEvent) String() string       { return response.Val(m) }

type MessageDocumentEvent struct {
	*Components
}

func (*MessageDocumentEvent) GetType() wpp.EventType { return wpp.EventTypeMessageDocument }
func (m *MessageDocumentEvent) String() string       { return response.Val(m) }

type MessageImageEvent struct {
	*Components
}

func (*MessageImageEvent) GetType() wpp.EventType { return wpp.EventTypeMessageImage }
func (m *MessageImageEvent) String() string       { return response.Val(m) }

type MessageInteractiveEvent struct {
	*Components
}

func (*MessageInteractiveEvent) GetType() wpp.EventType { return wpp.EventTypeMessageInteractive }
func (m *MessageInteractiveEvent) String() string       { return response.Val(m) }

type MessageOrderEvent struct {
	*Components
}

func (*MessageOrderEvent) GetType() wpp.EventType { return wpp.EventTypeMessageOrder }
func (m *MessageOrderEvent) String() string       { return response.Val(m) }

type MessageStickerEvent struct {
	*Components
}

func (*MessageStickerEvent) GetType() wpp.EventType { return wpp.EventTypeMessageSticker }
func (m *MessageStickerEvent) String() string       { return response.Val(m) }

type MessageSystemEvent struct {
	*Components
}

func (*MessageSystemEvent) GetType() wpp.EventType { return wpp.EventTypeMessageSystem }
func (m *MessageSystemEvent) String() string       { return response.Val(m) }

type MessageTextEvent struct {
	*Components
}

func (*MessageTextEvent) GetType() wpp.EventType { return wpp.EventTypeMessageText }
func (m *MessageTextEvent) String() string       { return response.Val(m) }

type MessageUnknownEvent struct {
	*Components
}

func (*MessageUnknownEvent) GetType() wpp.EventType { return wpp.EventTypeMessageUnknown }
func (m *MessageUnknownEvent) String() string       { return response.Val(m) }

type MessageVideoEvent struct {
	*Components
}

func (*MessageVideoEvent) GetType() wpp.EventType { return wpp.EventTypeMessageVideo }
func (m *MessageVideoEvent) String() string       { return response.Val(m) }

type MessageContactEvent struct {
	*Components
}

func (*MessageContactEvent) GetType() wpp.EventType { return wpp.EventTypeMessageContact }
func (m *MessageContactEvent) String() string       { return response.Val(m) }

type MessageReactionEvent struct {
	*Components
}

func (*MessageReactionEvent) GetType() wpp.EventType { return wpp.EventTypeMessageReaction }
func (m *MessageReactionEvent) String() string       { return response.Val(m) }

type MessageLocationEvent struct {
	*Components
}

func (*MessageLocationEvent) GetType() wpp.EventType { return wpp.EventTypeMessageLocation }
func (m *MessageLocationEvent) String() string       { return response.Val(m) }
