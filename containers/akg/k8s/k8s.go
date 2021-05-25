package k8s

import (
	"context"
	"fmt"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s struct {
	config    *rest.Config
	client    *kubernetes.Clientset
	InCluster bool
}

type Instance struct {
	Name string
}

type App struct {
	Name      string
	Instances []Instance
}

func (k *K8s) Configure() {
	if k.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		k.config = config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", "/auth/.kube/config")
		if err != nil {
			panic(err.Error())
		}
		k.config = config
	}
}

func (k *K8s) Connect() {
	clientset, err := kubernetes.NewForConfig(k.config)
	if err != nil {
		panic(err.Error())
	}

	k.client = clientset
}

func (k *K8s) Apps() ([]App, error) {
	apps := []App{}

	deployments, err := k.client.AppsV1().Deployments("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list deployments: %s", err)
		return []App{}, err
	}

	for _, deployment := range deployments.Items {
		name := deployment.Name
		app := App{
			Name: name,
		}

		pods, err := k.client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("failed to list pods: %s", err)
			return []App{}, err
		}

		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, fmt.Sprintf("%s-", name)) {
				instance := Instance{
					Name: pod.Name,
				}
				app.Instances = append(app.Instances, instance)
			}
		}

		apps = append(apps, app)
	}

	return apps, nil
}
