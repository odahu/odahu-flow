package vault

import (
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"path"
	"time"

	bank_vaults "github.com/banzaicloud/bank-vaults/pkg/sdk/vault"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahuflow_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
)

const (
	ConnKey          = "odahuflow_conn"
	ConnectionIdsKey = "keys"
)

var (
	MaxSize   = 500
	FirstPage = 0
)

// Could be used for storing sensitive data
type vaultConnRepository struct {
	vaultClient         *vaultapi.Client
	connectionVaultPath string
}

func NewRepository(
	vaultClient *vaultapi.Client,
	connectionVaultPath string,
) conn_repository.Repository {
	return &vaultConnRepository{
		vaultClient:         vaultClient,
		connectionVaultPath: connectionVaultPath,
	}
}

func NewRepositoryFromConfig(vaultConfig config.Vault) (conn_repository.Repository, error) {
	vConfig := vaultapi.DefaultConfig()
	vConfig.Address = vaultConfig.URL

	vClient, err := vaultapi.NewClient(vConfig)
	if err != nil {
		return nil, err
	}
	vClient.SetToken(vaultConfig.Token)
	// This client provides authentication using k8s api
	_, err = bank_vaults.NewClientFromRawClient(
		vClient,
		bank_vaults.ClientRole(vaultConfig.Role),
	)

	return NewRepository(
		vClient,
		vaultConfig.SecretEnginePath,
	), err
}

func (vcr *vaultConnRepository) GetConnection(connID string) (*connection.Connection, error) {
	conn, err := vcr.getConnectionFromVault(connID)
	if err != nil {
		return nil, err
	}

	return conn.DeleteSensitiveData(), nil
}

func (vcr *vaultConnRepository) GetDecryptedConnection(connID string) (*connection.Connection, error) {
	return vcr.getConnectionFromVault(connID)
}

func (vcr *vaultConnRepository) GetConnectionList(
	options ...conn_repository.ListOption,
) ([]connection.Connection, error) {
	list, err := vcr.vaultClient.Logical().List(vcr.connectionVaultPath)
	if err != nil {
		return nil, convertVaultErrToOdahuflowErr(err)
	}

	if list == nil {
		return []connection.Connection{}, nil
	}

	connectionIdsRaw, ok := list.Data[ConnectionIdsKey]
	if !ok {
		return nil, odahuflow_errors.SerializationError{}
	}

	connectionIds, ok := connectionIdsRaw.([]interface{})
	if !ok {
		return nil, odahuflow_errors.SerializationError{}
	}

	listOptions := &conn_repository.ListOptions{
		Filter: &conn_repository.Filter{},
		Page:   &FirstPage,
		Size:   &MaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	connResults := []connection.Connection{}
	startPosition := *listOptions.Page * (*listOptions.Size)

	// TODO: think about more effective way to extract list of connections from vault in future
	// We assume that a connection is usually changed rarely.
	// So connection can not be deleted during this operation.
	for _, connIDRaw := range connectionIds[startPosition:] {
		if len(connResults) >= *listOptions.Size {
			break
		}

		connID, ok := connIDRaw.(string)
		if !ok {
			return nil, odahuflow_errors.SerializationError{}
		}

		conn, err := vcr.GetConnection(connID)
		if err != nil {
			return nil, err
		}

		if len(listOptions.Filter.Type) == 0 {
			connResults = append(connResults, *conn)
		} else {
			for _, connType := range listOptions.Filter.Type {
				if connType == string(conn.Spec.Type) {
					connResults = append(connResults, *conn)
					continue
				}
			}
		}
	}

	return connResults, nil
}

func (vcr *vaultConnRepository) DeleteConnection(connID string) error {
	if _, err := vcr.GetConnection(connID); err != nil {
		return err
	}

	_, err := vcr.vaultClient.Logical().Delete(vcr.getFullPath(connID))

	return convertVaultErrToOdahuflowErr(err)
}

func (vcr *vaultConnRepository) UpdateConnection(conn *connection.Connection) error {
	existedConn, err := vcr.GetConnection(conn.ID)

	switch {
	case err == nil:
		// If err is not nil, then the connection already exists.
		existedConn.Spec = conn.Spec
		existedConn.Status.UpdatedAt = &metav1.Time{Time: time.Now()}
		err = vcr.createOrUpdateConnection(existedConn)
		if err != nil {
			conn.Status = existedConn.Status
		}
		return err
	case odahuflow_errors.IsNotFoundError(err):
		return odahuflow_errors.NotFoundError{Entity: conn.ID}
	default:
		return err
	}
}

func (vcr *vaultConnRepository) CreateConnection(conn *connection.Connection) error {
	_, err := vcr.GetConnection(conn.ID)

	switch {
	case err == nil:
		// If err is nil, then the connection already exists.
		return odahuflow_errors.AlreadyExistError{Entity: conn.ID}
	case odahuflow_errors.IsNotFoundError(err):
		conn.Status.CreatedAt = &metav1.Time{Time: time.Now()}
		conn.Status.UpdatedAt = &metav1.Time{Time: time.Now()}
		return vcr.createOrUpdateConnection(conn)
	default:
		return err
	}
}

func (vcr *vaultConnRepository) createOrUpdateConnection(conn *connection.Connection) error {
	_, err := vcr.vaultClient.Logical().Write(
		vcr.getFullPath(conn.ID),
		map[string]interface{}{
			ConnKey: conn,
		},
	)

	return convertVaultErrToOdahuflowErr(err)
}

func (vcr *vaultConnRepository) getConnectionFromVault(connID string) (*connection.Connection, error) {
	vaultEntity, err := vcr.vaultClient.Logical().Read(vcr.getFullPath(connID))
	if err != nil {
		return nil, convertVaultErrToOdahuflowErr(err)
	}

	// The vaultEntity is nil when binaryData is not located by specific connectionVaultPath
	if vaultEntity == nil {
		return nil, odahuflow_errors.NotFoundError{Entity: connID}
	}

	binaryData, ok := vaultEntity.Data[ConnKey]
	if !ok {
		return nil, odahuflow_errors.NotFoundError{Entity: connID}
	}

	connData, err := json.Marshal(binaryData)
	if err != nil {
		return nil, err
	}

	conn := &connection.Connection{}
	err = json.Unmarshal(connData, conn)
	if err != nil {
		return nil, odahuflow_errors.SerializationError{}
	}

	return conn, nil
}

// Construct the full vault connectionVaultPath for odahuflow connections
func (vcr *vaultConnRepository) getFullPath(connID string) (vaultPath string) {
	return path.Join(vcr.connectionVaultPath, connID)
}

func convertVaultErrToOdahuflowErr(err error) error {
	if verr, ok := err.(*vaultapi.ResponseError); ok && verr.StatusCode == http.StatusForbidden {
		return odahuflow_errors.ForbiddenError{}
	}

	return err
}
