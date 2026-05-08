package examplemessage

import (
	"fmt"

	clientpack "github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/client/ig"
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	type_ig "github.com/ecsavigne/client_wa_oficial/v2/types/ig"
	igpbv1 "github.com/ecsavigne/client_wa_oficial/v2/types/ig/gen/igpb/v1"
)

var NewClientIG = func() clientpack.Client {
	client, e := ig.NewClientIG(
		// ig.WithEnvFilePath("../../config_env"),
		ig.WithEnvFilePath("../config_env"),
		ig.WithToken("IGAAQ7JedvaElBZAGI3clhETVZAFa3hjeXRVd200cnZAucHplbkFrTVI3OUJwMTVfeVN2MHhzWjBPY2hvS2gwODROUTIybkhNWVlFcmdjdEFySmZAMU2VNcU1nMkVNeFlIZADJLcVBJVTBrT0o0NlFDdGtRZAGpadWFBRVU4ektSMWs5cwZDZD"),
		ig.WithUserID("17841448182209630"),
	)

	if e != nil {
		fmt.Printf("Error creating client IG: %s\n", e)
		return nil
	}

	return client
}

func SendMessageText(cl clientpack.Client, scope_id, msg string) response.Responser {
	msgIGText := new(igpbv1.InstagramTextMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)
	// recipient.SetId("2014023525994448")
	text := new(igpbv1.TextMessage)
	text.SetText(msg)
	msgIGText.SetRecipient(recipient)
	msgIGText.SetMessage(text)

	return cl.SendMessage(msgIGText)
}

// @sender_action can be "react" or "unreact"
// reaction can be "👍", "😂", "😮", "😢", "😡", or name e.g: "love"
func SendMessageReaction(cl clientpack.Client, scope_id, message_id, reaction string, sender_action ...string) response.Responser {
	var sender_act string
	if len(sender_action) > 0 {
		sender_act = sender_action[0]
	} else {
		sender_act = "react"
	}

	msgIGReaction := new(igpbv1.InstagramReactionMessage)
	msgIGReaction.SetSenderAction(sender_act)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)
	msgIGReaction.SetRecipient(recipient)

	payload := new(igpbv1.ReactionPayload)
	payload.SetReaction(reaction)
	payload.SetMessageId(message_id)
	msgIGReaction.SetPayload(payload)

	return cl.SendMessage(msgIGReaction)
}

