package router

import (
	"cpds/cpds-detector/internal/handlers"
	"cpds/cpds-detector/internal/middlewares"

	gormlogger "gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type resource struct {
	logger *zap.Logger
	db     *gorm.DB
}

func InitRouter(debug bool, logger *zap.Logger, db *gorm.DB) *gin.Engine {
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

	router.Group("/api/v1")

	return router
}
