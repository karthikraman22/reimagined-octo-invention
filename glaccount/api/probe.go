package api

import (
	"net/http"

	"achuala.in/ledger/broker"
	"github.com/gin-gonic/gin"
)

type HealthProbeResource struct {
	nc *broker.NatsClient
}

func NewHealthProbe(router *gin.Engine, nc *broker.NatsClient) HealthProbeResource {
	resource := HealthProbeResource{nc: nc}
	resource.setupHealthProbes(router)
	return resource
}

func (r HealthProbeResource) setupHealthProbes(router *gin.Engine) {
	health := router.Group("/health")
	rs := make(map[string]string)
	rs["status"] = "UP"
	health.GET("/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, rs)
	})
	health.GET("/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, rs)
	})
}
