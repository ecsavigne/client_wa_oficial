package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	clientpack "github.com/ecsavigne/client_wa_oficial/v2/client"

	"github.com/ecsavigne/client_wa_oficial/v2/example/examplemessage"
	"github.com/ecsavigne/client_wa_oficial/v2/types/general/response"
	"github.com/ecsavigne/client_wa_oficial/v2/types/wpp/response/event"
)

// EventHandler is a function that handles events
func EventHandler(data any) {
	fmt.Println("handler function")
	switch v := data.(type) {
	case []byte:
		fmt.Println("Data:\n", string(v))
	case *event.ErrorSocketConnectEvent:
		fmt.Println("ErrorSocketConnectEvent:\n", v)
	case *event.StatusMessageEvent:
		fmt.Println("StatusMessageEvent:\n", v.Entry[0].Changes[0].Value.Statuses[0].Status)
	case *event.MessageAudioEvent:
		fmt.Println("MessageAudioEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Audio)
	case *event.MessageButtonEvent:
		fmt.Println("MessageButtonEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Button)
	case *event.MessageContactEvent:
		fmt.Println("MessageContatEvent:\n", v.Entry[0].Changes[0].Value.Contacts[0])
	case *event.MessageDocumentEvent:
		fmt.Println("MessageDocumentEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Document)
	case *event.MessageImageEvent:
		fmt.Println("MessageImageEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Image)
	case *event.MessageInteractiveEvent:
		fmt.Println("MessageInteractiveEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Interactive)
	case *event.MessageOrderEvent:
		fmt.Println("MessageOrderEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Order)
	case *event.MessageStickerEvent:
		fmt.Println("MessageStickerEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Sticker)
	case *event.MessageSystemEvent:
		fmt.Println("MessageSystemEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].System)
	case *event.MessageTextEvent:
		fmt.Println("MessageTextEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Text.Body)
	case *event.MessageVideoEvent:
		fmt.Println("MessageVideoEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Video)
	case *event.MessageUnknownEvent:
		fmt.Println("MessageUnknownEvent:\n", v.Entry[0].Changes[0].Value.Messages[0])
	case *event.MessageLocationEvent:
		fmt.Println("MessageLocationEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Location)
	case *event.MessageReactionEvent:
		fmt.Println("MessageReactionEvent:\n", v.Entry[0].Changes[0].Value.Messages[0].Reaction)
	}
}

func testWpp(client clientpack.Client, rr response.Responser) {
	fmt.Println(strings.Repeat(".", 25), " Test Client Wpp ", strings.Repeat(".", 25))
	// Create one Client of WhatsApp Official
	// my_client := wpp.NewClientWA(
	// 	// wpp.WithEnvFilePath("../config_env"),
	// 	wpp.WithEnvFilePath("../../config_env"),
	// 	// wpp.WithWebhookSocket("wss://webhooks.xxx.com/wa_official/ws"),
	// 	// wpp.WithEventHandle(EventHandler),
	// 	// wpp.WithWabaID(""),
	// 	// wpp.WithToken(""),
	// )

	// if my_client.Error != nil {
	// 	// v, ok := my_client.Error.(*response.Error)
	// 	// if ok {
	// 	// 	fmt.Println("Errrr: \n", v)
	// 	// }
	// 	return
	// }

	// var m message.Messager = &message.MessageResponse{
	// 	MessagerKernel: message.MessagerKernel{
	// 		Type: "location",
	// 	},
	// 	Media: &message.Media{
	// 		Link: "https://example.com/audio.mp3",
	// 	},
	// }

	// fmt.Println("Nums: \n", my_client.GetInfoAllNumberInWA())

	// fmt.Println(my_client.GetInfoAllNumberInWA())
	// For handling the interruption of the application with Ctrl+C allows you to see the events handled by EventHandle
	// r := my_client.GetAllTplFromWaba("111_waba_id") // id = 1111111, name = hello_world
	// r := my_client.GetAllTplFromLibrary(wpp.QueryData{"language": str.Join([]string{"pt_BR"}, ",")}) // id = 1111111, name = hello_world
	// tplArr := r.GetTemplateResponse()
	// fmt.Println("Response: \n", len(tplArr.Data))
	// for _, v := range tplArr.Data {
	// 	fmt.Printf("Template Name: %s, Type: %s, LocateTemplate: %s\n", v.Name, v.GetTypeTpl(), v.GetLocatedTpl())
	// }
	// r := my_client.GetTemplateById("111111111111") // id = 1122206686609499, name = hello_world
	// r := my_client.GetTemplateByName("event_rsvp_reminder_2") // id = 1122206686622209499, name = hello_world

	// fmt.Println("Response: \n", r)

	// rr := my_client.SendReadNotification("wamid.")
	// tpl := &types.MockupTemplate{
	// 	Name:     "seguir_suporte_test_create",
	// 	Language: "pt_BR",
	// 	Category: "UTILITY",
	// 	Components: []types.MockupComponent{
	// 		// header
	// 		// {},
	// 		// body
	// 		{
	// 			Type: types.TC_BODY,
	// 			Text: "Olá {{1}}, estamos aguardando sua resposta sobre o caso de suporte {{2}}. Caso já tenha resolvido, por favor nos avise.",
	// 			Example: &types.Example{
	// 				BodyText: &[]types.PositionalParams{
	// 					{"Carlos", "SUP-9812"},
	// 				},
	// 			},
	// 		},
	// 		// bootons
	// 		{
	// 			Type: types.TC_BUTTONS,
	// 			ArrayButton: &types.ArrayButton{
	// 				{
	// 					Type: "QUICK_REPLY",
	// 					Text: "Caso resolvido",
	// 				},
	// 				{
	// 					Type: "QUICK_REPLY",
	// 					Text: "Ainda preciso de ajuda",
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// rr := my_client.DeleteTemplate(wpp.ParamDelete{
	// 	Name: "my_template_name_complete",
	// })
	// rr := my_client.UpdateTemplate(&types.MockupTemplate{
	// 	Name:     "my_template_header_footer_body",
	// 	Language: "en_US",
	// 	Category: "MARKETING",
	// 	Components: []types.MockupComponent{
	// 		{
	// 			Type: types.TC_BODY,
	// 			Text: "Olá TTT, estamos aguardando sua resposta sobre o caso de suporte PP. Caso já tenha resolvido, por favor nos avise.",
	// 		},
	// 	},
	// })

	// rr := my_client.SubscribedWabaInApps("1415160156667406")
	// rr := my_client.GetInfoSubscribedWaba("1415160156667406")
	// rr := my_client.UnSubscribedWaba("1415160156667406")
	// rr := my_client.RegisterForUseApi("984286011433245")
	// token := ""
	// rr := my_client.DebugToken(token)
	// rr := my_client.GetLimiteMsg(143390)
	// rr := my_client.GetLimiteMsg(143390)
	// rr := my_client.UnregisterNumber(94571)

	// fmt.Printf("Response: %s\n", rr)
	// fmt.Println("Permitions: ", types.GetPermission().Get("whatsapp_business_management", "Label"), "description: ", types.GetPermission().Get("whatsapp_business_management", "Description"))
	fmt.Println(strings.Repeat(".", 25), " Fin testing Client WhatsApp ", strings.Repeat(".", 25))
}

func testIg(client clientpack.Client, rr response.Responser) {
	fmt.Println(strings.Repeat(".", 25), " Test Client IG ", strings.Repeat(".", 25))
	client = examplemessage.NewClientIG() // fmt.Printf("Response Client IG: %s\n", client)
	// rr = examplemessage.SendMessageText(client, "2014023525994448", "Habla test ok, prueba de mensaje de texto") // 948876401474591
	// fmt.Printf("Response: %s\n", rr)
	// rr = examplemessage.SendMessageImage(client, "2014023525994448", "https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExeTNvcThmcG85OTRlMTF5MTZueWk4ZGMxcDRkMHcxMXA1Z2oxNzZlaSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/O3wIKp4DlubUl1Bezt/giphy.gif", "") // 948876401474591
	// rr = examplemessage.SendMessageImageCollage(client, "2014023525994448", []string{}, []string{"948876401474591", "948876401474591", "948876401474591", "1131758486677092", "1295230589210858", "1314539007200675"}) // 948876401474591
	// rr = examplemessage.SendMessageAudio(client, "2014023525994448", "https://oficial.crmsocialhub.com.br/static/test.aac", "817700310893199") //2741747206192358
	// rr = examplemessage.SendMessageVideo(client, "2014023525994448", "https://res.cloudinary.com/dczar4xfh/video/upload/v1776219555/samples/dance-2.mp4", "1942160596402518")
	// rr = examplemessage.SendMessageStickerLike(client, "2014023525994448")
	// rr = examplemessage.SendMessagePDFFile(client, "2014023525994448", "https://oficial.crmsocialhub.com.br/static/Fatura%20de%20novembro.pdf", "995658263029909")
	// rr = examplemessage.SendMessageReaction(client, "2014023525994448", "aWdfZAG1faXRlbToxOklHTWVzc2FnZAUlEOjE3ODQxNDQ4MTgyMjA5NjMwOjM0MDI4MjM2Njg0MTcxMDMwMTI0NDI3NjI0NDk2NzIxNDY0NzM0MzozMjc2NTU3MTU2MzMwMTY1NDg5MDI3MzEyMTExOTUwMjMzNgZDZD", "", "unreact")
	// rr = examplemessage.GetAccountInfoBusiness(client)
	// rr = examplemessage.SendMessageQuickReplies(client, "2014023525994448")
	// rr = examplemessage.SendPresence(client, "2014023525994448", "typing_on")
	// time.Sleep(10 * time.Second)
	// rr = examplemessage.SendPresence(client, "2014023525994448", "typing_off")
	// rr = examplemessage.SendMessagePersistentMenu(client)
	// rr = examplemessage.GetPersistentMenu(client)
	// rr = examplemessage.DeletePersistentMenu(client)
	// rr = examplemessage.SendMessageIceBreakers(client)
	// rr = examplemessage.GetIceBreakers(client)
	// rr = examplemessage.DeleteIceBreakers(client)
	// rr = examplemessage.GetInstagramLink(client)
	// rr = examplemessage.CreateInstagramWelcomeMessageFlowsADS(client)
	// rr = examplemessage.GetWelcomeMessageFlowsADS(client)
	rr = examplemessage.UpdateWelcomeMessageFlowsADS(client)
	// rr = examplemessage.DeleteWelcomeMessageFlowsADS(client, map[string]any{"flow_id": "2443201722808096"})
	fmt.Printf("Response: %v\n", rr)

	fmt.Println(strings.Repeat(".", 25), " Fin testing Client IG ", strings.Repeat(".", 25))
}

func main() {
	var (
		client clientpack.Client
		rr     response.Responser
	)

	// testWpp(client, rr)
	testIg(client, rr)

	fmt.Println("Press Ctrl+C to exit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
