package glaccount

import (
	"log"

	"achuala.in/ledger/broker"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type GLAccountProcessor struct {
	nc *broker.NatsClient
}

func NewGLAccountProcessor(nc *broker.NatsClient) *GLAccountProcessor {
	return &GLAccountProcessor{nc: nc}
}

func (p *GLAccountProcessor) Init() {
	p.RegisterHandler("glaccount.GetGLAccountById", "glacct-wrkr", getGLAccountById)
}

func (p *GLAccountProcessor) RegisterHandler(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {
	defer log.Printf("Subscribed to : %s", subject)
	return p.nc.Subscribe(subject, groupName, h)
}

func getGLAccountById(reqPayLoad []byte) (protoreflect.ProtoMessage, error) {
	req := &FindGLAByIdRq{}
	if err := proto.Unmarshal(reqPayLoad, req); err != nil {
		return nil, err
	}
	acc := &GeneralLedgerAccount{Id: uuid.New().String(), Name: req.Id}
	return &FindGLAByIdRs{Glaccount: acc, Status: "0"}, nil
}
