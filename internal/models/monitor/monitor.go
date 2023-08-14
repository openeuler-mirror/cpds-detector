package monitor

import (
	"cpds/cpds-detector/pkg/prometheus"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	stringutil "cpds/cpds-detector/pkg/utils/string"

	jsoniter "github.com/json-iterator/go"
)

type Operator interface {
	GetMonitorTargets(instance string) (*MonitorTargets, error)

	GetNodeInfo(instance string) ([]NodeInfo, error)

	GetNodeStatus(instance string) ([]NodeStatus, error)

	GetNodeResources(instance string, startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error)

	GetNodeContainerStatus(instance string) ([]prometheus.Metric, error)

	GetClusterResource(startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error)

	GetClusterContainerStatus(startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error)
}

type operator struct {
	prometheusConfig *prometheusConfig
}

type prometheusConfig struct {
	host string
	port int
}

func NewOperator(prometheusHost string, prometheusPort int) Operator {
	return &operator{
		prometheusConfig: &prometheusConfig{
			host: prometheusHost,
			port: prometheusPort,
		},
	}
}

func (o *operator) GetMonitorTargets(instance string) (*MonitorTargets, error) {
	url := fmt.Sprintf("http://%s:%d/api/v1/targets?scrapePool=cpds", o.prometheusConfig.host, o.prometheusConfig.port)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr promResponse
	if err := jsoniter.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	var pr1 promResponse
	for _, target :=range pr.Data.Targets{
		if strings.Contains(target.DiscoveredLabels.Address,instance) {
			pr1.Data.Targets=append(pr1.Data.Targets, target)
		}
	}
	var mt MonitorTargets
	for _, target := range pr1.Data.Targets {
		mt.addTargets(target.DiscoveredLabels.Address, target.Health)
	}

	return &mt, nil
}

func (o *operator) GetNodeInfo(instance string) ([]NodeInfo, error) {
	expr := "cpds_node_basic_info"
	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	var n []NodeInfo
	metrics := client.GetSingleMetric(expr, time.Now())
	for _, v := range metrics.MetricData.MetricValues {
		if instance != "" && stringutil.ExtractIP(v.Metadata["instance"]) != instance {
			continue
		}
		n = append(n, NodeInfo{
			Instance:      stringutil.ExtractIP(v.Metadata["instance"]),
			Arch:          v.Metadata["arch"],
			OSVersion:     v.Metadata["os_version"],
			KernelVersion: v.Metadata["kernel_version"],
		})
	}

	return n, nil
}

