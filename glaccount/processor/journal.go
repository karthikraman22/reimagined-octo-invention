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
	return &JournalProcessor{nc: nc, db: db}
}

func (p *JournalProcessor) Init() {
	p.RegisterHandler("glacct.jrnl.postnewgljournalentry", "glacct-wrkr", p.postNewGLJournalEntry)
}

func (p *JournalProcessor) RegisterHandler(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *JournalProcessor) postNewGLJournalEntry(reqPayLoad []byte) (protoreflect.ProtoMessage, error) {
	rq := &glaccount.PostNewJournalEntryRq{}
	if err := proto.Unmarshal(reqPayLoad, rq); err != nil {
		return nil, err
	}

	rs := &glaccount.PostNewJournalEntryRs{Status: "0", Id: uuid.NewString()}
	return rs, nil
}
