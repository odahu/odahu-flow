/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"path/filepath"
)

type APILocalBackendConfig struct {
	// Path to a dir with ODAHU CRDs
	LocalBackendCRDPath string `json:"localBackendCrdPath"`
}

type APIBackendConfig struct {
	// Type of the backend. Available values:
	//    * local
	//    * config
	Type string `json:"type"`
	// Local backend
	Local APILocalBackendConfig `json:"local"`
}

type APIConfig struct {
	Backend APIBackendConfig `json:"backend"`
	// API HTTP port
	Port int `json:"port"`
	// If true then only webserver will be setup.
	// Without background workers responsible to monitor storage and call services
	DisableWorkers bool `json:"disableWorkers"`
}

const (
	LocalBackendType  = "local"
	ConfigBackendType = "config"
)

func NewDefaultAPIConfig() APIConfig {
	return APIConfig{
		DisableWorkers: false,
		Backend: APIBackendConfig{
			Type: ConfigBackendType,
			Local: APILocalBackendConfig{
				LocalBackendCRDPath: filepath.Join("config", "crds"),
			},
		},
		Port: 5000,
	}
}
