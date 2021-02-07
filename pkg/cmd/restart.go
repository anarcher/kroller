package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type RestartConfig struct {
	rootCfg *RootConfig
	targets stringSlice
}

func NewRestartCmd(rootCfg *RootConfig) *ffcli.Command {
	cfg := &RestartConfig{
		rootCfg: rootCfg,
	}

	fs := flag.NewFlagSet("kroller restart", flag.ExitOnError)
	fs.String("config", "", "config file (optional)")
	fs.Var(&cfg.targets, "target", "only use the specified objects (Format: <namespace>/<type>/<name>)")
	rootCfg.RegisterFlags(fs)

	c := &ffcli.Command{
		Name:    "restart",
		FlagSet: fs,
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
	for _, t := range c.targets {
		fmt.Println("target:", t)
	}

	client := c.rootCfg.KubeClient

	deploys, err := client.Deployments(ctx)
	if err != nil {
		return err
	}

	for _, d := range deploys.Items {
		fmt.Printf("%s/%s", d.Namespace, d.Name)
		println("")
	}

	return nil
}
