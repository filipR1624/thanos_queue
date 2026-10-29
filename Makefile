.PHONY: cluster delete build deploy reload

cluster:
	@kind create cluster

delete:
	@kind delete cluster

build:
	@docker build . -t thanos_queue:dev
	@kind load docker-image thanos_queue:dev

deploy:
	@kubectl apply -f demo/prometheus
	@kubectl apply -f demo/thanos
	@kubectl apply -f demo/redpanda
	@kubectl rollout status sts/redpanda --timeout=2m30s
	@kubectl exec sts/redpanda -- rpk topic create remote-write --if-not-exists
	@kubectl apply -f demo/tq

reload: build
	@kubectl rollout restart deploy/tq

sever:
	@kubectl patch svc receive-service -p '{"spec":{"selector":{"app":"outage-test"}}}'
	@kubectl rollout restart deploy/tq
