package resource

import (
	"bytes"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"test.kubebuilder.io/project/api/v1alpha1"
)

func GetResourceYAMLs(ctx context.Context, c client.Client, eventList *v1.EventList) ([]string, error) {
	var yamlStrings []string
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true})

	for _, event := range eventList.Items {
		obj, err := getResource(ctx, c, event.InvolvedObject)
		if err != nil {
			return nil, err
		}

		yamlBytes, err := serializeToYAML(obj, s)
		if err != nil {
			return nil, err
		}

		yamlStrings = append(yamlStrings, string(yamlBytes))
	}

	return yamlStrings, nil
}

func getResource(ctx context.Context, c client.Client, ref v1.ObjectReference) (runtime.Object, error) {
	obj := &v1.Pod{} // 예시로 Pod를 사용합니다. 실제로는 ref.Kind에 따라 다른 타입을 사용해야 할 수 있습니다.
	err := c.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ref.Namespace}, obj)
	return obj, err
}

func serializeToYAML(obj runtime.Object, s *json.Serializer) ([]byte, error) {
	var yamlBuffer bytes.Buffer
	if err := s.Encode(obj, &yamlBuffer); err != nil {
		return nil, err
	}
	return yamlBuffer.Bytes(), nil
}

func GetResult(ctx context.Context, c client.Client) ([]v1alpha1.Result, error) {
	var results []v1alpha1.Result
	eventList := &v1.EventList{} // EventList의 포인터를 생성합니다.

	// EventList를 가져옵니다. 이때 List 메서드를 사용합니다.
	if err := c.List(ctx, eventList); err != nil {
		return nil, err
	}

	for _, event := range eventList.Items {
		result := v1alpha1.Result{
			Spec: v1alpha1.ResultSpec{
				Name:      event.Name,
				Namespace: event.Namespace,
				Kind:      event.InvolvedObject.Kind,
				Message:   event.Message,
			},
		}
		results = append(results, result)
	}

	return results, nil
}
