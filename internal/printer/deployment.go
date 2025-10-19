package printer

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func deploymentParams(o *unstructured.Unstructured, now time.Time) ([]string, error) {
	var deployment appsv1.Deployment
	err := convertUnstructured(&deployment, o)
	if err != nil {
		return []string{}, err
	}

	return []string{
		deployment.Name,
		fmt.Sprintf("%v/%v", deployment.Status.AvailableReplicas, deployment.Status.Replicas),
		calculateResourceDuration(o, now),
	}, nil
}
