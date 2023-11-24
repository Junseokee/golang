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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1alpha1 "test.kubebuilder.io/project/api/v1alpha1"
)

// KubegptReconciler reconciles a Kubegpt object
type KubegptReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.test.kubebuilder.io,resources=kubegpts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.test.kubebuilder.io,resources=kubegpts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.test.kubebuilder.io,resources=kubegpts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kubegpt object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *KubegptReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	//kubegptConfig := &corev1alpha1.Kubegpt{}
	//err := r.Get(ctx, req.NamespacedName, kubegptConfig)
	//if err != nil {
	//	return ctrl.Result{}, client.IgnoreNotFound(err)
	//}
	//
	//deployment := v1.Deployment{}
	//err = r.Get(ctx, client.ObjectKey{Namespace: kubegptConfig.Namespace,
	//	Name: "k8sgpt-deployment"}, &deployment)
	//if client.IgnoreNotFound(err) != nil {
	//	return ctrl.Result{}, client.IgnoreNotFound(err)
	//}
	//err = resources.Sync(ctx, r.Client, *k8sgptConfig, resources.SyncOp)
	//if err != nil {
	//	k8sgptReconcileErrorCount.Inc()
	//	return r.finishReconcile(err, false)
	//}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubegptReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Kubegpt{}).
		Complete(r)
}
