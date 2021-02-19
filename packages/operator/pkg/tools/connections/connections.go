package connections

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"strings"
)

const (
	sep = ":"
	defaultStorageNamePrefix = "odahu-"
)

type ConnectionClient interface {
	GetConnection(id string) (*connection.Connection, error)
}

func getConnNameAndRCloneStorageName (s string) (connName string, storageName string) {
	res := strings.Split(s, sep)
	connName = res[0]
	if len(res) > 1 {
		storageName = res[1]
	} else {
		storageName = defaultStorageNamePrefix + connName
	}

	return
}

func GenerateRClone(clusterValues []string, fileValues []string, connClient ConnectionClient ) (err error) {

	for _, value := range clusterValues {
		connName, storageName := getConnNameAndRCloneStorageName(value)
		conn, fetchErr := connClient.GetConnection(connName)

		if fetchErr != nil {
			err = multierr.Append(err, fetchErr)
			continue
		}

		if decodeErr := conn.DecodeBase64Fields(); decodeErr != nil {
			zap.S().Warnw(
				"There are some problems with decoding base64 fields. Maybe they are already decoded",
				zap.Error(decodeErr),
			)
		}

		_, genErr := rclone.NewObjectStorageWithName(storageName, &conn.Spec)
		if genErr != nil {
			err = multierr.Append(err, genErr)
			continue
		}
	}


	return err

}
