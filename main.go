package main

import (
	"github.com/anthdm/hollywood/actor"
)

func main() {
	e := actor.NewEngine()
	e.Spawn(NewServer(":6000"), "server") // retornar um PID
	<-make(chan struct{})
}
