package errors

import "fmt"

const (
	SUCCESS        = 0
	DATABASE_ERROR = 101
	SOCKET_ERROR   = 102

	ANALYSIS_RULE_UPDATED_ERROR = 1001

	PROMETHEUS_QUERY_ERROR       = 2001
	PROMETHEUS_QUERY_RANGE_ERROR = 2002

	MONITOR_GET_NODE_STATUS_ERROR              = 3001
	MONITOR_GET_NODE_INFO_ERROR                = 3002
	MONITOR_GET_NODE_RESOURCES_ERROR           = 3003
	MONITOR_GET_NODE_CONTAINER_STATUS_ERROR    = 3004
	MONITOR_GET_CLUSTER_RESOURCES_ERROR        = 3005
	MONITOR_GET_CLUSTER_CONTAINER_STATUS_ERROR = 3006
	MONITOR_GET_TARGET_ERROR                   = 3007
)

var DetectorResultCodeMap = map[uint16]string{
	SUCCESS:        "success",
	DATABASE_ERROR: "database error",
	SOCKET_ERROR:   "socket error",

	ANALYSIS_RULE_UPDATED_ERROR: "rule updated error",

	PROMETHEUS_QUERY_ERROR:       "prometheus query error",
	PROMETHEUS_QUERY_RANGE_ERROR: "prometheus query range error",
}

type Error struct {
	Err        error
	ResultCode uint16
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", DetectorResultCodeMap[e.ResultCode], e.Err.Error())
}

func NewError(resultCode uint16, err error) error {
	return &Error{
		ResultCode: resultCode,
		Err:        err,
	}
}

func IsErrorWithCode(err error, desiredResultCode uint16) bool {
	if err == nil {
		return false
	}

	serverError, ok := err.(*Error)
	if !ok {
		return false
	}

	return serverError.ResultCode == desiredResultCode
}
