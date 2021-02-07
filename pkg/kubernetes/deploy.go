package kubernetes

import (
	"context"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) Deployments(ctx context.Context) (*appv1.DeploymentList, error) {
	deploys, err := c.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load deployments: %w", err)
	}

	return deploys, nil
}
