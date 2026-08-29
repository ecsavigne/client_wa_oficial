package event

import (
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type MessageEvent struct {
	*igpbv1.InstagramWebhookEvent
}

// func (*MessageEvent) GetType() igpbv1.EventType { return igpbv1.EventTypeMessage }

// func (m *MessageEvent) String() string { return response.Val(m) }

type IGStatusMessageEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGStatusMessageEvent) GetType() string  { return IGEventTypeStatusMessage }
func (m *IGStatusMessageEvent) String() string { return response.Val(m) }
func (m *IGStatusMessageEvent) GetMessageStatus() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageAudioEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageAudioEvent) GetType() string  { return IGEventTypeMessageAudio }
func (m *IGMessageAudioEvent) String() string { return response.Val(m) }
func (m *IGMessageAudioEvent) GetMessageAudio() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageFileEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageFileEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageFileEvent) String() string { return response.Val(m) }
func (m *IGMessageFileEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageImageEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageImageEvent) GetType() string  { return IGEventTypeMessageImage }
func (m *IGMessageImageEvent) String() string { return response.Val(m) }
func (m *IGMessageImageEvent) GetMessageImage() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageInteractiveEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageInteractiveEvent) GetType() string  { return IGEventTypeMessageInteractive }
func (m *IGMessageInteractiveEvent) String() string { return response.Val(m) }
func (m *IGMessageInteractiveEvent) GetMessageInteractive() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageOrderEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageOrderEvent) GetType() string  { return IGEventTypeMessageOrder }
func (m *IGMessageOrderEvent) String() string { return response.Val(m) }
func (m *IGMessageOrderEvent) GetMessageOrder() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageStickerEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageStickerEvent) GetType() string  { return IGEventTypeMessageSticker }
func (m *IGMessageStickerEvent) String() string { return response.Val(m) }
func (m *IGMessageStickerEvent) GetMessageSticker() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageSystemEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageSystemEvent) GetType() string  { return IGEventTypeMessageSystem }
func (m *IGMessageSystemEvent) String() string { return response.Val(m) }
func (m *IGMessageSystemEvent) GetMessageSystem() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageTextEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageTextEvent) GetType() string  { return IGEventTypeMessageText }
func (m *IGMessageTextEvent) String() string { return response.Val(m) }
func (m *IGMessageTextEvent) GetMessageText() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageUnknownEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageUnknownEvent) GetType() string  { return IGEventTypeMessageUnknown }
func (m *IGMessageUnknownEvent) String() string { return response.Val(m) }
func (m *IGMessageUnknownEvent) GetMessageUnknown() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageVideoEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageVideoEvent) GetType() string  { return IGEventTypeMessageVideo }
func (m *IGMessageVideoEvent) String() string { return response.Val(m) }
func (m *IGMessageVideoEvent) GetMessageVideo() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageContactEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageContactEvent) GetType() string  { return IGEventTypeMessageContact }
func (m *IGMessageContactEvent) String() string { return response.Val(m) }
func (m *IGMessageContactEvent) GetMessageContact() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageReactionEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageReactionEvent) GetType() string  { return IGEventTypeMessageReaction }
func (m *IGMessageReactionEvent) String() string { return response.Val(m) }
func (m *IGMessageReactionEvent) GetMessageReaction() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageLocationEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageLocationEvent) GetType() string  { return IGEventTypeMessageLocation }
func (m *IGMessageLocationEvent) String() string { return response.Val(m) }
func (m *IGMessageLocationEvent) GetMessageLocation() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageButtonEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageButtonEvent) GetType() string  { return IGEventTypeMessageButton }
func (m *IGMessageButtonEvent) String() string { return response.Val(m) }
func (m *IGMessageButtonEvent) GetMessageButton() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageTemplateEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageTemplateEvent) GetType() string  { return IGEventTypeTemplateMessage }
func (m *IGMessageTemplateEvent) String() string { return response.Val(m) }
func (m *IGMessageTemplateEvent) GetMessageTemplate() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageMediaEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageMediaEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageMediaEvent) String() string { return response.Val(m) }
func (m *IGMessageMediaEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageShareEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageShareEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageShareEvent) String() string { return response.Val(m) }
func (m *IGMessageShareEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageIGPostEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageIGPostEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageIGPostEvent) String() string { return response.Val(m) }
func (m *IGMessageIGPostEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageStoryMentionEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageStoryMentionEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageStoryMentionEvent) String() string { return response.Val(m) }
func (m *IGMessageStoryMentionEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageIGReelEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageIGReelEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageIGReelEvent) String() string { return response.Val(m) }
func (m *IGMessageIGReelEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageReelEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageReelEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageReelEvent) String() string { return response.Val(m) }
func (m *IGMessageReelEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageStoryEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageStoryEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageStoryEvent) String() string { return response.Val(m) }
func (m *IGMessageStoryEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageIGStoryEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageIGStoryEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageIGStoryEvent) String() string { return response.Val(m) }
func (m *IGMessageIGStoryEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageEditEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageEditEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageEditEvent) String() string { return response.Val(m) }
func (m *IGMessageEditEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessageEphemeralEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessageEphemeralEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessageEphemeralEvent) String() string { return response.Val(m) }
func (m *IGMessageEphemeralEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessagingPostbacksEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessagingPostbacksEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessagingPostbacksEvent) String() string { return response.Val(m) }
func (m *IGMessagingPostbacksEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMessagingReferralEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMessagingReferralEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMessagingReferralEvent) String() string { return response.Val(m) }
func (m *IGMessagingReferralEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGCommentsEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGCommentsEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGCommentsEvent) String() string { return response.Val(m) }
func (m *IGCommentsEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGLiveCommentsEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGLiveCommentsEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGLiveCommentsEvent) String() string { return response.Val(m) }
func (m *IGLiveCommentsEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}

type IGMentionsEvent struct {
	*igpbv1.InstagramWebhookEvent
	IGAccountID string
}

func (*IGMentionsEvent) GetType() string  { return IGEventTypeMessageDocument }
func (m *IGMentionsEvent) String() string { return response.Val(m) }
func (m *IGMentionsEvent) GetMessageDocument() *igpbv1.InstagramWebhookEvent {
	return m.InstagramWebhookEvent
}
