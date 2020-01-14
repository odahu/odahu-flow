/*
 * Copyright 2019 EPAM Systems
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

package api

import (
	"github.com/spf13/viper"
	"path/filepath"
)

const (
	// Type of the backend. Available values:
	//    * local
	//    * config
	BackendType = "api.backend.type"
	// Path to a dir with Legion CRDs
	LocalBackendCRDPath = "api.backend.local.crd_path"
	// WEB Server port
	Port = "api.port"
)

const (
	LocalBackendType  = "local"
	ConfigBackendType = "config"
)

func init() {
	viper.SetDefault(BackendType, ConfigBackendType)
	viper.SetDefault(LocalBackendCRDPath, filepath.Join("config", "crds"))

	viper.SetDefault(Port, 5000)
}
