package handler

import (
	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
)

type AccountHandler struct {
	nc *broker.NatsClient
}

func NewAccountHandler(nc *broker.NatsClient) *AccountHandler {
	svc := &AccountHandler{nc: nc}
	return svc
}

func (h *AccountHandler) GetGLAccountById(rq *glaccount.GetGLAByIdRq) (*glaccount.GetGLAByIdRs, error) {
	rs := &glaccount.GetGLAByIdRs{}
	if err := h.nc.Send("glacct.getglaccountbyid", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
