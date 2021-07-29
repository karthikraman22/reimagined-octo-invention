package processor

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/repository"
	"achuala.in/ledger/util"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type JournalProcessor struct {
	nc       *broker.NatsClient
	db       *sql.DB
	acctRepo *repository.AccountRepository
}

func NewJournalProcessor(nc *broker.NatsClient, db *sql.DB) *JournalProcessor {
	ar := repository.NewAccountRepository(db)
	p := &JournalProcessor{nc: nc, db: db, acctRepo: ar}
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

	debtors := rq.Entry.Debtors
	creditors := rq.Entry.Creditors

	if err := validateSumTotal(debtors, creditors); err != nil {
		rs.Status = err.Error()
		return rs
	}

	sql := `INSERT INTO gl_journal(id,tranaction_id,ext_ref,account_number,entry_date, amount,entry_type, notes,
		reversal_id, manual_entry, tsignature, tags) VALUES %s `
	valueStrings := []string{}
	valueArgs := []interface{}{}
	i := 0
	for _, d := range debtors {
		valueStrings, valueArgs = buildStatment(valueStrings, valueArgs, i, rq.Entry, d.AccountNumber, d.Amount.Value)
		i = i + 1
	}
	for _, c := range creditors {
		valueStrings, valueArgs = buildStatment(valueStrings, valueArgs, i, rq.Entry, c.AccountNumber, c.Amount.Value)
		i = i + 1
	}
	stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
	log.Printf("%s", stmt)
	_, err := p.db.Exec(stmt, valueArgs...)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	rs.Status = "0"
	return rs
}

func validateSumTotal(debtors []*glaccount.Debtor, creditors []*glaccount.Creditor) error {

	if len(debtors) == 0 && len(creditors) == 0 {
		// TODO: Load the gl debit accounts
		return errors.New("no debtors or creditors")
	}
	var drAmt = 0.0
	var crAmt = 0.0
	for _, d := range debtors {
		if d.Amount == nil || d.Amount.Value <= 0 {
			return errors.New("invalid_debit_amount")
		}
		drAmt = drAmt + d.Amount.Value
	}
	for _, c := range creditors {
		if c.Amount == nil || c.Amount.Value <= 0 {
			return errors.New("invalid_credit_amount")
		}
		crAmt = crAmt + c.Amount.Value
	}
	if drAmt != crAmt {
		return errors.New("debits_credits_mismatch")
	}
	return nil
}

func buildStatment(valueStrings []string, valueArgs []interface{}, i int, e *glaccount.JournalEntry, acctNumber string, amount float64) ([]string, []interface{}) {
	id, _ := uuid.NewUUID()
	reversalId := util.ToUuid(e.ReversalId)
	signature := sha256.Sum256([]byte(fmt.Sprintf("%s.%s.%s.%g.%s", id, acctNumber, e.TransactionId, amount, e.Type)))
	tagsJson, _ := json.Marshal(e.Tags)
	entryDate, _ := time.Parse("20060102", e.EntryDate)
	valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d,$%d, $%d, $%d,$%d, $%d, $%d,$%d, $%d, $%d)",
		i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6, i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12))
	valueArgs = append(valueArgs, id)
	valueArgs = append(valueArgs, e.TransactionId)
	valueArgs = append(valueArgs, e.ExternalRef)
	valueArgs = append(valueArgs, acctNumber)
	valueArgs = append(valueArgs, entryDate)
	valueArgs = append(valueArgs, amount)
	valueArgs = append(valueArgs, e.Type)
	valueArgs = append(valueArgs, e.Notes)
	valueArgs = append(valueArgs, reversalId)
	valueArgs = append(valueArgs, e.ManualEntry)
	valueArgs = append(valueArgs, signature[:])
	valueArgs = append(valueArgs, tagsJson)
	return valueStrings, valueArgs
}
