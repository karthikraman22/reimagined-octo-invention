package main

import (
	"database/sql"
	"log"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/config"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.NewConfig("conf.yaml")

	s := service.NewService(cfg)

	nc := broker.NewBroker(cfg.String("nats_hosts"))
	nc.Connect()
	defer nc.Disconnect()

	db, err := sql.Open(cfg.String("db.driver"), cfg.String("db.uri"))
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()

	glaccount.NewGLAccountResource(s.Router, nc)

	// Regiser the subscriber
	p := glaccount.NewGLAccountProcessor(nc)
	p.Init()

	s.Run()

}
