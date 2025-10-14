package kube

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func Watch(
	namespace string,
	dataCh chan [2]interface{},
	stopCh chan struct{},
	getInformer func(factory informers.SharedInformerFactory) cache.SharedIndexInformer,
) {
	// Build Informer
	factory := informers.NewSharedInformerFactoryWithOptions(
		client,
		time.Minute,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(opt *metav1.ListOptions) {
			opt.FieldSelector = fields.Everything().String()
		}),
	)

	informer := getInformer(factory)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			go func() { dataCh <- [2]interface{}{nil, obj} }()
		},
		DeleteFunc: func(obj interface{}) {
			go func() { dataCh <- [2]interface{}{obj, nil} }()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			go func() { dataCh <- [2]interface{}{oldObj, newObj} }()
		},
	})

	go factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	<-stopCh
}
