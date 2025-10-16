package resources

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

var Deployment = KubeResource{
	Scheme: appsv1.SchemeGroupVersion.WithResource("deployment"),

	Key: func(obj interface{}) string {
		deployment := obj.(*appsv1.Deployment)
		return deployment.Name
	},

	Params: func(obj interface{}, now time.Time) []string {
		deployment := obj.(*appsv1.Deployment)
		return []string{
			deployment.Name,
			fmt.Sprintf("%v/%v", deployment.Status.AvailableReplicas, deployment.Status.Replicas),
		}
	},
}
