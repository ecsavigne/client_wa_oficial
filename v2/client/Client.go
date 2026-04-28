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
	/*
		@typeInfo can be "account_business" for get info of account business, for example in IG, in WhatsApp can be "phone_number_id" for get info of phone number id
	*/
	Get(typeInfo string) response.Responser

	/*
		@typeDelete can be "message" for delete a message, for example in IG, in WhatsApp can be "message" for delete a message.

		@data is the data needed to delete the message, for example in IG, can be the id of the message, in WhatsApp can be the id of the message and the phone number id
	*/
	Delete(typeDelete string, data ...map[string]any) response.Responser

	/*
		@typeUpdate can be "message" for update a message, for example in IG, in WhatsApp can be "message" for update a message.

		@data is the data needed to update the message, for example in IG, can be the id of the message, in WhatsApp can be the id of the message and the phone number id

	*/
	Update(typeUpdate string, data ...map[string]any) response.Responser
	/*
		@recipient_id is the id of the recipient of the presence
		@action can be "typing_on", "typing_off", or "recording_on", "recording_off" depending on the client, for example in IG can be "typing_on" or "typing_off", in WhatsApp can be "typing_on", "typing_off", "recording_on", "recording_off"
	*/
	SendPresence(recipient_id, presence string) response.Responser
	MultipartRequest(method string, data proto.Message, ePoint string) (*http.Request, error)
}
