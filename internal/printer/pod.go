package printer

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func podParams(o *unstructured.Unstructured, now time.Time) ([]string, error) {
	var pod corev1.Pod
	err := convertUnstructured(&pod, o)
	if err != nil {
		return []string{}, err
	}

	return []string{
		pod.Name,
		string(pod.Status.Phase),
		calculateResourceDuration(o, now),
		pod.Status.PodIP,
	}, nil
}
