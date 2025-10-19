package printer

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ParamGenerator func(o *unstructured.Unstructured, now time.Time) ([]string, error)

var paramsFuncs = map[schema.GroupVersionKind]ParamGenerator{
	corev1.SchemeGroupVersion.WithKind("Pod"): podParams,
	corev1.SchemeGroupVersion.WithKind("Node"): nodeParams,
}

func GetParams(t schema.GroupVersionKind, o *unstructured.Unstructured, now time.Time) []string {
	gen, ok := paramsFuncs[t]
	if !ok {
		gen = defaultPrinter
	}

	params, err := gen(o, now)
	if err != nil || len(params) == 0 {
		params, _ = defaultPrinter(o, now)
	}

	return params
}
