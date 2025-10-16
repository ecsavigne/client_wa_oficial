package clientoficial

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type Client interface {
	SendMessage(msg message.Messager) response.Responser
}
