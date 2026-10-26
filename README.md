# thanos_queue

## Description

This project was developed as the final year project for SE603 Software Development Curriculum in SETU 2026/2027.

It's a service placed between Prometheus and Thanos that implements a Kafka queue to persist data while accelerating
write acknowledgements.

## Requirements
- Docker
- kind (local Kubernetes cluster)
- kubectl

## How to run
- Ensure Docker is running
- Create a local Kubernetes cluster: `kind create cluster`
- Apply the manifests: `kubectl apply -f demo`
  - run this from the repository's root
- Verify it's operational: `kubectl port-forward svc/query-service 9090:9090`
  - Port-forwarding will continue until manually cancelled using Ctrl + C
  - Open `localhost:9090` in your browser, you should see the Thanos Query UI
  - Run a query like `up` and execute it.
  - You should see the query returned with `receive_replica="0"` and "`tenant_id="default-tenant"`
- Cleanup: `kubectl delete -f demo && kind delete cluster`
