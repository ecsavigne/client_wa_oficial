package response

type Business struct {
	KernelResponser
	ResponseType                              string `json:"response_type,omitempty"`
	Id                                        string `json:"id,omitempty"`
	BlockOfflineAnalytics                     bool   `json:"block_offline_analytics,omitempty"`
	CollaborativeAdsManagedPartnerEligibility string `json:"collaborative_ads_managed_partner_eligibility,omitempty"`
	CreatedBy                                 string `json:"created_by,omitempty"`
	CreatedTime                               string `json:"created_time,omitempty"`
	ExtendedUpdatedTime                       string `json:"extended_updated_time,omitempty"`
	IsHidden                                  bool   `json:"is_hidden,omitempty"`
	Link                                      string `json:"link,omitempty"`
	Name                                      string `json:"name,omitempty"`
	PaymentAccountId                          string `json:"payment_account_id,omitempty"`
	PrimaryPage                               string `json:"primary_page,omitempty"`
	ProfilePictureUri                         string `json:"profile_picture_uri,omitempty"`
	TimezoneId                                uint32 `json:"timezone_id,omitempty"`
	TwoFactorType                             string `json:"two_factor_type,omitempty"`
	UpdatedTime                               string `json:"updated_time,omitempty"`
	VerificationStatus                        string `json:"verification_status,omitempty"`
	Vertical                                  string `json:"vertical,omitempty"`
	VerticalId                                uint32 `json:"vertical_id,omitempty"`
	Success                                   bool   `json:"success,omitempty"`
	Message                                   string `json:"message,omitempty"`
}

func NewBusiness(config ResponserRequest) *Business {
	if v, ok := config.(*Business); ok {
		v.ResponseType = ResponseBusiness
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Business")
}
