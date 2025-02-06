package cmd

import (
	"context"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
)

const Version = "v0.3.0"

func NewVersionCmd() *ffcli.Command {
	c := &ffcli.Command{
		Name:       "version",
		ShortUsage: "show version of kroller",
		ShortHelp:  "show version of kroller",
		Exec: func(context.Context, []string) error {
			fmt.Printf("version: %s\n", Version)
			return nil
		},
	}
	return c
}
