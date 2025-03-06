package types

type MessageVideo struct {
	Messager `json:"messager,omitempty"`
	Header
	Media `json:"video"`
}

func NewMessageVideo(m MessageVideo) Messager {
	mk := &messagerKernel{
		Type: MessageTypeVideo,
		m:    m,
	}

	m.Messager = mk
	return &m
}
