/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package connection

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"go.uber.org/multierr"
	"time"
)

// A layer on top of connection repository that prepares data, e.g. base64 decoding
type Service interface {
	GetConnection(id string, encrypted bool) (*connection.Connection, error)
	GetConnectionList(options ...conn_repository.ListOption) ([]connection.Connection, error)
	DeleteConnection(id string) error
	UpdateConnection(connection connection.Connection) (*connection.Connection, error)
	CreateConnection(connection connection.Connection) (*connection.Connection, error)
}

type serviceImpl struct {
	repo conn_repository.Repository
}

func NewService(repo conn_repository.Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetConnection(id string, encrypted bool) (*connection.Connection, error) {
	conn, err := s.repo.GetConnection(id)
	if err != nil {
		return nil, err
	}

	if encrypted {
		conn.DeleteSensitiveData()
	}
	conn.EncodeBase64Fields()

	return conn, nil
}

func (s *serviceImpl) GetConnectionList(options ...conn_repository.ListOption) ([]connection.Connection, error) {
	connections, err := s.repo.GetConnectionList(options...)
	if err != nil {
		return connections, err
	}

	for i := range connections {
		connections[i].DeleteSensitiveData()
		connections[i].EncodeBase64Fields()
	}

	return connections, err
}

func (s *serviceImpl) DeleteConnection(id string) error {
	if _, err := s.repo.GetConnection(id); err != nil {
		return err
	}

	return s.repo.DeleteConnection(id)
}

func (s *serviceImpl) UpdateConnection(connection connection.Connection) (*connection.Connection, error) {
	existedConn, err := s.repo.GetConnection(connection.ID)

	switch {
	case errors.IsNotFoundError(err):
		return nil, errors.NotFoundError{Entity: connection.ID}
	case err != nil:
		return nil, err
	}

	if err := connection.DecodeBase64Fields(); err != nil {
		return nil, errors.InvalidEntityError{
			Entity:           fmt.Sprintf("Connection %s", connection.ID),
			ValidationErrors: multierr.Errors(err),
		}
	}

	existedConn.Spec = connection.Spec
	existedConn.UpdatedAt = time.Now()
	err = s.repo.UpdateConnection(existedConn)
	if err != nil {
		connection.Status = existedConn.Status
	}

	existedConn.DeleteSensitiveData()
	existedConn.EncodeBase64Fields()
	return existedConn, err
}

func (s *serviceImpl) CreateConnection(connection connection.Connection) (*connection.Connection, error) {
	_, err := s.repo.GetConnection(connection.ID)

	switch {
	case err == nil:
		// If err is nil, then the connection already exists.
		return nil, errors.AlreadyExistError{Entity: connection.ID}
	case errors.IsNotFoundError(err):
		// Good path
	default:
		return nil, err
	}

	connection.CreatedAt = time.Now()
	connection.UpdatedAt = connection.CreatedAt

	if err := connection.DecodeBase64Fields(); err != nil {
		return nil, errors.InvalidEntityError{
			Entity:           fmt.Sprintf("Connection %s", connection.ID),
			ValidationErrors: multierr.Errors(err),
		}
	}

	err = s.repo.SaveConnection(&connection)
	if err != nil {
		return nil, err
	}

	connection.DeleteSensitiveData()
	connection.EncodeBase64Fields()
	return &connection, err
}