// send media message with url or attachment_id, if url is not empty, send with url, else send with attachment_id
/*
	type image file: jpeg, png, gif. max-size:8MB
*/
func SendMessageImage(cl clientpack.Client, scope_id, url, attacment_id string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	attachment := new(igpbv1.Attachment)
	payload := new(igpbv1.AttachmentPayload)
	if attacment_id != "" {
		payload.SetAttachmentId(attacment_id)
	} else {
		payload.SetUrl(url)
	}
	attachment.SetPayload(payload)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_TEMPLATE.String())
	msg.SetAttachment(attachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

func SendMessageImageCollage(cl clientpack.Client, scope_id string, urls, attachment_ids []string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	// attachment := new(igpbv1.Attachment)
	arrAttachment := make([]*igpbv1.Attachment, len(attachment_ids)) //new
	if len(attachment_ids) > 0 {
		for i, attId := range attachment_ids {
			arrAttachment[i] = new(igpbv1.Attachment)
			payload := new(igpbv1.AttachmentPayload)
			payload.SetAttachmentId(attId)
			arrAttachment[i].SetPayload(payload)
			arrAttachment[i].SetType(type_ig.IG_ATTACHMENT_TYPE_IMAGE.String())
		}
	} else {
		for i, url := range urls {
			arrAttachment[i] = new(igpbv1.Attachment)
			payload := new(igpbv1.AttachmentPayload)
			payload.SetUrl(url)
			arrAttachment[i].SetPayload(payload)
			arrAttachment[i].SetType(type_ig.IG_ATTACHMENT_TYPE_IMAGE.String())
		}
		// payload.SetUrl(url)
	}

	msg.SetAttachments(arrAttachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

// send media message with url or attachment_id, if url is not empty, send with url, else send with attachment_id
/*
type audio file: aac, m4a, wav, mp4. max-size: 25MB
*/
func SendMessageAudio(cl clientpack.Client, scope_id, url, attacment_id string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	attachment := new(igpbv1.Attachment)
	payload := new(igpbv1.AttachmentPayload)
	if attacment_id != "" {
		payload.SetAttachmentId(attacment_id)
	} else {
		payload.SetUrl(url)
	}
	attachment.SetPayload(payload)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_AUDIO.String())
	msg.SetAttachment(attachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

// send media message with url or attachment_id, if url is not empty, send with url, else send with attachment_id
/*
	type video file: mp4, ogg, avi, mov, webm. max-size:25MB
*/
func SendMessageVideo(cl clientpack.Client, scope_id, url, attacment_id string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	attachment := new(igpbv1.Attachment)
	payload := new(igpbv1.AttachmentPayload)
	if attacment_id != "" {
		payload.SetAttachmentId(attacment_id)
	} else {
		payload.SetUrl(url)
	}
	attachment.SetPayload(payload)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_VIDEO.String())
	msg.SetAttachment(attachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

/*
type pdf file: pdf. max-size:25MB
*/
func SendMessagePDFFile(cl clientpack.Client, scope_id, url, attacment_id string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	attachment := new(igpbv1.Attachment)
	payload := new(igpbv1.AttachmentPayload)
	if attacment_id != "" {
		payload.SetAttachmentId(attacment_id)
	} else {
		payload.SetUrl(url)
	}
	attachment.SetPayload(payload)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_FILE.String())
	msg.SetAttachment(attachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

func SendMessageStickerLike(cl clientpack.Client, scope_id string) response.Responser {
	msgIGMedia := new(igpbv1.InstagramMediaMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msg := new(igpbv1.MediaMessage)
	attachment := new(igpbv1.Attachment)
	payload := new(igpbv1.AttachmentPayload)

	attachment.SetPayload(payload)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_LIKE_HEART.String())
	msg.SetAttachment(attachment)
	msgIGMedia.SetRecipient(recipient)
	msgIGMedia.SetMessage(msg)

	return cl.SendMessage(msgIGMedia)
}

func SendMessageQuickReplies(cl clientpack.Client, scope_id string) response.Responser {
	igQuickRepliesMessage := new(igpbv1.InstagramQuickRepliesMessage)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)
	igQuickRepliesMessage.SetRecipient(recipient)

	qrMessage := new(igpbv1.QuickRepliesMessage)
	arrQReply := new([]*igpbv1.QuickReply)

	// section 1
	qr := new(igpbv1.QuickReply)
	qr.SetContentType("text")
	qr.SetTitle("Option 1")
	qr.SetPayload("click button option_1")
	*arrQReply = append(*arrQReply, qr)

	// section 2
	qr = new(igpbv1.QuickReply)
	qr.SetContentType("user_email")
	qr.SetTitle("email")
	qr.SetPayload("click button email")
	*arrQReply = append(*arrQReply, qr)

	// section 3
	qr = new(igpbv1.QuickReply)
	qr.SetContentType("user_phone_number")
	qr.SetTitle("phone")
	qr.SetPayload("click button phone")
	*arrQReply = append(*arrQReply, qr)

	qrMessage.SetQuickReplies(*arrQReply)
	qrMessage.SetText("Choose an option:")
	igQuickRepliesMessage.SetMessage(qrMessage)

	return cl.SendMessage(igQuickRepliesMessage)
}

// Solo un menu persistente por cuenta
func SendMessagePersistentMenu(cl clientpack.Client) response.Responser {
	igPersistentMenuMessage := new(igpbv1.InstagramPersistentMenuMessage)
	arrPersistentsItemMenu := new([]*igpbv1.PersistentMenuItem)

	pItemMenu := new(igpbv1.PersistentMenuItem)
	pItemMenu.SetLocale("default")
	call_to_actions := new([]*igpbv1.CallToAction)

	// call_to_actions 1
	call_to_action := new(igpbv1.CallToAction)
	call_to_action.SetType("postback")
	call_to_action.SetTitle("Habla con un agente")
	call_to_action.SetPayload("AYUDA_CONSULTA")
	*call_to_actions = append(*call_to_actions, call_to_action)

	// call_to_actions 2
	call_to_action = new(igpbv1.CallToAction)
	call_to_action.SetType("postback")
	call_to_action.SetTitle("Sugerencias de atuendos")
	call_to_action.SetPayload("CURACION")
	*call_to_actions = append(*call_to_actions, call_to_action)

	// call_to_actions 3
	call_to_action = new(igpbv1.CallToAction)
	call_to_action.SetType("web_url")
	call_to_action.SetTitle("SocialHub")
	call_to_action.SetUrl("https://oficial.crmsocialhub.com.br/")
	*call_to_actions = append(*call_to_actions, call_to_action)

	pItemMenu.SetCallToActions(*call_to_actions)
	*arrPersistentsItemMenu = append(*arrPersistentsItemMenu, pItemMenu)
	igPersistentMenuMessage.SetPersistentMenu(*arrPersistentsItemMenu)

	return cl.SendMessage(igPersistentMenuMessage)
}

func SendMessageIceBreakers(cl clientpack.Client) response.Responser {
	igIceBreakersMessage := new(igpbv1.InstagramIceBreakersMessage)
	itemsIceBreakers := new([]*igpbv1.IceBreakersItem)
	iceBreakersItem := new(igpbv1.IceBreakersItem)
	iceBreakersItem.SetLocale("default")

	igIceBreakersMessage.SetPlatform("instagram")

	call_to_actions := new([]*igpbv1.CallToAction)

	// call_to_actions 1
	call_to_action := new(igpbv1.CallToAction)
	call_to_action.SetQuestion("Pregunta de rompehielos 1")
	call_to_action.SetPayload("PAYLOAD_ROMPEHIELOS_1")
	*call_to_actions = append(*call_to_actions, call_to_action)

	// call_to_actions 2
	call_to_action = new(igpbv1.CallToAction)
	call_to_action.SetQuestion("Pregunta de rompehielos 2")
	call_to_action.SetPayload("PAYLOAD_ROMPEHIELOS_2")
	*call_to_actions = append(*call_to_actions, call_to_action)

	// call_to_actions 3
	call_to_action = new(igpbv1.CallToAction)
	call_to_action.SetQuestion("Pregunta de rompehielos 3")
	call_to_action.SetPayload("PAYLOAD_ROMPEHIELOS_3")
	*call_to_actions = append(*call_to_actions, call_to_action)

	iceBreakersItem.SetCallToActions(*call_to_actions)
	*itemsIceBreakers = append(*itemsIceBreakers, iceBreakersItem)
	igIceBreakersMessage.SetIceBreakers(*itemsIceBreakers)

	return cl.SendMessage(igIceBreakersMessage)
}

// Só aprarece se estão mostrado a um anúncio do Facebook ou Instagram
func CreateInstagramWelcomeMessageFlowsADS(cl clientpack.Client) response.Responser {
	igWelcomeMessageFlows := new(igpbv1.InstagramWelcomeMessageFlows)
	igWelcomeMessageFlows.SetEligiblePlatforms([]string{"instagram"})
	igWelcomeMessageFlows.SetName("welcome_message_flows_example")

	welcome_message_flows := make([]*igpbv1.WelcomeMessageFlowItem, 0)
	welcome_message_flow := new(igpbv1.WelcomeMessageFlowItem)

	message := new(igpbv1.WelcomeMessageFlowItem_Message)
	message.SetText("¡Bienvenido a nuestra tienda! ¿En qué podemos ayudarte hoy?")

	quick_replies := new([]*igpbv1.QuickReply)
	// quick_reply 1
	quick_reply := new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 1")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_1")
	*quick_replies = append(*quick_replies, quick_reply)

	// quick_reply 2
	quick_reply = new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 2")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_2")
	*quick_replies = append(*quick_replies, quick_reply)

	// quick_reply 3
	quick_reply = new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 3")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_3")
	*quick_replies = append(*quick_replies, quick_reply)

	message.SetQuickReplies(*quick_replies)

	welcome_message_flow.SetMessage(message)
	welcome_message_flows = append(welcome_message_flows, welcome_message_flow)

	igWelcomeMessageFlows.SetWelcomeMessageFlow(welcome_message_flows)

	return cl.SendMessage(igWelcomeMessageFlows)
}

func SendPresence(cl clientpack.Client, scope_id string, action string) response.Responser {
	return cl.SendPresence(scope_id, action)
}

func MarkRead(cl clientpack.Client, scope_id string) response.Responser {
	return cl.MarkRead(scope_id)
}

// gets
/*
id:						The app user's app-scoped ID
user_id:				The Instagram professional acount ID, <IG_ID>, for your app user. This ID is value of the id field received in webhook notifications for this account.
username:				The app user's Instagram username.
name:					The app user's name
account_type:			The app user's account type. Can be Business or Media_Creator.
profile_picture_url:	The URL for the app user's profile picture.
followers_count:		The number of followers of the app user's Instagram professional account
follows_count:			The number of Instagram accounts the app user's Instagram professional account follows
media_count:			The number of Media object on the User
*/
func GetAccountInfoBusiness(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_INFO_ACCOUNT_BUSINESS)
}

func GetPersistentMenu(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_INFO_PERSISTENT_MENU)
}

