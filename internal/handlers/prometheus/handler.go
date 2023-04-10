package prometheus

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cpds/cpds-detector/internal/models/prometheus"
	cpdserr "cpds/cpds-detector/internal/pkg/errors"
	"cpds/cpds-detector/internal/pkg/response"
	"cpds/cpds-detector/pkg/cpds-detector/config"
	prometheusutils "cpds/cpds-detector/pkg/utils/prometheus"
	timeutils "cpds/cpds-detector/pkg/utils/time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	Query() gin.HandlerFunc

	QueryRange() gin.HandlerFunc
}

type handler struct {
	logger   *zap.Logger
	operator prometheus.Operator
}

func New(config *config.Config, logger *zap.Logger) Handler {
	return &handler{
		logger:   logger,
		operator: prometheus.NewOperator(config.PrometheusOptions.Host, config.PrometheusOptions.Port),
	}
}

func (h *handler) Query() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		opt, err := parseQueryParams(ctx)
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, err))
			return
		}

		responseData, err := h.operator.Query(opt.expression, opt.time)
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, err))
			return
		} else if len(responseData.MetricValues) == 0 {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, errors.New("no metric data")))
			return
		}

		response.HandleOK(ctx, responseData)
	}
}

func (h *handler) QueryRange() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		opt, err := parseRangeQueryParams(ctx)
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_RANGE_ERROR, err))
			return
		}

		responseData, err := h.operator.QueryRange(opt.expression, opt.startTime, opt.endTime, time.Duration(opt.step)*time.Second)
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, err))
			return
		} else if len(responseData.MetricValues) == 0 {
			response.HandleError(ctx, http.StatusInternalServerError, cpdserr.NewError(cpdserr.PROMETHEUS_QUERY_ERROR, errors.New("no metric data")))
			return
		}

		response.HandleOK(ctx, responseData)
	}
}

func parseQueryParams(ctx *gin.Context) (*queryParams, error) {
	expression := ctx.Query("query")
	if expression == "" {
		return nil, errors.New("query cannot be empty")
	}

	timeStr := ctx.DefaultQuery("time", strconv.FormatInt(time.Now().Unix(), 10))
	time, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil && !timeutils.IsTimestamp(time) {
		return nil, fmt.Errorf("invalid timestamp format: %d", time)
	}

	p := &queryParams{
		expression: expression,
		time:       time,
	}
	return p, nil
}

func parseRangeQueryParams(ctx *gin.Context) (*rangeQueryParams, error) {
	expression := ctx.Query("query")
	if expression == "" && !prometheusutils.IsExprValid(expression) {
		return nil, fmt.Errorf("invalid query expression: %s", expression)
	}

	startTime, err := strconv.ParseInt(ctx.Query("start_time"), 10, 64)
	if err != nil && !timeutils.IsTimestamp(startTime) {
		return nil, fmt.Errorf("invalid timestamp format: %d", startTime)
	}

	endTime, err := strconv.ParseInt(ctx.Query("end_time"), 10, 64)
	if err != nil && !timeutils.IsTimestamp(endTime) {
		return nil, fmt.Errorf("invalid timestamp format: %d", endTime)
	}

	step, err := strconv.Atoi(ctx.Query("step"))
	if err != nil {
		return nil, fmt.Errorf("invalid step format: %s", ctx.Query("step"))
	}

	p := &rangeQueryParams{
		expression: expression,
		startTime:  startTime,
		endTime:    endTime,
		step:       step,
	}
	return p, nil
}
