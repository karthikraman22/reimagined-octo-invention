package glaccount

import "achuala.in/ledger/broker"

type GLAccountJournalHandler struct {
	nc *broker.NatsClient
}

func (h *GLAccountJournalHandler) PostNewGLJournalEntry(rq *PostNewGLAcctJrnlEntryRq) (*PostNewGLAcctJrnlEntryRs, error) {
	rs := &PostNewGLAcctJrnlEntryRs{}
	if err := h.nc.Send("glaccount.PostNewGLAcctJrnlEntry", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}