package clientoficial

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type Client interface {
	GetType() string
	SendMessage(msg message.Messager) response.Responser
}
