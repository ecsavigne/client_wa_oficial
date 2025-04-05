package message

import (
	"encoding/json"

	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type MessageResponse struct {
	MessagerKernel
	*Media            `json:",omitempty"`
	*Text             `json:"text,omitempty"`
	*InteractiveProto `json:"interactive,omitempty"`
	*Reaction         `json:"reaction,omitempty"`
	*Location         `json:"location,omitempty"`
	*Contact          `json:"contact,omitempty"`
	*Template         `json:"template,omitempty"`
}

func (*MessageResponse) NewResponseMessage(config Messager) *MessageResponse {
	switch v := any(config).(type) {
	case *MessageResponse:
		v.MessagerKernel.parent = v
		return v
	}
	panic("Invalid protocol type, expected *MessageResponse")
}

func (m *MessageResponse) MarshalJSON() ([]byte, error) {
	type Alias MessageResponse // Evitar recursión
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Serializer all except of Media
	data, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	switch {
	case m.MessagerKernel.Type != "audio" && m.MessagerKernel.Type != "Response" &&
		m.MessagerKernel.Type != "video" && m.MessagerKernel.Type != "document" && m.MessagerKernel.Type != "sticker":
		return data, nil
	case m.Media == nil, m.MessagerKernel.Type == "":
		return data, nil // without dinamic field
	}

	// Convertir a mapa para agregar la clave dinámica
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Serializar el Media y agregarlo bajo Header.Type
	mediaData, _ := json.Marshal(m.Media)
	var mediaMap map[string]any
	if err := json.Unmarshal(mediaData, &mediaMap); err != nil {
		return nil, err
	}

	result[m.MessagerKernel.Type] = mediaMap
	var mapKey map[string]any
	b, err := json.Marshal(m.Media)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &mapKey); err != nil {
		return nil, err
	}

	for k := range mapKey {
		delete(result, k)
	}

	return json.Marshal(result)
}

func (m *MessageResponse) String() string {
	return response.Val(m)
}

func (m MessageResponse) IsTypeResponse() bool {
	switch m.MessagerKernel.Type {
	case "audio", "image", "video", "document", "sticker", "interactive", "location", "contact", "text", "template", "reaction":
		return true
	}

	return false
}
