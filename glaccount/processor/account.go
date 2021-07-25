package processor

import (
	"database/sql"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type AccountProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewAccountProcessor(nc *broker.NatsClient, db *sql.DB) *AccountProcessor {
	return &AccountProcessor{nc: nc, db: db}
}

func (p *AccountProcessor) Init() {
	p.RegisterHandler("glacct.getglaccountbyid", "glacct-wrkr", p.getGLAccountById)
}

func (p *AccountProcessor) RegisterHandler(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *AccountProcessor) getGLAccountById(reqPayLoad []byte) (protoreflect.ProtoMessage, error) {
	req := &glaccount.GetGLAByIdRq{}
	if err := proto.Unmarshal(reqPayLoad, req); err != nil {
		return nil, err
	}
	acc := &glaccount.GeneralLedgerAccount{Id: uuid.New().String(), Name: req.Id}
	return &glaccount.GetGLAByIdRs{Glaccount: acc, Status: "0"}, nil
}
