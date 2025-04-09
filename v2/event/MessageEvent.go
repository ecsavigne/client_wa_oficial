package event

import "github.com/ecsavigne/client_wa_oficial/v2/types"

type MessageEvent struct {
	*MessageWebhook
}

func (*MessageEvent) GetType() types.EventType { return types.EventTypeMessage }

type StatusMessageEvent struct {
	*MessageWebhook
}

func (*StatusMessageEvent) GetType() types.EventType { return types.EventTypeStatusMessage }

type MessageAudioEvent struct {
	*MessageWebhook
}

func (*MessageAudioEvent) GetType() types.EventType { return types.EventTypeMessageAudio }

type MessageButtonEvent struct {
	*MessageWebhook
}

func (*MessageButtonEvent) GetType() types.EventType { return types.EventTypeMessageButton }

type MessageDocumentEvent struct {
	*MessageWebhook
}

func (*MessageDocumentEvent) GetType() types.EventType { return types.EventTypeMessageDocument }

type MessageImageEvent struct {
	*MessageWebhook
}

func (*MessageImageEvent) GetType() types.EventType { return types.EventTypeMessageImage }

type MessageInteractiveEvent struct {
	*MessageWebhook
}

func (*MessageInteractiveEvent) GetType() types.EventType { return types.EventTypeMessageInteractive }

type MessageOrderEvent struct {
	*MessageWebhook
}

func (*MessageOrderEvent) GetType() types.EventType { return types.EventTypeMessageOrder }

type MessageStickerEvent struct {
	*MessageWebhook
}

func (*MessageStickerEvent) GetType() types.EventType { return types.EventTypeMessageSticker }

type MessageSystemEvent struct {
	*MessageWebhook
}

func (*MessageSystemEvent) GetType() types.EventType { return types.EventTypeMessageSystem }

type MessageTextEvent struct {
	*MessageWebhook
}

func (*MessageTextEvent) GetType() types.EventType { return types.EventTypeMessageText }

type MessageUnknownEvent struct {
	*MessageWebhook
}

func (*MessageUnknownEvent) GetType() types.EventType { return types.EventTypeMessageUnknown }

type MessageVideoEvent struct {
	*MessageWebhook
}

func (*MessageVideoEvent) GetType() types.EventType { return types.EventTypeMessageVideo }

type MessageContatEvent struct {
	*MessageWebhook
}

func (*MessageContatEvent) GetType() types.EventType { return types.EventTypeMessageContact }

type MessageReactionEvent struct {
	*MessageWebhook
}

func (*MessageReactionEvent) GetType() types.EventType { return types.EventTypeMessageReaction }

type MessageLocationEvent struct {
	*MessageWebhook
}

func (*MessageLocationEvent) GetType() types.EventType { return types.EventTypeMessageLocation }
