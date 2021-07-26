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

func (h *AccountHandler) CreateNewAccount(rq *glaccount.CreateNewAcctRq) (*glaccount.CreateNewAcctRs, error) {
	rs := &glaccount.CreateNewAcctRs{}
	if err := h.nc.Send("glacct.createnewaccount", rq, rs); err != nil {
		rs.Status = err.Error()
		return rs, err
	}
	return rs, nil
}

func (h *AccountHandler) GetAccountById(rq *glaccount.GetGLAByIdRq) (*glaccount.GetGLAByIdRs, error) {
	rs := &glaccount.GetGLAByIdRs{}
	if err := h.nc.Send("glacct.getaccountbyid", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
