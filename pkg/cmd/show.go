package cmd

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewShowCmd(rootCfg *RootConfig) *ffcli.Command {
	showNodesCmd := NewShowNodesCmd(rootCfg)

	c := &ffcli.Command{
		Name:       "show",
		ShortUsage: "show details of nodes or resources",
		ShortHelp:  "show details of nodes or resources",
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}
	c.Subcommands = []*ffcli.Command{
		showNodesCmd,
	}
	return c
}
