# Makefile
.PHONY: users-generate users-dev users-build-image users-minikube-build-image

users-generate:
	@echo "Running generate scripts..."
	@sh ./apps/users/scripts/generate-proto.sh
	@sh ./apps/users/scripts/generate-wire.sh

users-dev:
	@go run ./apps/users/cmd/main.go

users-build-image:
	@docker build -t users-service:dev -f apps/users/deploy/docker/Dockerfile .

users-minikube-build-image:
	@echo "Setting up Minikube Docker environment..."
	eval "$$(minikube docker-env)" && make users-build-image

users-minikube-deploy:
	kubectl apply -k apps/users/deploy/k8s/overlays/dev/

users-minikube-destroy:
	kubectl delete -k apps/users/deploy/k8s/overlays/dev/