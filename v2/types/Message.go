package types

import (
	"encoding/json"
	"fmt"
	"strings"
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
}

type messagerKernel struct {
	Type             string `json:"type,omitempty"`
	Link             string `json:"link,omitempty"`
	MessagingProduct string `json:"messaging_product,omitempty"`
	m                Messager
}

func (m *messagerKernel) GetType() string {
	return m.Type
}

func (m *messagerKernel) IsTypeResponse() bool {
	return false
}

func toJonReader(m Messager) *strings.Reader {
	databytes, e := json.Marshal(m)
	if e != nil {
		log := fmt.Sprintf("Error in function [types.Message.go -> toJonReader] creating json reader. Error is: %s", e.Error())
		panic(log)
	}
	return strings.NewReader(string(databytes))
}

func (m *messagerKernel) ToJSONReader() *strings.Reader {
	return toJonReader(m.m)
}

func (m *messagerKernel) GetMessagingProduct() string {
	return m.MessagingProduct
}

func (m *messagerKernel) GetMessageLink() string {
	switch v := m.m.(type) {
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
	default:
		return ""
	}
}

func (m *messagerKernel) GetMessageId() string {
	switch v := m.m.(type) {
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
	default:
		return ""
	}
}

func (m *messagerKernel) SetLink(link string) {
	switch v := m.m.(type) {
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
	}
}

func (m *messagerKernel) SetId(id string) {
	switch v := m.m.(type) {
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
	}
}
