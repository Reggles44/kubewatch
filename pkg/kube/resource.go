package kube

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

func GetResource(args []string, configFlags *genericclioptions.ConfigFlags) (*resource.Result, error) {
	r := resource.NewBuilder(configFlags).
		Unstructured().
		NamespaceParam(*configFlags.Namespace).DefaultNamespace().AllNamespaces(*configFlags.Namespace == "").
		FilenameParam(true, &resource.FilenameOptions{}).
		LabelSelectorParam("").
		FieldSelectorParam("").
		RequestChunksOf(500).
		ResourceTypeOrNameArgs(true, args...).
		SingleResourceType().
		Latest().
		TransformRequests().
		Do()

	if err := r.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
