package types

type CATEGORY = string
type PARAMETER_FORMAT = string
type LANGUAGE = string

const (
	HEADER CATEGORY = "header"
	BODY   CATEGORY = "body"
	FOOTER CATEGORY = "footer"
)

const (
	TEXT   PARAMETER_FORMAT = "text"
	NUMBER PARAMETER_FORMAT = "number"
	DATE   PARAMETER_FORMAT = "date"
)

const (
	ENGLISH    LANGUAGE = "en_US"
	SPANISH    LANGUAGE = "es_ES"
	PORTUGUESE LANGUAGE = "pt_BR"
)

type ArrayButton = []Button

type PositionalParams = []string
type NamedParam struct {
	ParamName string `json:"param_name" validate:"required"`
	Example   string `json:"example" validate:"required"`
}

type HeaderText = []PositionalParams
type HeaderHandle = []string
type HeaderTextNamedParam = []NamedParam
type BodyTextNamedParam = []NamedParam
type BodyText = []PositionalParams

type Example struct {
	*HeaderHandle         `json:"header_handle,omitempty"`
	*HeaderText           `json:"header_text,omitempty"`
	*BodyText             `json:"body_text,omitempty"`
	*BodyTextNamedParam   `json:"body_text_named_params,omitempty"`
	*HeaderTextNamedParam `json:"header_text_named_params,omitempty"`
}

type Header struct {
	// type of assets of media content. Configurable to "IMAGE", "VIDEO" o "DOCUMENT"
	Format string `json:"format" validate:"required"`
}

type MockupComponent struct {
	Type string `json:"type" validate:"required"`
	// accept parameters for programming ex:
	// 		Posicional: se incluye una matriz de parámetros posicionales numerados que corresponden a posiciones numéricas en el texto del cuerpo con ejemplos.
	// 		Por ejemplo: “Hello {{1}}, your account balance is {{2}}” | [ “John”, “$1,000” ]
	// 		Nominal: se incluyen objetos JSON que contengan un parámetro con nombre y ejemplos.
	// Footer not accept parameters
	// Por ejemplo: { "param_name": "order_id", "example": "335628"}
	// less than 60 characters in header and footer, 1024 characters in body,
	Text         string `json:"text,omitempty"`
	*Example     `json:"example,omitempty"`
	*Header      `json:",omitempty"`
	*ArrayButton `json:"buttons,omitempty"`
}

type ResponseMT struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
}

// All template have limit of one body component
type MockupTemplate struct {
	// less than 512 characters
	Name             string            `json:"name" validate:"required"`
	Category         CATEGORY          `json:"category" validate:"required"`
	PreviousCategory string            `json:"previous_category,omitempty"`
	SubCategory      string            `json:"sub_category,omitempty"`
	ParameterFormat  PARAMETER_FORMAT  `json:"parameter_format,omitempty"`
	Language         LANGUAGE          `json:"language" validate:"required"`
	Components       []MockupComponent `json:"components,omitempty"`
	*ResponseMT      `json:",omitempty"`
}
