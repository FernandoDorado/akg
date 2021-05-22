package main

import (
	"log"

	"akg/k8s"
	"akg/webserver"
)

func main() {
	log.Print("debug begin")

	// k8s
	k := &k8s.K8s{
		InCluster: true,
	}
	k.Configure()
	k.Connect()
	apps := k.Apps()
	log.Printf("Apps: %s", apps)

	// gin
	ws := &webserver.WebServer{}
	ws.Start()
}
