package v1

import corev1 "k8s.io/api/core/v1"

type NginxSpec struct {
	Image string `json:"image,omitempty"`

	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:default:=80
	ServicePort int32 `json:"servicePort,omitempty"`

	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	PodLabels map[string]string `json:"podLabels,omitempty"`

	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	Env []corev1.EnvVar `json:"env,omitempty"`
}
