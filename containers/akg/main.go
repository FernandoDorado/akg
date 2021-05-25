package main

import (
	"log"

	"akg/k8s"
	"akg/webserver"
)

func main() {
	// k8s
	k := &k8s.K8s{
		InCluster: true,
	}
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
