/*
Copyright 2023.

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

	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1alpha1 "github.com/Azanul/adapter-operator/api/v1alpha1"
	adapterpkg "github.com/Azanul/adapter-operator/pkg"
	"github.com/go-logr/logr"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/layer5io/meshkit/utils"
)

// AdapterReconciler reconciles a Adapter object
type AdapterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=api.my.domain,resources=adapters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.my.domain,resources=adapters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.my.domain,resources=adapters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AdapterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log

	log = log.WithValues("controller", "Adapter")
	log = log.WithValues("namespace", req.NamespacedName)
	log.Info("Reconciling Adapter")

	baseResource := &apiv1alpha1.Adapter{}

	// Check if resource exists
	err := r.Get(ctx, req.NamespacedName, baseResource)
	if err != nil {
		if kubeerror.IsNotFound(err) {
			baseResource.Name = req.Name
			baseResource.Namespace = req.Namespace
			return r.reconcileAdapter(ctx, false, baseResource, req)
		}
		return ctrl.Result{}, err
	}

	// Check if Adapter controller running
	result, err := r.reconcileAdapter(ctx, true, baseResource, req)
	if err != nil {
		return ctrl.Result{}, ErrReconcileAdapter(err)
	}

	// Patch the adapter resource
	patch, err := utils.Marshal(baseResource)
	if err != nil {
		return ctrl.Result{}, ErrUpdateAdapter(err)
	}

	err = r.Status().Patch(ctx, baseResource, client.RawPatch(types.MergePatchType, []byte(patch)))
	if err != nil {
		return ctrl.Result{}, ErrUpdateAdapter(err)
	}

	return result, nil
}

func (r *AdapterReconciler) reconcileAdapter(ctx context.Context, enable bool, baseResource *apiv1alpha1.Adapter, req ctrl.Request) (ctrl.Result, error) {
	object := adapterpkg.GetObjects(baseResource)[adapterpkg.ServerObject]
	err := r.Get(ctx,
		types.NamespacedName{
			Name:      baseResource.Name,
			Namespace: baseResource.Namespace,
		},
		object,
	)

	if err != nil && kubeerror.IsNotFound(err) && enable {
		r.Log.Info("Creating")
		_ = util.SetControllerReference(baseResource, object, r.Scheme)
		er := r.Create(ctx, object)
		if er != nil {
			return ctrl.Result{}, ErrCreateAdapter(er)
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil && enable {
		return ctrl.Result{}, ErrGetAdapter(err)
	} else if err == nil {
		if enable {
			r.Log.Info("Updating")
			er := r.Update(ctx, object)
			if er != nil {
				return ctrl.Result{}, ErrUpdateAdapter(er)
			}
		} else {
			r.Log.Info("Updating")
			er := r.Delete(ctx, object)
			if er != nil {
				return ctrl.Result{}, ErrDeletingAdapter(er)
			}
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AdapterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Adapter{}).
		Complete(r)
}
