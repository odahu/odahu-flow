package validation_test

import (
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	. "github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		IDArgs  string
		wantErr bool
	}{
		{
			"Valid iD",
			"some-valid-id",
			false,
		},
		{
			"If ID contains 1 symbol than it is valid",
			"s",
			false,
		},
		{
			"If ID starts with '-' than it is not valid",
			"-some-valid-id",
			true,
		},
		{
			"If ID ends with '-' than it is not valid",
			"some-valid-id-",
			true,
		},
		{
			"Empty ID is not valid",
			"",
			true,
		},
		{
			"If ID contains more 64 symbols, than it is not valid",
			strings.Repeat("a", 64),
			true,
		},
		{
			"Check max length of id",
			strings.Repeat("a", 63),
			false,
		},
		{
			"Check 3 length of id",
			strings.Repeat("a", 3),
			false,
		},
		{
			"If ID contains upper symbols, than it is not valid",
			"SOME-id",
			true,
		},
		{
			"ID starts with number, invalid",
			"123-abc",
			true,
		},
		{
			"ID ends with number, valid",
			"abc-123",
			false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateID(tt.IDArgs); (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestValidateResources(t *testing.T) {

	const memTpl = "Invalid value: \"%s\": must be less than or equal to memory limit"
	const cpuTpl = "Invalid value: \"%s\": must be less than or equal to cpu limit"

	as := assert.New(t)

	testData := []struct {
		memoryRequests string
		memoryLimits string
		cpuRequests string
		cpuLimits string

		want []string
	}{
		{
			memoryRequests: "128Mi", memoryLimits:   "120Mi",
			cpuRequests:    "120m", cpuLimits:      "100m",
			want: []string{
				fmt.Sprintf(memTpl, "128Mi"),
				fmt.Sprintf(cpuTpl, "120m"),
			},
		},
		{
			memoryRequests: "128Mi", memoryLimits:   "128Mi",
			cpuRequests:    "128m", cpuLimits:      "128m",
			want: []string{},
		},
		{
			memoryRequests: "128Mi", memoryLimits:   "130Mi",
			cpuRequests:    "120m", cpuLimits:      "130m",
			want: []string{},
		},
	}

	for _, test := range testData {

		res := &odahuflowv1alpha1.ResourceRequirements{
			Requests: &odahuflowv1alpha1.ResourceList{
				CPU:    &test.cpuRequests,
				Memory: &test.memoryRequests,
			},
			Limits: &odahuflowv1alpha1.ResourceList{
				CPU:    &test.cpuLimits,
				Memory: &test.memoryLimits,
			},
		}

		err := ValidateResources(res, config.NvidiaResourceName)
		if len(test.want) > 0 {
			as.Error(err)
			errStr := err.Error()

			for _, expectedSubstr := range test.want {
				as.Contains(errStr, expectedSubstr)
			}

		} else {
			as.NoError(err)
		}

	}

}