package main

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

	"github.com/gin-gonic/gin"
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
}

type replica struct {
	Name string `json:"name"`
}

type application struct {
	Name     string    `json:"name"`
	Replicas []replica `json:"replicas"`
}

func main() {
	k := &k8s{}
	k.InCluster = (os.Getenv("IN_CLUSTER") == "true")
	k.CloudProvider = os.Getenv("CLOUD_PROVIDER")
	if ok, _ := k.configure(); ok {
		k.connect()
	}

	// http
	e := gin.Default()
	e.LoadHTMLGlob("templates/*")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	e.GET("/api/v1/apps", func(c *gin.Context) {
		apps, err := k.apps()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"apps": apps,
			})
		}
	})
	e.GET("/health/live", func(c *gin.Context) {})
	e.Run()
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
			log.Printf(msg)
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

	// call digital ocean api
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/kubernetes/clusters/%s/credentials", clusterId), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := httpClient.Do(req)
	if err != nil {
		msg := fmt.Sprintf("failed to call digitalocean api: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}
	if resp.StatusCode != http.StatusOK {
		msg := "non 2XX response from digitalocean api"
		log.Print(msg)
		return false, errors.New(msg)
	}

	// parse response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read digitalocean api response: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}
	body := string(respBytes)
	var data map[string]string
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal digitalocean api response json: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}

	// get cert
	cert, err := base64.StdEncoding.DecodeString(data["certificate_authority_data"])
	if err != nil {
		msg := fmt.Sprintf("failed to decode cert: %s", err)
		log.Print(msg)
		return false, errors.New(msg)
	}

	// config
	k.Config.TLSClientConfig.CAData = []byte(cert)

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
