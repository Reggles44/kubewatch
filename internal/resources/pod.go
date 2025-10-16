package resources

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var Pod = KubeResource{
	Scheme: corev1.SchemeGroupVersion.WithResource("pods"),

	Key: func(obj interface{}) string {
		pod := obj.(*corev1.Pod)
		return pod.Name
	},

	Params: func(obj interface{}, now time.Time) []string {
		pod := obj.(*corev1.Pod)

		var duration time.Duration
		startTime := pod.Status.StartTime
		if startTime != nil {
			duration = now.Sub(pod.Status.StartTime.Time)
		}

		return []string{
			pod.Name,
			string(pod.Status.Phase),
			fmt.Sprintf("%-8v", duration.Round(time.Second)),
			pod.Status.PodIP,
		}
	},
}
