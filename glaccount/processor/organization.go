package processor

import (
	"database/sql"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type OrganizationProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewOrganizationProcessor(nc *broker.NatsClient, db *sql.DB) *OrganizationProcessor {
	p := &OrganizationProcessor{nc: nc, db: db}
	p.RegisterHandlers()
	return p
}

func (p *OrganizationProcessor) RegisterHandlers() {
	p.RegisterHandler("glacct.org.createorg", "glacct-wrkr", p.createOrg)
}

func (p *OrganizationProcessor) RegisterHandler(subject string, groupName string, h func([]byte) protoreflect.ProtoMessage) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *OrganizationProcessor) createOrg(reqPayLoad []byte) protoreflect.ProtoMessage {
	rq := &glaccount.CreateNewOrgRq{}
	rs := &glaccount.CreateNewOrgRs{}
	if err := proto.Unmarshal(reqPayLoad, rq); err != nil {
		rs.Status = err.Error()
		return rs
	}

	sql := `INSERT INTO gl_organization(id,code,name) VALUES($1,$2,$3)`
	stmt, err := p.db.Prepare(sql)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	defer stmt.Close()
	newOrgId, _ := uuid.NewUUID()
	_, err = stmt.Exec(newOrgId, rq.OrgDetails.Code, rq.OrgDetails.Name)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	rs.Status = "0"
	rs.OrgId = newOrgId.String()
	return rs
}
