package v1

import (
	"gitee.com/cpds/cpds-detector/pkgs/rules"
)

type Handler struct {
	rules rules.Interface
}

func newRulesHandler() *Handler {
	return &Handler{}
}
