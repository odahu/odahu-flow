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

package trainer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	conn_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	train_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
)

const (
	modelTrainingFile                     = "mt.json"
	unsupportedConnectionTypeErrorMessage = "unexpected connection type: %s. Supported types: git, s3, gcs, azure blob"
	odahuflowProjectFile                  = "odahuflow.project.yaml"
)

type ModelTrainer struct {
	trainClient     train_api_client.Client
	connClient      conn_api_client.Client
	modelTrainingID string
	log             logr.Logger
	trainerConfig   config.TrainerConfig
}

func NewModelTrainer(
	trainAPIClient train_api_client.Client,
	connAPIClient conn_api_client.Client,
	trainerConfig config.TrainerConfig) *ModelTrainer {

	trainingLogger := logf.Log.WithName("trainer").
		WithValues(odahuflow.ModelTrainingIDLogPrefix, trainerConfig.ModelTrainingID)

	return &ModelTrainer{
		trainClient:     trainAPIClient,
		connClient:      connAPIClient,
		modelTrainingID: trainerConfig.ModelTrainingID,
		log:             trainingLogger,
		trainerConfig:   trainerConfig,
	}
}

// This function prepares a training environment. To do this, it performs the following steps:
//   1) It extracts the training entity from repository storage, for example, from the API server.
//   2) It downloads the source code of model training.
//   3) The setup function downloads all training data.
//   4) Finally, it saves the training entity to allow an ML toolchain to use it.
func (mt *ModelTrainer) Setup() (err error) {
	// Extracts the training entity
	k8sTraining, err := mt.getTraining()
	if err != nil {
		mt.log.Error(err, "Can not construct the training entity")

		return err
	}

	mt.log.Info("The training entity was constructed successfully")

	commitID := ""

	workDir := mt.trainerConfig.OutputDir
	if len(workDir) != 0 {
		mt.log.Info("Change current working dir", "new worker dir", workDir)

		if err := os.Chdir(workDir); err != nil {
			mt.log.Error(err, "Changing current working dir failed",
				"new worker dir", workDir,
			)

			return err
		}
	}

	connType := k8sTraining.AlgorithmSourceConnection.Conn.Spec.Type

	// Downloads a source code
	switch {
	case connType == connection.GITType:
		commitID, err = mt.cloneUserRepo(k8sTraining, mt.trainerConfig.OutputDir)
		if err != nil {
			mt.log.Error(err, "Error occurs during cloning project")

			return err
		}

		// Saves some data before starting a training
		if err := mt.trainClient.SaveModelTrainingResult(
			k8sTraining.ModelTraining.ID,
			&odahuflowv1alpha1.TrainingResult{
				CommitID: commitID,
			},
		); err != nil {
			mt.log.Error(err, "Cannot save the commit id")

			return err
		}

		mt.log.Info("The commit ID was saved", "commit_id", commitID)
	case connection.ObjectStorageTypesSet[connType]:
		if err := mt.downloadAlgorithm(k8sTraining); err != nil {
			mt.log.Error(err, "Downloading algorithm failed")

			return err
		}
	default:
		return errors.New(fmt.Sprintf(unsupportedConnectionTypeErrorMessage, k8sTraining.AlgorithmSourceConnection.Conn.Spec.Type))
	}

	mt.log.Info("The model source code was downloaded", "dir", workDir)

	if err := mt.downloadData(k8sTraining); err != nil {
		mt.log.Error(err, "Downloading training data failed")

		return err
	}

	mt.log.Info("The training data was downloaded")

	// TODO: We make available all connections to a toolchain script. Do we need it?
	mtBytes, err := json.Marshal(k8sTraining)
	if err != nil {
		mt.log.Error(err, "Marshaling of the training entity to JSON format failed.")

		return err
	}

	return ioutil.WriteFile(modelTrainingFile, mtBytes, 0644) //nolint file reads in another docker image
}

type trainingDescription struct {
	Output map[string]string `yaml:"output"`
}

