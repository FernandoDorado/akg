package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adamkgray/dok8cert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type k8s struct {
	Config        *rest.Config
	Client        *kubernetes.Clientset
	InCluster     bool
	CloudProvider string
	Ready         bool
}

type replica struct {
	Name string `json:"name"`
}

type application struct {
	Name     string    `json:"name"`
	Replicas []replica `json:"replicas"`
}

func (k *k8s) configure() (bool, error) {
	// configure in/out cluster basic credentials
	if k.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			msg := fmt.Sprintf("failed to create in-cluster config: %s", err)
			log.Print(msg)
			return false, errors.New(msg)
		}
		k.Config = config
	} else {
		homedir := homedir.HomeDir()
		config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/%s", homedir, ".kube/config"))
		if err != nil {
			msg := fmt.Sprintf("failed to create out-of-cluster config: %s", err)
			log.Print(msg)
			return false, errors.New(msg)
		}
		k.Config = config
	}

	if k.CloudProvider == "digitalocean" {
		_, err := k.configureForDigitalOcean()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// https://stackoverflow.com/questions/65042279/python-kubernetes-client-requests-fail-with-unable-to-get-local-issuer-certific
func (k *k8s) configureForDigitalOcean() (bool, error) {
	clusterId := os.Getenv("DO_CLUSTER_ID")
	if clusterId == "" {
		msg := "client wants to use digitalocean but DO_CLUSTER_ID is not provided"
		log.Print(msg)
		return false, errors.New(msg)
	}

	accessToken := os.Getenv("DO_ACCESS_TOKEN")
	if accessToken == "" {
		msg := "client wants to use digitalocean but DO_ACCESS_TOKEN is not provided"
		log.Print(msg)
		return false, errors.New(msg)
	}

	_, err := dok8cert.Update(clusterId, accessToken, k.Config)
	if err != nil {
		msg := fmt.Sprintf("failed to get custom TLS certificate from digitalocean: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}

	return true, nil
}

func (k *k8s) connect() (bool, error) {
	clientset, err := kubernetes.NewForConfig(k.Config)
	if err != nil {
		msg := fmt.Sprintf("failed to connect to kubernetes api: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}

	k.Client = clientset
	k.Ready = true

	return true, nil
}

func (k *k8s) apps() ([]application, error) {
	deployments, err := k.Client.AppsV1().Deployments("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		msg := fmt.Sprintf("failed to list deployments: %s", err)
		log.Print(msg)
		return []application{}, errors.New(msg)
	}

	apps := []application{}

	for _, deployment := range deployments.Items {
		name := deployment.Name
		app := application{
			Name: name,
		}

		pods, err := k.Client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			msg := fmt.Sprintf("failed to list pods: %s", err)
			log.Print(msg)
			return []application{}, errors.New(msg)
		}

		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, fmt.Sprintf("%s-", name)) {
				replica := replica{
					Name: pod.Name,
				}
				app.Replicas = append(app.Replicas, replica)
			}
		}

		apps = append(apps, app)
	}

	return apps, nil
}
