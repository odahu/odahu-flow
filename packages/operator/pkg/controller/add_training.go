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

package controller

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/modeltraining"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	// AddToManagerTrainingFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerTrainingFuncs = append(AddToManagerTrainingFuncs, modeltraining.Add)
}

// AddToManagerTrainingFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerTrainingFuncs []func(
	manager.Manager, config.ModelTrainingConfig, config.OperatorConfig, config.CommonConfig, string,
) error

// AddToManager adds all Controllers to the Manager
func AddTrainingToManager(m manager.Manager, trainingConfig config.ModelTrainingConfig, operatorConfig config.OperatorConfig, commonConfig config.CommonConfig, gpuResourceName string) error {
	for _, f := range AddToManagerTrainingFuncs {
		if err := f(m, trainingConfig, operatorConfig, commonConfig, gpuResourceName); err != nil {
			return err
		}
	}
	return nil
}
