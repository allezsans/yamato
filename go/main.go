package main

import (
	"github.com/allezsans/yamato/go/discord"
)

func main() {
	loop := make(chan bool)
	go discord.Start()
	<-loop
}
