package kubeclient

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const APIServer = "127.0.0.1:8080"

type Client struct {
	APIServer string
	Clientset *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	restConfig := &rest.Config{
		Host: APIServer,
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}

	client := &Client{
		APIServer: APIServer,
		Clientset: clientset,
	}
	return client, nil
}

func (client *Client) ListNode() (*v1.NodeList, error) {
	nodeList, err := client.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}

func (client *Client) ListPod() (*v1.PodList, error) {
	podList, err := client.Clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList, nil
}
