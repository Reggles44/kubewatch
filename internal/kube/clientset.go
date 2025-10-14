package kube

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var client *kubernetes.Clientset

func init() {
	// Setup clientset
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// cs, err := dynamic.NewForConfig(config)
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client = cs
}
