package handler

import (
	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
)

type OrganizationtHandler struct {
	nc *broker.NatsClient
}

func (h *OrganizationtHandler) CreateGLAOrg(rq *glaccount.CreateNewOrgRq) (*glaccount.CreateNewOrgRs, error) {
	rs := &glaccount.CreateNewOrgRs{}
	if err := h.nc.Send("glacct.org.createglaorg", rq, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
