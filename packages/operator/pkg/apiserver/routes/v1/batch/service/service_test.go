/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package service_test

import (
	"github.com/gin-gonic/gin"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	batch "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/batch/service"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/batch/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	router := gin.Default()
	service := &mocks.Service{}
	service.On("Get", mock.Anything, "tf-predictor").Return(api_types.InferenceService{
		ID:           "tf-predictor",
		DeletionMark: false,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		Spec:         api_types.InferenceServiceSpec{},
		Status:       api_types.InferenceServiceStatus{},
	}, nil)
	batch.SetupRoutes(router, service)


	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, strings.Replace(batch.GetURL, ":id", "tf-predictor", -1), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}
