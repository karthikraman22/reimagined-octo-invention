package handler

import (
	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
)

type OrganizationtHandler struct {
	nc *broker.NatsClient
}

func NewOrganizationtHandler(nc *broker.NatsClient) *OrganizationtHandler {
	return &OrganizationtHandler{nc: nc}
}

func (h *OrganizationtHandler) CreateOrg(rq *glaccount.CreateNewOrgRq) (*glaccount.CreateNewOrgRs, error) {
	rs := &glaccount.CreateNewOrgRs{}
	if err := h.nc.Send("glacct.org.createorg", rq, rs); err != nil {
		rs.Status = err.Error()
		return rs, err
	}
	return rs, nil
}
