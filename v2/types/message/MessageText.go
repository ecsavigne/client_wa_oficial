package message

type Text struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageText struct {
	MessagerKernel
	Text `json:"text"`
}

func (*MessageText) NewTextMessage(config Messager) *MessageText {
	switch v := any(config).(type) {
	case *MessageText:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageText")
}
