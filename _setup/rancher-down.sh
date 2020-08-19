#!/bin/sh
kubectl delete namespace $(kubectl get namespaces| grep p- | awk '{print $1}' | xargs)
kubectl delete namespace $(kubectl get namespaces| grep cattle | awk '{print $1}' | xargs)
kubectl delete namespace local cert-manager
