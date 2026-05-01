package response

type ResponseType = string

const (
	ResponseSuccess ResponseType = "response_success"
	ResponseError   ResponseType = "response_error"
	ResponseUnknow  ResponseType = "response_unknow"
)
