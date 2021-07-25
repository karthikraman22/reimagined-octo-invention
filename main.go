package main

import (
	"database/sql"
	"log"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/config"
	"achuala.in/ledger/glaccount/api"
	"achuala.in/ledger/glaccount/processor"
	"achuala.in/ledger/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.NewConfig("conf.yaml", "GLACCT_")

	s := service.NewService(cfg)

	nc := broker.NewBroker(cfg.String("nats.hosts"))
	nc.Connect()
	defer nc.Disconnect()

	db, err := sql.Open(cfg.String("db.driver"), cfg.String("db.uri"))
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()

	api.NewGLAccountResource(s.Router, nc)

	// Regiser the subscribers
	processor.NewAccountProcessor(nc, db).Init()
	processor.NewJournalProcessor(nc, db).Init()

	s.Run()

}
