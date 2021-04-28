package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Client) Deployments(ctx context.Context) (*appv1.DeploymentList, error) {
	deploys, err := c.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load deployments: %w", err)
	}

	return deploys, nil
}

func (c *Client) UpdateDeployment(ctx context.Context, d *appv1.Deployment) error {
	_, err := c.clientset.AppsV1().Deployments(d.Namespace).Update(ctx, d, metav1.UpdateOptions{})
	return err
}

func (c *Client) PatchDeployment(ctx context.Context, d *appv1.Deployment, patch map[string]interface{}) error {
	bs, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	pt := types.StrategicMergePatchType
	_, err = c.clientset.AppsV1().Deployments(d.Namespace).Patch(ctx, d.Name, pt, bs, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
