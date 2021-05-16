package k8s

import (
    "context"

	  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type K8s struct {
    config *rest.Config
    client *kubernetes.Clientset
    InCluster bool
}

func (k *K8s) Configure() {
    if k.InCluster {
        config, err := rest.InClusterConfig()
        if err != nil {
            panic(err.Error())
        }

        k.config = config
    } else {
        panic("not in cluster!")
    }
}

func (k *K8s) Connect() {
    clientset, err := kubernetes.NewForConfig(k.config)
    if err != nil {
        panic(err.Error())
    }

    k.client = clientset
}

func (k *K8s) Deployments() []string {
    pods, err := k.client.ExtensionsV1beta1().Deployments("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    var Deployments []string

    for _, deployment := range pods.Items {
        Deployments = append(Deployments, deployment.ObjectMeta.Name)
    }

    return Deployments
}
