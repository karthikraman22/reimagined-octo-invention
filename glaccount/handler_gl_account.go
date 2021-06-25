package glaccount

import (
	"achuala.in/ledger/broker"
)

type GLAccountHandler struct {
	nc *broker.NatsClient
}

func NewGLAccountHandler(nc *broker.NatsClient) *GLAccountHandler {
	svc := &GLAccountHandler{nc: nc}
	return svc
}

func (h *GLAccountHandler) GetGLAccountById(rq *FindGLAByIdRq) (*FindGLAByIdRs, error) {
	rs := &FindGLAByIdRs{}
	if err := h.nc.Send("glacct.getglaccountbyid", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
