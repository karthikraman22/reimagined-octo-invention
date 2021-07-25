package handler

import (
	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
)

type JournalHandler struct {
	nc *broker.NatsClient
}

func NewJournalHandler(nc *broker.NatsClient) *JournalHandler {
	return &JournalHandler{nc: nc}
}

func (h *JournalHandler) PostNewGLJournalEntry(rq *glaccount.PostNewJournalEntryRq) (*glaccount.PostNewJournalEntryRs, error) {
	rs := &glaccount.PostNewJournalEntryRs{}
	if err := h.nc.Send("glacct.jrnl.postnewgljournalentry", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
