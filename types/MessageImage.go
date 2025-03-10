package types

type MessageImage struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"image"`
}

func NewMessageImage(m MessageImage) Messager {
	mk := &messagerKernel{
		Type:             MessageTypeImage,
		m:                m,
		Link:             m.Link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return &m
}
