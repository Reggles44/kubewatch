package kube

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/resource"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/kubectl/pkg/util/interrupt"
)

func WatchRuntimeObject(r *resource.Result, dataCh chan watch.Event) error {
	obj, err := r.Object()
	if err != nil {
		return err
	}

	// watching from resourceVersion 0, starts the watch at ~now and
	// will return an initial watch event.  Starting form ~now, rather
	// the rv of the object will insure that we start the watch from
	// inside the watch window, which the rv of the object might not be.

	rv := "0"
	isList := meta.IsListType(obj)
	if isList {
		// the resourceVersion of list objects is ~now but won't return
		// an initial watch event
		rv, err = meta.NewAccessor().ResourceVersion(obj)
		if err != nil {
			return err
		}
	}

	// print the current object
	var objsToPrint []runtime.Object
	if isList {
		objsToPrint, _ = meta.ExtractList(obj)
	} else {
		objsToPrint = append(objsToPrint, obj)
	}

	// Existing events
	for _, obj = range objsToPrint {
		dataCh <- watch.Event{Type: watch.Added, Object: obj}
	}

	// print watched changes
	w, err := r.Watch(rv)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	intr := interrupt.New(nil, cancel)
	intr.Run(func() error {
		_, err := watchtools.UntilWithoutRetry(ctx, w,
			func(e watch.Event) (bool, error) {
				dataCh <- e
				return false, nil
			})
		return err
	})

	return nil
}
