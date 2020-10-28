package deployment_test

import (
	. "github.com/odahu/odahu-flow/packages/operator/pkg/deployment" //nolint
	"github.com/stretchr/testify/assert"
	"testing"
)

const roleName = "TestRole"

func TestReadDefaultPoliciesAndRender(t *testing.T) {
	data, err := ReadDefaultPoliciesAndRender(roleName)
	assert.NoError(t, err)

	assert.Len(t, data, 3)

	assert.Contains(t, data, "core.rego")
	assert.Contains(t, data, "mapper.rego")
	assert.Contains(t, data, "roles.rego")

	assert.Contains(t, data["core.rego"], roleName)
}