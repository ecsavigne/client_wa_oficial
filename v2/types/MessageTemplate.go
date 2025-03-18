package types

type Language struct {
	Code string `json:"code" validate:"required"` // code language in: https://developers.facebook.com/docs/whatsapp/api/messages/message-templates#language
}

type Currency struct {
	FallbackValue string  `json:"fallback_value" validate:"required"` // Value predeterminate if not work the localization
	Code          string  `json:"code" validate:"required"`           // El código de la divisa, como se define en la norma ISO 4217..Code currency in: https://developers.facebook.com/docs/whatsapp/api/messages/message-templates#currency
	Amount1000    float64 `json:"amount_1000" validate:"required"`    // Importe multiplicado por 1.000.
}

type DateTime struct {
	FallbackValue string `json:"fallback_value" validate:"required"` // Texto predeterminado. En la API de la nube, siempre usamos el valor alternativo y no intentamos localizar usando otros campos opcionales.
}

/*
En las plantillas basadas en texto, los únicos tipos de parámetros admitidos son currency, date_time y text
*/
type Parameter struct {
	Type     string   `json:"type" validate:"required"` // Value:text, image, currency, date_time, video, document, payload (if type of component is button)
	Image    Media    `json:"image,omitempty"`          // use obligatory with type=image
	Video    Media    `json:"video,omitempty"`          // use obligatory with type=video
	Document Media    `json:"document,omitempty"`       // use obligatory with type=document
	Text     string   `json:"text,omitempty"`           // use obligatory with type=text
	Currency Currency `json:"currency,omitempty"`       // use obligatory with type=currency
	DateTime DateTime `json:"date_time,omitempty"`      // use obligatory with type=date_time
	Payload  string   `json:"payload,omitempty"`        // use obligatory with type=payload
}

type Component struct {
	Type       string      `json:"type" validate:"required"` // Value:header, body, footer
	Parameters []Parameter `json:"parameters,omitempty"`
	SubType    string      `json:"sub_type,omitempty"` // use with type=button
	Index      string      `json:"index,omitempty"`    // use with type=button
}

type Template struct {
	/*
		Espacio de nombres de la plantilla.
		A partir de la versión v2.27.8, este debe ser el espacio de nombres asociado a la cuenta de WhatsApp Business que posee el número de teléfono correspondiente al cliente de la API de WhatsApp Business actual, o el mensaje no se enviará
	*/
	Namespace  string      `json:"namespace" validate:"required"`
	Name       string      `json:"name" validate:"required"`
	Language   Language    `json:"language" validate:"required"` // Language in can be represented the template, codes exemples: (pt_BR, en_US, etc)
	Components []Component `json:"components,omitempty"`
}

type MessageTemplate struct {
	Messager `json:"messager,omitempty"`
	Header
	Template `json:"template"`
}

func NewMessageTemplate(m *MessageTemplate) Messager {
	mk := &messagerKernel{
		Type: MessageTypeTemplate,
		m:    m,
	}

	m.Messager = mk
	return m
}
