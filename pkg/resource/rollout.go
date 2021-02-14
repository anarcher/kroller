package resource

import (
	"context"

	"github.com/anarcher/kroller/pkg/kubernetes"
)

type Rollout interface {
	Kind() string
	Name() string
	Namespace() string

	Restart(context.Context) error
}

type RolloutList []Rollout

func (rl *RolloutList) add(r Rollout) {
	*rl = append(*rl, r)
}

func GetRolloutList(ctx context.Context, client *kubernetes.Client) (RolloutList, error) {
	ds, err := client.Deployments(ctx)
	if err != nil {
		return nil, err
	}

	sts, err := client.StatefulSets(ctx)
	if err != nil {
		return nil, err
	}

	var rl RolloutList
	addDeploymentList(&rl, ds, client)
	addStatefulSetList(&rl, sts, client)

	return rl, nil
}
