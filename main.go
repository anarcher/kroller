package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anarcher/kroller/pkg/cmd"
	"github.com/anarcher/kroller/pkg/kubernetes"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	rootCmd, cfg := cmd.NewRootCmd()
	restartCmd := cmd.NewRestartCmd(cfg)
	drainCmd := cmd.NewDrainCmd(cfg)
	showCmd := cmd.NewShowCmd(cfg)
	versionCmd := cmd.NewVersionCmd()

	rootCmd.Subcommands = []*ffcli.Command{
		restartCmd,
		drainCmd,
		showCmd,
		versionCmd,
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error during Parse: %v\n", err)
		os.Exit(1)
	}

	client, err := kubernetes.NewClient(cfg.KubeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error constructing kubernetes client: %v\n", err)
		os.Exit(1)
	}
	cfg.KubeClient = client

	if err := rootCmd.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
