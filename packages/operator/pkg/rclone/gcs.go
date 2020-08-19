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

package rclone

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	_ "github.com/rclone/rclone/backend/googlecloudstorage" // s3 specific handlers
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"io/ioutil"
	"net/url"
)

func createGcsConfig(configName string, conn *v1alpha1.ConnectionSpec) (*FileDescription, error) {
	_, err := fs.Find("googlecloudstorage")
	if err != nil {
		log.Error(err, "")
		return nil, err
	}

	options := map[string]interface{}{
		"env_auth":   true,
		"location":   conn.Region,
		"bucket_acl": "private",
	}

	if len(conn.KeySecret) != 0 {
		options["service_account_credentials"] = conn.KeySecret
	}

	// TODO: remove debug prints
	dat, err := ioutil.ReadFile(config.ConfigPath)
	log.Info("RClone config content", "config", string(dat))
	if err := config.CreateRemote(configName, "googlecloudstorage", options, true, false); err != nil {
		return nil, err
	}

	parsedURI, err := url.Parse(conn.URI)
	if err != nil {
		log.Error(err, "Parsing data binding URI", "connection uri", conn.URI)

		return nil, err
	}

	return &FileDescription{
		FsName: fmt.Sprintf("%s:%s", configName, parsedURI.Host),
		Path:   parsedURI.Path,
	}, nil
}
