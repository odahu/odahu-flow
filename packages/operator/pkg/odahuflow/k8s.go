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

package odahuflow

import (
	"fmt"
)

const (
	LastAppliedHashAnnotation = "operator.odahuflow.org/last-applied-hash"
	PackagerSetupStep         = "setup"
	PackagerPackageStep       = "packager"
	PackagerResultStep        = "result"
	TrainerSetupStep          = "setup"
	TrainerTrainStep          = "trainer"
	TrainerValidationStep     = "validation"
	TrainerResultStep         = "result"
)

func GeneratePackageResultCMName(mpID string) string {
	return fmt.Sprintf("%s-mp-result", mpID)
}

func GenerateTrainingResultCMName(mtID string) string {
	return fmt.Sprintf("%s-mt-result", mtID)
}

func GenerateDeploymentConnectionSecretName(connName string) string {
	return fmt.Sprintf("%s-regsecret", connName)
}