func GetIceBreakers(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_INFO_ICE_BREAKERS)
}

func GetInstagramLink(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_INFO_LINK)
}

func GetWelcomeMessageFlowsADS(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_INFO_WELCOME_MESSAGE_FLOWS)
}

// deletes
func DeletePersistentMenu(cl clientpack.Client) response.Responser {
	return cl.Delete(type_ig.IG_DELETE_PERSISTENT_MENU)
}
func DeleteIceBreakers(cl clientpack.Client) response.Responser {
	return cl.Delete(type_ig.IG_DELETE_ICE_BREAKERS)
}

func DeleteWelcomeMessageFlowsADS(cl clientpack.Client, data map[string]any) response.Responser {
	return cl.Delete(type_ig.IG_DELETE_WELCOME_MESSAGE_FLOWS, data)
}

// updates
// Só é possível atualizar o fluxo de mensagens de boas-vindas, e para isso é necessário o id do fluxo que se deseja atualizar, e o fluxo atualizado
func UpdateWelcomeMessageFlowsADS(cl clientpack.Client) response.Responser {
	mapData := map[string]any{"flow_id": "1361283519166009"}

	igWelcomeMessageFlows := new(igpbv1.InstagramWelcomeMessageFlows)
	// igWelcomeMessageFlows.SetEligiblePlatforms([]string{"instagram"})
	igWelcomeMessageFlows.SetName("welcome_message_flows_example")

	welcome_message_flows := make([]*igpbv1.WelcomeMessageFlowItem, 0)
	welcome_message_flow := new(igpbv1.WelcomeMessageFlowItem)

	message := new(igpbv1.WelcomeMessageFlowItem_Message)
	message.SetText("¡Bienvenido a nuestra tienda! ¿En qué podemos ayudarte hoy?")

	quick_replies := new([]*igpbv1.QuickReply)
	// quick_reply 1
	quick_reply := new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 1")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_1")
	*quick_replies = append(*quick_replies, quick_reply)

	// quick_reply 2
	quick_reply = new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 2")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_2")
	*quick_replies = append(*quick_replies, quick_reply)

	// quick_reply 3
	quick_reply = new(igpbv1.QuickReply)
	quick_reply.SetTitle("Pregunta de rompehielos 3")
	quick_reply.SetContentType("text")
	quick_reply.SetPayload("PAYLOAD_ROMPEHIELOS_3")
	*quick_replies = append(*quick_replies, quick_reply)

	message.SetQuickReplies(*quick_replies)

	welcome_message_flow.SetMessage(message)
	welcome_message_flows = append(welcome_message_flows, welcome_message_flow)

	igWelcomeMessageFlows.SetWelcomeMessageFlow(welcome_message_flows)

	mapData["msg"] = igWelcomeMessageFlows

	return cl.Update(type_ig.IG_UPDATE_WELCOME_MESSAGE_FLOWS, mapData)
}

