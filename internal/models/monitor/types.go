package monitor

type promResponse struct {
	Status string `json:"status"`
	Data   struct {
		Targets []struct {
			DiscoveredLabels struct {
				Address string `json:"__address__"`
			} `json:"discoveredLabels"`
			Health string `json:"health"`
		} `json:"activeTargets"`
	} `json:"data"`
}

type MonitorTargets struct {
	Targets []struct {
		Instance string `json:"instance"`
		Status   string `json:"status"`
	} `json:"targets"`
}

type NodeInfo struct {
	Instance      string `json:"instance"`
	Arch          string `json:"arch"`
	KernelVersion string `json:"kernal_version"`
	OSVersion     string `json:"os_version"`
}

type NodeStatus struct {
	Instance  string `json:"instance"`
	Container struct {
		Running int `json:"running"`
		Total   int `json:"total"`
	} `json:"container"`
	Cpu struct {
		Usage     float64 `json:"usage"`
		UsedCore  float64 `json:"used_core"`
		TotalCore float64 `json:"total_core"`
	} `json:"cpu"`
	Memory struct {
		Usage      float64 `json:"usage"`
		UsedBytes  float64 `json:"used_bytes"`
		TotalBytes float64 `json:"total_bytes"`
	} `json:"memory"`
	Disk struct {
		Usage      float64 `json:"usage"`
		UsedBytes  float64 `json:"used_bytes"`
		TotalBytes float64 `json:"total_bytes"`
	} `json:"disk"`
}

func (mt *MonitorTargets) addTargets(instance, status string) error {
	mt.Targets = append(mt.Targets, struct {
		Instance string `json:"instance"`
		Status   string `json:"status"`
	}{
		Instance: instance,
		Status:   status,
	})
	return nil
}
