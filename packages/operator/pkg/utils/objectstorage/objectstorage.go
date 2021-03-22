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

package objectstorage

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// SyncArchiveOrDir sync object storage directory or archive
// If conn.URI + relPath refers to .zip / .tar.gz file then it will be unzipped to localPath
// If conn.URI + relPath refers to directory then it will be synced to localPath (localPath must ends with "/")
// If conn.URI + relPath refers to file then it will be copied to localPath (localPath must not ends with "/")
func SyncArchiveOrDir(conn v1alpha1.ConnectionSpec, relPath string, localPath string) error {

	log := zap.S()

	storage, err := rclone.NewObjectStorage(&conn)
	if err != nil {
		log.Error(err, "repository creation")
		return err
	}

	fullModelPath := path.Join(storage.RemoteConfig.Path, relPath)

	if err = os.Mkdir(localPath, 0777); err != nil {
		log.Error(err, "output dir creation")

		return err
	}

	if strings.HasSuffix(relPath,".zip") || strings.HasSuffix(relPath,".tar.gz") {
		file, err := ioutil.TempFile("", relPath)
		if err != nil {
			log.Error(err, "Unable to create temp model archive file to fetch into")
			return err
		}
		tempZip := file.Name()
		if err := storage.Download(tempZip, fullModelPath); err != nil {
			return err
		}

		if err := utils.Unzip(tempZip, localPath); err != nil {
			log.Error(err, "unzip training data")

			return err
		}

	} else {
		log.Info("We suppose that %s is not archive")
		if err := storage.Download(localPath, fullModelPath); err != nil {
			return err
		}
	}

	return nil

}