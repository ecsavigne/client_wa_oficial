package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type SmbAppStateSyncEvent struct {
	*Components
}

func (*SmbAppStateSyncEvent) GetType() wpp.EventType {
	return wpp.EventTypSmbAppStateSync
}

func (m *SmbAppStateSyncEvent) String() string { return response.Val(m) }

type SmbMessageEchoesEvent struct {
	*Components
}

func (*SmbMessageEchoesEvent) GetType() wpp.EventType {
	return wpp.EventTypeSmbMessageEchoes
}

func (m *SmbMessageEchoesEvent) String() string { return response.Val(m) }
