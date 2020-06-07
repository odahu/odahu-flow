package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
)

const (
	toolchainIntegrationTable   = "odahu_operator_toolchain_integration"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

type ToolchainRepository struct {
	DB *sql.DB
}

func (tr ToolchainRepository) GetToolchainIntegration(name string) (*training.ToolchainIntegration, error) {

	ti := new(training.ToolchainIntegration)

	err := tr.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status FROM %s WHERE id = $1", toolchainIntegrationTable),
		name,
	).Scan(&ti.ID, &ti.Spec, &ti.Status)

	switch {
	case err == sql.ErrNoRows:
		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		return nil, err
	default:
		return ti, nil
	}

}

func (tr ToolchainRepository) GetToolchainIntegrationList(options ...kubernetes.ListOption) ([]training.ToolchainIntegration, error) {
	rows, err := tr.DB.Query(
		fmt.Sprintf("SELECT id, spec, status FROM %s", toolchainIntegrationTable),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tis []training.ToolchainIntegration

	for rows.Next() {
		ti := new(training.ToolchainIntegration)
		err := rows.Scan(&ti.ID, &ti.Spec, &ti.Status)
		if err != nil {
			return nil, err
		}
		tis = append(tis, *ti)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tis, nil

}

func (tr ToolchainRepository) DeleteToolchainIntegration(name string) error {
	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", toolchainIntegrationTable)
	_, err := tr.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (tr ToolchainRepository) UpdateToolchainIntegration(md *training.ToolchainIntegration) error {
	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", toolchainIntegrationTable)
	_, err := tr.DB.Exec(sqlStatement, md.Spec, md.Status, md.ID)
	if err != nil {
		return err
	}
	return nil
}

func (tr ToolchainRepository) CreateToolchainIntegration(md *training.ToolchainIntegration) error {
	_, err := tr.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", toolchainIntegrationTable),
		md.ID, md.Spec, md.Status,
	)

	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: md.ID}
		}
		return err
	}
	return nil

}
