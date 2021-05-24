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

package memory

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"sync"
)

type repository struct {
	mu *sync.Mutex
	m  map[string]*connection.Connection
}

// Not for Production! Development / Testing purposes only
func NewRepository() conn_repository.Repository {
	return repository{
		mu: &sync.Mutex{},
		m:  make(map[string]*connection.Connection),
	}
}

func (r repository) GetConnection(id string) (*connection.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	conn, ok := r.m[id]
	if !ok {
		return nil, odahu_errors.NotFoundError{Entity: id}
	}
	return conn, nil
}

// Does not support filtering yet
func (r repository) GetConnectionList(options ...conn_repository.ListOption) (res []connection.Connection, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range r.m {
		res = append(res, *v)
	}
	return res, err
}

func (r repository) DeleteConnection(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.m, id)
	return nil
}

func (r repository) UpdateConnection(connection *connection.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[connection.ID] = connection
	return nil
}

func (r repository) SaveConnection(connection *connection.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[connection.ID] = connection
	return nil
}
