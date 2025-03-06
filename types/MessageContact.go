package types

type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
}

type Email struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

// Al menos uno de los parámetros opcionales debe incluirse junto con el parámetro formatted_name.
type Name struct {
	FormattedName string `json:"formatted_name" validate:"required"` // Required
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MiddleName    string `json:"middle_name"`
	Suffix        string `json:"suffix"`
	Prefix        string `json:"prefix"`
}

type Org struct {
	Company    string `json:"company"`
	Department string `json:"department"`
	Title      string `json:"title"`
}

type Phone struct {
	Phone string  `json:"phone"`
	Type  string  `json:"type"`
	WaID  *string `json:"wa_id,omitempty"`
}

type URL struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type Contact struct {
	Addresses []Address `json:"addresses"`
	Birthday  string    `json:"birthday"` // YYYY-MM-DD
	Emails    []Email   `json:"emails"`
	Name      Name      `json:"name" validate:"required"` // Required
	Org       Org       `json:"org"`
	Phones    []Phone   `json:"phones"`
	Urls      []URL     `json:"urls"`
}

type MessageContact struct {
	Messager `json:"messager,omitempty"`
	Header
	Contact []Contact `json:"contacts"`
}

func NewMessageContact(m MessageContact) Messager {
	mk := &messagerKernel{
		Type: MessageTypeContact,
		m:    m,
	}

	m.Messager = mk
	return &m
}
