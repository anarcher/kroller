package cmd

import (
	"context"
	"flag"

	"github.com/anarcher/kroller/pkg/kubernetes"
	"github.com/anarcher/kroller/pkg/ui"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type RootConfig struct {
	KubeConfig string
	Verbose    bool

	KubeClient *kubernetes.Client
}

func NewRootCmd() (*ffcli.Command, *RootConfig) {
	fs := flag.NewFlagSet("kroller", flag.ExitOnError)

	cfg := &RootConfig{}
	cfg.RegisterFlags(fs)

	c := &ffcli.Command{
		Name:       "kroller",
		ShortUsage: "kroller <subcommand>",
		ShortHelp:  "kroller",
		Options: []ff.Option{
			ff.WithEnvVarNoPrefix(),
		},
		FlagSet: fs,
		Exec:    cfg.Exec,
	}

	return c, cfg
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.KubeConfig, "kubeconfig", "", "kubeconfig file")
	fs.BoolVar(&c.Verbose, "v", false, "log verbose output")
}

func (c *RootConfig) Exec(context.Context, []string) error {
	ui.PrintBanner("kroller")
	return flag.ErrHelp
}
