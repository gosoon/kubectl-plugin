package kubeclient

import v1 "k8s.io/api/core/v1"

type Interface interface {
	ListNode() (*v1.NodeList, error)

	ListPod() (*v1.PodList, error)
}
