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

func (d *Deployment) Restart(ctx context.Context) error {
	obj := d.deploy
	if obj.Spec.Template.ObjectMeta.Annotations == nil {
		obj.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	obj.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	client := d.client
	return client.UpdateDeployment(ctx, obj)
}
