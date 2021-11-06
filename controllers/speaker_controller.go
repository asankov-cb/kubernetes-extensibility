/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	istaconv1 "github.com/asankov-cb/kubernetes-extensibility/api/v1"
)

// SpeakerReconciler reconciles a Speaker object
type SpeakerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=istacon.org,resources=speakers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=istacon.org,resources=speakers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=istacon.org,resources=speakers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SpeakerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// once we are in this method the Speaker object is already created in the database

	log := log.FromContext(ctx)

	// even though the object is already created, we do not get it as an argument.
	// instead we get its namespace and name(which Kubernetes uses to identify uniquely, so basically an id)
	// and we can use that info to fetch it from the database
	speaker := istaconv1.Speaker{}
	if err := r.Get(ctx, req.NamespacedName, &speaker); err != nil {
		log.Error(err, "error while fetching Speaker for database")
		// returning an error here, means that Kubernetes will requeue that request and this logic will be retried
		return ctrl.Result{}, err
	}

	log.Info("Created Speaker", "id", req.Name, "name", speaker.Spec.Name)

	// once we have the Speaker here we can use it to run some domain logic

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpeakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&istaconv1.Speaker{}).
		Complete(r)
}
