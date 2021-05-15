#!/bin/bash
echo "[ISTIO]"
istioctl upgrade -f istio.yml -y

echo "[POST INSTALL]"
kubectl apply -f gateway.yml
