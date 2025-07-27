/*
Permission in WebHook for received notifications about changes in templates components (message_template_components_update)
*/
package response

import "github.com/ecsavigne/client_wa_oficial/v2/types"

type WebHookTemplateResponse struct {
	KernelResponser
	ResponseType string `json:"response_type,omitempty"`
	// Template id
	MessageTemplateId uint64 `json:"message_template_id,omitempty"`
	// Template name
	MessageTemplateName string `json:"message_template_name,omitempty"`
	// Template language
	MessageTemplateLanguage string `json:"message_template_language,omitempty"`
	// New header of template before update, stay empty if user don't input header
	MessageTemplateTitle string `json:"message_template_title,omitempty"`
	// New body of the template before update, stay empty if user don't input body
	MessageTemplateElement string `json:"message_template_element,omitempty"`
	// New footer of the template before update, stay empty if user don't input footer
	MessageTemplateFooter string `json:"message_template_footer,omitempty"`
	// New buttons of the template before update, stay empty if user don't input
	MessageTemplateButtons []types.Button `json:"message_template_buttons,omitempty"`
}

func NewWebHookTemplateResponse(config Responser) *WebHookTemplateResponse {
	if v, ok := config.(*WebHookTemplateResponse); ok {
		v.ResponseType = ResponseWebHookTemplate
		v.KernelResponser.parent = v
		return v
	}
	panic("type Responser is not *WebHookTemplateResponse")
}
