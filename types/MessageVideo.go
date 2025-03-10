package types

type MessageVideo struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"video"`
}

func NewMessageVideo(m MessageVideo) Messager {
	mk := &messagerKernel{
		Type:             MessageTypeVideo,
		m:                m,
		Link:             m.Link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return &m
}
