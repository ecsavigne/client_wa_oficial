package main

import (
	"fmt"
	"os"
	"os/signal"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2/client"
	"github.com/ecsavigne/client_wa_oficial/v2/event"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
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

func main() {
	// Create one Client of WhatsApp Official
	my_client := clientoficial.NewClientWA(
		// clientoficial.WithEnvFilePath("../config_env"),
		clientoficial.WithEnvFilePath("../../config_env"),
		// clientoficial.WithWebhookSocket("wss://webhooks.xxx.com/wa_official/ws"),
		clientoficial.WithEventHandle(EventHandler),
		clientoficial.WithWabaID(""),
		clientoficial.WithToken(""),
	)

	if my_client.Error != nil {
		v, ok := my_client.Error.(*response.Error)
		if ok {
			fmt.Println("Errrr: \n", v)
		}
		return
	}

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
	// r := my_client.GetAllTplFromLibrary(clientoficial.QueryData{"language": str.Join([]string{"pt_BR"}, ",")}) // id = 1111111, name = hello_world
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
	// rr := my_client.DeleteTemplate(clientoficial.ParamDelete{
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
	// fmt.Printf("Response: %+v\n", rr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
