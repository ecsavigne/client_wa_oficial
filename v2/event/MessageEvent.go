package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type MessageEvent struct {
	*MessageWebhook
}

func (*MessageEvent) GetType() types.EventType { return types.EventTypeMessage }

func (m *MessageEvent) String() string { return response.Val(m) }

type StatusMessageEvent struct {
	*MessageWebhook
}

func (*StatusMessageEvent) GetType() types.EventType { return types.EventTypeStatusMessage }
func (m *StatusMessageEvent) String() string         { return response.Val(m) }

type MessageAudioEvent struct {
	*MessageWebhook
}

func (*MessageAudioEvent) GetType() types.EventType { return types.EventTypeMessageAudio }
func (m *MessageAudioEvent) String() string         { return response.Val(m) }

type MessageButtonEvent struct {
	*MessageWebhook
}

func (*MessageButtonEvent) GetType() types.EventType { return types.EventTypeMessageButton }
func (m *MessageButtonEvent) String() string         { return response.Val(m) }

type MessageDocumentEvent struct {
	*MessageWebhook
}

func (*MessageDocumentEvent) GetType() types.EventType { return types.EventTypeMessageDocument }
func (m *MessageDocumentEvent) String() string         { return response.Val(m) }

type MessageImageEvent struct {
	*MessageWebhook
}

func (*MessageImageEvent) GetType() types.EventType { return types.EventTypeMessageImage }
func (m *MessageImageEvent) String() string         { return response.Val(m) }

type MessageInteractiveEvent struct {
	*MessageWebhook
}

func (*MessageInteractiveEvent) GetType() types.EventType { return types.EventTypeMessageInteractive }
func (m *MessageInteractiveEvent) String() string         { return response.Val(m) }

type MessageOrderEvent struct {
	*MessageWebhook
}

func (*MessageOrderEvent) GetType() types.EventType { return types.EventTypeMessageOrder }
func (m *MessageOrderEvent) String() string         { return response.Val(m) }

type MessageStickerEvent struct {
	*MessageWebhook
}

func (*MessageStickerEvent) GetType() types.EventType { return types.EventTypeMessageSticker }
func (m *MessageStickerEvent) String() string         { return response.Val(m) }

type MessageSystemEvent struct {
	*MessageWebhook
}

func (*MessageSystemEvent) GetType() types.EventType { return types.EventTypeMessageSystem }
func (m *MessageSystemEvent) String() string         { return response.Val(m) }

type MessageTextEvent struct {
	*MessageWebhook
}

func (*MessageTextEvent) GetType() types.EventType { return types.EventTypeMessageText }
func (m *MessageTextEvent) String() string         { return response.Val(m) }

type MessageUnknownEvent struct {
	*MessageWebhook
}

func (*MessageUnknownEvent) GetType() types.EventType { return types.EventTypeMessageUnknown }
func (m *MessageUnknownEvent) String() string         { return response.Val(m) }

type MessageVideoEvent struct {
	*MessageWebhook
}

func (*MessageVideoEvent) GetType() types.EventType { return types.EventTypeMessageVideo }
func (m *MessageVideoEvent) String() string         { return response.Val(m) }

type MessageContactEvent struct {
	*MessageWebhook
}

func (*MessageContactEvent) GetType() types.EventType { return types.EventTypeMessageContact }
func (m *MessageContactEvent) String() string         { return response.Val(m) }

type MessageReactionEvent struct {
	*MessageWebhook
}

func (*MessageReactionEvent) GetType() types.EventType { return types.EventTypeMessageReaction }
func (m *MessageReactionEvent) String() string         { return response.Val(m) }

type MessageLocationEvent struct {
	*MessageWebhook
}

func (*MessageLocationEvent) GetType() types.EventType { return types.EventTypeMessageLocation }
func (m *MessageLocationEvent) String() string         { return response.Val(m) }
