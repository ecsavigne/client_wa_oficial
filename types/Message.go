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
	return m.Link
}
