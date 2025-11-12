package connect

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aide-family/magicbox/load"
)

const defaultKubeConfig = "~/.kube/config"

func NewKubernetesClientSet(kubeConfig string) (*kubernetes.Clientset, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		if kubeConfig == "" {
			kubeConfig = defaultKubeConfig
		}
		restConfig, err = clientcmd.BuildConfigFromFlags("", load.ExpandHomeDir(kubeConfig))
		if err != nil {
			return nil, err
		}
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
