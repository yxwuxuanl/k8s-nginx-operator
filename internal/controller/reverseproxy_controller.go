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

package controller

import (
	"context"
	nginxv1 "github.com/yxwuxuanl/k8s-nginx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ReverseProxyReconciler reconciles a ReverseProxy object
type ReverseProxyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=nginx.lin2ur.cn,resources=reverseproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nginx.lin2ur.cn,resources=reverseproxies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nginx.lin2ur.cn,resources=reverseproxies/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps;services,verbs=*
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=*
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReverseProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *ReverseProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reverseProxy := &nginxv1.ReverseProxy{}

	if err := r.Get(ctx, req.NamespacedName, reverseProxy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err := updateObject(ctx, r.Client, r.Scheme, reverseProxy)
	defer updateStatus(
		ctx,
		r.Client,
		reverseProxy,
		&err,
		func(n *nginxv1.ReverseProxy, condition metav1.Condition) {
			meta.SetStatusCondition(&n.Status.Conditions, condition)
		},
	)

	if err != nil {
		r.Recorder.Event(reverseProxy, "Warning", "ReconcileFailed", err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

var deleteOnlyPred = builder.WithPredicates(predicate.Funcs{
	DeleteFunc: func(event event.DeleteEvent) bool {
		return true
	},
	CreateFunc: func(createEvent event.CreateEvent) bool {
		return false
	},
	UpdateFunc: func(updateEvent event.UpdateEvent) bool {
		return false
	},
	GenericFunc: func(genericEvent event.GenericEvent) bool {
		return false
	},
})

var createOrUpdatePred = builder.WithPredicates(predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
	},
	CreateFunc:  func(e event.CreateEvent) bool { return true },
	DeleteFunc:  func(e event.DeleteEvent) bool { return false },
	GenericFunc: func(e event.GenericEvent) bool { return false },
})

// SetupWithManager sets up the controller with the Manager.
func (r *ReverseProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Recorder = mgr.GetEventRecorderFor("nginx-operator")

	enqueueForOwner := handler.EnqueueRequestForOwner(
		mgr.GetScheme(),
		mgr.GetRESTMapper(),
		&nginxv1.ReverseProxy{},
	)

	return ctrl.NewControllerManagedBy(mgr).
		For(&nginxv1.ReverseProxy{}, createOrUpdatePred).
		Watches(&corev1.ConfigMap{}, enqueueForOwner, deleteOnlyPred).
		Watches(&corev1.Service{}, enqueueForOwner, deleteOnlyPred).
		Watches(&appsv1.Deployment{}, enqueueForOwner, deleteOnlyPred).
		Complete(r)
}
