package event

type EventInterface interface {
	GetType() string
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
