package message

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
}

type MessageLocation struct {
	MessagerKernel
	Location `json:"location"`
}

func (*MessageLocation) NewLocationMessage(config Messager) *MessageLocation {
	switch v := any(config).(type) {
	case *MessageLocation:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageLocation")
}
