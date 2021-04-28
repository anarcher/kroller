package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Client) StatefulSets(ctx context.Context) (*appv1.StatefulSetList, error) {
	sts, err := c.clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load statefulsets: %w", err)
	}

	return sts, nil
}

func (c *Client) UpdateStatefulSet(ctx context.Context, s *appv1.StatefulSet) error {
	_, err := c.clientset.AppsV1().StatefulSets(s.Namespace).Update(ctx, s, metav1.UpdateOptions{})
	return err
}

func (c *Client) PatchStatefulSet(ctx context.Context, s *appv1.StatefulSet, patch map[string]interface{}) error {
	bs, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	pt := types.StrategicMergePatchType
	_, err = c.clientset.AppsV1().StatefulSets(s.Namespace).Patch(ctx, s.Name, pt, bs, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
