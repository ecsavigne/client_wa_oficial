# Code of example of how to use client official and handle of the event
```
package main

import (
	"fmt"
	"os"
	"os/signal"

	clientoficial "github.com/ecsavigne/client_wa_oficial"
	"github.com/ecsavigne/client_wa_oficial/event"
	"github.com/ecsavigne/client_wa_oficial/types"
)

func EventHandler(data any) {
	switch v := data.(type) {
	case []byte:
		fmt.Println("Data:\n", string(v))
	case *event.EventErrorSocketConnect:
		fmt.Println("EventErrorSocketConnect:\n", v)
	}
}

func main() {
	my_client := clientoficial.NewClientWA(clientoficial.Config{
		EnvFilePath:   "./config_env",
		WebhookSocket: "wss://webhooks.savcoe-services.com/wa_official/ws1",
		EventHandle:   EventHandler,
	})

	if my_client.Error != nil {
		v, ok := my_client.Error.(*types.Error)
		if ok {
			fmt.Println("Errrr: \n", v)
		}
		return
	}

	// For handling the interruption of the application with Ctrl+C allows you to see the events handled by EventHandle
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
```