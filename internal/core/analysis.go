package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"cpds/cpds-detector/pkg/cpds-detector/config"
	"cpds/cpds-detector/pkg/prometheus"
	prometheusutils "cpds/cpds-detector/pkg/utils/prometheus"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	r *resource
)

type resource struct {
	prometheusConfig *prometheusConfig
	rulesChan        chan []Rule
	logger           *zap.Logger
	db               *gorm.DB
}

type prometheusConfig struct {
	host string
	port int
}

func InitAnalysis(config *config.Config, logger *zap.Logger, db *gorm.DB) error {
	r = &resource{
		prometheusConfig: &prometheusConfig{
			host: config.PrometheusOptions.Host,
			port: config.PrometheusOptions.Port,
		},
		rulesChan: make(chan []Rule),
		logger:    logger,
		db:        db.Session(&gorm.Session{}),
	}
	var rules []Rule

	go manageRules(r.rulesChan)

	query := db.Session(&gorm.Session{})
	if err := query.Model(&Rule{}).Find(&rules).Error; err != nil {
		return err
	}
	r.rulesChan <- rules

	return nil
}

func RuleUpdated() error {
	var rules []Rule
	db := r.db.Session(&gorm.Session{})
	if err := db.Find(&rules).Error; err != nil {
		return err
	}
	r.rulesChan <- rules

	return nil
}

func manageRules(rulesChan chan []Rule) {
	listendRules := make(map[Rule]chan bool)

	for {
		select {
		case rules := <-rulesChan:
			for oldRule, stopChan := range listendRules {
				if len(rules) == 0{
					delete(listendRules,oldRule)
					stopChan <- true
				}
				for index, newRule := range rules {
					if oldRule == newRule {
						break
					}
					if index == len(rules)-1 {
						delete(listendRules,oldRule)
						stopChan <- true
					}
				}
			}

			for _, newRule := range rules {
				if _, exist := listendRules[newRule]; !exist {
					stopChan := make(chan bool)
					listendRules[newRule] = stopChan
					// Save a copy of the pointer value in a local variable,
					// so that each goroutine can work with its own unique pointer.
					go listenRule(r.prometheusConfig, newRule, stopChan)
				}
			}
		}
	}
}

func listenRule(promConf *prometheusConfig, rule Rule, stopChan chan bool) {
	p, err := prometheus.NewPrometheus(promConf.host, promConf.port)
	if err != nil {
		r.logger.Warn(fmt.Sprintf("faild to listen rule: %s", rule.Name))
		return
	}

	if !prometheusutils.IsExprValid(rule.Expression) {
		r.logger.Warn(fmt.Sprintf("faild to listen rule %s, invalid rule expression: %s", rule.Name, rule.Expression))
		return
	}

	subHealthRule := fmt.Sprintf("%s %s %s",
		rule.Expression,
		rule.SubhealthConditionType,
		strconv.FormatFloat(rule.SubhealthThresholds, 'f', -1, 64),
	)
	faultRule := fmt.Sprintf("%s %s %s",
		rule.Expression,
		rule.FaultConditionType,
		strconv.FormatFloat(rule.FaultThresholds, 'f', -1, 64),
	)

	r.logger.Info(fmt.Sprintf("start to listen rule: %s", rule.Name))

	for {
		select {
		case <-stopChan:
			return
		default:
			subMetric := p.GetSingleMetric(subHealthRule, time.Now())
			faultMetric := p.GetSingleMetric(faultRule, time.Now())

			if !isEmptyMetric(&faultMetric) {
				insertFaultRecord(rule.ID, rule.Name, &faultMetric)
			} else if !isEmptyMetric(&subMetric) {
				insertSubhealthRecord(rule.ID, rule.Name, &subMetric)
			}

			time.Sleep(time.Second * 5)
		}
	}
}

func insertSubhealthRecord(id uint, name string, metric *prometheus.Metric) error {
	return insertRecord(id, name, "subhealth", metric)
}

func insertFaultRecord(id uint, name string, metric *prometheus.Metric) error {
	return insertRecord(id, name, "fault", metric)
}

func insertRecord(id uint, name, kind string, metric *prometheus.Metric) error {
	var newRecord, existingRecord *Analysis
	analysisTime := metric.MetricData.MetricValues[0].Sample[0]
	newRecord = &Analysis{
		RuleID:     id,
		RuleName:   name,
		Status:     kind,
		Count:      1,
		CreateTime: int64(analysisTime),
		UpdateTime: int64(analysisTime),
	}

	db := r.db.Session(&gorm.Session{}).Model(&Analysis{})
	result := db.Where("rule_id = ?", id).Order("update_time DESC").Limit(1).Find(&existingRecord)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if (existingRecord.ID == 0) ||
		(time.Unix(int64(analysisTime), 0).Sub(time.Unix(existingRecord.CreateTime, 0)) > time.Hour*24) { //
		if err := db.Create(newRecord).Error; err != nil {
			return err
		}
	} else {
		if kind == "fault" {
			existingRecord.Status = kind
		}
		result := db.Model(&existingRecord).Updates(Analysis{RuleName: name,Count: existingRecord.Count + 1, UpdateTime: int64(analysisTime)})
		if result.Error != nil {
			return result.Error
		}

	}

	return nil
}

func isEmptyMetric(m *prometheus.Metric) bool {
	return len(m.MetricValues) == 0
}