func (o *operator) GetNodeStatus(instance string) ([]NodeStatus, error) {
	var n []NodeStatus
	var mtx sync.Mutex
	var wg sync.WaitGroup

	t, err := o.GetMonitorTargets(instance)
	if err != nil {
		return nil, err
	}

	for _, target := range t.Targets {
		wg.Add(1)
		go func(
			target struct {
				Instance string `json:"instance"`
				Status   string `json:"status"`
			},
			config *prometheusConfig,
		) {
			if target.Status != "up" {
				mtx.Lock()
				n = append(n, NodeStatus{
					Instance: stringutil.ExtractIP(target.Instance),
				})
				mtx.Unlock()
				wg.Done()
				return
			}

			exprMap := make(map[string]string)
			exprMap["container_total"] = fmt.Sprintf("sum (cpds_container_state{instance=~\"%s.*\"})", target.Instance)
			exprMap["container_running"] = fmt.Sprintf("sum (cpds_container_state{instance=~\"%s.*\",state=\"running\"})", target.Instance)
			exprMap["cpu_used_core"]=fmt.Sprintf("sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\", mode!=\"idle\",instance=~\"%s.*\"}[1m]))",target.Instance)
			exprMap["cpu_usage"] = fmt.Sprintf("1-(sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\", mode=\"idle\",instance=~\"%s.*\"}[1m]))/sum (rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}[1m])))", instance, instance)
			exprMap["cpu_total_core"]=fmt.Sprintf("sum (rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}[1m]))",target.Instance)
			exprMap["cpu_number_core"]=fmt.Sprintf("count(sum(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}) by (cpu))",target.Instance)
			exprMap["memory_usage"] = fmt.Sprintf("cpds_node_memory_usage_bytes{instance=~\"%s.*\"} / cpds_node_memory_total_bytes{instance=~\"%s.*\"}", target.Instance, target.Instance)
			exprMap["memory_used_bytes"] = fmt.Sprintf("cpds_node_memory_usage_bytes{instance=~\"%s.*\"}", target.Instance)
			exprMap["memory_total_bytes"] = fmt.Sprintf("cpds_node_memory_total_bytes{instance=~\"%s.*\"}", target.Instance)
			exprMap["disk_usage"] = fmt.Sprintf("sum(cpds_node_blk_total_bytes{instance=\"%s\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_usage_bytes{instance=\"%s\"})/(sum(cpds_node_blk_total_bytes{instance=\"%s\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_usage_bytes{instance=\"%s\"})+sum(cpds_node_blk_total_bytes{instance=\"%s\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_available_bytes{instance=\"%s\"}))", target.Instance, target.Instance, target.Instance, target.Instance, target.Instance, target.Instance)
			exprMap["disk_used_bytes"] = fmt.Sprintf("sum(cpds_node_blk_total_bytes{instance=\"%s\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_usage_bytes{instance=\"%s\"})", target.Instance, target.Instance)
			exprMap["disk_total_bytes"] = fmt.Sprintf("sum(cpds_node_blk_total_bytes{instance=\"%s\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_total_bytes{instance=\"%s\"})", target.Instance, target.Instance)

			client, _ := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
			metrics := client.GetMultiMetrics(exprMap, time.Now())
			ns := NodeStatus{
				Instance: target.Instance,
			}
			for _, metric := range metrics {
				switch metric.MetricName {
				case "container_total":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Container.Total = int(metric.MetricData.MetricValues[0].Sample[1])
					}else{
						ns.Container.Total = 0
					}
				case "container_running":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Container.Running = int(metric.MetricData.MetricValues[0].Sample[1])
					}else{
						ns.Container.Running = 0
					}
						
				case "cpu_usage":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Cpu.Usage = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Cpu.Usage = 0
					}
				case "cpu_used_core":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Cpu.UsedCore = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Cpu.UsedCore = 0
					}
				case "cpu_number_core":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Cpu.NumberCores = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Cpu.NumberCores = 0
					}
				case "cpu_total_core":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Cpu.TotalCore = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Cpu.TotalCore = 0
					}
				case "memory_usage":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Memory.Usage = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Memory.Usage = 0
					}
				case "memory_used_bytes":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Memory.UsedBytes = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Memory.UsedBytes = 0
					}
				case "memory_total_bytes":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Memory.TotalBytes = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Memory.TotalBytes = 0
					}
				case "disk_usage":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Disk.Usage = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Disk.Usage = 0
					}
				case "disk_used_bytes":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Disk.UsedBytes = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Disk.UsedBytes = 0
					}
				case "disk_total_bytes":
					if len(metric.MetricData.MetricValues) > 0 {
						ns.Disk.TotalBytes = metric.MetricData.MetricValues[0].Sample[1]
					}else{
						ns.Disk.TotalBytes = 0
					}
				}
			}
			mtx.Lock()
			n = append(n, ns)
			mtx.Unlock()

			wg.Done()
		}(target, o.prometheusConfig)
	}
	wg.Wait()

	return n, nil
}

