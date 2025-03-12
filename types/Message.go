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
	SetLink(string)
	SetId(string)
}

type messagerKernel struct {
	Type             string
	Link             string
	MessagingProduct string
	m                Messager
}

func (m *messagerKernel) GetType() string {
	return m.Type
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
	if v, ok := m.m.(MessageImage); ok {
		return v.Media.Link
	}

	if v, ok := m.m.(MessageVideo); ok {
		return v.Media.Link
	}

	if v, ok := m.m.(MessageAudio); ok {
		return v.Media.Link
	}

	if v, ok := m.m.(MessageDocument); ok {
		return v.Media.Link
	}

	if v, ok := m.m.(MessageSticker); ok {
		return v.Media.Link
	}

	return ""
}

func (m *messagerKernel) SetLink(link string) {
	if v, ok := m.m.(MessageImage); ok {
		v.Media.Link = link
		m.m = v
	}

	if v, ok := m.m.(MessageVideo); ok {
		v.Media.Link = link
		m.m = v
	}

	if v, ok := m.m.(MessageAudio); ok {
		v.Media.Link = link
		m.m = v
	}

	if v, ok := m.m.(MessageDocument); ok {
		v.Media.Link = link
		m.m = v
	}

	if v, ok := m.m.(MessageSticker); ok {
		v.Media.Link = link
		m.m = v
	}
}

func (m *messagerKernel) SetId(id string) {
	if v, ok := m.m.(MessageImage); ok {
		v.Media.Id = id
		m.m = v
	}

	if v, ok := m.m.(MessageVideo); ok {
		v.Media.Id = id
		m.m = v
	}

	if v, ok := m.m.(MessageAudio); ok {
		v.Media.Id = id
		m.m = v
	}

	if v, ok := m.m.(MessageDocument); ok {
		v.Media.Id = id
		m.m = v
	}

	if v, ok := m.m.(MessageSticker); ok {
		v.Media.Id = id
		m.m = v
	}
}
