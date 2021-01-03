package model

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes ...
type Routes struct {
	RouteEngine *gin.Engine
}

// BuildRoutes ...
func (md *Routes) BuildRoutes() error {
	v1 := md.RouteEngine.Group("/v1")
	{
		v1.GET("/model/:id", func(c *gin.Context) {
			model, err := GetModelByID(c.Param("id"))
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("error : %s", err))
			} else {
				c.JSON(http.StatusOK, model)
			}
		})
	}
	return nil
}
