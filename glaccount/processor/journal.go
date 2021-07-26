package processor

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/util"
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

	_, err := p.validateAccounts(rq.Entries)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}

	sql := `INSERT INTO gl_journal(id,tranaction_id,ext_ref,account_id,entry_date, amount,entry_type, notes,
		reversal_id, manual_entry, tsignature, tags) VALUES %s `
	valueStrings := []string{}
	valueArgs := []interface{}{}
	i := 0
	for _, e := range rq.Entries {
		id, _ := uuid.NewUUID()
		acctId := util.ToUuid(e.AccountId)
		reversalId := util.ToUuid(e.ReversalId)
		signature := sha256.Sum256([]byte(fmt.Sprintf("%s.%s.%s.%g.%s", id, e.AccountId, e.TransactionId, e.Amount, e.Type)))
		tagsJson, _ := json.Marshal(e.Tags)
		entryDate, _ := time.Parse("MM-DD-YYYY", e.EntryDate)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d,$%d, $%d, $%d,$%d, $%d, $%d,$%d, $%d, $%d)",
			i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6, i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12))
		valueArgs = append(valueArgs, id)
		valueArgs = append(valueArgs, e.TransactionId)
		valueArgs = append(valueArgs, e.ExternalRef)
		valueArgs = append(valueArgs, acctId)
		valueArgs = append(valueArgs, entryDate)
		valueArgs = append(valueArgs, e.Amount)
		valueArgs = append(valueArgs, e.Type)
		valueArgs = append(valueArgs, e.Notes)
		valueArgs = append(valueArgs, reversalId)
		valueArgs = append(valueArgs, e.ManualEntry)
		valueArgs = append(valueArgs, signature[:])
		valueArgs = append(valueArgs, tagsJson)
		i = i + 1
	}
	stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
	log.Printf("%s", stmt)
	_, err = p.db.Exec(stmt, valueArgs...)
	if err != nil {
		rs.Status = err.Error()
		return rs
	}
	rs.Status = "0"
	return rs
}

func (p *JournalProcessor) validateAccounts(entries []*glaccount.JournalEntry) (map[string]glaccount.GeneralLedgerAccount, error) {
	return nil, nil
}
