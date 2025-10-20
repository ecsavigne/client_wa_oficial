package clientoficial

import (
	"encoding/json"
	"sync"

	"github.com/ecsavigne/client_wa_oficial/v2/event"
)

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
	if !cl.messageIsForMe(msg) {
		return
	}

	switch {
	case len(msg.Entry) != 0 &&
		len(msg.Entry[0].Changes) != 0 &&
		len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
		msg.Entry[0].Changes[0].Value.Messages[0].Type != "":
		switch msg.Entry[0].Changes[0].Value.Messages[0].Type {
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
			cl.Config.EventHandle(message)
		}
	case len(msg.Entry) != 0 &&
		len(msg.Entry[0].Changes) != 0 &&
		len(msg.Entry[0].Changes[0].Value.Statuses) != 0:
		evt = &event.StatusMessageEvent{
			Components: msg,
		}

		recipientId := msg.Entry[0].Changes[0].Value.Statuses[0].RecipientID
		id := msg.Entry[0].Changes[0].Value.Statuses[0].ID

		mu := sync.Mutex{}
		mu.Lock()
		pair, ok := infoContacts[id]
		mu.Unlock()

		if msg.Entry[0].Changes[0].Value.Statuses[0].Status == "failed" &&
			msg.Entry[0].Changes[0].Value.Statuses[0].Errors[0].Message == "Message undeliverable" {
			if ok {
				pair.Channel <- InfoContact{
					ContactPhone: pair.Phone,
					RecipientID:  recipientId,
					IsOnWhats:    false,
				}
			}

		} else {
			if ok {
				pair.Channel <- InfoContact{
					ContactPhone: pair.Phone,
					RecipientID:  recipientId,
					IsOnWhats:    true,
				}
			}
		}
	case len(msg.Entry[0].Changes[0].Value.Messages) != 0 &&
		len(msg.Entry[0].Changes[0].Value.Messages[0].Contacts) != 0:
		evt = &event.MessageContactEvent{
			Components: msg,
		}
	default:
		cl.Config.EventHandle(message)
	}

	cl.Config.EventHandle(evt)

}
