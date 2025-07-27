package types

type Url_Button struct {
	// Admite one var max, ej: "https://www.luckyshrub.com/shop/", ej2: https://www.luckyshrub.com/shop?promo={{1}}. max size = 2.000 chars
	Url string `json:"url" validate:"required"`
	// Required if <URL> contains a variable
	Example []string `json:"example,omitempty"`
}

type Phone_Number_Button struct {
	// less than 20 chars
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type Copy_Code_Button struct {
	// 15 chars maximum. String that will be copied in Clipboard
	Example string `json:"example,omitempty"`
}

type Flow_Button struct {
	FlowId uint64 `json:"flow_id" validate:"required"`
	// name of proccess only allowed on Cloud Api
	FlowName string `json:"flow_name,omitempty" validate:"required"`
	//  string in format json ex: "{\"version\": \"3.1\", \"screens\": [...]}"
	FlowJson string `json:"flow_json" validate:"required"`
	// Allowed values: "navigate", "data_exchange". Use "navigate" for define first screen with part of template message.
	FlowAction string `json:"flow_action,omitempty"`
	// Optional only if <flow_action> is "navigate". The id of the first screen of process
	NavigateScreen string `json:"navigate_screen,omitempty"`
	//  Allowed values: "DOCUMENT", "PROMOTION", "REVIEW". Default "PROMOTION"
	Icon string `json:"icon,omitempty"`
}

type Quick_Reply_Button struct {
}

type Button struct {
	Type string `json:"type" validate:"required"`
	// 25 characters maximum
	Text string `json:"text" validate:"required"`
	*Url_Button
	*Phone_Number_Button
	*Copy_Code_Button
	*Flow_Button
	*Quick_Reply_Button
}
