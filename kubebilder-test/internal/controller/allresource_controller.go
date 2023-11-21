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
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	corev1alpha1 "test.kubebuilder.io/project/api/v1alpha1"
	"test.kubebuilder.io/project/pkg/resource"
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
	var events v1.EventList
	if err := r.Client.List(ctx, &events); err != nil {
		return ctrl.Result{}, err
	}

	for _, event := range events.Items {
		// JSON 형식으로 직렬화할 오브젝트를 설정합니다.
		var pod v1.Pod // 적절한 오브젝트 타입으로 변경

		jsonString, err := resource.SerializeObjectAsJSON(ctx, r.Client, client.ObjectKey{Name: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace}, &pod)
		if err != nil {
			// 에러 처리
			continue
		}

		l.Info("JSON for object", "json", jsonString)
	}

	results, err := resource.GetResult(ctx, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	for _, result := range results {
		l.Info("로그로그 :", "name", result)
	}

	return ctrl.Result{}, nil

	//// 이벤트 목록을 가져옵니다.
	//var events v1.EventList
	//if err := r.Client.List(ctx, &events); err != nil {
	//	return ctrl.Result{}, err
	//}
	//
	//// JSON to YAML 변환을 위한 Serializer 생성
	//s := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true})
	//
	//for _, event := range events.Items {
	//	// 연관된 오브젝트를 가져옵니다.
	//	var obj runtime.Object // 연관된 오브젝트 타입에 따라 변경
	//	if err := r.Client.Get(ctx, client.ObjectKey{Name: event.InvolvedObject.Name, Namespace: event.InvolvedObject.Namespace}, obj); err != nil {
	//		// 오브젝트 가져오기 실패
	//		continue
	//	}
	//
	//	// 오브젝트를 YAML 형태로 변환
	//	yamlBytes, err := serializeToYAML(obj, s)
	//	if err != nil {
	//		// YAML 변환 실패
	//		continue
	//	}
	//
	//	// 여기서 yamlBytes를 사용하여 필요한 작업을 수행합니다.
	//	// 예: 로깅, 저장, CRD 업데이트 등
	//	l.Info("YAML for object", "yaml", string(yamlBytes))
	//}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AllResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Result{}).
		Watches(&v1.Event{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1.Pod{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
