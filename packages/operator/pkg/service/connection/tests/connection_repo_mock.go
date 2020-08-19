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
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"github.com/stretchr/testify/mock"
)

type ConnectionRepositoryMock struct {
	mock.Mock
	// fields to remember passed pointers state
	UpdatedConnection connection.Connection
	CreatedConnection connection.Connection
}

func (c *ConnectionRepositoryMock) GetConnection(id string) (*connection.Connection, error) {
	mockedResult := c.Called(id)
	if conn, ok := mockedResult.Get(0).(*connection.Connection); ok {
		return conn, mockedResult.Error(1)
	}
	return nil, mockedResult.Error(1)
}

func (c *ConnectionRepositoryMock) GetConnectionList(options ...conn_repository.ListOption) ([]connection.Connection, error) {
	mockedResult := c.Called(options)
	return mockedResult.Get(0).([]connection.Connection), mockedResult.Error(1)
}

func (c *ConnectionRepositoryMock) DeleteConnection(id string) error {
	mockedResult := c.Called(id)
	return mockedResult.Error(0)
}

func (c *ConnectionRepositoryMock) UpdateConnection(conn *connection.Connection) error {
	// Remember the state on passed connection for assertions
	c.UpdatedConnection = *conn
	mockedResult := c.Called(conn)
	return mockedResult.Error(0)
}

func (c *ConnectionRepositoryMock) CreateConnection(conn *connection.Connection) error {
	// Remember the state on passed connection for assertions
	c.CreatedConnection = *conn
	mockedResult := c.Called(conn)
	return mockedResult.Error(0)
}
