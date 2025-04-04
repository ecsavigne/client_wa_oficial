package main

import (
	"fmt"
	"os"
	"os/signal"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2"
	"github.com/ecsavigne/client_wa_oficial/v2/event"
	"github.com/ecsavigne/client_wa_oficial/v2/types/message"
	"github.com/ecsavigne/client_wa_oficial/v2/types/response"
)

// EventHandler is a function that handles events
func EventHandler(data any) {
	switch v := data.(type) {
	case []byte:
		fmt.Println("Data:\n", string(v))
	case *event.EventErrorSocketConnect:
		fmt.Println("EventErrorSocketConnect:\n", v)
	}
}

// main
func main() {
	// Create one Client of WhatsApp Official
	my_client := clientoficial.NewClientWA(clientoficial.Config{
		EnvFilePath:   "../../config_env",
		WebhookSocket: "wss://ws.example.com/ws1",
		EventHandle:   EventHandler,
	})

	if my_client.Error != nil {
		v, ok := my_client.Error.(*response.Error)
		if ok {
			fmt.Println("Errrr: \n", v)
		}
		return
	}

	var m message.Messager = &message.MessageResponse{
		MessagerKernel: message.MessagerKernel{
			Type: "location",
		},
		Media: &message.Media{
			Link: "https://example.com/audio.mp3",
		},
	}

	var mText message.Messager = &message.MessageText{
		MessagerKernel: message.MessagerKernel{
			Type: "location",
		},
	}

	fmt.Printf("M: %s\n", m)
	fmt.Printf("M: %s\n", mText)

	// fmt.Println(my_client.GetInfoAllNumberInWA())
	// For handling the interruption of the application with Ctrl+C allows you to see the events handled by EventHandle
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
