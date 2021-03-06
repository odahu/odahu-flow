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
	_ "github.com/rclone/rclone/backend/azureblob" // s3 specific handlers
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/rc"
)


func createAzureBlobConfig(configName string, conn *v1alpha1.ConnectionSpec) (*FileDescription, error) {
	_, err := fs.Find("azureblob")
	if err != nil {
		log.Error(err, "")
		return nil, err
	}

	if err := config.CreateRemote(configName, "azureblob", rc.Params{
		"sas_url": conn.KeySecret,
	}, true, false); err != nil {
		return nil, err
	}

	bucketName, pathInsideBucket, err := GetBucketAndPath(conn)
	if err != nil {
		log.Error(err, "Parsing data binding URI", "connection uri", conn.URI)
		return nil, err
	}
	return &FileDescription{
		FsName: fmt.Sprintf("%s:%s", configName, bucketName),
		Path:   pathInsideBucket,
	}, nil
}