func (o *operator) GetNodeResources(instance string, startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error) {
	exprMap := func(string) map[string]string {
		exprMap := make(map[string]string)
		exprMap["node_network_recive_iops"]=fmt.Sprintf("sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_receive_packets_total{instance=~\"%s.*\",interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))",instance)
		exprMap["node_network_transmit_iops"]=fmt.Sprintf("sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_transmit_packets_total{instance=~\"%s.*\",interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))",instance)
		exprMap["node_network_recive_bytes"]=fmt.Sprintf("sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_receive_bytes_total{instance=~\"%s.*\",interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))", instance)
		exprMap["node_network_transmit_bytes"]=fmt.Sprintf("sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_transmit_bytes_total{instance=~\"%s.*\",interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))", instance)
		exprMap["node_network_recive_drop_rate"] = fmt.Sprintf("sum(increase(cpds_node_network_receive_drop_total{instance=~\"%s.*\"}[1m])) / sum(increase(cpds_node_network_receive_packets_total{instance=~\"%s.*\"}[1m])) or vector(0)",instance,instance)
		exprMap["node_network_transmit_drop_rate"] = fmt.Sprintf("sum(increase(cpds_node_network_transmit_drop_total{instance=~\"%s.*\"}[1m])) / sum(increase(cpds_node_network_transmit_packets_total{instance=~\"%s.*\"}[1m])) or vector(0)",instance,instance)
		exprMap["node_cpu_usage"] = fmt.Sprintf("1-(sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\", mode=\"idle\",instance=~\"%s.*\"}[1m]))/sum (rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}[1m])))", instance, instance)
		exprMap["node_cpu_iowait"]=fmt.Sprintf("sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",mode=\"iowait\",instance=~\"%s.*\"}[1m]))/sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}[1m]))", instance, instance)
		exprMap["node_network_recive_error_rate"] = fmt.Sprintf("sum(increase(cpds_node_network_receive_errors_total{instance=~\"%s.*\"}[1m])) / sum(increase(cpds_node_network_receive_packets_total{instance=~\"%s.*\"}[1m])) or vector(0)",instance,instance)
		exprMap["node_network_transmit_error_rate"] = fmt.Sprintf("sum(increase(cpds_node_network_transmit_errors_total{instance=~\"%s.*\"}[1m])) / sum(increase(cpds_node_network_transmit_packets_total{instance=~\"%s.*\"}[1m])) or vector(0)",instance,instance)
		exprMap["node_memory_usage"] = fmt.Sprintf("cpds_node_memory_usage_bytes{instance=~\"%s.*\"} / cpds_node_memory_total_bytes{instance=~\"%s.*\"}", instance, instance)
		exprMap["node_retransm_rate"]=fmt.Sprintf("sum(increase(cpds_node_netstat_tcp_retrans_segs{instance=~\"%s.*\"}[1m])) / sum(increase(cpds_node_netstat_tcp_out_segs{instance=~\"%s.*\"}[1m])) or vector(0)",instance,instance)
		exprMap["node_disk_usage"] = fmt.Sprintf("sum(cpds_node_blk_total_bytes{instance=~\"%s.*\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_usage_bytes{instance=~\"%s.*\"})/(sum(cpds_node_blk_total_bytes{instance=~\"%s.*\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_usage_bytes{instance=~\"%s.*\"})+sum(cpds_node_blk_total_bytes{instance=~\"%s.*\",mount=~\".+\"}*0 + on(mount) group_right cpds_node_fs_available_bytes{instance=~\"%s.*\"}))", instance, instance, instance, instance, instance, instance)
		exprMap["node_disk_written_bytes"] = fmt.Sprintf("sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\",instance=~\"%s.*\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device) group_right rate(cpds_node_disk_written_bytes_total{instance=~\"%s.*\"}[1m]))", instance, instance)
		exprMap["node_disk_read_bytes"] = fmt.Sprintf("sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\",instance=~\"%s.*\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device) group_right rate(cpds_node_disk_read_bytes_total{instance=~\"%s.*\"}[1m]))", instance, instance)
		exprMap["node_disk_written_complete"] = fmt.Sprintf("sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\",instance=~\"%s.*\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device) group_right rate(cpds_node_disk_writes_completed_total{instance=~\"%s.*\"}[1m]))", instance, instance)
		exprMap["node_disk_read_complete"] = fmt.Sprintf("sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\",instance=~\"%s.*\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device) group_right rate(cpds_node_disk_reads_completed_total{instance=~\"%s.*\"}[1m]))", instance, instance)
		// TODO: get network recive/transmit rate, get network drop rate, get network error rate

		return exprMap
	}(instance)

	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metrics := client.GetMultiMetricsOverTime(exprMap, startTime, endTime, step)
	return metrics, nil
}