func SendPublishImg(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"images_url": []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRN2z0ERwXQUqH29urPuzWueLXKhJAY6SMyAA&ss"},
	}

	return cl.Create(type_ig.IG_CREATE_POST, dataParam)
}

func SendPublishVideo(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"videos_url": []string{"https://res.cloudinary.com/dczar4xfh/video/upload/v1776219555/samples/dance-2.mp4"},
	}

	return cl.Create(type_ig.IG_CREATE_POST, dataParam)
}

func SendPublishMultiFile(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"images_url": []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRN2z0ERwXQUqH29urPuzWueLXKhJAY6SMyAA&ss"},
		"videos_url": []string{"https://res.cloudinary.com/dczar4xfh/video/upload/v1776219555/samples/dance-2.mp4"},
	}

	return cl.Create(type_ig.IG_CREATE_POST, dataParam)
}

func SendHistoryImg(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"images_url": []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRN2z0ERwXQUqH29urPuzWueLXKhJAY6SMyAA&ss"},
	}

	return cl.Create(type_ig.IG_CREATE_STORY, dataParam)
}

func SendHistoryVideo(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"videos_url": []string{"https://res.cloudinary.com/dczar4xfh/video/upload/v1776219555/samples/dance-2.mp4"},
	}

	return cl.Create(type_ig.IG_CREATE_STORY, dataParam)
}

