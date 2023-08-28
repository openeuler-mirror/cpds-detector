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
