package api

import (
	"net/http"

	"achuala.in/ledger/broker"
	"achuala.in/ledger/glaccount"
	"achuala.in/ledger/glaccount/handler"
	"github.com/gin-gonic/gin"
)

type GLAccountResource struct {
	hAccount *handler.AccountHandler
	hJournal *handler.JournalHandler
}

func NewGLAccountResource(router *gin.Engine, nc *broker.NatsClient) GLAccountResource {
	resource := GLAccountResource{hAccount: handler.NewAccountHandler(nc), hJournal: handler.NewJournalHandler(nc)}
	resource.setupGLAccountRoutes(router)
	return resource
}

func (r GLAccountResource) setupGLAccountRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	r.addV1Routes(v1)
}

func (r GLAccountResource) addV1Routes(rg *gin.RouterGroup) {
	gl := rg.Group("/glaccount")

	gl.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		gla, err := r.hAccount.GetGLAccountById(&glaccount.GetGLAByIdRq{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gla)
		}
	})

	gl.POST("/journal/new", func(c *gin.Context) {
		payLoad := &glaccount.NewJournalEntry{}
		if err := c.BindJSON(payLoad); err != nil {
			c.AbortWithStatus(400)
			return
		}
		rq := &glaccount.PostNewJournalEntryRq{Entry: payLoad}
		rs, err := r.hJournal.PostNewGLJournalEntry(rq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, rs)
		}
	})
}