func SendHistoryMultiFile(cl clientpack.Client) response.Responser {
	dataParam := map[string]any{
		"images_url": []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRN2z0ERwXQUqH29urPuzWueLXKhJAY6SMyAA&ss"},
		"videos_url": []string{"https://res.cloudinary.com/dczar4xfh/video/upload/v1776219555/samples/dance-2.mp4"},
	}

	return cl.Create(type_ig.IG_CREATE_STORY, dataParam)
}

func SendComment(cl clientpack.Client, media_id, comment string) response.Responser {
	dataParam := map[string]any{
		"ig_media_id": media_id,
		"message":     comment,
	}

	return cl.Create(type_ig.IG_CREATE_COMMENT, dataParam)
}

func SendReplyComment(cl clientpack.Client, comment_id, reply string) response.Responser {
	dataParam := map[string]any{
		"ig_comment_id": comment_id,
		"message":       reply,
	}

	return cl.Create(type_ig.IG_CREATE_REPLY_COMMENT, dataParam)
}

func HideComment(cl clientpack.Client, comment_id string, hide bool) response.Responser {
	dataParam := map[string]any{
		"ig_comment_id": comment_id,
		"hide":          hide,
	}

	return cl.Create(type_ig.IG_CREATE_HIDE_COMMENT, dataParam)
}

func EnableComment(cl clientpack.Client, ig_media_id string, comment_enabled bool) response.Responser {
	dataParam := map[string]any{
		"ig_media_id":     ig_media_id,
		"comment_enabled": comment_enabled,
	}

	return cl.Create(type_ig.IG_CREATE_ENABLE_COMMENT, dataParam)
}

func DeleteComment(cl clientpack.Client, ig_comment_id string) response.Responser {
	dataParam := map[string]any{
		"ig_comment_id": ig_comment_id,
	}

	return cl.Delete(type_ig.IG_DELETE_COMMENT, dataParam)
}

