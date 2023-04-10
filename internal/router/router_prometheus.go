package router

import (
	prometheusHandler "cpds/cpds-detector/internal/handlers/prometheus"

	"github.com/gin-gonic/gin"
)

func setPrometheusRouter(api *gin.RouterGroup, r *resource) {
	rulesApi := api.Group("prometheus")
	{
		rulesHandler := prometheusHandler.New(r.config, r.logger)
		rulesApi.GET("/query", rulesHandler.Query())
		rulesApi.GET("/query_range", rulesHandler.QueryRange())
	}
}
