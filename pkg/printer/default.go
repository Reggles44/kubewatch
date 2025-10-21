package printer

import (
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func defaultPrinter(o *unstructured.Unstructured, now time.Time) ([]string, error) {
	return []string{
		o.GetName(),
		o.GetNamespace(),
		o.GetKind(),
		calculateResourceDuration(o, now),
	}, nil
}
