package client

import (
	"net/http"

	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	"google.golang.org/protobuf/proto"
)

type CLIENT_TYPE string

const (
	CLIENT_WHATSAPP  CLIENT_TYPE = "whatsapp"
	CLIENT_FACEBOOK  CLIENT_TYPE = "facebook"
	CLIENT_IG        CLIENT_TYPE = "instagram"
	CLIENT_MESSENGER CLIENT_TYPE = "messenger"
)

func (c CLIENT_TYPE) String() string {
	return string(c)
}

type ConfigClient interface {
	GetToken() string
	GetVersion() string
	SetBaseUrl(string)
	GetBaseUrl() string
	String() string
	GetType() TYPE_CONFIG
	GetUserID() string
}

type Client interface {
	GetType() string
	GetConfig() ConfigClient
	SendMessage(msg proto.Message) response.Responser
	String() string
	MultipartRequest(method string, data proto.Message, ePoint string) (*http.Request, error)
}
