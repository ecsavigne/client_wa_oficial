package message

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
)

type Messager interface {
	GetType() string
	ToJSONReader() *strings.Reader
	GetMessagingProduct() string
	GetMessageLink() string
	GetMessageId() string
	SetLink(string)
	SetId(string)
	IsTypeResponse() bool
	GetFileHeader() *multipart.FileHeader
	ResetFileHeader()
	GetInteractiveMessage() *MessageInteractive
}

type Context struct {
	MessageId string `json:"message_id" valid:"required"`
}

type MessagerKernel struct {
	parent                Messager
	MessagingProduct      string `json:"messaging_product,omitempty"`
	RecipientType         string `json:"recipient_type,omitempty" validate:"required"` // "individual"
	To                    string `json:"to,omitempty" validate:"required"`
	Type                  string `json:"type" validate:"required"` // "text" | "image" | "audio" | "document" | "location" | "video" | "button" | "interactive" | "template" | "sticker" | "contacts" | "reaction"
	Status                string `json:"status,omitempty"`
	BizOpaqueCallbackData string `json:"biz_opaque_callback_data,omitempty"`
	*Context              `json:"context,omitempty"`
}

func NewMessage(config Messager) Messager {
	switch v := config.(type) {
	case *MessageAudio:
		v.MessagerKernel.parent = v
		return v
	case *MessageContact:
		v.MessagerKernel.parent = v
		return v
	case *MessageDocument:
		v.MessagerKernel.parent = v
		return v
	case *MessageImage:
		v.MessagerKernel.parent = v
		return v
	case *MessageInteractive:
		v.MessagerKernel.parent = v
		return v
	case *MessageLocation:
		v.MessagerKernel.parent = v
		return v
	case *MessageReaction:
		v.MessagerKernel.parent = v
		return v
	case *MessageResponse:
		v.MessagerKernel.parent = v
		return v
	case *MessageSticker:
		v.MessagerKernel.parent = v
		return v
	case *MessageTemplate:
		v.MessagerKernel.parent = v
		return v
	case *MessageText:
		v.MessagerKernel.parent = v
		return v
	case *MessageVideo:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageAudio, *MessageContact, *MessageDocument, *MessageImage, *MessageInteractive, *MessageLocation, *MessageReaction, *MessageResponse, *MessageSticker, *MessageTemplate, *MessageText, *MessageVideo")
}

func (m *MessagerKernel) GetType() string {
	switch m.parent.(type) {
	case *MessageImage:
		return wpp.MessageTypeImage
	case *MessageAudio:
		return wpp.MessageTypeAudio
	case *MessageVideo:
		return wpp.MessageTypeVideo
	case *MessageDocument:
		return wpp.MessageTypeDocument
	case *MessageSticker:
		return wpp.MessageTypeSticker
	case *MessageResponse:
		return wpp.MessageTypeResponse
	case *MessageLocation:
		return wpp.MessageTypeLocation
	case *MessageContact:
		return wpp.MessageTypeContact
	case *MessageText:
		return wpp.MessageTypeText
	case *MessageTemplate:
		return wpp.MessageTypeTemplate
	case *MessageReaction:
		return wpp.MessageTypeReaction
	case *MessageInteractive:
		return wpp.MessageTypeInteractive
	}
	panic("Invalid protocol type")
}

func (m *MessagerKernel) GetFileHeader() *multipart.FileHeader {
	switch v := m.parent.(type) {
	case *MessageImage:
		return v.Media.FileHeader
	case *MessageVideo:
		return v.Media.FileHeader
	case *MessageAudio:
		return v.Media.FileHeader
	case *MessageDocument:
		return v.Media.FileHeader
	case *MessageSticker:
		return v.Media.FileHeader
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			return v.Media.FileHeader
		}
		return nil
	default:
		return nil
	}
}

func (m *MessagerKernel) ResetFileHeader() {
	switch v := m.parent.(type) {
	case *MessageImage:
		v.Media.FileHeader = nil
	case *MessageVideo:
		v.Media.FileHeader = nil
	case *MessageAudio:
		v.Media.FileHeader = nil
	case *MessageDocument:
		v.Media.FileHeader = nil
	case *MessageSticker:
		v.Media.FileHeader = nil
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			v.Media.FileHeader = nil
		}
	}
}

func (m *MessagerKernel) IsTypeResponse() bool {
	return false
}

func toJonReader(m Messager) *strings.Reader {
	databytes, e := json.Marshal(m)
	if e != nil {
		log := fmt.Sprintf("Error in function [wpp.Message.go -> toJonReader] creating json reader. Error is: %s", e.Error())
		panic(log)
	}
	return strings.NewReader(string(databytes))
}

func (m *MessagerKernel) ToJSONReader() *strings.Reader {
	if m.parent == nil {
		return nil
	}
	return toJonReader(m.parent)
}

func (m *MessagerKernel) GetMessagingProduct() string {
	return m.MessagingProduct
}

func (m *MessagerKernel) GetMessageLink() string {
	switch v := m.parent.(type) {
	case *MessageImage:
		return v.Media.Link
	case *MessageVideo:
		return v.Media.Link
	case *MessageAudio:
		return v.Media.Link
	case *MessageDocument:
		return v.Media.Link
	case *MessageSticker:
		return v.Media.Link
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			return v.Media.Link
		}
		return ""
	default:
		return ""
	}
}

func (m *MessagerKernel) GetMessageId() string {
	switch v := m.parent.(type) {
	case *MessageImage:
		return v.Media.Id
	case *MessageVideo:
		return v.Media.Id
	case *MessageAudio:
		return v.Media.Id
	case *MessageDocument:
		return v.Media.Id
	case *MessageSticker:
		return v.Media.Id
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			return v.Media.Id
		}
		return ""
	default:
		return ""
	}
}

func (m *MessagerKernel) SetLink(link string) {
	switch v := m.parent.(type) {
	case *MessageImage:
		v.Media.Link = link
	case *MessageVideo:
		v.Media.Link = link
	case *MessageAudio:
		v.Media.Link = link
	case *MessageDocument:
		v.Media.Link = link
	case *MessageSticker:
		v.Media.Link = link
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			v.Media.Link = link
		}
	}
}

func (m *MessagerKernel) SetId(id string) {
	switch v := m.parent.(type) {
	case *MessageImage:
		v.Media.Id = id
	case *MessageVideo:
		v.Media.Id = id
	case *MessageAudio:
		v.Media.Id = id
	case *MessageDocument:
		v.Media.Id = id
	case *MessageSticker:
		v.Media.Id = id
	case *MessageResponse:
		switch v.MessagerKernel.Type {
		case "audio", "image", "video", "document", "sticker":
			v.Media.Id = id
		}
	}
}

func (m *MessagerKernel) GetInteractiveMessage() *MessageInteractive {
	switch v := m.parent.(type) {
	case *MessageInteractive:
		return v
	default:
		return nil
	}
}
