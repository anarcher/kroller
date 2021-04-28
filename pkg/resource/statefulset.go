package resource

import (
	"context"
	"time"

	"github.com/anarcher/kroller/pkg/kubernetes"
	appv1 "k8s.io/api/apps/v1"
)

type StatefulSet struct {
	client *kubernetes.Client
	sts    *appv1.StatefulSet
}

func addStatefulSetList(rl *RolloutList, list *appv1.StatefulSetList, client *kubernetes.Client) {
	for _, _s := range list.Items {
		s := _s
		rl.add(&StatefulSet{
			sts:    &s,
			client: client,
		})
	}
}

func (s *StatefulSet) Kind() string {
	return "statefulset"
}

func (s *StatefulSet) Name() string {
	return s.sts.Name
}

func (s *StatefulSet) Namespace() string {
	return s.sts.Namespace
}

func (s *StatefulSet) Restart(ctx context.Context) error {
	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"kubectl.kubernetes.io/restartedAt": time.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}

	return s.client.PatchStatefulSet(ctx, s.sts, patch)
}
