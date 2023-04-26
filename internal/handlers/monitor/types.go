package monitor

type clusterMonitorQueryParams struct {
	StartTime  int64 `json:"start_time"`
	EndTime    int64 `json:"end_time"`
	StepSecond int64 `json:"step"`
}

type nodeInfoQueryParams struct {
	Instance string `json:"instance"`
}

type nodeMonitorDataQueryParams struct {
	Instance   string `json:"instance"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	StepSecond int64  `json:"step"`
}
