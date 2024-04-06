/*
Copyright 2024.

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
	"github.com/yxwuxuanl/k8s-nginx-operator/internal/nginx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TCPReverseProxySpec defines the desired state of TCPReverseProxy
type TCPReverseProxySpec struct {
	NginxSpec `json:",inline"`

	nginx.CommonConfig          `json:",inline"`
	nginx.TCPReverseProxyConfig `json:",inline"`
}

// TCPReverseProxyStatus defines the observed state of TCPReverseProxy
type TCPReverseProxyStatus struct {
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName="ngxtcppxy"
//+kubebuilder:printcolumn:name="Reconciled",type="string",JSONPath=".status.conditions[?(@.type == 'Reconciled')].status"
//+kubebuilder:printcolumn:name="Servers",type="string",JSONPath=".spec.servers..server"
//+kubebuilder:printcolumn:name="NginxImage",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Replicas",type="number",JSONPath=".spec.replicas"

// TCPReverseProxy is the Schema for the tcpreverseproxies API
type TCPReverseProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TCPReverseProxySpec   `json:"spec,omitempty"`
	Status TCPReverseProxyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TCPReverseProxyList contains a list of TCPReverseProxy
type TCPReverseProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TCPReverseProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TCPReverseProxy{}, &TCPReverseProxyList{})
}

func (rp *TCPReverseProxy) GetNginxSpec() NginxSpec {
	return rp.Spec.NginxSpec
}

func (rp *TCPReverseProxy) GetNginxConfig() nginx.Config {
	return nginx.Config{
		CommonConfig:          rp.Spec.CommonConfig,
		TCPReverseProxyConfig: &rp.Spec.TCPReverseProxyConfig,
	}
}

func (rp *TCPReverseProxy) GetNamePrefix() string {
	return "ngx-tcpproxy-"
}

func (in *TCPReverseProxy) GetServiceProtocol() corev1.Protocol {
	return corev1.Protocol(in.Spec.Protocol)
}
