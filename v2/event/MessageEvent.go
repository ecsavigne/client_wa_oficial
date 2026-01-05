package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type MessageEvent struct {
	*Components
}

func (*MessageEvent) GetType() types.EventType { return types.EventTypeMessage }

func (m *MessageEvent) String() string { return response.Val(m) }

type StatusMessageEvent struct {
	*Components
}

func (*StatusMessageEvent) GetType() types.EventType { return types.EventTypeStatusMessage }
func (m *StatusMessageEvent) String() string         { return response.Val(m) }

type MessageAudioEvent struct {
	*Components
}

func (*MessageAudioEvent) GetType() types.EventType { return types.EventTypeMessageAudio }
func (m *MessageAudioEvent) String() string         { return response.Val(m) }

type MessageButtonEvent struct {
	*Components
}

func (*MessageButtonEvent) GetType() types.EventType { return types.EventTypeMessageButton }
func (m *MessageButtonEvent) String() string         { return response.Val(m) }

type MessageDocumentEvent struct {
	*Components
}

func (*MessageDocumentEvent) GetType() types.EventType { return types.EventTypeMessageDocument }
func (m *MessageDocumentEvent) String() string         { return response.Val(m) }

type MessageImageEvent struct {
	*Components
}

func (*MessageImageEvent) GetType() types.EventType { return types.EventTypeMessageImage }
func (m *MessageImageEvent) String() string         { return response.Val(m) }

type MessageInteractiveEvent struct {
	*Components
}

func (*MessageInteractiveEvent) GetType() types.EventType { return types.EventTypeMessageInteractive }
func (m *MessageInteractiveEvent) String() string         { return response.Val(m) }

type MessageOrderEvent struct {
	*Components
}

func (*MessageOrderEvent) GetType() types.EventType { return types.EventTypeMessageOrder }
func (m *MessageOrderEvent) String() string         { return response.Val(m) }

type MessageStickerEvent struct {
	*Components
}

func (*MessageStickerEvent) GetType() types.EventType { return types.EventTypeMessageSticker }
func (m *MessageStickerEvent) String() string         { return response.Val(m) }

type MessageSystemEvent struct {
	*Components
}

func (*MessageSystemEvent) GetType() types.EventType { return types.EventTypeMessageSystem }
func (m *MessageSystemEvent) String() string         { return response.Val(m) }

type MessageTextEvent struct {
	*Components
}

func (*MessageTextEvent) GetType() types.EventType { return types.EventTypeMessageText }
func (m *MessageTextEvent) String() string         { return response.Val(m) }

type MessageUnknownEvent struct {
	*Components
}

func (*MessageUnknownEvent) GetType() types.EventType { return types.EventTypeMessageUnknown }
func (m *MessageUnknownEvent) String() string         { return response.Val(m) }

type MessageVideoEvent struct {
	*Components
}

func (*MessageVideoEvent) GetType() types.EventType { return types.EventTypeMessageVideo }
func (m *MessageVideoEvent) String() string         { return response.Val(m) }

type MessageContactEvent struct {
	*Components
}

func (*MessageContactEvent) GetType() types.EventType { return types.EventTypeMessageContact }
func (m *MessageContactEvent) String() string         { return response.Val(m) }

type MessageReactionEvent struct {
	*Components
}

func (*MessageReactionEvent) GetType() types.EventType { return types.EventTypeMessageReaction }
func (m *MessageReactionEvent) String() string         { return response.Val(m) }

type MessageLocationEvent struct {
	*Components
}

func (*MessageLocationEvent) GetType() types.EventType { return types.EventTypeMessageLocation }
func (m *MessageLocationEvent) String() string         { return response.Val(m) }

type MessageTemplateEvent struct {
	*Components
}

func (*MessageTemplateEvent) GetType() types.EventType { return types.EventTypeTemplateMessage }
func (m *MessageTemplateEvent) String() string         { return response.Val(m) }
