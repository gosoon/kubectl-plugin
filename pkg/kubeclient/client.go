package kubeclient

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfig = "/root/.kube/config"

type Client struct {
	Clientset *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	client := &Client{
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
	podList, err := client.Clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList, nil
}
