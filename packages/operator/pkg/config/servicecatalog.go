//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package config

import "time"

type ServiceCatalog struct {
	Auth AuthConfig `json:"auth"`
	BaseURL string `json:"baseUrl"`
	FetchTimeout time.Duration `json:"fetchTimeout"`
	HandleTimeout time.Duration `json:"handleTimeout"`
	// ServiceCatalog uses EdgeURL to call MLServer by adding ModelRoute prefix to EdgeURL path
	EdgeURL string `json:"edgeURL"`
	// ServiceCatalog set EdgeHost as Host header in requests to ML servers
	EdgeHost string `json:"edgeHost"`
}

func NewDefaultServiceCatalogConfig() ServiceCatalog {
	return ServiceCatalog{
		FetchTimeout: 2 * time.Second,
		HandleTimeout: 2 * time.Second,
	}
}
