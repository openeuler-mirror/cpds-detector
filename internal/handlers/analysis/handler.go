/* 
 *  Copyright 2023 CPDS Author
 *  
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *  
 *       https://www.apache.org/licenses/LICENSE-2.0
 *  
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
