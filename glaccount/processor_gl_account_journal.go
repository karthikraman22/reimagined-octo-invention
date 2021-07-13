package glaccount

import (
	"database/sql"

	"achuala.in/ledger/broker"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type GLAccountJournalProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewGLAccountJournalProcessor(nc *broker.NatsClient, db *sql.DB) *GLAccountJournalProcessor {
	return &GLAccountJournalProcessor{nc: nc, db: db}
}

func (p *GLAccountJournalProcessor) Init() {
	p.RegisterHandler("glacct.jrnl.postnewgljournalentry", "glacct-wrkr", p.postNewGLJournalEntry)
}

func (p *GLAccountJournalProcessor) RegisterHandler(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *GLAccountJournalProcessor) postNewGLJournalEntry(reqPayLoad []byte) (protoreflect.ProtoMessage, error) {
	rq := &PostNewGLAcctJrnlEntryRq{}
	if err := proto.Unmarshal(reqPayLoad, rq); err != nil {
		return nil, err
	}

	rs := &PostNewGLAcctJrnlEntryRs{Status: "0", Id: uuid.NewString()}
	return rs, nil
}
