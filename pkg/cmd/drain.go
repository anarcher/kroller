package cmd

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/anarcher/kroller/pkg/aws"
	"github.com/anarcher/kroller/pkg/kubernetes"
	"github.com/anarcher/kroller/pkg/ui"

	"github.com/fatih/color"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rodaine/table"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DrainConfig struct {
	rootCfg                  *RootConfig
	awsRegion                string
	gracePeriod              time.Duration
	node                     string
	isTerminateNode          bool
	decrementDesiredCapacity bool
}

func NewDrainCmd(rootCfg *RootConfig) *ffcli.Command {
	cfg := &DrainConfig{
		rootCfg: rootCfg,
	}

	fs := flag.NewFlagSet("kroller drain", flag.ExitOnError)
	fs.String("config", "", "config file (optional)")
	fs.StringVar(&cfg.awsRegion, "aws-region", "ap-northeast-2", "The region to use for node")
	fs.DurationVar(&cfg.gracePeriod, "grace-period", (30 * time.Second), "Pod grace-period")
	fs.StringVar(&cfg.node, "node", "", "The node that should drain")
	fs.BoolVar(&cfg.isTerminateNode, "terminate-node", false, "Terminate the AWS instance in the autoscaling group")
	fs.BoolVar(&cfg.decrementDesiredCapacity, "decr-desired-capacity", false, "Decrement desired capacity of the autoscaling group")
	rootCfg.RegisterFlags(fs)

	c := &ffcli.Command{
		Name:       "drain",
		ShortUsage: "drain node",
		ShortHelp:  "drain node",
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

func (c *DrainConfig) Exec(ctx context.Context, args []string) error {
	if c.node == "" {
		return fmt.Errorf("node is required")
	}

	if err := c.drainNode(ctx); err != nil {
		return err
	}

	if c.isTerminateNode == true {
		if err := c.terminateNode(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *DrainConfig) drainNode(ctx context.Context) error {
	verbose := c.rootCfg.Verbose
	kubeClient := c.rootCfg.KubeClient

	node, err := kubeClient.Node(ctx, c.node)
	if err != nil {
		return err
	}

	allPods, err := kubeClient.PodsOnNode(ctx, c.node)
	if err != nil {
		return err
	}

	pods := filterRollPods(allPods.Items)

	ui.PodList(pods)
	fmt.Println("")
	fmt.Printf(color.GreenString("Do you want to continue and drain?"))
	ok, err := ui.AskForConfirm()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	if _, err := kubeClient.CordonNode(ctx, node); err != nil {
		return err
	}

	ui.Print("", verbose)
	ui.PrintTitle("Cordon\n", verbose)
	ui.Print(fmt.Sprintf("[✓] %s cordoned\n\n", node.ObjectMeta.Name), verbose)

	ui.PrintTitle("Evict Pods\n", verbose)
	rollPods(ctx, kubeClient, pods, c.gracePeriod, verbose)

	return nil
}

func (c *DrainConfig) terminateNode(ctx context.Context) error {
	verbose := c.rootCfg.Verbose

	fmt.Println("")
	fmt.Printf(color.RedString("Do you want to continue and terminate the node? "))
	ok, err := ui.AskForConfirm()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	ui.Print("", verbose)
	ui.PrintTitle("Node termination:\n", verbose)

	client, err := aws.NewClient(c.awsRegion)
	if err != nil {
		return err
	}

	instanceID, err := client.GetInstanceID(c.node)
	if err != nil {
		return err
	}

	ui.Print(fmt.Sprintf("%-25s %s", "Private DNS:", c.node), verbose)
	ui.Print(fmt.Sprintf("%-25s %s", "Instance ID:", instanceID), verbose)
	ui.Print(fmt.Sprintf("Decrement desired capacity: %v", c.decrementDesiredCapacity), verbose)

	if err := client.TerminateInstance(instanceID, c.decrementDesiredCapacity); err != nil {
		return err
	}

	ui.Print("\n", verbose)
	ui.Print("[✓] Node has been terminated!\n", true)
	return nil
}

func filterRollPods(pods []v1.Pod) []v1.Pod {
	var res []v1.Pod
	for _, p := range pods {
		controllerRef := metav1.GetControllerOf(&p)
		if controllerRef == nil {
			continue
		}

		if controllerRef.Kind == "DaemonSet" {
			continue
		}
		res = append(res, p)
	}
	return res
}

func rollPods(ctx context.Context, kubeClient *kubernetes.Client, pods []v1.Pod, gracePeriod time.Duration, verbose bool) error {

	graceP := int64(gracePeriod.Seconds())
	deleteOptions := metav1.DeleteOptions{GracePeriodSeconds: &graceP}

	fmt.Println("")
	tbl := table.New(" ", "Evict pod", "New pod", "New node")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, pod := range pods {
		err := kubeClient.DeletePod(ctx, pod, deleteOptions)
		if err != nil {
			return err
		}
		newPod, err := kubeClient.DetermineNewPod(ctx, pod)
		if err != nil {
			return err
		}
		if newPod != nil {
			if err := kubeClient.WaitForPodToBeReady(ctx, newPod); err != nil {
				return err
			}
			tbl.AddRow("[✓]", pod.Name, newPod.Name, newPod.Spec.NodeName)
		} else {
			tbl.AddRow("[✓]", pod.Name, "?", "?")
		}

		if verbose {
			fmt.Printf("Evicting pod: %s\n", pod.Name)
		}
	}
	tbl.Print()

	return nil
}
