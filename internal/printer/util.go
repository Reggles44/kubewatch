package printer

import (
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func calculateResourceDuration(o *unstructured.Unstructured, now time.Time) string {
	duration := now.Sub(o.GetCreationTimestamp().Time)
	return fmt.Sprintf("%-8v", duration.Round(time.Second))
}

func convertUnstructured(dest any, obj *unstructured.Unstructured) error {
	raw, err := json.Marshal(obj.Object)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, dest)
}
