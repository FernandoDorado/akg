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

func (k *K8s) AkgPods() []string {
    pods, err := k.client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    var AkgPods []string

    for _, pod := range pods.Items {
        name := pod.ObjectMeta.Name
        if name[0:4] == "akg-" {
            AkgPods = append(AkgPods, name)
        }
    }

    return AkgPods
}
