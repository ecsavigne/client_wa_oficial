package event

import (
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response"
)

type TYPE_ACCOUNT_UPDATE_EVENT string

var (
	ACCOUNT_DELETED, ACCOUNT_RESTRICTION, ACCOUNT_VIOLATION, AD_ACCOUNT_LINKED                 TYPE_ACCOUNT_UPDATE_EVENT
	AUTH_INTL_PRICE_ELIGIBILITY_UPDATE, BUSINESS_PRIMARY_LOCATION_COUNTRY_UPDATE               TYPE_ACCOUNT_UPDATE_EVENT
	DISABLED_UPDATE, MM_LITE_TERMS_SIGNED, PARTNER_ADDED                                       TYPE_ACCOUNT_UPDATE_EVENT
	PARTNER_APP_INSTALLED, PARTNER_APP_UNINSTALLED, PARTNER_CLIENT_CERTIFICATION_STATUS_UPDATE TYPE_ACCOUNT_UPDATE_EVENT
	PARTNER_REMOVED, VOLUME_BASED_PRICING_TIER_UPDATE, ACCOUNT_OFFBOARDED                      TYPE_ACCOUNT_UPDATE_EVENT
	ACCOUNT_RECONNECTED                                                                        TYPE_ACCOUNT_UPDATE_EVENT
)

type AccountUpdateEvent struct {
	*Components
}

func (*AccountUpdateEvent) GetType() wpp.EventType { return wpp.EventTypeAccountUpdate }

func (m *AccountUpdateEvent) String() string { return response.Val(m) }

func (self *AccountUpdateEvent) GetEvent() string {
	// return self.GetEntry()[0].GetChange()[0].GetValue().Event
	return self.Components.GetEvent()
}

func (self *AccountUpdateEvent) GetWabaInfo() (vWabaInfo *WabaInfo) {
	defer func() {
		if r := recover(); r != nil {
			vWabaInfo = nil
		}
	}()

	return self.GetEntry()[0].Changes[0].Value.WabaInfo
}

type AccountAlertsEvent struct {
	*Components
}

func (*AccountAlertsEvent) GetType() wpp.EventType { return wpp.EventTypeAccountAlerts }

func (m *AccountAlertsEvent) String() string { return response.Val(m) }

type AccountReviewUpdateEvent struct {
	*Components
}

func (*AccountReviewUpdateEvent) GetType() wpp.EventType { return wpp.EventTypeAccountAlerts }

func (m *AccountReviewUpdateEvent) String() string { return response.Val(m) }
