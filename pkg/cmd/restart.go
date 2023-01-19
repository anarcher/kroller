package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/anarcher/kroller/pkg/resource"
	"github.com/anarcher/kroller/pkg/target"
	"github.com/anarcher/kroller/pkg/ui"
	"github.com/fatih/color"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"k8s.io/apimachinery/pkg/labels"
)

type RestartConfig struct {
	rootCfg      *RootConfig
	targets      stringSlice
	nodeSelector string
}

func NewRestartCmd(rootCfg *RootConfig) *ffcli.Command {
	cfg := &RestartConfig{
		rootCfg: rootCfg,
	}

	fs := flag.NewFlagSet("kroller restart", flag.ExitOnError)
	fs.String("config", "", "config file (optional)")
	fs.Var(&cfg.targets, "target", "only use the specified objects (Format: <namespace>/<type>/<name>)")
	fs.StringVar(&cfg.nodeSelector, "node-selector", "", "node label selector used to filter nodes, if empty all nodes are selected (group=nodegroup)")
	rootCfg.RegisterFlags(fs)

	c := &ffcli.Command{
		Name:       "restart",
		ShortUsage: "restart all rollout resources (deployment,statefulset)",
		ShortHelp:  "restart all rollout resources",
		FlagSet:    fs,
		Options: []ff.Option{
			ff.WithEnvVarNoPrefix(),
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
		},
		Exec: cfg.Exec,
	}

	return c
}

func (c *RestartConfig) Exec(ctx context.Context, args []string) error {
	client := c.rootCfg.KubeClient

	rl, err := resource.GetRolloutList(ctx, client)
	if err != nil {
		return err
	}

	if c.nodeSelector != "" {
		rl = c.matchNodeSelector(rl)
	}

	if len(c.targets) > 0 {

		exprs, err := target.StrExps(c.targets...)
		if err != nil {
			return err
		}

		rl = target.Filter(rl, exprs)
	}

	ui.RolloutList(rl)

	fmt.Println("")
	fmt.Printf(color.GreenString("Do you want to continue and restart? "))
	ok, err := ui.AskForConfirm()
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	for _, r := range rl {
		if err := r.Restart(ctx); err != nil {
			return err
		}
		fmt.Println(color.YellowString("Restarting %s/%s/%s...", r.Namespace(), r.Kind(), r.Name()))
	}

	return nil
}

func (c *RestartConfig) matchNodeSelector(rl resource.RolloutList) resource.RolloutList {
	var nodeSelector labels.Selector
	if ns, err := labels.Parse(c.nodeSelector); err != nil {
		log.Fatalf("parsing node selector: %s", err)
	} else {
		nodeSelector = ns
	}

	out := make(resource.RolloutList, 0, len(rl))

	for _, r := range rl {
		set := labels.Set(r.NodeSelector())
		if nodeSelector.Matches(set) {
			out = append(out, r)
		}
	}

	return out
}
