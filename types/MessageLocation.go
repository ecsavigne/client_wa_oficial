package types

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
}

type MessageLocation struct {
	Messager `json:"messager,omitempty"`
	Header
	Location `json:"location,omitempty"`
}

func NewMessageLocation(m *MessageLocation) Messager {
	mk := &messagerKernel{
		Type: MessageTypeLocation,
		m:    m,
	}

	m.Messager = mk
	return m
}
