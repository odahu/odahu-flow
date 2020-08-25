package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelDeploymentTable = "odahu_operator_deployment"
	ModelRouteTable = "odahu_operator_route"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

var (
	log = logf.Log.WithName("model-deployment--repository--postgres")
	MaxSize   = 500
	FirstPage = 0
)

type DeploymentPostgresRepo struct {
	DB *sql.DB
}

func (repo DeploymentPostgresRepo) GetModelDeployment(name string) (*deployment.ModelDeployment, error) {

	mt := new(deployment.ModelDeployment)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status FROM %s WHERE id = $1", ModelDeploymentTable),
		name,
	).Scan(&mt.ID, &mt.Spec, &mt.Status)

	switch {
	case err == sql.ErrNoRows:

		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		log.Error(err, "error during sql query")

		return nil, err
	default:

		return mt, nil
	}

}

func (repo DeploymentPostgresRepo) GetModelDeploymentList(options ...filter.ListOption) (
	[]deployment.ModelDeployment, error,
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

	sb := sq.Select("*").From("odahu_operator_deployment").
		OrderBy("id").
		Offset(uint64(offset)).
		Limit(uint64(*listOptions.Size)).PlaceholderFormat(sq.Dollar)

	sb = utils.TransformFilter(sb, listOptions.Filter)
	stmt, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "error during rows.Close()")
		}
		if err := rows.Err(); err != nil {
			log.Error(err, "error during rows iterating")
		}
	}()

	mts := make([]deployment.ModelDeployment, 0)

	for rows.Next() {
		mt := new(deployment.ModelDeployment)
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status)
		if err != nil {
			return nil, err
		}
		mts = append(mts, *mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mts, nil

}

func (repo DeploymentPostgresRepo) DeleteModelDeployment(name string) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelDeployment(name)
	if err != nil {
		return err
	}

	// If exists, delete it

	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ModelDeploymentTable)
	_, err = repo.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (repo DeploymentPostgresRepo) UpdateModelDeployment(mt *deployment.ModelDeployment) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelDeployment(mt.ID)
	if err != nil {
		return err
	}

	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", ModelDeploymentTable)
	_, err = repo.DB.Exec(sqlStatement, mt.Spec, mt.Status, mt.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo DeploymentPostgresRepo) CreateModelDeployment(mt *deployment.ModelDeployment) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelDeploymentTable),
		mt.ID, mt.Spec, mt.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mt.ID}
		}
		return err
	}
	return nil

}

func (repo DeploymentPostgresRepo) GetModelRoute(name string) (*deployment.ModelRoute, error) {

	mt := new(deployment.ModelRoute)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status FROM %s WHERE id = $1", ModelRouteTable),
		name,
	).Scan(&mt.ID, &mt.Spec, &mt.Status)

	switch {
	case err == sql.ErrNoRows:

		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		log.Error(err, "error during sql query")

		return nil, err
	default:

		return mt, nil
	}

}

func (repo DeploymentPostgresRepo) GetModelRouteList(options ...filter.ListOption) (
	[]deployment.ModelRoute, error,
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

	sb := sq.Select("*").From("odahu_operator_route").
		OrderBy("id").
		Offset(uint64(offset)).
		Limit(uint64(*listOptions.Size)).PlaceholderFormat(sq.Dollar)

	sb = utils.TransformFilter(sb, listOptions.Filter)
	stmt, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "error during rows.Close()")
		}
		if err := rows.Err(); err != nil {
			log.Error(err, "error during rows iterating")
		}
	}()

	mts := make([]deployment.ModelRoute, 0)

	for rows.Next() {
		mt := new(deployment.ModelRoute)
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status)
		if err != nil {
			return nil, err
		}
		mts = append(mts, *mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mts, nil

}

func (repo DeploymentPostgresRepo) DeleteModelRoute(name string) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelRoute(name)
	if err != nil {
		return err
	}

	// If exists, delete it

	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ModelRouteTable)
	_, err = repo.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (repo DeploymentPostgresRepo) UpdateModelRoute(mt *deployment.ModelRoute) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelRoute(mt.ID)
	if err != nil {
		return err
	}

	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", ModelRouteTable)
	_, err = repo.DB.Exec(sqlStatement, mt.Spec, mt.Status, mt.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo DeploymentPostgresRepo) CreateModelRoute(mt *deployment.ModelRoute) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelRouteTable),
		mt.ID, mt.Spec, mt.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mt.ID}
		}
		return err
	}
	return nil

}
