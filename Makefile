build:
	@echo "==> building image"
	@docker build . -t thanos_queue:dev

cluster:
	@kind create cluster

clean:
	@kind delete cluster

load:
	@kind load docker-image thanos_queue:dev

deploy:
	@kubectl apply -f demo

full:
	make cluster
	make build
	make load
	make deploy
