//
//    Copyright 2019 EPAM Systems
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

package training_test

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

const (
	testNamespace               = "default"
	testMtID                    = "test-mt"
	testMtID1                   = "test-mt-id-1"
	testMtID2                   = "test-mt-id-2"
	testModelVersion1           = "1"
	testModelVersion2           = "2"
	testModelName               = "test_name"
	testVcsReference            = "origin/develop123"
	testToolchainIntegrationID  = "ti"
	testToolchainIntegrationID1 = "ti-1"
	testToolchainIntegrationID2 = "ti-2"
	testMtEntrypoint            = "script.py"
	testMtVCSID                 = "odahu-flow-test"
	testToolchainMtImage        = "toolchain-image-test:123"
	testMtImage                 = "image-test:123"
	testMtReference             = "feat/123"
	testModelNameFilter         = "model_name"
	testModelVersionFilter      = "model_version"
	testMtDataPath              = "data/path"
	testMtOutConn               = "some-output-connection"
	testMtOutConnDefault        = "default-output-connection"
	testMpOutConnNotFound       = "out-conn-not-found"
)

var (
	db *sql.DB
	kubeClient client.Client
	cfg *rest.Config
)

func testMainWrapper(m *testing.M) int {
	// Setup Test DB

	var closeDB func() error
	var err error
	db, _, closeDB, err = testenvs.SetupTestDB()
	defer func() {
		if err := closeDB(); err != nil {
			log.Print("Error during release test DB resources")
		}
	}()
	if err != nil {
		return -1
	}

	var closeKube func() error
	kubeClient, cfg, closeKube, _, err = testenvs.SetupTestKube(
		filepath.Join("..", "..", "..", "..", "..", "config", "crds"),
		filepath.Join("..", "..", "..", "..", "..", "hack", "tests", "thirdparty_crds"),
	)
	defer func() {
		if err := closeKube(); err != nil {
			log.Print("Error during release test Kube Environment resources")
		}
	}()
	if err != nil {
		return -1
	}

	return m.Run()
}

func TestMain(m *testing.M) {
	utils.SetupLogger()

	os.Exit(testMainWrapper(m))
}
