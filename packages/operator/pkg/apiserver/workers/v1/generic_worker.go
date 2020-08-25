package v1

import (
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/types"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	log = logf.Log.WithName("generic-worker")
)

type GenericWorker struct {
	mgr          ctrl.Manager
	name         string
	launchPeriod time.Duration
	syncer       types.StorageServiceSyncer
}

func NewGenericWorker(
	mgr ctrl.Manager,
	name string,
	launchPeriod time.Duration,
	syncer types.StorageServiceSyncer,
) GenericWorker {
	return GenericWorker{
		mgr:          mgr,
		name:         name,
		launchPeriod: launchPeriod,
		syncer:       syncer,
	}
}

// Return name of runner
func (r *GenericWorker) String() string {
	return r.name
}

func (r *GenericWorker) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	ID := request.Name
	eLog := log.WithValues("ID", ID, "Flow", "Storage <- Service", "worker-name", r.String())
	// We need periodically check kube entities even if there are have no updates
	// So we can remove them if corresponding entity in storage is deleted yet
	periodicCheck := ctrl.Result{RequeueAfter: time.Second * 10}

	kubeEntity, err := r.syncer.GetFromService(ID)
	if odahu_errors.IsNotFoundError(err) {
		return periodicCheck, nil
	}
	if err != nil {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return ctrl.Result{}, err
	}

	if err := kubeEntity.ReportStatus(); err != nil {

		if odahu_errors.IsNotFoundError(err) {
			return periodicCheck, nil
		}
		eLog.Error(err, "Unable to update storage entity")
		return ctrl.Result{}, err
	}
	eLog.Info("Storage entity was maybe updated")
	return periodicCheck, nil

}

func (r *GenericWorker) syncSpecs(_ context.Context) error {

	eLog := log.WithValues("Flow", "Storage -> Service", "worker-name", r.String())

	create, update, del, delDB, err := r.getNotSynced()

	if len(create) > 0 || len(update) > 0 || len(del) > 0|| len(delDB) > 0 {
		log.Info(
			"Get not synced entities",
			"create", len(create), "update", len(update), "del", len(del), "delDB", len(delDB),
		)
	}

	if err != nil {
		return err
	}

	for _, storeEn := range create {
		crErr := storeEn.CreateInService()
		if crErr != nil {
			eLog.Error(err, "Unable to create service entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("Service entity is created", "ID", storeEn.GetID())
		}

	}

	for _, storeEn := range update {
		upErr := storeEn.UpdateInService()
		if upErr != nil {
			eLog.Error(err, "Unable to update service entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("Service entity is updated", "ID", storeEn.GetID())
		}
	}

	for _, servEn := range del {
		delErr := servEn.Delete()
		if delErr != nil {
			eLog.Error(err, "Unable to delete service entity", "ID", servEn.GetID())
		} else {
			eLog.Info("Service entity is deleted", "ID", servEn.GetID())
		}

	}

	for _, storeEn := range delDB {
		delErr := storeEn.DeleteInDB()
		if delErr != nil {
			eLog.Error(err, "Unable to delete DB entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("DB entity is deleted", "ID", storeEn.GetID())
		}

	}

	return nil
}

func (r *GenericWorker) getNotSynced() (
	toCreateInService []types.StorageEntity,
	toUpdateInService []types.StorageEntity,
	toDeleteInService []types.ExternalServiceEntity,
	toDeleteInDB []types.StorageEntity, err error, ) {

	servEnsList, err := r.syncer.ListService()
	if err != nil {
		return nil, nil, nil, toDeleteInDB, err
	}

	storeEnsList, err := r.syncer.ListStorage()
	if err != nil {
		return nil, nil, nil, toDeleteInDB, err
	}

	// Find all entities that are exist only in storage or their spec in storage is different with service
	servEnsIndex := make(map[string]types.ExternalServiceEntity)
	for _, en := range servEnsList {
		servEnsIndex[en.GetID()] = en
	}
	for _, storeEn := range storeEnsList {

		servEn, existsInService := servEnsIndex[storeEn.GetID()]

		if storeEn.IsFinished(){
			continue
		}

		if !existsInService && !storeEn.HasDeletionMark() {
			toCreateInService = append(toCreateInService, storeEn)
			continue
		}
		if !existsInService && storeEn.HasDeletionMark() {
			toDeleteInDB = append(toDeleteInDB, storeEn)
			continue
		}
		if existsInService && !servEn.IsDeleting() && storeEn.HasDeletionMark() {
			toDeleteInService = append(toDeleteInService, servEn)
			continue
		}
		if existsInService && servEn.IsDeleting() && storeEn.HasDeletionMark() {
			continue
		}

		eq, eqErr := specsEqual(servEn, storeEn)
		if eqErr != nil {
			log.Error(
				eqErr,
				fmt.Sprintf("unable compare specs storage: %+v; service: %+v", storeEn, servEn),
			)
			continue
		}

		if !eq {
			toUpdateInService = append(toUpdateInService, storeEn) // toUpdateInService
			continue
		}
	}

	// Find all zombies (entities that exist only in service)
	storeEnsIndex := make(map[string]types.StorageEntity)
	for _, en := range storeEnsList {
		storeEnsIndex[en.GetID()] = en
	}
	for _, servEn := range servEnsList {
		if _, ok := storeEnsIndex[servEn.GetID()]; !ok {
			toDeleteInService  = append(toDeleteInService, servEn)
		}
	}

	return toCreateInService, toUpdateInService, toDeleteInService, toDeleteInDB, err

}

func (r *GenericWorker) Run(ctx context.Context) (err error) {

	log.Info(fmt.Sprintf("%v is running", r.String()))

	// Attach controller to manager
	// This controller get updates from K8S resource and update state in DB
	err = r.syncer.AttachController(r.mgr, r)
	if err != nil {
		log.Error(err, "unable to initialize controller")
		return err
	}
	log.Info("Persistence controller for attached to kube manager")

	// We sync entities in service from storage every `launchPeriod` seconds
	t := time.NewTicker(r.launchPeriod)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = r.syncSpecs(ctx)
			if err != nil {
				log.Error(err, "Error while syncSpecs")
			}
			continue

		case <-ctx.Done():
			log.Info(fmt.Sprintf("Cancellation signal was received in %v", r.String()))
		}
		break
	}

	return nil
}
