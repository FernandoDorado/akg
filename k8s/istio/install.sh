#!/bin/bash
echo "[ISTIO]"
istioctl install -f istio.yml -y

echo "[POST INSTALL"
kubectl apply -f gateway.yml
