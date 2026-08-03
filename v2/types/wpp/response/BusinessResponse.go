package response

type MarketingMessagesOnboardingStatus struct {
	Status string `json:"status,omitempty"`
}

type CreatedBy struct {
	Id       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Business *Business `json:"business,omitempty"`
}

type Business struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	Id           string `json:"id,omitempty"`
	// Whether offline analytics is blocked for this business.
	BlockOfflineAnalytics                     bool   `json:"block_offline_analytics,omitempty"`
	CollaborativeAdsManagedPartnerEligibility string `json:"collaborative_ads_managed_partner_eligibility,omitempty"`
	// The creator of this business.
	CreatedBy   *CreatedBy `json:"created_by,omitempty"`
	CreatedTime string     `json:"created_time,omitempty"`
	// The last update time for this business's extended credits.
	ExtendedUpdatedTime string `json:"extended_updated_time,omitempty"`
	IsHidden            bool   `json:"is_hidden,omitempty"`
	// The URI of the business's profile page.
	Link string `json:"link,omitempty"`
	Name string `json:"name,omitempty"`
	// The payment account ID of this business.
	PaymentAccountId string `json:"payment_account_id,omitempty"`
	// The primary Facebook Page associated with this business.
	PrimaryPage string `json:"primary_page,omitempty"`
	// The profile picture URI of this business.
	ProfilePictureUri string `json:"profile_picture_uri,omitempty"`
	// The time zone ID of this business.
	TimezoneId uint32 `json:"timezone_id,omitempty"`
	// The two-factor authentication method used by this business.
	TwoFactorType string `json:"two_factor_type,omitempty"`
	// The name of the user who last updated this business.
	UpdatedBy string `json:"updated_by,omitempty"`
	// The time when this business was last updated.
	UpdatedTime string `json:"updated_time,omitempty"`
	/* The verification status of this business. verification_status. enum {expired, failed, ineligible, not_verified, pending, pending_need_more_info, pending_submission, rejected, revoked, verified}*/
	VerificationStatus string `json:"verification_status,omitempty"`
	// The industry sector associated with this business.
	Vertical string `json:"vertical,omitempty"`
	//Business industry identifier.
	VerticalId uint32 `json:"vertical_id,omitempty"`
	Success    bool   `json:"success,omitempty"`
	Message    string `json:"message,omitempty"`
	/*Maximum number of unique WhatsApp users your Business Manager can message outside the 24-hour customer service window. This limit is shared across all WhatsApp phone numbers associated with your Business Manager.. enum {TIER_100K, TIER_10K, TIER_250, TIER_2K, TIER_UNLIMITED, UNTIERED} */
	WhatsappBusinessManagerMessagingLimit string                             `json:"whatsapp_business_manager_messaging_limit,omitempty"`
	MarketingMessagesOnboardingStatus     *MarketingMessagesOnboardingStatus `json:"marketing_messages_onboarding_status,omitempty"`
}

func NewBusiness(config Responser) *Business {
	if v, ok := config.(*Business); ok {
		v.ResponseType = ResponseBusiness
		v.KernelResponser.parent = v
		return v
	}
	panic("type ResponserRequest is not *Business")
}
