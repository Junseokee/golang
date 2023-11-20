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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	corev1alpha1 "test.kubebuilder.io/project/api/v1alpha1"
)

// AllResourceReconciler reconciles a AllResource object
type AllResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=allresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=allresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=allresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AllResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *AllResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	event := &corev1.Event{}
	result := &corev1alpha1.Result{}
	//pod := &corev1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, event); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	//l.Info("pod", "name", pod.Name, "namespace", pod.Namespace)

	result.Spec.Name = event.Name
	result.Spec.Namespace = event.Namespace
	result.Spec.Message = event.Message
	result.Spec.Kind = event.InvolvedObject.Kind
	///*	if err := r.Get(ctx, client.ObjectKey{Namespace: result.Spec.Namespace, Name: result.Name}, result); err != nil {
	//		return ctrl.Result{}, client.IgnoreNotFound(err)
	//	} else {
	//		l.Info("pod",
	//			"pod :", pod)
	//	}*/

	l.Info("Event",
		"이름 :", result.Spec.Name,
		"네임스페이스 :", result.Spec.Namespace,
		"메세지 : ", result.Spec.Message,
		"리소스 :", result.Spec.Kind)

	//result.Spec.Event = append(result.Spec.Event, corev1alpha1.Event{
	//	Name:      event.Name,
	//	Namespace: event.Namespace,
	//	Kind:      event.Kind,
	//	Message:   event.Message,
	//})
	//l.Info("event_log", "이름 :", event.Name, "종류 :", event.Kind, "메세지 :", event.Message)
	//l.Info("Result",
	//	"Result", result.Spec.Event)

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *AllResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Event{}).
		Complete(r)
}