// This function saves a training result. To do this, it performs the following steps:
//   1) It extracts the training entity from repository storage, for example, from the API server.
//   2) It creates the training zip archive.
//   2) It uploads the training zip archive to object storage.
//   4) Finally, it saves the training entity results.
func (mt *ModelTrainer) SaveResult() error {
	k8sTraining, err := mt.getTraining()
	if err != nil {
		return err
	}

	mt.log.Info("The training entity was constructed successfully")

	outputZipName, err := odahuflow.ProduceTrainingZipName(
		k8sTraining.ModelTraining.Spec.Model.ArtifactNameTemplate,
		&odahuflow.TrainingZipNameConfig{
			Name:    k8sTraining.ModelTraining.Spec.Model.Name,
			Version: k8sTraining.ModelTraining.Spec.Model.Version,
		},
	)
	if err != nil {
		return err
	}

	outputTrainingDir := mt.trainerConfig.OutputDir
	mt.log.Info("Run to zip the dir", "dir", outputTrainingDir, "archive_name",
		outputZipName)

	jsonFile, err := os.Open(filepath.Join(outputTrainingDir, odahuflowProjectFile))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var trainingDesc trainingDescription
	err = yaml.Unmarshal(byteValue, &trainingDesc)
	if err != nil {
		return err
	}

	err = utils.ZipDir(outputTrainingDir, outputZipName)
	if err != nil {
		mt.log.Info("Zipping the dir failed", "dir", outputTrainingDir, "archive name",
			outputZipName)
		return err
	}

	mt.log.Info("Run to zip the dir", "dir", outputTrainingDir, "archive_name",
		outputZipName)

	storage, err := rclone.NewObjectStorage(&k8sTraining.OutputConn.Spec)
	if err != nil {
		return err
	}

	if err := storage.Upload(outputZipName, path.Join(storage.RemoteConfig.Path, outputZipName)); err != nil {
		return err
	}

	if err := mt.trainClient.SaveModelTrainingResult(
		k8sTraining.ModelTraining.ID,
		&odahuflowv1alpha1.TrainingResult{
			RunID:        trainingDesc.Output["run_id"],
			ArtifactName: outputZipName,
		},
	); err != nil {
		mt.log.Error(err, "Cannot save the training result")
	}

	return nil
}

func (mt *ModelTrainer) downloadData(k8sTraining *training.K8sTrainer) error {
	if len(k8sTraining.InputData) == 0 {
		mt.log.Info("Model k8sTraining data is empty. Skip downloading")

		return nil
	}
	for _, mtData := range k8sTraining.InputData {
		mt.log.Info("Run download k8sTraining data",
			"remote_path", mtData.RemotePath,
			"local_path", mtData.LocalPath,
			"connection_type", mtData.DataBinding.Type,
			"connection_uri", mtData.DataBinding.URI,
		)

		storage, err := rclone.NewObjectStorage(&mtData.DataBinding)
		if err != nil {
			return err
		}

		if err := storage.Download(mtData.LocalPath, mtData.RemotePath); err != nil {
			return err
		}
	}

	return nil
}

func (mt *ModelTrainer) downloadAlgorithm(k8sTraining *training.K8sTrainer) error {
	mt.log.Info("Run download k8sTraining algorithm",
		"remote_path", k8sTraining.AlgorithmSourceConnection.Path,
		"connection_type", k8sTraining.AlgorithmSourceConnection.Conn.Spec.Type,
		"connection_uri", k8sTraining.AlgorithmSourceConnection.Conn.Spec.URI,
	)

	storage, err := rclone.NewObjectStorage(&k8sTraining.AlgorithmSourceConnection.Conn.Spec)
	if err != nil {
		return err
	}

	localDir := mt.trainerConfig.OutputDir

	if !strings.HasSuffix(localDir, "/") {
		localDir += "/"
	}

	if err := storage.Download(localDir, k8sTraining.AlgorithmSourceConnection.Path); err != nil {
		return err
	}

	return nil
}
