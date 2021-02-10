package resource

import (
	"context"

	appv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	deploy *appv1.Deployment
}

func addDeploymentList(rl *RolloutList, list *appv1.DeploymentList) {
	for _, _d := range list.Items {
		d := _d
		rl.add(&Deployment{&d})
	}
}

func (d *Deployment) Kind() string {
	return "deployment"
}

func (d *Deployment) Name() string {
	return d.deploy.Name
}

func (d *Deployment) Namespace() string {
	return d.deploy.Namespace
}

func (d *Deployment) Restart(ctx context.Context) error {
	return nil
}
