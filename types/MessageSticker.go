package types

/*
Un objeto media que contiene un sticker.
se admiten stickers de salida estáticos y animados de terceros,
además de todos los tipos de stickers de entrada. Un sticker estático debe
contener 512 x 512 píxeles y no puede exceder los 100 KB. Un sticker
animado debe contener 512 x 512 píxeles y no puede exceder los 500 KB.
*/
type MessageSticker struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"sticker"`
}

func NewMessageSticker(m *MessageSticker) Messager {
	mk := &messagerKernel{
		Type:             MessageTypeSticker,
		m:                m,
		Link:             m.Link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return m
}
