#!/bin/bash
echo "[CERT-MANAGER]"
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm upgrade \
  -i \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.3.1 \
  --set installCRDs=true \
  -f values.yml

echo "[POST INSTALL]"
kubectl apply -f digitalocean-dns.yml
kubectl apply -f cluster-issuer.yml
kubectl apply -f certificate.yml

