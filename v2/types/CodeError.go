package types

// Message string
type MessageError = string
type TypeError = string
type CodeError = int64

// MessageErrors
const (
	MsgErrorEnvNotFound     MessageError = "env_not_found or env_not_set"
	MsgErrorBaseUrlEmpty    MessageError = "BaseUrl (WA_BASE_URL) is empty in .env file"
	MsgErrorApiVersionEmpty MessageError = "Cloud API Version (CLOUD_API_VERSION) is empty in .env file"
	MsgErrorTlsInternal     MessageError = "remote error: tls: internal error"
	MsgErrorBadHandshake    MessageError = "websocket: bad handshake"
	MsgErrorUnexpectedClose MessageError = "Is Unexpected Close Error in socket"
	MsgErrorPhoneIdEmpty    MessageError = "WA Phone number ID (WA_PHONE_NUMBER_ID) is empty in .env file"
	MsgErrorBadRequest      MessageError = "Bad Request"
	MsgErrorUrlNotFound     MessageError = "Url Not Found"
	MsgErrorUnauthorized    MessageError = "Unauthorized"
)

// CodeErrors
const (
	CodeErrorEnvNotFound     CodeError = 1
	CodeErrorTlsInternal     CodeError = 2
	CodeErrorBadHandshake    CodeError = 3
	CodeErrorUnexpectedClose CodeError = 4
	CodeErrorBaseUrlEmpty    CodeError = 5
	CodeErrorApiVersionEmpty CodeError = 6
	CodeErrorPhoneIdEmpty    CodeError = 7
	CodeErrorBadRequest      CodeError = 400
	CodeErrorUnauthorized    CodeError = 401
	CodeErrorUrlNotFound     CodeError = 404
)

// TypeErrors
const (
	TypeErrorConfig          TypeError = "env_not_found"
	TypeErrorTlsInternal     TypeError = "tls_internal"
	TypeErrorBadHandshake    TypeError = "bad_handshake"
	TypeErrorUnexpectedClose TypeError = "unexpected_close"
	TypeErrorBaseUrlEmpty    TypeError = "base_url_empty"
	TypeErrorApiVersionEmpty TypeError = "api_version_empty"
	TypeErrorPhoneIdEmpty    TypeError = "phone_id_empty"
	TypeErrorBadRequest      TypeError = "bad_request"
	TypeErrorUnauthorized    TypeError = "unauthorized"
	TypeErrorUrlNotFound     TypeError = "url_not_found"
)
