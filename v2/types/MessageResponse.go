package types

import (
	"encoding/json"
)

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
	case m.Header.Type != "audio" && m.Header.Type != "image" &&
		m.Header.Type != "video" && m.Header.Type != "document" && m.Header.Type != "sticker":
		return data, nil
	case m.Media == nil, m.Header.Type == "":
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

	result[m.Header.Type] = mediaMap
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

func NewMessageResponse(m *MessageResponse) Messager {
	link := ""
	switch m.Header.Type {
	case "audio", "image", "video", "document", "sticker":
		link = m.Media.Link
	}
	mk := &messagerKernel{
		Type:             "response",
		m:                m,
		Link:             link,
		MessagingProduct: m.MessagingProduct,
	}

	m.Messager = mk
	return m
}

func (m *MessageResponse) String() string {
	return val(m)
}

func (m *MessageResponse) IsTypeResponse() bool {
	switch m.Header.Type {
	case "audio", "image", "video", "document", "sticker", "interactive", "location", "contact", "text", "template", "reaction":
		return true
	}

	return false
}
