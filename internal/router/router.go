package router

import (
	"cpds/cpds-detector/internal/handlers"
	"cpds/cpds-detector/internal/middlewares"
	"cpds/cpds-detector/pkg/cpds-detector/config"

	gormlogger "gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type resource struct {
	config *config.Config
	logger *zap.Logger
	db     *gorm.DB
}

func InitRouter(debug bool, config *config.Config, logger *zap.Logger, db *gorm.DB) *gin.Engine {
	r := &resource{
		config: config,
		logger: logger,
		db:     db,
	}

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		db.Logger.LogMode(gormlogger.Silent)
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middlewares.LoggerMiddleware(logger))

	// test route
	router.GET("/ping", handlers.GetPing)

	apiv1 := router.Group("/api/v1")
	setAnalysisRouter(apiv1, r)
	setPrometheusRouter(apiv1, r)
	setMonitorRouter(apiv1, r)

	return router
}
