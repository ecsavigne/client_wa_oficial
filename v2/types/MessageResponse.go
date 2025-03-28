package types

type MessageResponse struct {
	Messager `json:"messager,omitempty"`
	Header
	*Media            `json:",omitempty"`
	*Text             `json:"text,omitempty"`
	*InteractiveProto `json:"interactive,omitempty"`
	*Reaction         `json:"reaction,omitempty"`
	*Location         `json:"location,omitempty"`
	*Contact          `json:"contact,omitempty"`
	*Template         `json:"template,omitempty"`
}

func NewMessageResponse(m *MessageResponse) Messager {
	mk := &messagerKernel{
		Type: "response",
		m:    m,
	}

	m.Messager = mk
	return m
}

func (m *MessageResponse) IsTypeResponse() bool {
	switch m.Header.Type {
	case "audio", "image", "video", "document", "sticker", "interactive", "location", "contact", "text", "template", "reaction":
		return true
	}

	return false
}
