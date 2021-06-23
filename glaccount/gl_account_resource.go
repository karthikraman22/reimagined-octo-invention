package glaccount

import (
	"net/http"

	"achuala.in/ledger/broker"
	"github.com/gin-gonic/gin"
)

type GLAccountResource struct {
	glaHandler *GLAccountHandler
}

func NewGLAccountResource(router *gin.Engine, nc *broker.NatsClient) GLAccountResource {
	handler := NewGLAccountHandler(nc)
	resource := GLAccountResource{glaHandler: handler}
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
}
