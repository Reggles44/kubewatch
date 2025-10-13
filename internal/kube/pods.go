package kube

import (
	"context"
	"log"
	"time"

	"github.com/reggles44/kubewatch/internal/display"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func PodList(namespace string) []*corev1.Pod {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	ppods := []*corev1.Pod{}
	for _, pod := range pods.Items {
		ppods = append(ppods, &pod)
	}

	return ppods
}

func WatchPods(display *display.Model[corev1.Pod], namespace string) {
	// Build Informer
	factory := informers.NewSharedInformerFactoryWithOptions(
		client,
		time.Minute,
		// informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(opt *metav1.ListOptions) {
			opt.FieldSelector = fields.Everything().String()
		}),
	)

	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			// fmt.Printf("[+] Pod added in %s: %s\n", pod.Namespace, pod.GetName())
			go func() { display.ChangeChannel <- [2]*corev1.Pod{nil, pod} }()
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			// fmt.Printf("[-] Pod deleted from %s: %s\n", pod.Namespace, pod.GetName())
			go func() { display.ChangeChannel <- [2]*corev1.Pod{pod, nil} }()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)
			// if oldPod.ResourceVersion != newPod.ResourceVersion {
			// 	fmt.Printf("[~] Pod updated in %s: %s\n", namespace, newPod.GetName())
			// }
			go func() { display.ChangeChannel <- [2]*corev1.Pod{oldPod, newPod} }()
		},
	})

	// Create a channel to signal when to stop the informer.
	// When this channel is closed, the informer will gracefully shut down.
	stopCh := make(chan struct{})
	defer close(stopCh)

	// Start the informer factory. This begins the process of listing and watching events.
	// This also runs in a goroutine, so it doesn't block the current goroutine.
	go factory.Start(stopCh)

	// Wait for the informer's caches to be synced. This is important!
	// It ensures the informer has retrieved the initial state of all Pods
	// before it starts processing real-time events. This prevents missing initial events.
	// It will block until caches are synced or stopCh is closed.
	factory.WaitForCacheSync(stopCh)
	log.Printf("Cache synced for namespace: %s. Ready to watch events.\n", namespace)

	// This line keeps the goroutine running indefinitely.
	// It will block until the 'stopCh' channel is closed, allowing the informer to run in the background.
	<-stopCh
	// If `stopCh` is closed, this goroutine will exit.
}
