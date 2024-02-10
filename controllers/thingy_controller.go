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

package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	crapv1 "github.com/tvs/cou-test/api/v1"
)

// ThingyReconciler reconciles a Thingy object
type ThingyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crap.tvs.io,resources=thingies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crap.tvs.io,resources=thingies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crap.tvs.io,resources=thingies/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ThingyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	thingy := &crapv1.Thingy{}

	if err := r.Client.Get(ctx, req.NamespacedName, thingy); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
	}

	// CreateOrUpdate a ConfigMap with nothing of value inside
	cm := &corev1.ConfigMap{}
	cm.Name = thingy.Name
	cm.Namespace = thingy.Namespace

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.Data = map[string]string{
			"poop": "dook",
		}

		return nil
	})

	if err != nil {
		logger.Error(err, "Thingy reconcile failed")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ThingyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crapv1.Thingy{}).
		Watches(
			source.NewKindWithCache(&metav1.PartialObjectMetadata{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
			}, mgr.GetCache()),
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}
