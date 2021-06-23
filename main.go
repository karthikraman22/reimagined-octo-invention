package main

import (
	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/service"
)

func main() {

	s := service.NewService("ledger")

	nc := broker.NewBroker("localhost:4222")
	nc.Connect()
	defer nc.Disconnect()

	glaccount.NewGLAccountResource(s.Router, nc)

	// Regiser the subscriber
	p := glaccount.NewGLAccountProcessor(nc)
	p.Init()

	s.Run()

}
