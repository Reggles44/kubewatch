package resources

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	corev1 "k8s.io/api/core/v1"
)

var Node = KubeResource{
	Scheme: schema.GroupVersionResource{Group: "core", Version: "v1", Resource: "nod"},

	Key: func(obj interface{}) string {
		node := obj.(*corev1.Node)
		return node.Name
	},

	Params: func(obj interface{}, now time.Time) []string {
		node := obj.(*corev1.Node)
		params := []string{
			node.Name,
			fmt.Sprintf("%v Pods", node.Status.Allocatable.Pods()),
			fmt.Sprintf("%v Cores", node.Status.Allocatable.Cpu()),
			fmt.Sprintf("%v Mem", node.Status.Allocatable.Memory()),
		}
		return params
	},
}
