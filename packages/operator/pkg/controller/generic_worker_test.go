package controller_test

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types/mocks"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestID = "TestID"

type tAsserts struct {
	DeleteInDBCalled      bool
	CreateInRuntimeCalled bool
	UpdateInRuntimeCalled bool
	DeleteInRuntimeCalled bool
}

type tRuntime struct {
	exists bool
	spec uint64
	deleting bool
}

type tStorage struct {
	exists bool
	spec uint64
	delMark bool
	isFinished bool
}

type TestData struct {
	asserts tAsserts

	storage tStorage
	runtime tRuntime
}

func initMocks(td TestData) (*mocks.RuntimeAdapter, *mocks.StorageEntity, *mocks.RuntimeEntity) {

	se := new(mocks.StorageEntity)
	se.On("GetID").Return(TestID)
	se.On("GetSpecHash").Return(td.storage.spec, nil)
	se.On("HasDeletionMark").Return(td.storage.delMark)
	se.On("IsFinished").Return(td.storage.isFinished)
	se.On("DeleteInDB").Return(nil)
	se.On("UpdateInRuntime").Return(nil)
	se.On("CreateInRuntime").Return(nil)


	re := new(mocks.RuntimeEntity)
	re.On("GetID").Return(TestID)
	re.On("GetSpecHash").Return(td.runtime.spec, nil)
	re.On("IsDeleting").Return(td.runtime.deleting)
	re.On("Delete").Return(nil)

	a := new(mocks.RuntimeAdapter)

	if td.storage.exists {
		a.On("GetFromStorage", TestID).Return(se, nil)
		a.On("ListStorage").Return([]types.StorageEntity{se}, nil)
	} else {
		a.On("GetFromStorage", TestID).Return(nil, odahu_errs.NotFoundError{Entity: TestID})
		a.On("ListStorage").Return([]types.StorageEntity{}, nil)
	}

	if td.runtime.exists {
		a.On("GetFromRuntime", TestID).Return(re, nil)
		a.On("ListRuntime").Return([]types.RuntimeEntity{re}, nil)
	} else {
		a.On("GetFromRuntime", TestID).Return(nil, odahu_errs.NotFoundError{Entity: TestID})
		a.On("ListRuntime").Return([]types.RuntimeEntity{}, nil)
	}

	return a, se, re
}

func TestGenericWorker_SyncSpecs(t *testing.T) {

	as := assert.New(t)

	for i, td := range []TestData{
		// Cases about no actions are required
		{
			asserts: tAsserts{},
			storage: tStorage{exists: false},
			runtime: tRuntime{exists: false},
		},
		{
			asserts: tAsserts{},
			storage: tStorage{exists: true},
			runtime: tRuntime{exists: true},

		},
		{  // We do nothing if process in tRuntime is finished (TODO: consider delete it in tRuntime because of garbage)
			asserts: tAsserts{},
			storage: tStorage{exists: true, isFinished: true},
			runtime: tRuntime{exists: false},

		},
		{ // We do nothing if entity in tRuntime is deleting now
			asserts: tAsserts{},
			storage: tStorage{exists: true},
			runtime: tRuntime{exists: true, deleting: true},

		},

		// Cases about some actions are required
		{  // We delete entity in DB if there is deletion mark and corresponding process in tRuntime was already deleted
			asserts: tAsserts{DeleteInDBCalled: true},
			storage: tStorage{exists: true, delMark: true},
			runtime: tRuntime{exists: false},

		},
		{ // We delete process in tRuntime if it is not deleting now but we have deletion mark in tStorage
			asserts: tAsserts{DeleteInRuntimeCalled: true},
			storage: tStorage{exists: true, delMark: true},
			runtime: tRuntime{exists: true},

		},
		{  // We delete zombie process in tRuntime (that does not have corresponding entity in tStorage)
			asserts: tAsserts{DeleteInRuntimeCalled: true},
			storage: tStorage{exists: false},
			runtime: tRuntime{exists: true},

		},
		{  // We create not finished and not marked to delete entities in tRuntime
			asserts: tAsserts{CreateInRuntimeCalled: true},
			storage: tStorage{exists: true, delMark: false, isFinished: false},
			runtime: tRuntime{exists: false},

		},
		{  // We update in tRuntime if spec is changed
			asserts: tAsserts{UpdateInRuntimeCalled: true},
			storage: tStorage{exists: true, spec: 2, delMark: false, isFinished: false},
			runtime: tRuntime{exists: true, spec: 1},
		},
		{  // If process is finished but user change spec we need rerun it
			asserts: tAsserts{UpdateInRuntimeCalled: true},
			storage: tStorage{exists: true, isFinished: true, spec: 2},
			runtime: tRuntime{exists: true, spec: 1},
		},

	} {

		adapter, se, re := initMocks(td)
		worker := controller.NewGenericWorker("", time.Hour, adapter)

		as.NoError(worker.SyncSpecs(context.TODO()))

		ok := true

		if td.asserts.DeleteInDBCalled {
			ok = ok && se.AssertCalled(t, "DeleteInDB")
		} else {
			ok = ok && se.AssertNotCalled(t, "DeleteInDB")
		}
		if td.asserts.CreateInRuntimeCalled {
			ok = ok && se.AssertCalled(t, "CreateInRuntime")
		} else {
			ok = ok && se.AssertNotCalled(t, "CreateInRuntime")
		}
		if td.asserts.UpdateInRuntimeCalled {
			ok = ok && se.AssertCalled(t, "UpdateInRuntime")
		} else {
			ok = ok && se.AssertNotCalled(t, "UpdateInRuntime")
		}
		if td.asserts.DeleteInRuntimeCalled {
			ok = ok && re.AssertCalled(t, "Delete")
		} else {
			ok = ok && re.AssertNotCalled(t, "Delete")
		}

		if !ok {
			t.Logf("#%d Test Data\nStorage: %+v\nRuntime: %+v\n\n\n",  i, td.storage, td.runtime)
		}


	}

}