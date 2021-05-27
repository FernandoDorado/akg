package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8s struct {
	Config        *rest.Config
	Client        *kubernetes.Clientset
	InCluster     bool
	CloudProvider string
}

type Replica struct {
	Name string
}

type App struct {
	Name     string
	Replicas []Replica
}

func (k *K8s) Configure() {
	// configure in/out cluster basic credentials
	if k.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		k.Config = config
	} else {
		homedir := homedir.HomeDir()
		config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/%s", homedir, ".kube/config"))
		if err != nil {
			panic(err.Error())
		}
		k.Config = config
	}

	if k.CloudProvider == "digitalocean" {
		k.ConfigureForDigitalOcean()
	}
}

func (k *K8s) ConfigureForDigitalOcean() {
	// fetch credentials from env
	clusterId := os.Getenv("DO_CLUSTER_ID")
	if clusterId == "" {
		panic(errors.New("client wants to use digitalocean but DO_CLUSTER_ID is not provided"))
	}
	accessToken := os.Getenv("DO_ACCESS_TOKEN")
	if accessToken == "" {
		panic(errors.New("client wants to use digitalocean but DO_ACCESS_TOKEN is not provided"))
	}

	// call digital ocean api
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/kubernetes/clusters/%s/credentials", clusterId), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("failed to call digitalocean api: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(errors.New("non 2XX response from digitalocean api"))
	}

	// parse response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read digitalocean api response: %s", err)
		panic(err)
	}
	body := string(respBytes)
	var data map[string]string
	json.Unmarshal([]byte(body), &data)

	// get cert
	cert, err := base64.StdEncoding.DecodeString(data["certificate_authority_data"])
	if err != nil {
		log.Printf("failed to decode cert: %s", err)
		panic(err)
	}

	// config
	k.Config.TLSClientConfig.CAData = []byte(cert)
}

func (k *K8s) Connect() {
	clientset, err := kubernetes.NewForConfig(k.Config)
	if err != nil {
		panic(err.Error())
	}

	k.Client = clientset
}

func (k *K8s) Replicas() ([]Replica, error) {
	pods, err := k.Client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list pods: %s", err)
		return []Replica{}, err
	}

	instances := []Replica{}

	for _, pod := range pods.Items {
		replica := Replica{
			Name: pod.Name,
		}
		instances = append(instances, replica)
	}

	return instances, nil
}

func (k *K8s) Apps() ([]App, error) {
	deployments, err := k.Client.AppsV1().Deployments("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list deployments: %s", err)
		return []App{}, err
	}

	apps := []App{}

	for _, deployment := range deployments.Items {
		name := deployment.Name
		app := App{
			Name: name,
		}

		pods, err := k.Client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("failed to list pods: %s", err)
			return []App{}, err
		}

		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, fmt.Sprintf("%s-", name)) {
				replica := Replica{
					Name: pod.Name,
				}
				app.Replicas = append(app.Replicas, replica)
			}
		}

		apps = append(apps, app)
	}

	return apps, nil
}