func (o *operator) GetNodeContainerStatus(instance string) ([]prometheus.Metric, error) {
	exprMap := func(string) map[string]string {
		exprMap := make(map[string]string)
		exprMap["node_container_status"] = fmt.Sprintf("cpds_container_state{instance=~\"%s.*\"}", instance)
		exprMap["node_container_cpu_usage"] = fmt.Sprintf("((sum(rate(cpds_container_cpu_usage_seconds_total{instance=~\"%s.*\"}[1m])) by (container)) /scalar(sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"} [1m]))))*scalar(count(sum(cpds_node_cpu_seconds_total{cpu!=\"cpu\",instance=~\"%s.*\"}) by (cpu)))",instance,instance,instance)
		exprMap["node_container_memory_used"] = fmt.Sprintf("cpds_container_memory_usage_bytes{instance=~\"%s.*\"}", instance)
		exprMap["node_container_inbound_traffic"] = fmt.Sprintf("sum(cpds_container_network_receive_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\",instance=~\"%s.*\",network_mode!~\"host|container.*|none\"}) by (container)", instance)
		exprMap["node_container_outbound_traffic"] = fmt.Sprintf("sum(cpds_container_network_transmit_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\",instance=~\"%s.*\",network_mode!~\"host|container.*|none\"}) by (container)", instance)
		return exprMap
	}(instance)

	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metrics := client.GetMultiMetrics(exprMap, time.Now())
	return metrics, nil
}

func (o *operator) GetClusterResource(startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error) {
	exprMap := func() map[string]string {
		exprMap := make(map[string]string)
		exprMap["cluster_cpu_usage"] = "1-(sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\", mode=\"idle\"}[1m]))/sum (rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\"}[1m])))"
		exprMap["cluster_memory_usage"] = "sum(cpds_node_memory_usage_bytes)/scalar(sum(cpds_node_memory_total_bytes))"
		exprMap["cluster_network_recive_bytes"] = "sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_receive_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))"
		exprMap["cluster_network_transmit_bytes"] = "sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_transmit_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))"
		exprMap["cluster_disk_usage"] = "sum(cpds_node_blk_total_bytes{mount=~\".+\"}*0 + on(mount,instance) group_right cpds_node_fs_usage_bytes)/(sum(cpds_node_blk_total_bytes{mount=~\".+\"}*0 + on(mount,instance) group_right cpds_node_fs_usage_bytes)+sum(cpds_node_blk_total_bytes{mount=~\".+\"}*0 + on(mount,instance) group_right cpds_node_fs_available_bytes))"
		exprMap["cluster_disk_written_complete"] = "sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device,instance) group_right rate(cpds_node_disk_writes_completed_total[1m]))"
		exprMap["cluster_disk_read_complete"] = "sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device,instance) group_right rate(cpds_node_disk_reads_completed_total[1m]))"
		exprMap["cluster_disk_written_bytes"] = "sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device,instance) group_right rate(cpds_node_disk_written_bytes_total[1m]))"
		exprMap["cluster_cpu_iowait"] = "sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\",mode=\"iowait\"}[1m]))/sum(rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\"}[1m]))"
		exprMap["cluster_network_recive_drop_rate"] = "sum(increase(cpds_node_network_receive_drop_total[1m])) / sum(increase(cpds_node_network_receive_packets_total[1m]))"
		exprMap["cluster_network_recive_error_rate"] = "sum(increase(cpds_node_network_receive_errors_total[1m])) / sum(increase(cpds_node_network_receive_packets_total[1m]))"
		exprMap["cluster_network_transmit_drop_rate"] = "sum(increase(cpds_node_network_transmit_drop_total[1m])) / sum(increase(cpds_node_network_transmit_packets_total[1m]))"
		exprMap["cluster_network_transmit_error_rate"] = "sum(increase(cpds_node_network_transmit_errors_total[1m])) / sum(increase(cpds_node_network_transmit_packets_total[1m]))"
		exprMap["cluster_retransm_rate"]="sum(increase(cpds_node_netstat_tcp_retrans_segs[1m])) / sum(increase(cpds_node_netstat_tcp_out_segs[1m]))"
		exprMap["cluster_network_receive_iops"]="sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_receive_packets_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))"
		exprMap["cluster_network_transmit_iops"]="sum(cpds_node_network_info{mask=~\".+\"}*0+on(interface,instance) group_right rate(cpds_node_network_transmit_packets_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*|tunl.*\"}[1m]))"
		exprMap["cluster_disk_read_bytes"] = "sum((label_replace(cpds_node_blk_total_bytes{type=\"disk\"}, \"device\", \"$1\", \"name\", \"(.+)\")*0)+ on(device,instance) group_right rate(cpds_node_disk_read_bytes_total[1m]))"
		return exprMap
	}()

	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metrics := client.GetMultiMetricsOverTime(exprMap, startTime, endTime, step)
	return metrics, nil
}

func (o *operator) GetClusterContainerStatus(startTime time.Time, endTime time.Time, step time.Duration) ([]prometheus.Metric, error) {
	exprMap := func() map[string]string {
		exprMap := make(map[string]string)
		exprMap["cluster_container_cpu_usage"] = "sum(rate(cpds_container_cpu_usage_seconds_total[1m]))/sum (rate(cpds_node_cpu_seconds_total{cpu!=\"cpu\"}[1m]))"
		exprMap["cluster_container_memory_usage"] = "sum(cpds_container_memory_usage_bytes) / sum(cpds_node_memory_total_bytes)"
		exprMap["cluster_container_recive_bytes"] = "sum(cpds_container_network_receive_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*\",network_mode!~\"host|container.*|none\"})"
		exprMap["cluster_container_write_bytes"] = "sum(cpds_container_network_transmit_bytes_total{interface!~\"lo|bond[0-9]|cbr[0-9]|veth.*|vir.*|docker.*|vnet.*|br.*|tap.*\",network_mode!~\"host|container.*|none\"})"
		exprMap["cluster_container_disk_usage"] = "sum(cpds_container_disk_usage_bytes)/((sum(cpds_node_fs_usage_bytes{fs=~\"/.*\"})+sum(cpds_node_fs_available_bytes{fs=~\"/.*\"})))"
		return exprMap
	}()

	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metrics := client.GetMultiMetricsOverTime(exprMap, startTime, endTime, step)
	return metrics, nil
}
