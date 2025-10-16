package resources

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	appsv1 "k8s.io/api/apps/v1"
)

var Deployment = KubeResource{
	Scheme: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployment"},

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
