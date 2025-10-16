package resources

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type KubeResource struct {
	Scheme schema.GroupVersionResource
	Key    func(obj interface{}) string
	Params func(obj interface{}, now time.Time) []string
}

var resMap = map[string]KubeResource{
	"pod":  Pod,
	"node": Node,
}

func Keys() []string {
	var keys []string

	for k := range resMap {
		keys = append(keys, k)
	}

	return keys
}

func Get(name string) (KubeResource, bool) {
	r, ok := resMap[name]
	return r, ok
}
