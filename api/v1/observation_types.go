/*
Copyright 2025.

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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/kubernetes/pkg/apis/networking"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type HTTPRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	// Header http.Header `json:"header"`
}

type HTTPResponse struct {
	Status string `json:"status,omitempty"`
	Body   string `json:"body,omitempty"`
}

type HTTPEndpoint struct {
	Request  HTTPRequest  `json:"request"`
	Response HTTPResponse `json:"response,omitempty"`
}

// ObservationSpec defines the desired state of Observation
type ObservationSpec struct {
	// +optional
	DaemonSets appsv1.DaemonSetList `json:"daemonSets,omitempty"`

	// +optional
	Deployments appsv1.DeploymentList `json:"deployments,omitempty"`

	// +optional
	StatefulSets appsv1.StatefulSetList `json:"statefulSet,omitempty"`

	// +optional
	Pods corev1.PodList `json:"pods,omitempty"`

	// +optional
	Services corev1.ServiceList `json:"service,omitempty"`

	// +optional
	PersistentVolumes corev1.PersistentVolumeList `json:"persistentVolume,omitempty"`

	// +optional
	Ingress corev1.ObjectReference `json:"ingress,omitempty"`

	// +optional
	HTTPEnpoints []HTTPEndpoint `json:"httpEndpoints,omitempty"`
}

// ObservationStatus defines the observed state of Observation.
type ObservationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the Observation resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Observation is the Schema for the observations API
type Observation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Observation
	// +required
	Spec ObservationSpec `json:"spec"`

	// status defines the observed state of Observation
	// +optional
	Status ObservationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ObservationList contains a list of Observation
type ObservationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Observation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Observation{}, &ObservationList{})
}
