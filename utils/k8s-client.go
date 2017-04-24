package client

import (
	k8s "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	watch "k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	clientset *k8s.Clientset
}

func NewK8sClint() *K8sClient {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientSet
	clientSet, err := k8s.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client := &K8sClient{
		clientset: clientSet,
	}
	return client
}

func (c *K8sClient) WatchNodes(listOption v1.ListOptions) (watch.Interface, error) {
	return c.clientset.CoreV1().Nodes().Watch(listOption)
}
