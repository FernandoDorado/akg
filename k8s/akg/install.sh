#!/bin/bash
kubectl apply -f vs.yml
kubectl apply -f cluster-role.yml
kubectl apply -f cluster-role-binding.yml
kubectl apply -f sa.yml
kubectl apply -f ksvc.yml
