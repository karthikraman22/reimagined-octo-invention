package glaccount

import (
	"net/http"

	"achuala.in/ledger/broker"
	"github.com/gin-gonic/gin"
)

type GLAccountResource struct {
	glaHandler     *GLAccountHandler
	glaJrnlHandler *GLAccountJournalHandler
}

func NewGLAccountResource(router *gin.Engine, nc *broker.NatsClient) GLAccountResource {
	handler := NewGLAccountHandler(nc)
	resource := GLAccountResource{glaHandler: handler, glaJrnlHandler: &GLAccountJournalHandler{nc: nc}}
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
		gla, err := r.glaHandler.GetGLAccountById(&FindGLAByIdRq{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gla)
		}
	})

	gl.POST("/journal/new", func(c *gin.Context) {
		payLoad := &NewGLAcctJournalEntry{}
		if err := c.BindJSON(payLoad); err != nil {
			c.AbortWithStatus(400)
			return
		}
		rq := &PostNewGLAcctJrnlEntryRq{Entry: payLoad}
		rs, err := r.glaJrnlHandler.PostNewGLJournalEntry(rq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, rs)
		}
	})
}