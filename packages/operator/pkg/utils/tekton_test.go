package utils_test

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	. "github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	k8s_resource "k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"
)

func TestCalculateHelperContainerResources(t *testing.T) {
	resourceCPU := *k8s_resource.NewQuantity(42, k8s_resource.DecimalSI)
	resourceMemory := *k8s_resource.NewQuantity(911, k8s_resource.DecimalSI)
	resourceGPU := *k8s_resource.NewQuantity(73, k8s_resource.DecimalSI)

	gpuResourceName := corev1.ResourceName(config.NvidiaResourceName)
	type args struct {
		res             corev1.ResourceRequirements
		gpuResourceName string
	}

	tests := []struct {
		name string
		args args
		want corev1.ResourceRequirements
	}{
		{
			"Empty Resources",
			args{
				corev1.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
				config.NvidiaResourceName,
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Override the requests resources",
			args{
				corev1.ResourceRequirements{
					Limits: nil,
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: *k8s_resource.NewQuantity(2, k8s_resource.DecimalSI),
						corev1.ResourceCPU:    *k8s_resource.NewQuantity(3, k8s_resource.DecimalSI),
					},
				},
				config.NvidiaResourceName,
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Empty CPU limit resources",
			args{
				corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceMemory: resourceCPU,
					},
					Requests: nil,
				},
				config.NvidiaResourceName,
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
			args{
				corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU: resourceCPU,
					},
					Requests: nil,
				},
				config.NvidiaResourceName,
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
			args{
				corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						gpuResourceName: resourceGPU,
					},
					Requests: nil,
				},
				config.NvidiaResourceName,
			},
			corev1.ResourceRequirements{
				Limits:   DefaultHelperLimits.DeepCopy(),
				Requests: EmptyHelperContainerRequestRes.DeepCopy(),
			},
		},
		{
			"Main workflow",
			args{
				corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resourceCPU,
						corev1.ResourceMemory: resourceMemory,
						gpuResourceName:       resourceGPU,
					},
					Requests: nil,
				},
				config.NvidiaResourceName,
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
			if got := CalculateHelperContainerResources(tt.args.res, tt.args.gpuResourceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHelperContainerResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
