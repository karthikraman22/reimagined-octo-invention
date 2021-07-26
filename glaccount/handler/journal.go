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

func (h *JournalHandler) PostEntry(rq *glaccount.PostJournalEntryRq) (*glaccount.PostJournalEntryRs, error) {
	rs := &glaccount.PostJournalEntryRs{}
	if err := h.nc.Send("glacct.jrnl.postentry", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
