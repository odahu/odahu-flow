package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
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

type DeploymentRepo struct {
	DB *sql.DB
}

func (repo DeploymentRepo) GetModelDeployment(name string) (*deployment.ModelDeployment, error) {

	mt := new(deployment.ModelDeployment)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status, deletionmark FROM %s WHERE id = $1", ModelDeploymentTable),
		name,
	).Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark)

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

func (repo DeploymentRepo) GetModelDeploymentList(options ...filter.ListOption) (
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

	sb := sq.Select("id, spec, status, deletionmark").From("odahu_operator_deployment").
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
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark)
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

func (repo DeploymentRepo) DeleteModelDeployment(name string) error {

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

func (repo DeploymentRepo) SetDeletionMark(id string, value bool) error {
	return utils.SetDeletionMark(repo.DB, ModelDeploymentTable, id, value)
}

func (repo DeploymentRepo) UpdateModelDeployment(md *deployment.ModelDeployment) error {

	md.Status.State = ""

	stmt, args, err := sq.Update(ModelDeploymentTable).
		Set("spec", md.Spec).
		Set("status", md.Status).
		Where(sq.Eq{"id": md.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := repo.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: md.ID}
	}

	return nil
}

func (repo DeploymentRepo) UpdateModelDeploymentStatus(id string, s v1alpha1.ModelDeploymentStatus) error {

	stmt, args, err := sq.Update(ModelDeploymentTable).
		Set("status", s).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := repo.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: id}
	}

	return nil
}

func (repo DeploymentRepo) CreateModelDeployment(md *deployment.ModelDeployment) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelDeploymentTable),
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

func (repo DeploymentRepo) GetModelRoute(name string) (*deployment.ModelRoute, error) {

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

func (repo DeploymentRepo) GetModelRouteList(options ...filter.ListOption) (
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

	sb := sq.Select("id, spec, status").From("odahu_operator_route").
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

func (repo DeploymentRepo) DeleteModelRoute(name string) error {

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

func (repo DeploymentRepo) UpdateModelRoute(mr *deployment.ModelRoute) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelRoute(mr.ID)
	if err != nil {
		return err
	}

	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", ModelRouteTable)
	_, err = repo.DB.Exec(sqlStatement, mr.Spec, mr.Status, mr.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo DeploymentRepo) CreateModelRoute(mr *deployment.ModelRoute) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelRouteTable),
		mr.ID, mr.Spec, mr.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mr.ID}
		}
		return err
	}
	return nil

}
