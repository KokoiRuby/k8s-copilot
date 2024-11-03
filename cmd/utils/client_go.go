package utils

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientGo struct {
	ClientSet       *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
}

func NewClientGo(kubeconfig string) (*ClientGo, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	return &ClientGo{
		ClientSet:       clientSet,
		DynamicClient:   dynamicClient,
		DiscoveryClient: discoveryClient,
	}, nil
}
