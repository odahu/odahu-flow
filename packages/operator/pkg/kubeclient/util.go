package kubeclient

import (
	"fmt"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"reflect"

	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	resLisRequestName = "request"
	resLisLimitName   = "limit"
)

func TransformFilter(filter interface{}, tagKey string) (selector labels.Selector, err error) {
	if filter == nil {
		return nil, nil
	}

	labelSelector := labels.NewSelector()

	elem := reflect.ValueOf(filter).Elem()
	for i := 0; i < elem.NumField(); i++ {
		value := elem.Field(i).Interface().([]string)
		var operator selection.Operator

		switch {
		case len(value) > 1:
			operator = selection.In
		case len(value) == 1:
			if value[0] == "*" {
				continue
			}
			operator = selection.Equals
		default:
			continue
		}

		name := elem.Type().Field(i).Tag.Get(tagKey)
		requirement, err := labels.NewRequirement(name, operator, value)
		if err != nil {
			return nil, err
		}

		labelSelector = labelSelector.Add(*requirement)
	}

	return labelSelector, nil
}

func IsResourcePresent(value *string) bool {
	return value != nil && *value != ""
}

func ConvertOdahuflowResourcesToK8s(
	odahuResources *odahuv1alpha1.ResourceRequirements, gpuResourceName string) (
	k8sResources corev1.ResourceRequirements, err error,
) {
	if odahuResources == nil {
		return k8sResources, err
	}

	var validationError error
	k8sResources.Requests, validationError = convertResourceList(odahuResources.Requests, resLisRequestName)
	err = multierr.Append(err, validationError)
	k8sResources.Limits, validationError = convertResourceList(odahuResources.Limits, resLisLimitName)
	err = multierr.Append(err, validationError)

	// Read more about GPU resources https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/#using-device-plugins
	// From the link above:
	//   * You can specify GPU limits without specifying requests because Kubernetes
	//     will use the limit as the request value by default.
	//   * You can specify GPU in both limits and requests but these two values must be equal.
	//   * You cannot specify GPU requests without specifying limits.
	// So we required from user to specify limits or requests. Limits has priority over requests.
	// We fill in only GPU limits in the k8s resource struct.
	if odahuResources.Limits != nil && IsResourcePresent(odahuResources.Limits.GPU) {
		err = multierr.Append(err, convertResource(
			odahuResources.Limits.GPU, k8sResources.Limits,
			corev1.ResourceName(gpuResourceName),
			"gpu", resLisLimitName,
		))
	} else if odahuResources.Requests != nil && IsResourcePresent(odahuResources.Requests.GPU) {
		if odahuResources.Limits == nil {
			k8sResources.Limits = corev1.ResourceList{}
		}

		err = multierr.Append(err, convertResource(
			odahuResources.Requests.GPU, k8sResources.Limits,
			corev1.ResourceName(gpuResourceName),
			"gpu", resLisRequestName,
		))
	}

	return k8sResources, err
}

func convertResourceList(odahuResourceList *odahuv1alpha1.ResourceList, resourceListType string) (
	k8sResourcesList corev1.ResourceList, err error,
) {
	if odahuResourceList != nil {
		k8sResourcesList = corev1.ResourceList{}
		err = multierr.Append(
			err, convertResource(odahuResourceList.CPU, k8sResourcesList, corev1.ResourceCPU,
				"cpu", resourceListType,
			))
		err = multierr.Append(err, convertResource(
			odahuResourceList.Memory, k8sResourcesList, corev1.ResourceMemory,
			"memory", resourceListType,
		))
	}

	return k8sResourcesList, err
}

func convertResource(odahuResource *string, k8sResource corev1.ResourceList,
	resourceName corev1.ResourceName, resourceType string, resourceListType string) (err error) {
	if IsResourcePresent(odahuResource) {
		var validationErr error
		k8sResource[resourceName], validationErr = resource.ParseQuantity(*odahuResource)

		if validationErr != nil {
			err = multierr.Append(err, fmt.Errorf(
				"validation of %s %s is failed: %s",
				resourceType, resourceListType, validationErr.Error(),
			))
		}
	}

	return err
}

func ConvertK8sErrToOdahuflowErr(err error) error {
	errStatus, ok := err.(*errors.StatusError)
	if ok {
		switch errStatus.Status().Code {
		case http.StatusNotFound:
			return odahu_errs.NotFoundError{Entity: errStatus.ErrStatus.Details.Name}
		case http.StatusConflict:
			return odahu_errs.AlreadyExistError{Entity: errStatus.ErrStatus.Details.Name}
		}
	}

	return nil
}
