package kubernetes

import (
	"context"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

func (c *Client) Node(ctx context.Context, nodeName string) (*v1.Node, error) {
	node, err := c.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the node")
	}
	return node, nil
}

func (c *Client) PodsOnNode(ctx context.Context, nodeName string) (*v1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{
			"spec.nodeName": nodeName,
			"status.phase":  "Running",
		}).String()})

	if err != nil {
		return nil, errors.Wrap(err, "failed to load pods on the node "+nodeName)
	}
	return pods, nil
}

func (c *Client) CordonNode(ctx context.Context, node *v1.Node) (*v1.Node, error) {
	node, err := c.clientset.CoreV1().Nodes().Get(ctx, node.ObjectMeta.Name, metav1.GetOptions{})
	node.Spec.Unschedulable = true
	n, err := c.clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to cordon node: "+node.Name)
	}
	return n, nil
}
