package message

type MessageImage struct {
	MessagerKernel
	Media `json:"image"`
}

func (*MessageImage) NewImageMessage(config Messager) *MessageImage {
	switch v := any(config).(type) {
	case *MessageImage:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageImage")
}
