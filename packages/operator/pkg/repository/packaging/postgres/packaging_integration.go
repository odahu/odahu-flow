//
//    Copyright 2020 EPAM Systems
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

package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	packagingIntegrationTable   = "odahu_operator_packaging_integration"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

var (
	MaxSize   = 500
	FirstPage = 0
)

type PackagingIntegrationRepository struct {
	DB *sql.DB
}

func (pir *PackagingIntegrationRepository) GetPackagingIntegration(name string) (
	*packaging.PackagingIntegration, error,
) {

	pi := new(packaging.PackagingIntegration)

	err := pir.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status FROM %s WHERE id = $1", packagingIntegrationTable),
		name,
	).Scan(&pi.ID, &pi.Spec, &pi.Status)

	switch {
	case err == sql.ErrNoRows:
		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		return nil, err
	default:
		return pi, nil
	}

}

func (pir *PackagingIntegrationRepository) GetPackagingIntegrationList(options ...filter.ListOption) (
	[]packaging.PackagingIntegration, error,
) {

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &FirstPage,
		Size:   &MaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	offset := *listOptions.Size * (*listOptions.Page)

	stmt := "SELECT id, spec, status FROM odahu_operator_packaging_integration ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := pir.DB.Query(stmt, *listOptions.Size, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pis []packaging.PackagingIntegration

	for rows.Next() {
		pi := new(packaging.PackagingIntegration)
		err := rows.Scan(&pi.ID, &pi.Spec, &pi.Status)
		if err != nil {
			return nil, err
		}
		pis = append(pis, *pi)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return pis, nil

}

func (pir *PackagingIntegrationRepository) DeletePackagingIntegration(name string) error {
	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", packagingIntegrationTable)
	_, err := pir.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (pir *PackagingIntegrationRepository) UpdatePackagingIntegration(pi *packaging.PackagingIntegration) error {
	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", packagingIntegrationTable)
	_, err := pir.DB.Exec(sqlStatement, pi.Spec, pi.Status, pi.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pir *PackagingIntegrationRepository) SavePackagingIntegration(pi *packaging.PackagingIntegration) error {

	_, err := pir.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", packagingIntegrationTable),
		pi.ID, pi.Spec, pi.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: pi.ID}
		}
		return err
	}
	return nil

}
