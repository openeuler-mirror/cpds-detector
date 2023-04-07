package router

import (
	analysisHandler "cpds/cpds-detector/internal/handlers/analysis"

	"github.com/gin-gonic/gin"
)

func setAnalysisRouter(api *gin.RouterGroup, r *resource) {

	rulesHandler := analysisHandler.New(r.logger)
	api.GET("/rule_updated", rulesHandler.RuleUpdated())
}
