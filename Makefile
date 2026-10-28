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
	@kubectl apply -f demo/tq

reload: build
	@kubectl rollout restart deploy/tq
