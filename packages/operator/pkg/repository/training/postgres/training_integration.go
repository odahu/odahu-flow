package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	trainingIntegrationTable    = "odahu_operator_training_integration"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

var (
	MaxSize   = 500
	FirstPage = 0
)

type TrainingIntegrationRepo struct {
	DB *sql.DB
}

func (tr TrainingIntegrationRepo) GetTrainingIntegration(name string) (*training.TrainingIntegration, error) {

	ti := new(training.TrainingIntegration)

	err := tr.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status, created, updated FROM %s WHERE id = $1", trainingIntegrationTable),
		name,
	).Scan(&ti.ID, &ti.Spec, &ti.Status, &ti.CreatedAt, &ti.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		return nil, err
	default:
		return ti, nil
	}

}

func (tr TrainingIntegrationRepo) GetTrainingIntegrationList(options ...filter.ListOption) (
	[]training.TrainingIntegration, error,
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

	stmt := "SELECT id, spec, status, created, updated " +
		"FROM odahu_operator_training_integration ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := tr.DB.Query(stmt, *listOptions.Size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tis []training.TrainingIntegration

	log.Info("!DEBUGGING1!", "result", tis, "&result", &tis)

	for rows.Next() {
		ti := new(training.TrainingIntegration)
		err := rows.Scan(&ti.ID, &ti.Spec, &ti.Status, &ti.CreatedAt, &ti.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tis = append(tis, *ti)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	log.Info("!DEBUGGING2", "result", tis, "&result", &tis)
	return tis, nil

}

func (tr TrainingIntegrationRepo) DeleteTrainingIntegration(name string) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := tr.GetTrainingIntegration(name)
	if err != nil {
		return err
	}

	// If exists, delete it

	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", trainingIntegrationTable)
	_, err = tr.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (tr TrainingIntegrationRepo) UpdateTrainingIntegration(md *training.TrainingIntegration) error {

	// First try to check that row exists otherwise raise exception to fit interface
	oldTi, err := tr.GetTrainingIntegration(md.ID)
	if err != nil {
		return err
	}

	md.Status = oldTi.Status

	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2, updated = $3 WHERE id = $4",
		trainingIntegrationTable)
	_, err = tr.DB.Exec(sqlStatement, md.Spec, md.Status, md.UpdatedAt, md.ID)
	if err != nil {
		return err
	}
	return nil
}

func (tr TrainingIntegrationRepo) SaveTrainingIntegration(md *training.TrainingIntegration) error {

	_, err := tr.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status, created, updated) VALUES($1, $2, $3, $4, $5)",
			trainingIntegrationTable),
		md.ID, md.Spec, md.Status, md.CreatedAt, md.UpdatedAt,
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
