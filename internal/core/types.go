package core

type Analysis struct {
	ID         uint   `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	RuleID     uint   `json:"rule_id" gorm:"not null"`
	RuleName   string `json:"rule_name" gorm:"not null"`
	Status     string `json:"status" gorm:"not null"`
	Count      uint   `json:"count" gorm:"not null"`
	CreateTime int64  `json:"create_time" gorm:"not null"`
	UpdateTime int64  `json:"update_time" gorm:"not null"`
}

type Rule struct {
	ID                     uint    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name                   string  `json:"name" gorm:"unique;not null"`
	Expression             string  `json:"expression" gorm:"not null;type:varchar(512)"`
	SubhealthConditionType string  `json:"subhealth_condition_type"`
	SubhealthThresholds    float64 `json:"subhealth_thresholds"`
	FaultConditionType     string  `json:"fault_condition_type"`
	FaultThresholds        float64 `json:"fault_thresholds"`
	Severity               string  `json:"severity" gorm:"not null"`
	Duration               string  `json:"duration" gorm:"not null"`
	CreateTime             int64   `json:"create_time" gorm:"not null"`
	UpdateTime             int64   `json:"update_time" gorm:"not null"`
}
