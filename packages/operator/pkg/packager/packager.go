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

package packager

import (
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	conn_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	pack_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"io/ioutil"
	"os"
	"path"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	resultFileName     = "result.json"
	modelPackagingFile = "mp.json"
)

type Packager struct {
	packagingClient       pack_api_client.Client
	connClient            conn_api_client.Client
	log                   logr.Logger
	modelPackagingID      string
	packagerConfig        config.PackagerConfig
}

func NewPackager(
	packagingClient pack_api_client.Client,
	connClient conn_api_client.Client,
	config config.PackagerConfig,
) *Packager {
	return &Packager{
		packagingClient: packagingClient,
		connClient: connClient,
		log: logf.Log.WithName("packager").WithValues(
			odahuflow.ModelPackagingIDLogPrefix, config.ModelPackagingID,
		),
		modelPackagingID: config.ModelPackagingID,
		packagerConfig:   config,
	}
}

// This function prepares a packaging environment. To do this, it performs the following steps:
//   1) It extracts the packaging entity from repository storage, for example, from the API server.
//   2) It downloads a trained artifact from object storage.
//   4) Finally, it saves the packaging entity to allow a packager to use it.
func (p *Packager) SetupPackager() (err error) {
	k8sPackaging, err := p.getPackaging()
	if err != nil {
		return err
	}

	if err := p.downloadData(k8sPackaging); err != nil {
		p.log.Error(err, "Downloading packaging data failed")
		return err
	}

	mtBytes, err := json.Marshal(k8sPackaging)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(modelPackagingFile, mtBytes, 0644) //nolint file reads in another docker image
}

// This function saves a packaging result. To do this, it performs the following steps:
//   1) It extracts the packaging entity from repository storage, for example, from the API server.
//   2) It reads the packaging result file from workspace.
//   4) Finally, it saves the packaging results.
func (p *Packager) SaveResult() error {
	k8sPackaging, err := p.getPackaging()
	if err != nil {
		return err
	}

	resultFile, err := os.Open(resultFileName)
	if err != nil {
		p.log.Error(err, "Open result file")
		return err
	}
	defer resultFile.Close()

	byteResult, err := ioutil.ReadAll(resultFile)
	if err != nil {
		p.log.Error(err, "Read from result file")
		return err
	}

	var result map[string]string
	if err := json.Unmarshal(byteResult, &result); err != nil {
		p.log.Error(err, "Unmarshal result")
		return err
	}

	packResult := make([]v1alpha1.ModelPackagingResult, 0, len(result))
	for name, value := range result {
		packResult = append(packResult, v1alpha1.ModelPackagingResult{
			Name:  name,
			Value: value,
		})
	}

	return p.packagingClient.SaveModelPackagingResult(k8sPackaging.ModelPackaging.ID, packResult)
}

func (p *Packager) downloadData(packaging *packaging.K8sPackager) (err error) {
	storage, err := rclone.NewObjectStorage(&packaging.ModelHolder.Spec)
	if err != nil {
		p.log.Error(err, "repository creation")

		return err
	}

	file, err := os.Create(packaging.TrainingZipName)
	if err != nil {
		p.log.Error(err, "zip creation")

		return err
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			p.log.Error(err, "Error during closing of the file")
		}

		if err == nil {
			err = closeErr
		}
	}()

	if err := storage.Download(
		packaging.TrainingZipName,
		path.Join(storage.RemoteConfig.Path, packaging.TrainingZipName),
	); err != nil {
		p.log.Error(err, "download training zip")

		return err
	}

	if err = os.Mkdir(p.packagerConfig.OutputDir, 0777); err != nil {
		p.log.Error(err, "output dir creation")

		return err
	}

	if err := utils.Unzip(packaging.TrainingZipName, p.packagerConfig.OutputDir); err != nil {
		p.log.Error(err, "unzip training data")

		return err
	}

	return nil
}
