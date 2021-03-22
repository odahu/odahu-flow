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
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
"os"
"path"
	"path/filepath"
	"strings"

)

// SyncArchiveOrDir sync object storage directory or archive
// If localPath ends with "/" slash then we assume that it's a directory otherwise file
// If (conn.URI + relPath) ends with .tar.gz or .zip then we assume that it's compressed tarball (.zip is legacy).
// If local path is directory but remote path is archive then we not only fetch archive but unzip it to local path
// Please take a look at ObjectStorage.Download
func SyncArchiveOrDir(conn v1alpha1.ConnectionSpec, relPath string, localPath string) error {

	log := zap.S()


	storage, err := rclone.NewObjectStorage(&conn)
	if err != nil {
		log.Error(err, "repository creation")
		return err
	}

	fullModelPath := path.Join(storage.RemoteConfig.Path, relPath)

	remoteFileIsArchive := strings.HasSuffix(relPath,".zip") || strings.HasSuffix(relPath,".tar.gz")

	// Automatically unzip archive to target location directory
	if remoteFileIsArchive {

		file, err := ioutil.TempFile("", relPath)
		if err != nil {
			log.Error(err, "Unable to create temp model archive file to fetch into")
			return err
		}
		tempZip := file.Name()
		if err := storage.Download(tempZip, fullModelPath); err != nil {
			return err
		}

		// Ensure that target directory for utils.Unzip destination exists
		if err = os.MkdirAll(localPath, 0777); err != nil {
			log.Error(err, "output dir creation")
			return err
		}
		// Unzip archive
		if err := utils.Unzip(tempZip, localPath); err != nil {
			log.Error(err, "unzip training data")

			return err
		}

	} else {
		if err := storage.Download(localPath, fullModelPath); err != nil {
			return err
		}
	}

	return nil

}

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
// Model can be whether directory or .tar.gz / .zip archive
// If localPath == "" then temp directory will be created and returned in newLocalPath
// otherwise newLocalPath == localPath
func (mr ModelRegistry) SyncModel(connName string, relPath string, localPath string) (newLocalPath string, err error) {

	newLocalPath = localPath

	conn, err := mr.connGetter.GetConnection(connName)
	if err != nil {
		return newLocalPath, err
	}

	if err = conn.DecodeBase64Fields(); err != nil {
		return newLocalPath, err
	}

	if localPath == "" {
		tempDir, err := ioutil.TempDir("", "model-files")
		if err != nil {
			return newLocalPath, err
		}
		newLocalPath = tempDir
		defer func() {_ = os.Remove(tempDir)}()
	}

	if !strings.HasSuffix(newLocalPath, "/") {
		zap.S().Infof("Adding suffix '/' to %s to mark it as directory " +
			"(model can be synced only to directory)", newLocalPath)
		newLocalPath = newLocalPath + "/"
	}

	return newLocalPath, SyncArchiveOrDir(conn.Spec, relPath, newLocalPath)
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
	localPath, err := mr.SyncModel(connName, path, "")
	if err != nil {
		return "", "", err
	}

	return getModelNameVersion(localPath)

}