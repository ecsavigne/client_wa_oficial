package event

import "github.com/ecsavigne/client_wa_oficial/v2/types"

type EventInterface interface {
	GetType() types.EventType
	String() string
	/*
	 * GetContactPhone() string
	 * GetPhoneID() string
	 * GetContactName() string
	 * GetTimestamp() int64
	 * GetMessageID() string
	 * GetResponseMessageID() string
	 * GetButtonText() string
	 * GetContextMessage()
	 */
}
