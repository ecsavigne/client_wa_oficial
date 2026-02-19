package response

type GranularScope struct {
	Scope     string   `json:"scope,omitempty"`
	TargetIDS []string `json:"target_ids,omitempty"`
}

type Data struct {
	AppID               string          `json:"app_id,omitempty"`
	Type                string          `json:"type,omitempty"`
	Application         string          `json:"application,omitempty"`
	DataAccessExpiresAt int64           `json:"data_access_expires_at,omitempty"`
	ExpiresAt           int64           `json:"expires_at,omitempty"`
	IsValid             bool            `json:"is_valid,omitempty"`
	IssuedAt            int64           `json:"issued_at,omitempty"`
	Scopes              []string        `json:"scopes,omitempty"`
	GranularScopes      []GranularScope `json:"granular_scopes,omitempty"`
	UserID              string          `json:"user_id,omitempty"`
}

type DebugTokenResponse struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	Data         Data   `json:"data"`
}

func NewDebugTokenResponse(config Responser) *DebugTokenResponse {
	if v, ok := config.(*DebugTokenResponse); ok {
		v.ResponseType = ResponseDebugToken
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *DebugTokenResponse")
}
