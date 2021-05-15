#!/bin/bash
kubectl apply -f gateway.yml
kubectl apply -f vs.yml
kubectl apply -f role.yml
kubectl apply -f rolebinding.yml
kubectl apply -f sa.yml
kubectl apply -f ksvc.yml
