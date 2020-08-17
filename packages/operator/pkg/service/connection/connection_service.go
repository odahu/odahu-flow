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
	} else {
		conn.EncodeBase64Secrets()
	}

	return conn, nil
}

func (s *serviceImpl) GetConnectionList(options ...conn_repository.ListOption) ([]connection.Connection, error) {
	connections, err := s.repo.GetConnectionList(options...)
	if err != nil {
		return connections, err
	}

	for i := range connections {
		connections[i].DeleteSensitiveData()
		//conn.DeleteSensitiveData()
	}

	return connections, err
}

func (s *serviceImpl) DeleteConnection(id string) error {
	return s.repo.DeleteConnection(id)
}

func (s *serviceImpl) UpdateConnection(connection connection.Connection) (*connection.Connection, error) {
	if err := connection.DecodeBase64Secrets(); err != nil {
		return nil, errors.InvalidEntityError{
			Entity:           fmt.Sprintf("Connection %s", connection.ID),
			ValidationErrors: multierr.Errors(err),
		}
	}

	updatedConn, err := s.repo.UpdateConnection(&connection)
	if err != nil {
		return updatedConn, err
	}
	return updatedConn.DeleteSensitiveData(), err
}

func (s *serviceImpl) CreateConnection(connection connection.Connection) (*connection.Connection, error) {
	if err := connection.DecodeBase64Secrets(); err != nil {
		return nil, errors.InvalidEntityError{
			Entity:           fmt.Sprintf("Connection %s", connection.ID),
			ValidationErrors: multierr.Errors(err),
		}
	}

	createdConn, err := s.repo.CreateConnection(&connection)
	if err != nil {
		return createdConn, err
	}
	return createdConn.DeleteSensitiveData(), err
}
