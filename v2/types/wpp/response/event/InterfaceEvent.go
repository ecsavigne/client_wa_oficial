package event

import "github.com/ecsavigne/client_wa_oficial/v2/types/wpp"

type EventInterface interface {
	GetType() wpp.EventType
	String() string

	/*
	 * GetContactPhone() string
	 * GetPhoneID() string
	 * GetContactName() string
	 * GetTimestamp() int64
	 * GetMessageID() string
	 * GetResponseMessageID() string
	 * GetButtonText() string
	 * GetContext()
	 */
}
