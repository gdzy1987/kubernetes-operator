package main

import (
	"flag"
	"time"

	"github.com/gosoon/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	informers "github.com/gosoon/kubernetes-operator/pkg/client/informers/externalversions"
	"github.com/gosoon/kubernetes-operator/pkg/controller"
	"github.com/resouer/k8s-controller-custom-resource/pkg/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	kubernetesClusterClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubernetesClusterInformerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, time.Second*30)

	controller := controller.NewController(kubeClient, kubernetesClusterClient,
		kubernetesClusterInformerFactory.Ecs().V1().KubernetesClusters())

	go kubernetesClusterInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}