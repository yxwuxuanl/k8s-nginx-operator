package controller

import (
	"context"
	"fmt"
	nginxv1 "github.com/yxwuxuanl/k8s-nginx-operator/api/v1"
	"gomodules.xyz/jsonpatch/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

const Finalizer = "nginx.lin2ur.cn/finalizer"

func updateObject(ctx context.Context, cli client.Client, scheme *runtime.Scheme, ngxObject NgxObject) error {
	objects, err := buildObjects(ngxObject)
	if err != nil {
		log.FromContext(ctx).Error(err, "failed to build objects")
		return err
	}

	for _, object := range objects {
		if err := controllerutil.SetControllerReference(ngxObject, object, scheme); err != nil {
			return err
		}

		if err := createOrUpdate(ctx, cli, object); err != nil {
			return err
		}
	}

	return nil
}

func deleteObject(ctx context.Context, cli client.Client, ngxObject NgxObject) error {
	objects, err := buildObjects(ngxObject)
	if err != nil {
		return err
	}

	for _, object := range objects {
		if err := cli.Delete(ctx, object); err != nil {
			if errors.IsNotFound(err) {
				continue
			}

			return fmt.Errorf("failed to delete object: %w", err)
		}

		log.FromContext(ctx).Info(
			"resource has been deleted",
			"objectKind", fmt.Sprintf("%T", object),
			"objectName", object.GetName(),
		)
	}

	return setFinalizer(ctx, cli, ngxObject, true)
}

func setFinalizer(ctx context.Context, cli client.Client, object client.Object, remove bool) error {
	if controllerutil.ContainsFinalizer(object, Finalizer) {
		if !remove {
			return nil
		}
		controllerutil.RemoveFinalizer(object, Finalizer)
	} else {
		if remove {
			return nil
		}
		controllerutil.AddFinalizer(object, Finalizer)
	}

	var patches []jsonpatch.JsonPatchOperation
	patches = append(patches, jsonpatch.JsonPatchOperation{
		Operation: "replace",
		Path:      "/metadata/finalizers",
		Value:     object.GetFinalizers(),
	})

	jsonPatches, _ := json.Marshal(patches)

	return cli.Patch(
		ctx,
		object,
		client.RawPatch(types.JSONPatchType, jsonPatches),
	)
}

func createOrUpdate(ctx context.Context, cli client.Client, object client.Object) error {
	existing := object.DeepCopyObject().(client.Object)

	l := log.FromContext(ctx,
		"objectKind", fmt.Sprintf("%T", object),
		"objectName", object.GetName(),
	)

	if err := cli.Get(ctx, client.ObjectKeyFromObject(object), existing); err != nil {
		if !errors.IsNotFound(err) {
			l.Error(err, "failed to get resource")
			return err
		}

		if err := cli.Create(ctx, object); err != nil {
			l.Error(err, "failed to create resource")
			return err
		}

		l.Info("resource has been created")
		return nil
	}

	if err := cli.Update(ctx, object); err != nil {
		l.Error(err, "failed to update resource")
		return err
	}

	l.Info("resource has been updated")
	return nil
}

func updateStatus[T client.Object](
	ctx context.Context,
	cli client.Client,
	object T,
	reconcileErr *error,
	mutateFn func(T, metav1.Condition),
) {
	condition := metav1.Condition{
		Type:               nginxv1.Reconciled,
		LastTransitionTime: metav1.NewTime(time.Now()),
		ObservedGeneration: object.GetGeneration(),
	}

	if *reconcileErr != nil {
		condition.Status = metav1.ConditionFalse
		condition.Reason = metav1.StatusFailure
		condition.Message = (*reconcileErr).Error()
	} else {
		condition.Status = metav1.ConditionTrue
		condition.Reason = metav1.StatusSuccess
	}

	newObj := object.DeepCopyObject().(T)
	mutateFn(newObj, condition)

	if err := cli.Status().Update(ctx, newObj); err != nil {
		log.FromContext(ctx).Error(err, "failed to update status")
	}
}
