package types

type GeneralResponse struct {
	Type       string `json:"type,omitempty"`
	*PhonesWA  `json:",omitempty"`
	*Error     `json:",omitempty"`
	*Success   `json:",omitempty"`
	*MediaInfo `json:",omitempty"`
}

func (a *GeneralResponse) GetType() string {
	return ResponseGeneralResponse
}
func (a *GeneralResponse) String() string {
	return val(a)
}

func (a *GeneralResponse) GetGeneralResponse() *GeneralResponse {
	return a
}
func (a *GeneralResponse) GetResponseError() *Error {
	return nil
}

func (a *GeneralResponse) GetResponseSuccess() *Success {
	return nil
}

func (a *GeneralResponse) GetResponseMediaInfo() *MediaInfo {
	return nil
}

func (a *GeneralResponse) GetResponsePhonesWA() *PhonesWA {
	return nil
}

func (a *GeneralResponse) GetResponseType() ResponserRequest {
	switch {
	case a.PhonesWA != nil:
		return a
	case a.Error != nil:
		return a.Error
	case a.Success != nil:
		return a.Success
	case a.MediaInfo != nil:
		return a.MediaInfo
	}

	return nil
}

func (a *GeneralResponse) IsType(t ResponseType) bool {
	return ResponseGeneralResponse == t
}