func GetComment(cl clientpack.Client, ig_media_id string) response.Responser {
	dataParam := map[string]any{
		"ig_media_id": ig_media_id,
	}

	return cl.Get(type_ig.IG_GET_COMMENT, dataParam)
}

func GetRepliesComments(cl clientpack.Client, ig_comment_id string) response.Responser {
	dataParam := map[string]any{
		"ig_comment_id": ig_comment_id,
	}

	return cl.Get(type_ig.IG_GET_REPLIES_COMMENTS, dataParam)
}

func SendButtonTemplateMessage(cl clientpack.Client, scope_id string) response.Responser {
	msgIGButtonTemplate := new(igpbv1.InstagramTemplateButtonTemplate)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msgIGButtonTemplate.SetRecipient(recipient)

	msg := new(igpbv1.ButtonMessage)
	// igpbv1.templateAttachment_Button
	attachment := new(igpbv1.TemplateAttachment)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_TEMPLATE.String())

	payload := new(igpbv1.TemplatePayload)
	payload.SetTemplateType(type_ig.IG_TEMPLATE_BUTTON.String())
	payload.SetText("Este é um exemplo de mensagem com template de botão.")

	buttons := new([]*igpbv1.TemplateButton)

	button1 := new(igpbv1.TemplateButton)
	button1.SetType(type_ig.IG_TYPE_BUTTON_POSTBACK.String())
	button1.SetTitle("Botão para postback")
	button1.SetPayload("data send by postback button for webhook")
	*buttons = append(*buttons, button1)

	button2 := new(igpbv1.TemplateButton)
	button2.SetType(type_ig.IG_TYPE_BUTTON_WEB_URL.String())
	button2.SetTitle("Botão para web_url")
	button2.SetUrl("https://app.socialhub.com/")
	*buttons = append(*buttons, button2)

	payload.SetButtons(*buttons)
	attachment.SetPayload(payload)
	msg.SetAttachment(attachment)
	msgIGButtonTemplate.SetMessage(msg)

	return cl.SendMessage(msgIGButtonTemplate)
}

func SendGenericTemplateMessage(cl clientpack.Client, scope_id string) response.Responser {
	msgIGButtonTemplate := new(igpbv1.InstagramTemplateButtonTemplate)
	recipient := new(igpbv1.Recipient)
	recipient.SetId(scope_id)

	msgIGButtonTemplate.SetRecipient(recipient)

	msg := new(igpbv1.ButtonMessage)
	// igpbv1.templateAttachment_Button
	attachment := new(igpbv1.TemplateAttachment)
	attachment.SetType(type_ig.IG_ATTACHMENT_TYPE_TEMPLATE.String())

	payload := new(igpbv1.TemplatePayload)
	payload.SetTemplateType(type_ig.IG_TEMPLATE_GENERIC.String())

	elements := new([]*igpbv1.Element)

	element1 := new(igpbv1.Element)
	// element1.SetType(type_ig.IG_TYPE_BUTTON_POSTBACK.String())
	element1.SetTitle("Botão para postback")
	element1.SetImageUrl("https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExeTNvcThmcG85OTRlMTF5MTZueWk4ZGMxcDRkMHcxMXA1Z2oxNzZlaSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/O3wIKp4DlubUl1Bezt/giphy.gif")
	element1.SetSubtitle("Este é um exemplo de mensagem com template genérico.")

	defaultAction := new(igpbv1.DefaultAction)
	defaultAction.SetType(type_ig.IG_TYPE_BUTTON_WEB_URL.String())
	defaultAction.SetUrl("https://app.socialhub.com/")
	element1.SetDefaultAction(defaultAction)

	buttons := new([]*igpbv1.TemplateButton)

	button1 := new(igpbv1.TemplateButton)
	button1.SetType(type_ig.IG_TYPE_BUTTON_POSTBACK.String())
	button1.SetTitle("Botão para postback")
	button1.SetPayload("data send by postback button for webhook")
	*buttons = append(*buttons, button1)

	button2 := new(igpbv1.TemplateButton)
	button2.SetType(type_ig.IG_TYPE_BUTTON_WEB_URL.String())
	button2.SetTitle("Botão para web_url")
	button2.SetUrl("https://app.socialhub.com/")
	*buttons = append(*buttons, button2)

	element1.SetButtons(*buttons)
	*elements = append(*elements, element1)

	payload.SetElements(*elements)
	attachment.SetPayload(payload)
	msg.SetAttachment(attachment)
	msgIGButtonTemplate.SetMessage(msg)

	return cl.SendMessage(msgIGButtonTemplate)
}

