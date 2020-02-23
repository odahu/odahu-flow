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

package kubernetes

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	train_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/training"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"
)

var (
	reqGPU        = "1"
	reqCPU        = "2"
	reqMem        = "3"
	limitGPU      = "4"
	limitCPU      = "5"
	limitMem      = "6"
	emptyResource = ""
)

func TestConvertOdahuflowResourcesToK8s(t *testing.T) {
	var ResourceGPU = v1.ResourceName(viper.GetString(train_conf.ResourceGPUName))

	g := NewGomegaWithT(t)
	wantReqGPU, err := resource.ParseQuantity(reqGPU)
	g.Expect(err).Should(BeNil())
	wantReqCPU, err := resource.ParseQuantity(reqCPU)
	g.Expect(err).Should(BeNil())
	wantReqMem, err := resource.ParseQuantity(reqMem)
	g.Expect(err).Should(BeNil())
	wantLimitGPU, err := resource.ParseQuantity(limitGPU)
	g.Expect(err).Should(BeNil())
	wantLimitCPU, err := resource.ParseQuantity(limitCPU)
	g.Expect(err).Should(BeNil())
	wantLimitMem, err := resource.ParseQuantity(limitMem)
	g.Expect(err).Should(BeNil())

	tests := []struct {
		name             string
		requirements     *v1alpha1.ResourceRequirements
		wantDepResources v1.ResourceRequirements
		wantErr          bool
	}{
		{
			name: "Only requirements",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					GPU:    &reqGPU,
					CPU:    &reqCPU,
					Memory: &reqMem,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU:    wantReqCPU,
					v1.ResourceMemory: wantReqMem,
				},
				Limits: v1.ResourceList{
					ResourceGPU: wantReqGPU,
				},
			},
		},
		{
			name: "Only limits",
			requirements: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					GPU:    &limitGPU,
					CPU:    &limitCPU,
					Memory: &limitMem,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					ResourceGPU:       wantLimitGPU,
					v1.ResourceCPU:    wantLimitCPU,
					v1.ResourceMemory: wantLimitMem,
				},
			},
		},
		{
			name: "Limits and requirements",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					GPU:    &reqGPU,
					CPU:    &reqCPU,
					Memory: &reqMem,
				},
				Limits: &v1alpha1.ResourceList{
					GPU:    &limitGPU,
					CPU:    &limitCPU,
					Memory: &limitMem,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU:    wantReqCPU,
					v1.ResourceMemory: wantReqMem,
				},
				Limits: v1.ResourceList{
					ResourceGPU:       wantLimitGPU,
					v1.ResourceCPU:    wantLimitCPU,
					v1.ResourceMemory: wantLimitMem,
				},
			},
		},
		{
			name: "If GPU limits is nill then apply GPU requests",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					GPU: &reqGPU,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{},
				Limits: v1.ResourceList{
					ResourceGPU: wantReqGPU,
				},
			},
		},
		{
			name: "Only memory limits",
			requirements: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					Memory: &limitMem,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceMemory: wantLimitMem,
				},
			},
		},
		{
			name: "Only cpu requests",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					CPU: &reqCPU,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU: wantReqCPU,
				},
			},
		},
		{
			name: "Only GPU requests",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					GPU: &reqGPU,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{},
				Limits: v1.ResourceList{
					ResourceGPU: wantReqGPU,
				},
			},
		},
		{
			name:             "Empty resources",
			requirements:     &v1alpha1.ResourceRequirements{},
			wantDepResources: v1.ResourceRequirements{},
		},
		{
			name: "Every resources is empty string",
			requirements: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					GPU:    &emptyResource,
					CPU:    &emptyResource,
					Memory: &emptyResource,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Limits: v1.ResourceList{},
			},
		},
		{
			name: "Every resources is nil",
			requirements: &v1alpha1.ResourceRequirements{
				Requests: &v1alpha1.ResourceList{
					GPU:    nil,
					CPU:    nil,
					Memory: nil,
				},
			},
			wantDepResources: v1.ResourceRequirements{
				Requests: v1.ResourceList{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			gotDepResources, err := ConvertOdahuflowResourcesToK8s(tt.requirements)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertOdahuflowResourcesToK8s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDepResources, tt.wantDepResources) {
				t.Errorf("ConvertOdahuflowResourcesToK8s() gotDepResources = %v, want %v", gotDepResources, tt.wantDepResources)
			}
		})
	}
}
