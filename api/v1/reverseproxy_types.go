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
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deepEqual = equality.Semantic.DeepEqual

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReverseProxySpec defines the desired state of ReverseProxy
type ReverseProxySpec struct {
	NginxSpec `json:",inline"`

	nginx.CommonConfig       `json:",inline"`
	nginx.ReverseProxyConfig `json:",inline"`
}

const Reconciled = "Reconciled"

// ReverseProxyStatus defines the observed state of ReverseProxy
type ReverseProxyStatus struct {
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName="ngxpxy"
//+kubebuilder:printcolumn:name="Reconciled",type="string",JSONPath=".status.conditions[?(@.type == 'Reconciled')].status"
//+kubebuilder:printcolumn:name="ProxyPass",type="string",JSONPath=".spec.proxyPass"
//+kubebuilder:printcolumn:name="NginxImage",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Replicas",type="number",JSONPath=".spec.replicas"

// ReverseProxy is the Schema for the reverseproxies API
type ReverseProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReverseProxySpec   `json:"spec,omitempty"`
	Status ReverseProxyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReverseProxyList contains a list of ReverseProxy
type ReverseProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReverseProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReverseProxy{}, &ReverseProxyList{})
}

func (rp *ReverseProxy) GetNginxSpec() NginxSpec {
	return rp.Spec.NginxSpec
}

func (rp *ReverseProxy) GetNginxConfig() nginx.Config {
	return nginx.Config{
		CommonConfig:       rp.Spec.CommonConfig,
		ReverseProxyConfig: &rp.Spec.ReverseProxyConfig,
	}
}

func (rp *ReverseProxy) GetNamePrefix() string {
	return "ngx-proxy-"
}
