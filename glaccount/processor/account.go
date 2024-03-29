package processor

import (
	"database/sql"
	"encoding/json"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	u "achuala.in/ledger/util"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type AccountProcessor struct {
	nc *broker.NatsClient
	db *sql.DB
}

func NewAccountProcessor(nc *broker.NatsClient, db *sql.DB) *AccountProcessor {
	p := &AccountProcessor{nc: nc, db: db}
	p.RegisterHandlers()
	return p
}

func (p *AccountProcessor) RegisterHandlers() {
	p.RegisterHandler("glacct.getaccountbyid", "glacct-wrkr", p.getAccountById)
	p.RegisterHandler("glacct.createnewaccount", "glacct-wrkr", p.createNewAccount)
}

func (p *AccountProcessor) RegisterHandler(subject string, groupName string, h func([]byte) protoreflect.ProtoMessage) error {
	return p.nc.Subscribe(subject, groupName, h)
}

func (p *AccountProcessor) createNewAccount(reqPayLoad []byte) protoreflect.ProtoMessage {
	rq := &glaccount.CreateNewAcctRq{}
	rs := &glaccount.CreateNewAcctRs{}
	if err := proto.Unmarshal(reqPayLoad, rq); err != nil {
		rs.Status = err.Error()
		return rs
	}

	sql := `INSERT INTO gl_account(id, number, code, description, parent_id, org_id, type, usage, 
			disabled, allow_manual_entries, tags) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	stmt, err := p.db.Prepare(sql)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	defer stmt.Close()
	newAcctId, _ := uuid.NewUUID()
	parentId := u.ToUuid(rq.AcctDetails.ParentId)
	tagsJson, _ := json.Marshal(rq.AcctDetails.Tags)
	orgId := u.ToUuid(rq.AcctDetails.OrganizationId)
	_, err = stmt.Exec(newAcctId, rq.AcctDetails.Number, rq.AcctDetails.Code,
		rq.AcctDetails.Description, parentId,
		orgId, rq.AcctDetails.Type, rq.AcctDetails.Usage,
		rq.AcctDetails.Disabled, rq.AcctDetails.AllowManualEntries, tagsJson)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	rs.Status = "0"
	rs.AccountId = newAcctId.String()
	return rs
}

func (p *AccountProcessor) getAccountById(reqPayLoad []byte) protoreflect.ProtoMessage {
	req := &glaccount.GetGLAByIdRq{}
	if err := proto.Unmarshal(reqPayLoad, req); err != nil {
		return nil
	}
	acc := &glaccount.GeneralLedgerAccount{Id: uuid.New().String()}
	return &glaccount.GetGLAByIdRs{AcctDetails: acc, Status: "0"}
}
