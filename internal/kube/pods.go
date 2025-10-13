package kube

import (
	"context"
	"time"

	"github.com/reggles44/kubewatch/internal/display"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func PodList(namespace string) []*v1.Pod {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	ppods := []*v1.Pod{}
	for _, pod := range pods.Items {
		ppods = append(ppods, &pod)
	}

	return ppods
}

func WatchPods(display *display.Model[v1.Pod], namespaces []string) {
	for _, ns := range namespaces {
		go func(namespace string) {
			// Build Informer
			factory := informers.NewSharedInformerFactoryWithOptions(client, time.Minute, informers.WithNamespace(namespace))
			informer := factory.Core().V1().Pods().Informer()
			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					pod := obj.(*v1.Pod)
					display.ChangeChannel <- [2]*v1.Pod{pod, nil}
				},
				UpdateFunc: func(oldObj interface{}, newObj interface{}) {
					oldPod := oldObj.(*v1.Pod)
					newPod := newObj.(*v1.Pod)
					display.ChangeChannel <- [2]*v1.Pod{oldPod, newPod}
				},
				DeleteFunc: func(obj interface{}) {
					pod := obj.(*v1.Pod)
					display.ChangeChannel <- [2]*v1.Pod{nil, pod}
				},
			})

			// Setup stop
			stopCh := make(chan struct{})
			defer close(stopCh)

			go factory.Start(stopCh)

			factory.WaitForCacheSync(stopCh)

			<-stopCh
		}(ns)
	}

	// select {}
}
