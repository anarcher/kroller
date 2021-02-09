package resource

import "context"

type Rollout interface {
	Kind() string
	Name() string
	Namespace() string

	Restart(context.Context) error
}

type RolloutList []Rollout
