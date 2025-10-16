package kube

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type WatchEvent struct {
	Event  string
	Object interface{}
}

func Watch(
	resource schema.GroupVersionResource,
	namespace string,
	dataCh chan WatchEvent,
) error {
	client.AppsV1().Deployments("")
	// Build Informer
	factory := informers.NewSharedInformerFactoryWithOptions(
		client,
		time.Minute,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(opt *metav1.ListOptions) {
			opt.FieldSelector = fields.Everything().String()
		}),
	)

	gi, err := factory.ForResource(resource)
	if err != nil {
		panic(err)
	}

	infor := gi.Informer()
	infor.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			go func() {
				dataCh <- WatchEvent{"add", obj}
			}()
		},
		DeleteFunc: func(obj interface{}) {
			go func() {
				dataCh <- WatchEvent{"delete", obj}
			}()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			go func() {
				dataCh <- WatchEvent{"delete", oldObj}
				dataCh <- WatchEvent{"add", newObj}
			}()
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	go factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	<-stopCh
	return nil
}
