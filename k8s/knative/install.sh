#!/bin/bash
echo "[KNATIVE]"
kubectl apply -f https://github.com/knative/operator/releases/download/v0.22.0/operator.yaml
kubectl apply -f namespace.yml
kubectl apply -f serving.yml
