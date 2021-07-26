package api

import (
	"net/http"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/glaccount/handler"
	"github.com/gin-gonic/gin"
)

type GLAccountResource struct {
	hAccount      *handler.AccountHandler
	hJournal      *handler.JournalHandler
	hOrganization *handler.OrganizationtHandler
}

func NewGLAccountResource(router *gin.Engine, nc *broker.NatsClient) GLAccountResource {
	resource := GLAccountResource{hAccount: handler.NewAccountHandler(nc),
		hJournal:      handler.NewJournalHandler(nc),
		hOrganization: handler.NewOrganizationtHandler(nc)}
	resource.setupGLAccountRoutes(router)
	return resource
}

func (r GLAccountResource) setupGLAccountRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	r.addV1Routes(v1)
}

func (r GLAccountResource) addV1Routes(rg *gin.RouterGroup) {
	gl := rg.Group("/glaccount")

	r.addAccountRoutesV1(gl)
	r.addOrgRoutesV1(gl)

	gl.POST("/journal/post", func(c *gin.Context) {
		rq := &glaccount.PostJournalEntryRq{}
		if err := c.BindJSON(rq); err != nil {
			c.AbortWithStatusJSON(400, err)
			return
		}
		rs, err := r.hJournal.PostEntry(rq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, rs)
		}
	})
}

func (r GLAccountResource) addAccountRoutesV1(rg *gin.RouterGroup) {
	rg.POST("/new", func(c *gin.Context) {
		rq := &glaccount.CreateNewAcctRq{}
		if err := c.BindJSON(rq); err != nil {
			c.AbortWithStatusJSON(400, err)
			return
		}
		rs, err := r.hAccount.CreateNewAccount(rq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, rs)
		}
	})

	rg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		gla, err := r.hAccount.GetAccountById(&glaccount.GetGLAByIdRq{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gla)
		}
	})

}

func (r GLAccountResource) addOrgRoutesV1(rg *gin.RouterGroup) {
	rg.POST("/org/new", func(c *gin.Context) {
		rq := &glaccount.CreateNewOrgRq{}
		if err := c.BindJSON(rq); err != nil {
			c.AbortWithStatus(400)
			return
		}
		rs, err := r.hOrganization.CreateOrg(rq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, rs)
		}
	})
}
