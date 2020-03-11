package utils

import (
	training_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/training"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	k8s_resource "k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"
)

func TestCalculateHelperContainerResources(t *testing.T) {
	resourceCPU := *k8s_resource.NewQuantity(42, k8s_resource.DecimalSI)
	resourceMemory := *k8s_resource.NewQuantity(911, k8s_resource.DecimalSI)
	resourceGPU := *k8s_resource.NewQuantity(73, k8s_resource.DecimalSI)

	gpuResourceName := corev1.ResourceName(viper.GetString(training_conf.ResourceGPUName))

	tests := []struct {
		name string
		args corev1.ResourceRequirements
		want corev1.ResourceRequirements
	}{
		{
			"Empty Resources",
			corev1.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Override the requests resources",
			corev1.ResourceRequirements{
				Limits: nil,
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: *k8s_resource.NewQuantity(2, k8s_resource.DecimalSI),
					corev1.ResourceCPU:    *k8s_resource.NewQuantity(3, k8s_resource.DecimalSI),
				},
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Empty CPU limit resources",
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resourceCPU,
				},
				Requests: nil,
			},
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resourceCPU,
					corev1.ResourceCPU:    DefaultHelperLimits.Cpu().DeepCopy(),
				},
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Empty Memory limit resources",
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU: resourceCPU,
				},
				Requests: nil,
			},
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resourceCPU,
					corev1.ResourceMemory: DefaultHelperLimits.Memory().DeepCopy(),
				},
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"GPU removing from limit",
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					gpuResourceName: resourceGPU,
				},
				Requests: nil,
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Main workflow",
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resourceCPU,
					corev1.ResourceMemory: resourceMemory,
					gpuResourceName:       resourceGPU,
				},
				Requests: nil,
			},
			corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resourceCPU,
					corev1.ResourceMemory: resourceMemory,
				},
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateHelperContainerResources(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHelperContainerResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
