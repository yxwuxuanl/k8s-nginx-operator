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
	"fmt"
	"github.com/yxwuxuanl/k8s-nginx-operator/internal/nginx"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var tcpreverseproxylog = logf.Log.WithName("tcpreverseproxy-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *TCPReverseProxy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-nginx-lin2ur-cn-v1-tcpreverseproxy,mutating=true,failurePolicy=fail,sideEffects=None,groups=nginx.lin2ur.cn,resources=tcpreverseproxies,verbs=create;update,versions=v1,name=mtcpreverseproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &TCPReverseProxy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *TCPReverseProxy) Default() {
	if !r.DeletionTimestamp.IsZero() {
		return
	}

	tcpreverseproxylog.Info("default", "name", r.Name)

	if r.Spec.Image == "" {
		r.Spec.Image = os.Getenv("DEFAULT_NGINX_IMAGE")
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-nginx-lin2ur-cn-v1-tcpreverseproxy,mutating=false,failurePolicy=fail,sideEffects=None,groups=nginx.lin2ur.cn,resources=tcpreverseproxies,verbs=create;update,versions=v1,name=vtcpreverseproxy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &TCPReverseProxy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TCPReverseProxy) ValidateCreate() (admission.Warnings, error) {
	if !r.DeletionTimestamp.IsZero() {
		return nil, nil
	}

	tcpreverseproxylog.Info("validate create", "name", r.Name)

	if r.Spec.Image == "" {
		return nil, field.Invalid(
			field.NewPath("spec").Child("image"),
			"",
			"not be empty",
		)
	}

	config, err := nginx.BuildConfig(r.GetNginxConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to build nginx config: %w", err)
	}

	if err := nginx.TestConfig(config); err != nil {
		return nil, fmt.Errorf("bad nginx config: %w", err)
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TCPReverseProxy) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	if !r.DeletionTimestamp.IsZero() {
		return nil, nil
	}

	tcpreverseproxylog.Info("validate update", "name", r.Name)

	return r.ValidateCreate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TCPReverseProxy) ValidateDelete() (admission.Warnings, error) {
	tcpreverseproxylog.Info("validate delete", "name", r.Name)

	return nil, nil
}
