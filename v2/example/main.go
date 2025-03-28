package main

import (
	"fmt"
	"os"
	"os/signal"

	clientoficial "github.com/ecsavigne/client_wa_oficial/v2"
	"github.com/ecsavigne/client_wa_oficial/v2/event"
	"github.com/ecsavigne/client_wa_oficial/v2/types"
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
		v, ok := my_client.Error.(*types.Error)
		if ok {
			fmt.Println("Errrr: \n", v)
		}
		return
	}

	var m types.Messager = types.NewMessageResponse(&types.MessageResponse{
		Header: types.Header{
			Type: "location",
		},
		Media: &types.Media{
			Link: "https://example.com/audio.mp3",
		},
	})

	fmt.Printf("M: %s\n", m)

	// fmt.Println(my_client.GetInfoAllNumberInWA())
	// For handling the interruption of the application with Ctrl+C allows you to see the events handled by EventHandle
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
