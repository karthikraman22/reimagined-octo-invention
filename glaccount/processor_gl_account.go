package glaccount

import (
	"database/sql"

	"achuala.in/ledger/broker"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type GLAccountProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewGLAccountProcessor(nc *broker.NatsClient, db *sql.DB) *GLAccountProcessor {
	return &GLAccountProcessor{nc: nc, db: db}
}

func (p *GLAccountProcessor) Init() {
	p.RegisterHandler("glacct.getglaccountbyid", "glacct-wrkr", p.getGLAccountById)
}

func (p *GLAccountProcessor) RegisterHandler(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *GLAccountProcessor) getGLAccountById(reqPayLoad []byte) (protoreflect.ProtoMessage, error) {
	req := &FindGLAByIdRq{}
	if err := proto.Unmarshal(reqPayLoad, req); err != nil {
		return nil, err
	}
	acc := &GeneralLedgerAccount{Id: uuid.New().String(), Name: req.Id}
	return &FindGLAByIdRs{Glaccount: acc, Status: "0"}, nil
}
