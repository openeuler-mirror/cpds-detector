package analysis

import (
	"cpds/cpds-detector/internal/core"
	cpdserr "cpds/cpds-detector/internal/pkg/errors"
	"cpds/cpds-detector/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	RuleUpdated() gin.HandlerFunc
}

type handler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) RuleUpdated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := core.RuleUpdated(); err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, err))
			return
		}

		response.HandleOK(ctx, nil)
	}
}
