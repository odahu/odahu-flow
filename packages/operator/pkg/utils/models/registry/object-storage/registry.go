/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

// Package object_storage provides API to fetch model files and metadata from s3/azureblob/gcs
// with minimum requirements to model files.
// Generally speaking object storage registry is not true model registry but just a way to sync
// models from object storage with minimum ODAHU conventions about files
// We expect that model files are stored whether in a directory in bucket or as .tar.gz / zip archive
// We expect that this directory/archive contain odahuflow.project.yaml file in root that describes next
// metadata about model:
// --
// model:
//   name: wine
//	 version: 1.2
// --
package object_storage

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/objectstorage"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConnGetter interface {
	GetConnection(id string) (*connection.Connection, error)
}

type ModelRegistry struct {
	connGetter ConnGetter
}

func NewModelRegistry(cg ConnGetter) *ModelRegistry {
	return &ModelRegistry{connGetter: cg}
}


// SyncModel fetch model files from registry to localPath. If localPath is nil
// then SyncModel set localPath to temp directory and fetch model files there
func (mr ModelRegistry) SyncModel(connName string, path string, localPath *string) error {

	conn, err := mr.connGetter.GetConnection(connName)
	if err != nil {
		return err
	}

	if localPath == nil {
		tempDir, err := ioutil.TempDir("", "model-files")
		if err != nil {
			return err
		}
		localPath = &tempDir
		defer func() {_ = os.Remove(*localPath)}()
	}
	return objectstorage.SyncArchiveOrDir(conn.Spec, path, *localPath)
}

func getModelNameVersion(modelDir string) (string, string, error) {

	type Model struct {
		Name string `json:"name"`
		Version string `json:"version"`
	}

	type ModelFile struct {
		Model Model `json:"model"`
	}

	modelFiles := []string{"odahuflow.project.yaml", "odahuflow.project.yml"}

	items, err := ioutil.ReadDir(modelDir)
	if err != nil {
		return "", "", err

	}
	for _, item := range items {
		file := item.Name()
		for _, mFile := range modelFiles {
			if mFile == file {

				fp := filepath.Join(modelDir, file)
				zap.S().Infof("Model metadata file is found %s", fp)
				data, err := ioutil.ReadFile(fp)
				if err != nil {
					return "", "", err
				}
				mf := ModelFile{}
				if err := yaml.Unmarshal(data, &mf); err != nil {
					return "", "", err
				}

				return mf.Model.Name, mf.Model.Version, nil
			}
		}
	}
	return "", "", fmt.Errorf("unable to find model metadata file")
}

// Meta fetch model name and version using model path
func (mr ModelRegistry) Meta(connName string, path string) (name string, version string, err error) {
	var localPath string
	if err = mr.SyncModel(connName, path, &localPath); err != nil {
		return
	}

	return getModelNameVersion(localPath)

}