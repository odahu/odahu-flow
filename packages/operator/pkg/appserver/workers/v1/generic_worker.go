package v1

import (
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1/types"
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
	result := ctrl.Result{}

	kubeEntity, err := r.syncer.GetFromService(ID)

	if err != nil {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return result, err
	}

	if err := kubeEntity.SaveResultInStorage(); err != nil {
		eLog.Error(err, "Unable to update storage entity")
		return result, err
	}
	eLog.Info("Storage entity was maybe updated")
	return result, err

}

func (r *GenericWorker) syncSpecs(_ context.Context) error {

	eLog := log.WithValues("Flow", "Storage -> Service", "worker-name", r.String())

	created, specChanged, servZombies, err := r.getNotSynced()

	if err != nil {
		return err
	}

	for _, storeEn := range created {
		crErr := storeEn.CreateInService()
		if crErr != nil {
			eLog.Error(err, "Unable to create service entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("Service entity is created", "ID", storeEn.GetID())
		}

	}

	for _, storeEn := range specChanged {
		upErr := storeEn.UpdateInService()
		if upErr != nil {
			eLog.Error(err, "Unable to update service entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("Service entity is updated", "ID", storeEn.GetID())
		}
	}

	for _, servEn := range servZombies {
		delErr := servEn.DeleteInService()
		if delErr != nil {
			eLog.Error(err, "Unable to delete service entity", "ID", servEn.GetID())
		} else {
			eLog.Info("Service entity is deleted", "ID", servEn.GetID())
		}

	}

	return nil
}

// return not synced entities
// created – entities that are exist only in storage
// specChanged – entities that are exist in service and in storage but specs differ
// servZombies – entities that are exist only in service
func (r *GenericWorker) getNotSynced() (
	created []types.StorageEntity,
	specChanged []types.StorageEntity,
	servZombies []types.ExternalServiceEntity, err error, ) {

	servEnsList, err := r.syncer.ListService()
	if err != nil {
		return nil, nil, nil, err
	}

	storeEnsList, err := r.syncer.ListStorage()
	if err != nil {
		return nil, nil, nil, err
	}

	// Find all entities that are exist only in storage or their spec in storage is different with service
	servEnsIndex := make(map[string]types.ExternalServiceEntity)
	for _, en := range servEnsList {
		servEnsIndex[en.GetID()] = en
	}
	for _, storeEn := range storeEnsList {

		servEn, ok := servEnsIndex[storeEn.GetID()]
		if !ok {
			created = append(created, storeEn)
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
			specChanged = append(specChanged, storeEn)
			continue
		}
	}

	// Find all entities that are exist only in service
	storeEnsIndex := make(map[string]types.StorageEntity)
	for _, en := range storeEnsList {
		storeEnsIndex[en.GetID()] = en
	}
	for _, servEn := range servEnsList {
		if _, ok := storeEnsIndex[servEn.GetID()]; !ok {
			servZombies = append(servZombies, servEn)
		}
	}

	return created, specChanged, servZombies, err

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
