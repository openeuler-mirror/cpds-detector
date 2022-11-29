package v1

import (
	"gitee.com/cpds/cpds-detector/pkgs/rules"
	"github.com/emicklei/go-restful"
)

func AddToContainer(container *restful.Container, r *rules.Rules) {
	webservice := new(restful.WebService)
	webservice.Path("/rules")

	container.Add(webservice)
}
