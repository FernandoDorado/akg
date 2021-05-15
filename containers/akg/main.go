package main

import (
    "log"

    "akg/webserver"
    "akg/k8s"
)

func main() {
    log.Print("debug begin")

    // k8s
    k := &k8s.K8s{
        InCluster: true,
    }
    k.Configure()
    k.Connect()
    akgPods := k.AkgPods()
    log.Printf("Akg Pods: %s", akgPods)

    // gin
    ws := &webserver.WebServer{}
    ws.Start()
}

