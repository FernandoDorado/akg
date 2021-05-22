#!/bin/bash
kubectl apply -f vs.yml
kubectl apply -f svc.yml
kubectl apply -f deployment.yml
