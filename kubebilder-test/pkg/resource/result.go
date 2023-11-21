package resource

import (
	"context"
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"test.kubebuilder.io/project/api/v1alpha1"
	// import to yaml
	//"k8s.io/apimachinery/pkg/runtime"
	////"k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"bytes"
	//"k8s.io/apimachinery/pkg/types"
)

// SerializeObjectAsJSON 함수는 주어진 Kubernetes 오브젝트를 JSON 형식으로 직렬화합니다.
func SerializeObjectAsJSON(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object) (string, error) {
	// 오브젝트를 Kubernetes API 서버로부터 가져옵니다.
	if err := c.Get(ctx, key, obj); err != nil {
		return "", err
	}

	// Pod 타입의 리소스 처리
	if pod, ok := obj.(*v1.Pod); ok {
		// 이미지와 라벨 정보만 추출
		result := struct {
			Images []string          `json:"images"`
			Labels map[string]string `json:"labels"`
		}{
			Images: extractImagesFromPod(pod),
			Labels: pod.Labels,
		}

		// 결과를 JSON으로 직렬화
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return "", err
		}

		return string(jsonBytes), nil
	}

	// 다른 타입의 리소스에 대한 처리 (필요한 경우)

	return "", nil
}
func extractImagesFromPod(pod *v1.Pod) []string {
	var images []string
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	for _, initContainer := range pod.Spec.InitContainers {
		images = append(images, initContainer.Image)
	}
	return images
}

func GetResult(ctx context.Context, c client.Client) ([]v1alpha1.Result, error) {
	var results []v1alpha1.Result
	eventList := &v1.EventList{} // EventList의 포인터를 생성합니다.

	// EventList를 가져옵니다. 이때 List 메서드를 사용합니다.
	if err := c.List(ctx, eventList); err != nil {
		return nil, err
	}

	return results, nil
}

//func GetResourceYAMLs(ctx context.Context, c client.Client, eventList *v1.EventList) ([]string, error) {
//	var yamlStrings []string
//	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true})
//
//	for _, event := range eventList.Items {
//		obj, err := getResource(ctx, c, event.InvolvedObject)
//		if err != nil {
//			return nil, err
//		}
//
//		yamlBytes, err := serializeToYAML(obj, s)
//		if err != nil {
//			return nil, err
//		}
//
//		yamlStrings = append(yamlStrings, string(yamlBytes))
//	}
//
//	return yamlStrings, nil
//}
//
//func getResource(ctx context.Context, c client.Client, ref v1.ObjectReference) (runtime.Object, error) {
//	obj := &v1.Pod{} // 예시로 Pod를 사용합니다. 실제로는 ref.Kind에 따라 다른 타입을 사용해야 할 수 있습니다.
//	err := c.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ref.Namespace}, obj)
//	return obj, err
//}
//
//func serializeToYAML(obj runtime.Object, s *json.Serializer) ([]byte, error) {
//	var yamlBuffer bytes.Buffer
//	if err := s.Encode(obj, &yamlBuffer); err != nil {
//		return nil, err
//	}
//	return yamlBuffer.Bytes(), nil
//}
