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

package trainer

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
)

// The function extracts data from a repository and creates the training entity.
func (mt *ModelTrainer) getTraining() (*training.K8sTrainer, error) {
	modelTraining, err := mt.trainClient.GetModelTraining(mt.modelTrainingID)
	if err != nil {
		return nil, err
	}

	vcs, err := mt.connClient.GetConnection(modelTraining.Spec.VCSName)
	if err != nil {
		return nil, err
	}

	// Since connRepo here is actually an HTTP client, it returns connection with base64-encoded secrets
	if err := vcs.DecodeBase64Fields(); err != nil {
		return nil, err
	}

	inputData := make([]training.InputDataBindingDir, 0, len(modelTraining.Spec.Data))
	for _, trainData := range modelTraining.Spec.Data {
		var trainDataConnSpec odahuflowv1alpha1.ConnectionSpec

		trainDataConn, err := mt.connClient.GetConnection(trainData.Connection)
		if err != nil {
			mt.log.Error(err, "Get train data", odahuflow.ConnectionIDLogPrefix, trainData.Connection)

			return nil, err
		}
		// Since connRepo here is actually an HTTP client, it returns connection with base64-encoded secrets
		if err := trainDataConn.DecodeBase64Fields(); err != nil {
			return nil, err
		}

		trainDataConnSpec = trainDataConn.Spec

		inputData = append(inputData, training.InputDataBindingDir{
			LocalPath:   trainData.LocalPath,
			RemotePath:  trainData.RemotePath,
			DataBinding: trainDataConnSpec,
		})
	}

	outputConn, err := mt.connClient.GetConnection(modelTraining.Spec.OutputConnection)
	if err != nil {
		return nil, err
	}
	// Since connRepo here is actually an HTTP client, it returns connection with base64-encoded secrets
	if err := outputConn.DecodeBase64Fields(); err != nil {
		return nil, err
	}

	ti, err := mt.trainClient.GetToolchainIntegration(modelTraining.Spec.Toolchain)
	if err != nil {
		return nil, err
	}

	return &training.K8sTrainer{
		VCS:        vcs,
		InputData:  inputData,
		OutputConn: outputConn,
		ModelTraining: &training.ModelTraining{
			ID:   modelTraining.ID,
			Spec: modelTraining.Spec,
		},
		ToolchainIntegration: ti,
	}, nil
}
