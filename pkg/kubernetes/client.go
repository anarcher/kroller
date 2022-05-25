package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	roclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
)

type Client struct {
	clientset         *kubernetes.Clientset
	rolloutsClientset roclientset.Interface
	kubeConfig        string
}

func NewClient(kubeConfig string) (*Client, error) {
	if kubeConfig == "" {
		kubeConfig = clientcmd.RecommendedHomeFile
	}

	clientConfig := makeClientConfig(kubeConfig)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("%w kubeconfig: %s", err, kubeConfig)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w kubeconfig: %s", err, kubeConfig)
	}
	rolloutsClientset, err := roclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("rolloutsClient: %w kubeconfig: %s", err, kubeConfig)
	}

	c := &Client{
		clientset:         clientset,
		rolloutsClientset: rolloutsClientset,
		kubeConfig:        kubeConfig,
	}

	return c, nil
}

func makeClientConfig(kubeConfig string) clientcmd.ClientConfig {
	cc := strings.Split(kubeConfig, ":")
	c := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: cc},
		&clientcmd.ConfigOverrides{},
	)

	return c
}
