package resource

import (
	"context"
	"time"

	"github.com/anarcher/kroller/pkg/kubernetes"
	appv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	client *kubernetes.Client
	deploy *appv1.Deployment
}

func addDeploymentList(rl *RolloutList, list *appv1.DeploymentList, client *kubernetes.Client) {
	for _, _d := range list.Items {
		d := _d
		rl.add(&Deployment{
			deploy: &d,
			client: client,
		})
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

func (d *Deployment) NodeSelector() map[string]string {
	return d.deploy.Spec.Template.Spec.NodeSelector
}

func (d *Deployment) Restart(ctx context.Context) error {
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

	return d.client.PatchDeployment(ctx, d.deploy, patch)
}
