package resource

import (
	"context"
	"time"

	"github.com/anarcher/kroller/pkg/kubernetes"
	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

type ArgoRollout struct {
	client  *kubernetes.Client
	rollout *v1alpha1.Rollout
}

func addArgoRollouts(rl *RolloutList, list *v1alpha1.RolloutList, client *kubernetes.Client) {
	for _, _i := range list.Items {
		i := _i
		rl.add(&ArgoRollout{
			rollout: &i,
			client:  client,
		})
	}
}

func (r *ArgoRollout) Kind() string {
	return "Rollout"
}

func (r *ArgoRollout) Name() string {
	return r.rollout.Name
}

func (r *ArgoRollout) Namespace() string {
	return r.rollout.Namespace
}

func (r *ArgoRollout) Restart(ctx context.Context) error {
	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"restartAt": time.Now().Format(time.RFC3339),
		},
	}

	return r.client.PatchRollout(ctx, r.rollout, patch)
}
