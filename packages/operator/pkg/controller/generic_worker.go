package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/utils"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"time"
)

type GenericWorker struct {
	name         string
	launchPeriod time.Duration
	syncer       types.RuntimeAdapter
}

func NewGenericWorker(
	name string,
	launchPeriod time.Duration,
	syncer types.RuntimeAdapter,
) GenericWorker {
	return GenericWorker{
		name:         name,
		launchPeriod: launchPeriod,
		syncer:       syncer,
	}
}

// Return name of runner
func (r *GenericWorker) String() string {
	return r.name
}

func (r *GenericWorker) SyncStatus(re types.RuntimeEntity, se types.StorageEntity) error {

	eLog := log.WithValues("ID", re.GetID(), "Flow", "Storage <- Service", "worker-name", r.String())

	eq, err := utils.SpecsEqual(re, se)
	if err != nil {
		return err
	}
	if !eq {
		return errors.New("spec in runtime not equal spec in storage. Skip status delivery")
	}

	if err := re.ReportStatus(); err != nil {
		if odahu_errors.IsNotFoundError(err) {
			return nil
		}
		if odahu_errors.IsSpecWasChangedError(err) {
			return nil
		}
		eLog.Error(err, "Unable to report status")
		return err
	}

	eLog.Info("Storage entity was maybe updated")
	return nil
}

func (r *GenericWorker) SyncSpecs(_ context.Context) error {

	eLog := log.WithValues("Flow", "Storage -> Service", "worker-name", r.String())

	create, update, del, delDB, err := r.GetNotSynced()

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
		crErr := storeEn.CreateInRuntime()
		if crErr != nil {
			eLog.Error(err, "Unable to create service entity", "ID", storeEn.GetID())
		} else {
			eLog.Info("Service entity is created", "ID", storeEn.GetID())
		}

	}

	for _, storeEn := range update {
		upErr := storeEn.UpdateInRuntime()
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

func (r *GenericWorker) GetNotSynced() (
	toCreateInService []types.StorageEntity,
	toUpdateInService []types.StorageEntity,
	toDeleteInService []types.RuntimeEntity,
	toDeleteInDB []types.StorageEntity, err error, ) {

	servEnsList, err := r.syncer.ListRuntime()
	if err != nil {
		return toCreateInService, toUpdateInService, toDeleteInService, toDeleteInDB, err
	}

	storeEnsList, err := r.syncer.ListStorage()
	if err != nil {
		return toCreateInService, toUpdateInService, toDeleteInService, toDeleteInDB, err
	}

	// Find all entities that are exist only in tStorage or their spec in tStorage is different with service
	servEnsIndex := make(map[string]types.RuntimeEntity)
	for _, en := range servEnsList {
		servEnsIndex[en.GetID()] = en
	}
	for _, storeEn := range storeEnsList {

		servEn, existsInService := servEnsIndex[storeEn.GetID()]

		if !existsInService && !storeEn.HasDeletionMark() && !storeEn.IsFinished() {
			toCreateInService = append(toCreateInService, storeEn)
			continue
		}

		if !existsInService && !storeEn.HasDeletionMark() && storeEn.IsFinished() {
			continue
		}
		if !existsInService && storeEn.HasDeletionMark() {
			toDeleteInDB = append(toDeleteInDB, storeEn)
			continue
		}

		eq, eqErr := utils.SpecsEqual(servEn, storeEn)
		if eqErr != nil {
			log.Error(
				eqErr,
				fmt.Sprintf("unable compare specs tStorage: %+v; service: %+v", storeEn, servEn),
			)
			continue
		}
		if storeEn.IsFinished() && eq{
			continue
		}

		if existsInService && !servEn.IsDeleting() && storeEn.HasDeletionMark() {
			toDeleteInService = append(toDeleteInService, servEn)
			continue
		}
		if existsInService && servEn.IsDeleting() && storeEn.HasDeletionMark() {
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

	// Register callback for status sync on tRuntime updates
	if err := r.syncer.SubscribeRuntimeUpdates(r.SyncStatus); err != nil {
		log.Error(err, "Unable to run status polling")
		return err
	}

	// We sync entities in service from tStorage every `launchPeriod` seconds
	t := time.NewTicker(r.launchPeriod)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = r.SyncSpecs(ctx)
			if err != nil {
				log.Error(err, "Error while SyncSpecs")
			}
			continue

		case <-ctx.Done():
			log.Info(fmt.Sprintf("Cancellation signal was received in %v", r.String()))
		}
		break
	}

	return nil
}
