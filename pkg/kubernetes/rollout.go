package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

func (c *Client) Rollouts(ctx context.Context) (*v1alpha1.RolloutList, error) {
	rls, err := c.rolloutsClientset.ArgoprojV1alpha1().Rollouts("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load argo rollouts: %w", err)
	}
	return rls, nil
}

func (c *Client) PatchRollout(ctx context.Context, r *v1alpha1.Rollout, patch map[string]interface{}) error {
	bs, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	pt := types.MergePatchType
	_, err = c.rolloutsClientset.ArgoprojV1alpha1().Rollouts(r.Namespace).Patch(ctx, r.Name, pt, bs, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
