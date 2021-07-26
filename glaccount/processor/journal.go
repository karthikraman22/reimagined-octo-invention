package processor

import (
	"database/sql"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type JournalProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewJournalProcessor(nc *broker.NatsClient, db *sql.DB) *JournalProcessor {
	p := &JournalProcessor{nc: nc, db: db}
	p.RegisterHandlers()
	return p
}

func (p *JournalProcessor) RegisterHandlers() {
	p.RegisterHandler("glacct.jrnl.postentry", "glacct-wrkr", p.postEntry)
}

func (p *JournalProcessor) RegisterHandler(subject string, groupName string, h func([]byte) protoreflect.ProtoMessage) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *JournalProcessor) postEntry(reqPayLoad []byte) protoreflect.ProtoMessage {
	rq := &glaccount.PostJournalEntryRq{}
	rs := &glaccount.PostJournalEntryRs{}
	if err := proto.Unmarshal(reqPayLoad, rq); err != nil {
		rs.Status = err.Error()
		return rs
	}
	p.db.Begin()
	sql := `INSERT INTO gl_journal`
	stmt, err := p.db.Prepare(sql)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	defer stmt.Close()
	newId, _ := uuid.NewUUID()
	_, err = stmt.Exec()
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	rs.Status = "0"
	rs.Id = newId.String()
	return rs
}