func SendPrivateReplyMessage(cl clientpack.Client, comment_id, msg string) response.Responser {
	msgIGPrivateReply := new(igpbv1.InstagramPrivateReplyMessage)

	recipient := new(igpbv1.CommentRecipient)
	recipient.SetCommentId(comment_id)
	msgIGPrivateReply.SetRecipient(recipient)

	msgPR := new(igpbv1.PrivateReplyMessage)
	msgPR.SetText(msg)
	msgIGPrivateReply.SetMessage(msgPR)

	return cl.SendMessage(msgIGPrivateReply)
}

// Meta tem que revisar e aprovar o app pra ig
func SendHumanAgentMessage(cl clientpack.Client, recipient_id, msg string) response.Responser {
	msgIGHumanAgent := new(igpbv1.InstagramHumanAgentMessage)
	msgIGHumanAgent.SetTag("HUMAN_AGENT")

	recipient := new(igpbv1.Recipient)
	recipient.SetId(recipient_id)
	msgIGHumanAgent.SetRecipient(recipient)

	msgHA := new(igpbv1.HumanAgentMessage)
	msgHA.SetText(msg)
	msgIGHumanAgent.SetMessage(msgHA)

	return cl.SendMessage(msgIGHumanAgent)
}

func SubscribeWebHook(cl clientpack.Client) response.Responser {
	return cl.SubscribeWebHook()
}

func UnsubscribeWebHook(cl clientpack.Client) response.Responser {
	return cl.UnsubscribeWebHook()
}

func GetSubscibeWebhookField(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_SUBSCRIBE_WEBHOOK_FIELD)
}

func GetMetricsMedia(cl clientpack.Client, igMediaID string) response.Responser {
	dataParam := map[string]any{
		"ig_media_id": igMediaID,
	}

	return cl.Get(type_ig.IG_GET_METRICS_MEDIA, dataParam)
}

func GetMetricsMediaInsight(cl clientpack.Client, igMediaID string) response.Responser {
	dataParam := map[string]any{
		"ig_media_id": igMediaID,
	}

	return cl.Get(type_ig.IG_GET_METRICS_MEDIA_INSIGHT, dataParam)
}

func GetMetricsAccountInsight(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_METRICS_USER_INSIGHT)
}

func GetConversations(cl clientpack.Client) response.Responser {
	return cl.Get(type_ig.IG_GET_LIST_CONVERSATION)
}

func GetUserConversation(cl clientpack.Client, igUserID string) response.Responser {
	dataParam := map[string]any{
		"ig_user_id": igUserID,
	}

	return cl.Get(type_ig.IG_GET_USER_CONVERSATION, dataParam)
}

func GetMessagesOfConversation(cl clientpack.Client, conversationID string) response.Responser {
	dataParam := map[string]any{
		"conversation_id": conversationID,
	}

	return cl.Get(type_ig.IG_GET_MESSAGES_CONVERSATION, dataParam)
}

func GetMessageInfo(cl clientpack.Client, messageID string) response.Responser {
	dataParam := map[string]any{
		"message_id": messageID,
	}

	return cl.Get(type_ig.IG_GET_INFO_MESSAGE, dataParam)
}
