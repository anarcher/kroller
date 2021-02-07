package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset  *kubernetes.Clientset
	kubeConfig string
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

	c := &Client{
		clientset:  clientset,
		kubeConfig: kubeConfig,
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
