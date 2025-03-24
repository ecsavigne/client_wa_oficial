package types

type Another struct {
	Type      string `json:"type,omitempty"`
	*PhonesWA `json:",omitempty"`
}

func (a *Another) GetType() string {
	return ResponseAnother
}
func (a *Another) String() string {
	return val(a)
}

func (a *Another) GetResponseAnother() *Another {
	return a
}
func (a *Another) GetResponseError() *Error {
	return nil
}

func (a *Another) GetResponseSuccess() *Success {
	return nil
}

func (a *Another) GetResponseMediaInfo() *MediaInfo {
	return nil
}

func (a *Another) IsType(t ResponseType) bool {
	return ResponseAnother == t
}
