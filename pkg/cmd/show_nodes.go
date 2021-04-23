package cmd

import (
	"context"
	"flag"

	"github.com/fatih/color"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rodaine/table"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type NodeInfo struct {
	Node          v1.Node
	Deployments   []v1.Pod
	Statusfulsets []v1.Pod
}

type ShowNodesConfig struct {
	rootCfg       *RootConfig
	labelSelector string
}

func NewShowNodesCmd(rootCfg *RootConfig) *ffcli.Command {
	cfg := &ShowNodesConfig{
		rootCfg: rootCfg,
	}

	fs := flag.NewFlagSet("kroller show nodes", flag.ExitOnError)
	fs.String("config", "", "config file (optional)")
	fs.StringVar(&cfg.labelSelector, "l", "", "selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")

	rootCfg.RegisterFlags(fs)

	c := &ffcli.Command{
		Name:       "nodes",
		ShortUsage: "show details of nodes",
		ShortHelp:  "show details of nodes",
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

func (c *ShowNodesConfig) Exec(ctx context.Context, args []string) error {
	//verbose := c.rootCfg.Verbose
	kubeCli := c.rootCfg.KubeClient

	selector, err := fields.ParseSelector(c.labelSelector)
	if err != nil {
		return err
	}

	nodes, err := kubeCli.Nodes(ctx, selector)
	if err != nil {
		return err
	}

	var nodeInfos []*NodeInfo
	for _, n := range nodes.Items {
		pods, err := kubeCli.PodsOnNode(ctx, n.Name)
		if err != nil {
			return err
		}

		var (
			deploys []v1.Pod
			sts     []v1.Pod
		)
		for _, p := range pods.Items {
			if isDeploymentPod(&p) {
				deploys = append(deploys, p)
				continue
			}
			if isStatefulSetPod(&p) {
				sts = append(sts, p)
			}
		}

		ni := &NodeInfo{
			Node:          n,
			Deployments:   deploys,
			Statusfulsets: sts,
		}
		nodeInfos = append(nodeInfos, ni)
	}

	tbl := table.New("Node", "Deployment", "Statusfulsets")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, n := range nodeInfos {
		tbl.AddRow(n.Node.Name, len(n.Deployments), len(n.Statusfulsets))
	}

	tbl.Print()

	return nil
}

func isDeploymentPod(p *v1.Pod) bool {
	controllerRef := metav1.GetControllerOf(p)
	if controllerRef == nil {
		return false
	}
	if controllerRef.Kind == "ReplicaSet" {
		return true
	}

	return false
}

func isStatefulSetPod(p *v1.Pod) bool {
	controllerRef := metav1.GetControllerOf(p)
	if controllerRef == nil {
		return false
	}
	if controllerRef.Kind == "StatefulSet" {
		return true
	}
	return false
}
