package processor

import (
	"database/sql"

	"achuala.in/ledger/broker"
)

type OrganizationProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}
