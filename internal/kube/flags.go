package kube

import "k8s.io/cli-runtime/pkg/genericclioptions"

func NewConfigFlags() *genericclioptions.ConfigFlags {
	return genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
}

// func NewPrintFlags() *genericclioptions.PrintFlags {
// 	return genericclioptions.NewPrintFlags("")
// }
