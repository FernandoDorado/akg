#!/bin/bash
echo "[ISTIO]"
istioctl upgrade -f istio.yml -y

echo "[CUSTOM]"
kubectl apply -f custom.yml
