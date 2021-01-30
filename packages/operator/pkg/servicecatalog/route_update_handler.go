package servicecatalog

import (
	"errors"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/deployment"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/inspectors"
	"go.uber.org/zap"
)

type UpdateHandler struct {
	Inspectors       map[string]inspectors.ModelServerInspector
	Catalog          Catalog
	DeploymentClient deployment.Client
}

func extractModelDeploymentID(modelRoute deployment_types.ModelRoute) string {
	if modelRoute.Default && len(modelRoute.Spec.ModelDeploymentTargets) > 0 {
		target := modelRoute.Spec.ModelDeploymentTargets[0]
		return target.Name
	}
	return ""
}

func (r UpdateHandler) inspectModel(
	prefix string, predictor string, log *zap.SugaredLogger) (model model_types.ServedModel, err error) {
	log.Debugw("Inspecting deployed model...", "predictor", predictor, "prefix", prefix)

	inspector, ok := r.Inspectors[predictor]
	if !ok {
		errorStr := "no inspector for predictor: " + predictor
		log.Error(errorStr)
		return model, errors.New(errorStr)
	}

	return inspector.Inspect(prefix, log)
}

func (r UpdateHandler) Handle(object interface{}, log *zap.SugaredLogger) (err error) {

	event, ok := object.(event_types.RouteEvent)
	if !ok {
		return fmt.Errorf("unable to assert received object as event_types.RouteEvent: %+v", object)
	}

	log = log.With("EventType", event.EventType, "EventDatetime", event.Datetime.String())

	if event.EventType == event_types.ModelRouteDeletedEventType {
		r.Catalog.Delete(event.EntityID, log)
		return nil
	}

	if !event.Payload.Default {
		log.Debug("Not default route will not be processed")
		return nil
	}

	if event.Payload.Status.State != odahuflowv1alpha1.ModelRouteStateReady {
		return fmt.Errorf("ModelRoute %s has not a Ready state", event.EntityID)
	}

	route := Route{
		ID:        event.EntityID,
		Prefix:    event.Payload.Spec.URLPrefix,
		IsDefault: event.Payload.Default,
	}

	mdID := extractModelDeploymentID(event.Payload)
	log.Infow("Handling route event for deployment", "md id", mdID)
	md, err := r.DeploymentClient.GetModelDeployment(mdID)
	if err != nil {
		log.Error(err, "failed to GET model deployment from API")
		return err
	}

	var servedModel model_types.ServedModel
	servedModel, err = r.inspectModel(route.Prefix, md.Spec.Predictor, log)
	if err != nil {
		return err
	}
	deployedModel := model_types.DeployedModel{
		DeploymentID: mdID,
		ServedModel:  servedModel,
	}
	route.Model = deployedModel

	log.Infow("Adding a model to catalog", "route", route)
	if err != r.Catalog.CreateOrUpdate(route) {
		log.Error(err, "adding to catalog")
		return err
	}

	return nil

}
