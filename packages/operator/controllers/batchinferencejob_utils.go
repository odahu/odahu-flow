package controllers

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
)

// getBucketNames return bucket names for apitypes.InferenceService in next order:
// inputConnection bucket, outputConnection bucket, modelConnection bucket
func getBucketNames(job *v1alpha1.BatchInferenceJob, connAPI ConnGetter) (
	iBucket string, oBucket string, mBucket string, errs error,
	) {


	conn, err := connAPI.GetConnection(job.Spec.InputConnection)
	if err != nil {
		errs = multierr.Append(errs, err)
	} else {
		iBucket, _, err = rclone.GetBucketAndPath(&conn.Spec)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	conn, err = connAPI.GetConnection(job.Spec.OutputConnection)
	if err != nil {
		errs = multierr.Append(errs, err)
	} else {
		oBucket, _, err = rclone.GetBucketAndPath(&conn.Spec)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	conn, err = connAPI.GetConnection(job.Spec.ModelConnection)
	if err != nil {
		errs = multierr.Append(errs, err)
	} else {
		mBucket, _, err = rclone.GetBucketAndPath(&conn.Spec)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return iBucket, oBucket, mBucket, errs
}


// BatchJobToTaskSpec generate tektoncd TaskSpec based on v1alpha1.BatchInferenceJob
func BatchJobToTaskSpec(job *v1alpha1.BatchInferenceJob,
	connAPI ConnGetter,
	gpuResourceName string, rcloneImage string,
	toolsSecret string,
	toolsImage string,
) (ts *tektonv1beta1.TaskSpec, err error) {

	jobRes, err := kubernetes.ConvertOdahuflowResourcesToK8s(job.Spec.Resources, gpuResourceName)
	if err != nil {
		return ts, fmt.Errorf("unable to convert odahu resources to kubernetes resources: %s", err)
	}

	helpContainerRes := utils.CalculateHelperContainerResources(jobRes, gpuResourceName)

	iBucket, oBucket, mBucket, err := getBucketNames(job, connAPI)
	if err != nil {
		return ts, err
	}


	ts = &tektonv1beta1.TaskSpec{
		Steps:   []tektonv1beta1.Step{
			GetConfigureRCloneStep(
				toolsImage, job.Spec.InputConnection, job.Spec.OutputConnection, job.Spec.ModelConnection, helpContainerRes),
			GetSyncDataStep(rcloneImage, iBucket, job.Spec.InputPath, helpContainerRes),
			GetSyncModelStep(rcloneImage, mBucket, job.Spec.ModelPath, helpContainerRes),
			GetValidateInputStep(toolsImage, helpContainerRes),
			GetLogInputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes),
			GetUserContainer(job.Spec.Image, job.Spec.Command, job.Spec.Args, jobRes),
			GetValidateOutputStep(toolsImage, helpContainerRes),
			GetLogOutputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes),
			GetSyncOutputStep(rcloneImage, oBucket, job.Spec.OutputPath, helpContainerRes),
		},
		Volumes: []corev1.Volume{
			{
				Name: toolsConfigVolume,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: toolsSecret,
					},
				},
			},
		},
	}

	return ts, err
}