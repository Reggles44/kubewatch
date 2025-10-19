package printer

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func nodeParams(o *unstructured.Unstructured, now time.Time) ([]string, error) {
	var node corev1.Node
	err := convertUnstructured(&node, o)
	if err != nil {
		return []string{}, err
	}

	return []string{
		node.Name,
		string(node.Status.Phase),
		calculateResourceDuration(o, now),
		fmt.Sprintf("%v Pods", node.Status.Allocatable.Pods()),
		fmt.Sprintf("%v Cores", node.Status.Allocatable.Cpu()),
		fmt.Sprintf("%v Mem", node.Status.Allocatable.Memory()),
		// string(node.Status.Addresses[0].Type),
		// node.Status.Addresses[0].Address,
		// string(node.Status.Addresses[1].Type),
		// node.Status.Addresses[1].Address,
		// string(node.Status.Addresses[2].Type),
		// node.Status.Addresses[2].Address,
	}, nil
}
