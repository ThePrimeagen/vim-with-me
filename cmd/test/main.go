package main

import (
	"fmt"
	"log"

	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		fmt.Printf("msg: %+v\n", msg)
	})

	client.Join("theprimeagen")
	err := client.Connect()
	if err != nil {
		log.Fatal("I AM DEAD AND DYLAN RAIDED ME", err)
	}

	for {
	}
}
