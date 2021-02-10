// +build !ignore_autogenerated

/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Connection) DeepCopyInto(out *Connection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Connection.
func (in *Connection) DeepCopy() *Connection {
	if in == nil {
		return nil
	}
	out := new(Connection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Connection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionList) DeepCopyInto(out *ConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Connection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionList.
func (in *ConnectionList) DeepCopy() *ConnectionList {
	if in == nil {
		return nil
	}
	out := new(ConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionSpec) DeepCopyInto(out *ConnectionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionSpec.
func (in *ConnectionSpec) DeepCopy() *ConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(ConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionStatus) DeepCopyInto(out *ConnectionStatus) {
	*out = *in
	if in.SecretName != nil {
		in, out := &in.SecretName, &out.SecretName
		*out = new(string)
		**out = **in
	}
	if in.ServiceAccountName != nil {
		in, out := &in.ServiceAccountName, &out.ServiceAccountName
		*out = new(string)
		**out = **in
	}
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionStatus.
func (in *ConnectionStatus) DeepCopy() *ConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(ConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataBindingDir) DeepCopyInto(out *DataBindingDir) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataBindingDir.
func (in *DataBindingDir) DeepCopy() *DataBindingDir {
	if in == nil {
		return nil
	}
	out := new(DataBindingDir)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DefaultRouteTemplate) DeepCopyInto(out *DefaultRouteTemplate) {
	*out = *in
	if in.Attempts != nil {
		in, out := &in.Attempts, &out.Attempts
		*out = new(int32)
		**out = **in
	}
	if in.PerTryTimeout != nil {
		in, out := &in.PerTryTimeout, &out.PerTryTimeout
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DefaultRouteTemplate.
func (in *DefaultRouteTemplate) DeepCopy() *DefaultRouteTemplate {
	if in == nil {
		return nil
	}
	out := new(DefaultRouteTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariable) DeepCopyInto(out *EnvironmentVariable) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariable.
func (in *EnvironmentVariable) DeepCopy() *EnvironmentVariable {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JsonSchema) DeepCopyInto(out *JsonSchema) {
	*out = *in
	if in.Required != nil {
		in, out := &in.Required, &out.Required
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JsonSchema.
func (in *JsonSchema) DeepCopy() *JsonSchema {
	if in == nil {
		return nil
	}
	out := new(JsonSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeployment) DeepCopyInto(out *ModelDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeployment.
func (in *ModelDeployment) DeepCopy() *ModelDeployment {
	if in == nil {
		return nil
	}
	out := new(ModelDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentList) DeepCopyInto(out *ModelDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentList.
func (in *ModelDeploymentList) DeepCopy() *ModelDeploymentList {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentSpec) DeepCopyInto(out *ModelDeploymentSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		*out = new(int32)
		**out = **in
	}
	if in.MaxReplicas != nil {
		in, out := &in.MaxReplicas, &out.MaxReplicas
		*out = new(int32)
		**out = **in
	}
	if in.LivenessProbeInitialDelay != nil {
		in, out := &in.LivenessProbeInitialDelay, &out.LivenessProbeInitialDelay
		*out = new(int32)
		**out = **in
	}
	if in.ReadinessProbeInitialDelay != nil {
		in, out := &in.ReadinessProbeInitialDelay, &out.ReadinessProbeInitialDelay
		*out = new(int32)
		**out = **in
	}
	if in.RoleName != nil {
		in, out := &in.RoleName, &out.RoleName
		*out = new(string)
		**out = **in
	}
	if in.ImagePullConnectionID != nil {
		in, out := &in.ImagePullConnectionID, &out.ImagePullConnectionID
		*out = new(string)
		**out = **in
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ContainerConcurrency != nil {
		in, out := &in.ContainerConcurrency, &out.ContainerConcurrency
		*out = new(int64)
		**out = **in
	}
	if in.DefaultRoute != nil {
		in, out := &in.DefaultRoute, &out.DefaultRoute
		*out = new(DefaultRouteTemplate)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentSpec.
func (in *ModelDeploymentSpec) DeepCopy() *ModelDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentStatus) DeepCopyInto(out *ModelDeploymentStatus) {
	*out = *in
	if in.LastCredsUpdatedTime != nil {
		in, out := &in.LastCredsUpdatedTime, &out.LastCredsUpdatedTime
		*out = (*in).DeepCopy()
	}
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentStatus.
func (in *ModelDeploymentStatus) DeepCopy() *ModelDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentTarget) DeepCopyInto(out *ModelDeploymentTarget) {
	*out = *in
	if in.Weight != nil {
		in, out := &in.Weight, &out.Weight
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentTarget.
func (in *ModelDeploymentTarget) DeepCopy() *ModelDeploymentTarget {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelIdentity) DeepCopyInto(out *ModelIdentity) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelIdentity.
func (in *ModelIdentity) DeepCopy() *ModelIdentity {
	if in == nil {
		return nil
	}
	out := new(ModelIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelPackaging) DeepCopyInto(out *ModelPackaging) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelPackaging.
func (in *ModelPackaging) DeepCopy() *ModelPackaging {
	if in == nil {
		return nil
	}
	out := new(ModelPackaging)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelPackaging) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelPackagingList) DeepCopyInto(out *ModelPackagingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelPackaging, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelPackagingList.
func (in *ModelPackagingList) DeepCopy() *ModelPackagingList {
	if in == nil {
		return nil
	}
	out := new(ModelPackagingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelPackagingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelPackagingResult) DeepCopyInto(out *ModelPackagingResult) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelPackagingResult.
func (in *ModelPackagingResult) DeepCopy() *ModelPackagingResult {
	if in == nil {
		return nil
	}
	out := new(ModelPackagingResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelPackagingSpec) DeepCopyInto(out *ModelPackagingSpec) {
	*out = *in
	if in.ArtifactName != nil {
		in, out := &in.ArtifactName, &out.ArtifactName
		*out = new(string)
		**out = **in
	}
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]Target, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelPackagingSpec.
func (in *ModelPackagingSpec) DeepCopy() *ModelPackagingSpec {
	if in == nil {
		return nil
	}
	out := new(ModelPackagingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelPackagingStatus) DeepCopyInto(out *ModelPackagingStatus) {
	*out = *in
	if in.ExitCode != nil {
		in, out := &in.ExitCode, &out.ExitCode
		*out = new(int32)
		**out = **in
	}
	if in.Reason != nil {
		in, out := &in.Reason, &out.Reason
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		*out = make([]ModelPackagingResult, len(*in))
		copy(*out, *in)
	}
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelPackagingStatus.
func (in *ModelPackagingStatus) DeepCopy() *ModelPackagingStatus {
	if in == nil {
		return nil
	}
	out := new(ModelPackagingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelRoute) DeepCopyInto(out *ModelRoute) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelRoute.
func (in *ModelRoute) DeepCopy() *ModelRoute {
	if in == nil {
		return nil
	}
	out := new(ModelRoute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelRoute) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelRouteList) DeepCopyInto(out *ModelRouteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelRoute, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelRouteList.
func (in *ModelRouteList) DeepCopy() *ModelRouteList {
	if in == nil {
		return nil
	}
	out := new(ModelRouteList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelRouteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelRouteSpec) DeepCopyInto(out *ModelRouteSpec) {
	*out = *in
	if in.Mirror != nil {
		in, out := &in.Mirror, &out.Mirror
		*out = new(string)
		**out = **in
	}
	if in.ModelDeploymentTargets != nil {
		in, out := &in.ModelDeploymentTargets, &out.ModelDeploymentTargets
		*out = make([]ModelDeploymentTarget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Attempts != nil {
		in, out := &in.Attempts, &out.Attempts
		*out = new(int32)
		**out = **in
	}
	if in.PerTryTimeout != nil {
		in, out := &in.PerTryTimeout, &out.PerTryTimeout
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelRouteSpec.
func (in *ModelRouteSpec) DeepCopy() *ModelRouteSpec {
	if in == nil {
		return nil
	}
	out := new(ModelRouteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelRouteStatus) DeepCopyInto(out *ModelRouteStatus) {
	*out = *in
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelRouteStatus.
func (in *ModelRouteStatus) DeepCopy() *ModelRouteStatus {
	if in == nil {
		return nil
	}
	out := new(ModelRouteStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelTraining) DeepCopyInto(out *ModelTraining) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelTraining.
func (in *ModelTraining) DeepCopy() *ModelTraining {
	if in == nil {
		return nil
	}
	out := new(ModelTraining)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelTraining) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelTrainingList) DeepCopyInto(out *ModelTrainingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelTraining, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelTrainingList.
func (in *ModelTrainingList) DeepCopy() *ModelTrainingList {
	if in == nil {
		return nil
	}
	out := new(ModelTrainingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelTrainingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelTrainingSpec) DeepCopyInto(out *ModelTrainingSpec) {
	*out = *in
	out.Model = in.Model
	if in.CustomEnvs != nil {
		in, out := &in.CustomEnvs, &out.CustomEnvs
		*out = make([]EnvironmentVariable, len(*in))
		copy(*out, *in)
	}
	if in.HyperParameters != nil {
		in, out := &in.HyperParameters, &out.HyperParameters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EntrypointArguments != nil {
		in, out := &in.EntrypointArguments, &out.EntrypointArguments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make([]DataBindingDir, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelTrainingSpec.
func (in *ModelTrainingSpec) DeepCopy() *ModelTrainingSpec {
	if in == nil {
		return nil
	}
	out := new(ModelTrainingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelTrainingStatus) DeepCopyInto(out *ModelTrainingStatus) {
	*out = *in
	if in.ExitCode != nil {
		in, out := &in.ExitCode, &out.ExitCode
		*out = new(int32)
		**out = **in
	}
	if in.Reason != nil {
		in, out := &in.Reason, &out.Reason
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
	if in.Artifacts != nil {
		in, out := &in.Artifacts, &out.Artifacts
		*out = make([]TrainingResult, len(*in))
		copy(*out, *in)
	}
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelTrainingStatus.
func (in *ModelTrainingStatus) DeepCopy() *ModelTrainingStatus {
	if in == nil {
		return nil
	}
	out := new(ModelTrainingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Modifiable) DeepCopyInto(out *Modifiable) {
	*out = *in
	if in.CreatedAt != nil {
		in, out := &in.CreatedAt, &out.CreatedAt
		*out = (*in).DeepCopy()
	}
	if in.UpdatedAt != nil {
		in, out := &in.UpdatedAt, &out.UpdatedAt
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Modifiable.
func (in *Modifiable) DeepCopy() *Modifiable {
	if in == nil {
		return nil
	}
	out := new(Modifiable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagingIntegration) DeepCopyInto(out *PackagingIntegration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagingIntegration.
func (in *PackagingIntegration) DeepCopy() *PackagingIntegration {
	if in == nil {
		return nil
	}
	out := new(PackagingIntegration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackagingIntegration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagingIntegrationList) DeepCopyInto(out *PackagingIntegrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackagingIntegration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagingIntegrationList.
func (in *PackagingIntegrationList) DeepCopy() *PackagingIntegrationList {
	if in == nil {
		return nil
	}
	out := new(PackagingIntegrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackagingIntegrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagingIntegrationSpec) DeepCopyInto(out *PackagingIntegrationSpec) {
	*out = *in
	in.Schema.DeepCopyInto(&out.Schema)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagingIntegrationSpec.
func (in *PackagingIntegrationSpec) DeepCopy() *PackagingIntegrationSpec {
	if in == nil {
		return nil
	}
	out := new(PackagingIntegrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagingIntegrationStatus) DeepCopyInto(out *PackagingIntegrationStatus) {
	*out = *in
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagingIntegrationStatus.
func (in *PackagingIntegrationStatus) DeepCopy() *PackagingIntegrationStatus {
	if in == nil {
		return nil
	}
	out := new(PackagingIntegrationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceList) DeepCopyInto(out *ResourceList) {
	*out = *in
	if in.GPU != nil {
		in, out := &in.GPU, &out.GPU
		*out = new(string)
		**out = **in
	}
	if in.CPU != nil {
		in, out := &in.CPU, &out.CPU
		*out = new(string)
		**out = **in
	}
	if in.Memory != nil {
		in, out := &in.Memory, &out.Memory
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceList.
func (in *ResourceList) DeepCopy() *ResourceList {
	if in == nil {
		return nil
	}
	out := new(ResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequirements) DeepCopyInto(out *ResourceRequirements) {
	*out = *in
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = new(ResourceList)
		(*in).DeepCopyInto(*out)
	}
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = new(ResourceList)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequirements.
func (in *ResourceRequirements) DeepCopy() *ResourceRequirements {
	if in == nil {
		return nil
	}
	out := new(ResourceRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchemaValidation) DeepCopyInto(out *SchemaValidation) {
	*out = *in
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]TargetSchema, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Arguments.DeepCopyInto(&out.Arguments)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchemaValidation.
func (in *SchemaValidation) DeepCopy() *SchemaValidation {
	if in == nil {
		return nil
	}
	out := new(SchemaValidation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Target) DeepCopyInto(out *Target) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Target.
func (in *Target) DeepCopy() *Target {
	if in == nil {
		return nil
	}
	out := new(Target)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetSchema) DeepCopyInto(out *TargetSchema) {
	*out = *in
	if in.ConnectionTypes != nil {
		in, out := &in.ConnectionTypes, &out.ConnectionTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetSchema.
func (in *TargetSchema) DeepCopy() *TargetSchema {
	if in == nil {
		return nil
	}
	out := new(TargetSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolchainIntegration) DeepCopyInto(out *ToolchainIntegration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolchainIntegration.
func (in *ToolchainIntegration) DeepCopy() *ToolchainIntegration {
	if in == nil {
		return nil
	}
	out := new(ToolchainIntegration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ToolchainIntegration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolchainIntegrationList) DeepCopyInto(out *ToolchainIntegrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ToolchainIntegration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolchainIntegrationList.
func (in *ToolchainIntegrationList) DeepCopy() *ToolchainIntegrationList {
	if in == nil {
		return nil
	}
	out := new(ToolchainIntegrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ToolchainIntegrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolchainIntegrationSpec) DeepCopyInto(out *ToolchainIntegrationSpec) {
	*out = *in
	if in.AdditionalEnvironments != nil {
		in, out := &in.AdditionalEnvironments, &out.AdditionalEnvironments
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolchainIntegrationSpec.
func (in *ToolchainIntegrationSpec) DeepCopy() *ToolchainIntegrationSpec {
	if in == nil {
		return nil
	}
	out := new(ToolchainIntegrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolchainIntegrationStatus) DeepCopyInto(out *ToolchainIntegrationStatus) {
	*out = *in
	in.Modifiable.DeepCopyInto(&out.Modifiable)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolchainIntegrationStatus.
func (in *ToolchainIntegrationStatus) DeepCopy() *ToolchainIntegrationStatus {
	if in == nil {
		return nil
	}
	out := new(ToolchainIntegrationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TrainingResult) DeepCopyInto(out *TrainingResult) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TrainingResult.
func (in *TrainingResult) DeepCopy() *TrainingResult {
	if in == nil {
		return nil
	}
	out := new(TrainingResult)
	in.DeepCopyInto(out)
	return out
}
