package event

import (
	"strconv"

	"github.com/ecsavigne/client_wa_oficial/v2/types"
	"github.com/ecsavigne/client_wa_oficial/v2/types/internal"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

type MessageButtonEvent struct {
	*Components
}

func (*MessageButtonEvent) GetType() types.EventType { return types.EventTypeMessageButton }
func (m *MessageButtonEvent) String() string         { return response.Val(m) }

func (bEvt *MessageButtonEvent) GetContactPhone() string {
	return internal.FirstNotEmpty(bEvt.Entry[0].Changes[0].Value.Contacts[0].WaId,
		bEvt.Entry[0].Changes[0].Value.Messages[0].From)
}

func (bEvt *MessageButtonEvent) GetPhoneID() string {
	return bEvt.Entry[0].Changes[0].Value.Metadata.PhoneNumberID
}

func (bEvt *MessageButtonEvent) GetContactName() string {
	return bEvt.Entry[0].Changes[0].Value.Contacts[0].Name
}

func (bEvt *MessageButtonEvent) GetTimestamp() uint64 {
	timeStamp, err := strconv.ParseUint(bEvt.Entry[0].Changes[0].Value.Messages[0].Timestamp, 10, 64)
	if err != nil {
		return 0
	}

	return timeStamp
}

func (bEvt *MessageButtonEvent) GetMessageID() string {
	return bEvt.Entry[0].Changes[0].Value.Messages[0].ID
}

func (bEvt *MessageButtonEvent) GetResponseMessageID() string {
	ctx := bEvt.GetContext()
	if ctx != nil {
		return ctx.ID
	}

	return ""
}

func (bEvt *MessageButtonEvent) GetButtonText() string {
	return internal.FirstNotEmpty(bEvt.Entry[0].Changes[0].Value.Messages[0].Button.Text,
		bEvt.Entry[0].Changes[0].Value.Messages[0].Button.Payload,
	)
}

func (bEvt *MessageButtonEvent) GetContext() *Context {
	return bEvt.Entry[0].Changes[0].Value.Messages[0].Context
}
