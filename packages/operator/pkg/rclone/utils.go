package rclone

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

// GetBucketAndPath return Connection bucket name and path
func GetBucketAndPath(c *v1alpha1.ConnectionSpec) (string, string, error) {
	if c.Type == connection.S3Type || c.Type == connection.GcsType {
		parsedURI, err := url.Parse(c.URI)
		if err != nil {
			return "", "", fmt.Errorf("unable to parse conn URI: %s", err)
		}
		return parsedURI.Host, parsedURI.Path, nil
	} else if c.Type == connection.AzureBlobType {
		parsedURI, err := url.Parse(c.URI)
		if err != nil {
			return "", "", fmt.Errorf("unable to parse conn URI: %s", err)
		}

		uriPath := parsedURI.Path
		uriPath= strings.TrimPrefix(uriPath, "/")

		pathParts := strings.Split(uriPath, "/")
		if len(pathParts) == 0 {
			return "", "", errors.New("azure URI must contain at least a bucket name")
		}
		bucketName := pathParts[0]
		pathInsideBucket := "/" + strings.Join(pathParts[1:], "/")
		return bucketName, pathInsideBucket, nil
	} else {
		return "", "", fmt.Errorf("not available for connection type: %s", c.Type)
	}
}
