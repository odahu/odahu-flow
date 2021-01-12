package servicecatalog_test

import (
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"testing"
)


// Compares jsons after unmarshalling into empty interface{}
// So ignores indents and whitespaces and compare only keys and values
func assertJSONEqual(t *testing.T, rawExpected []byte, rawActual []byte) {
	var raw1Interface interface{}
	err := json.Unmarshal(rawExpected, &raw1Interface)
	assert.NoError(t, err)

	var raw2Interface interface{}
	err = json.Unmarshal(rawActual, &raw2Interface)
	assert.NoError(t, err)

	assert.Equal(t, raw1Interface, raw2Interface)
}

// Compares jsons after unmarshalling into empty interface{}
// So ignores indents and whitespaces and compare only keys and values
func assertJSONEqualFile(t *testing.T, rawActual []byte, path string) {
	var rawExpected []byte
	rawExpected, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assertJSONEqual(t, rawExpected, rawActual)
}

func TestPrefixSwaggerUrls(t *testing.T) {

	raw, err := ioutil.ReadFile("testdata/swagger.json")
	assert.NoError(t, err)
	prefixed, err := servicecatalog.PrefixSwaggerUrls("/model/simple", model.Swagger2{
		Raw: raw,
	})
	assert.NoError(t, err)

	assertJSONEqualFile(t, prefixed.Raw, "testdata/prefixed_swagger.json")
}

func TestTagSwaggerMethods(t *testing.T) {

	raw, err := ioutil.ReadFile("testdata/swagger.json")
	assert.NoError(t, err)
	tagged, err := servicecatalog.TagSwaggerMethods([]string{"simple"}, model.Swagger2{
		Raw: raw,
	})
	assert.NoError(t, err)
	assertJSONEqualFile(t, tagged.Raw, "testdata/tagged_swagger.json")
}

func TestProcessSwaggerJSON(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/swagger.json")
	assert.NoError(t, err)

	log, err := zap.NewDevelopment()
	assert.NoError(t, err)
	catalog := servicecatalog.NewModelRouteCatalog(log.Sugar())


	assert.NoError(t, catalog.CreateOrUpdate(servicecatalog.Route{
		ID:     "simple",
		Prefix: "/model/simple",
		Model: model.DeployedModel{
			DeploymentID: "simple",
			ServedModel: model.ServedModel{
				Metadata: model.Metadata{},
				MLServer: "",
				Swagger:  model.Swagger2{Raw: raw},
			},
		},
	}))

	combinedSwagger, err := catalog.ProcessSwaggerJSON()
	assert.NoError(t, err)

	assertJSONEqualFile(t, []byte(combinedSwagger), "testdata/combined_swagger.json")
}
