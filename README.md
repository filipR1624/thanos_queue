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

`make` is used to automated most functionality. Targets include:
- cluster: Creates a kind cluster
- delete: Deletes the kind cluster (including any resources within)
- reload: Builds an image of this service and loads it into kind
- deploy: Deploys all relevant components (Thanos, Prometheus, Redpanda) and initializes a topic in Redpanda
