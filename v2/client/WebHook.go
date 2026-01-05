package clientoficial

import (
	"encoding/json"

	"github.com/ecsavigne/client_wa_oficial/v2/event"
	evt_types "github.com/ecsavigne/client_wa_oficial/v2/event/types"
)

func getTypeMessage(msg *event.Components) (typ string) {
	defer func() {
		if r := recover(); r != nil {
			typ = ""
		}
	}()

	return msg.Entry[0].Changes[0].Value.Messages[0].Type
}

func getSatusMessage(msg *event.Components) (status string) {
	defer func() {
		if r := recover(); r != nil {
			status = ""
		}
	}()

	return msg.Entry[0].Changes[0].Value.Statuses[0].Status
}

// sent, delivered, read, failed, deleted, warning
func isVailidStatusMessage(status string) bool {
	if status == "read" || status == "delivered" || status == "sent" || status == "failed" ||
		status == "deleted" || status == "warning" {
		return true
	}

	return false
}

func (cl *ClientWA) Broadcast(data map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	var evt event.EventInterface

	// Listener message of the server way WebHook
	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	msg := codeWebHook(message)
	isForme, typeNotification := cl.messageIsForMe(msg)
	if !isForme {
		return
	}

	switch typeNotification {
	// case len(msg.Entry) != 0 &&
	// 	len(msg.Entry[0].Changes) != 0 &&
	// 	len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
	// 	msg.Entry[0].Changes[0].Value.Messages[0].Type != ""
	// :
	case evt_types.WEBHOOK_NOTIFICATION_MESSAGE:
		switch getTypeMessage(msg) {
		case "audio":
			evt = &event.MessageAudioEvent{
				Components: msg,
			}
		case "button":
			evt = &event.MessageButtonEvent{
				Components: msg,
			}
		case "document":
			evt = &event.MessageDocumentEvent{
				Components: msg,
			}
		case "text":
			evt = &event.MessageTextEvent{
				Components: msg,
			}
		case "image":
			evt = &event.MessageImageEvent{
				Components: msg,
			}
		case "interactive":
			evt = &event.MessageInteractiveEvent{
				Components: msg,
			}
		case "order":
			evt = &event.MessageOrderEvent{
				Components: msg,
			}
		case "sticker":
			evt = &event.MessageStickerEvent{
				Components: msg,
			}
		case "system":
			evt = &event.MessageSystemEvent{
				Components: msg,
			}
		case "video":
			evt = &event.MessageVideoEvent{
				Components: msg,
			}
		case "reaction":
			evt = &event.MessageReactionEvent{
				Components: msg,
			}
		case "location":
			evt = &event.MessageLocationEvent{
				Components: msg,
			}
		case "contacts":
			evt = &event.MessageContactEvent{
				Components: msg,
			}
		case "unknown":
			evt = &event.MessageUnknownEvent{
				Components: msg,
			}
		default:
			// can be status message or another notification about message
			status := getSatusMessage(msg)

			switch {
			case status != "":
				if isVailidStatusMessage(status) {
					evt = &event.StatusMessageEvent{
						Components: msg,
					}
				}
			default:
				cl.Config.EventHandle(message)
			}
		}
	// case len(msg.Entry) != 0 &&
	// 	len(msg.Entry[0].Changes) != 0 &&
	// 	len(msg.Entry[0].Changes[0].Value.Statuses) != 0:
	// 	evt = &event.StatusMessageEvent{
	// 		Components: msg,
	// 	}

	// 	recipientId := msg.Entry[0].Changes[0].Value.Statuses[0].RecipientID
	// 	id := msg.Entry[0].Changes[0].Value.Statuses[0].ID

	// 	mu := sync.Mutex{}
	// 	mu.Lock()
	// 	pair, ok := infoContacts[id]
	// 	mu.Unlock()

	// 	if msg.Entry[0].Changes[0].Value.Statuses[0].Status == "failed" &&
	// 		msg.Entry[0].Changes[0].Value.Statuses[0].Errors[0].Message == "Message undeliverable" {
	// 		if ok {
	// 			pair.Channel <- InfoContact{
	// 				ContactPhone: pair.Phone,
	// 				RecipientID:  recipientId,
	// 				IsOnWhats:    false,
	// 			}
	// 		}

	// 	} else {
	// 		if ok {
	// 			pair.Channel <- InfoContact{
	// 				ContactPhone: pair.Phone,
	// 				RecipientID:  recipientId,
	// 				IsOnWhats:    true,
	// 			}
	// 		}
	// 	}
	// case len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
	// 	len(msg.Entry[0].Changes[0].Value.Messages[0].Contacts) != 0:
	// 	evt = &event.MessageContactEvent{
	// 		Components: msg,
	// 	}
	case evt_types.WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_CATEGORY,
		evt_types.WEBHOOK_NOTIFICATION_TEMPLATE_UPDATE_STATUS:
		evt = &event.MessageTemplateEvent{
			Components: msg,
		}
	default:
		cl.Config.EventHandle(message)
		return
	}

	cl.Config.EventHandle(evt)

}
