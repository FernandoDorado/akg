package main

import (
	"log"
	"os"

	"akg/k8s"
	"akg/webserver"
)

func main() {
	// k8s
	k := &k8s.K8s{}
	k.InCluster = (os.Getenv("IN_CLUSTER") == "true")
	k.CloudProvider = os.Getenv("CLOUD_PROVIDER")
	k.Configure()
	k.Connect()

	apps, err := k.Apps()
	if err != nil {
		log.Printf("failed to get apps: %s", err)
	} else {
		log.Printf("found apps: %s", apps)
	}

	// gin
	ws := &webserver.WebServer{}
	ws.Start()
}
