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
	_ "github.com/rclone/rclone/backend/s3" // s3 specific handlers
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
)

func createS3config(configName string, conn *v1alpha1.ConnectionSpec) (*FileDescription, error) {
	_, err := fs.Find("s3")
	if err != nil {
		log.Error(err, "")
		return nil, err
	}

	if err := config.CreateRemote(configName, "s3", map[string]interface{}{
		fs.ConfigProvider:   "AWS",
		"env_auth":          false,
		"region":            conn.Region,
		"bucket_acl":        "private",
		"access_key_id":     conn.KeyID,
		"secret_access_key": conn.KeySecret,
		// https://github.com/rclone/rclone/issues/1824
		// Workaround can be replaced after rclone v1.54 release
		// with option server_side_encryption: aws:kms
		"upload_cutoff": 0,
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
