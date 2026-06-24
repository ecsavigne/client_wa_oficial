package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type PartnerSolutionsEvent struct {
	*Components
}

func (*PartnerSolutionsEvent) GetType() wpp.EventType {
	return wpp.EventTypePartnerSolutions
}

func (m *PartnerSolutionsEvent) String() string { return response.Val(m) }
