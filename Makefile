# Makefile

# Available services
SERVICES := users tasks

# Default service for single-service commands
SERVICE ?= users

# Generic targets that work for any service
%-dev:
	@go run ./apps/$*/cmd/api/main.go

%-generate:
	@echo "Running generate scripts for $*..."
	@sh ./apps/$*/scripts/generate-proto.sh
	@cd ./apps/$* && sh ./scripts/generate-wire.sh

%-build-image:
	@docker build -t $*-service:dev -f apps/$*/deploy/docker/Dockerfile .

%-minikube-build-image:
	@echo "Setting up Minikube Docker environment for $*..."
	@eval "$$(minikube docker-env)" && $(MAKE) $*-build-image

%-minikube-deploy:
	@kubectl apply -k apps/$*/deploy/k8s/overlays/dev/

%-minikube-destroy:
	@kubectl delete -k apps/$*/deploy/k8s/overlays/dev/

# Convenience targets for all services
.PHONY: dev generate build-image minikube-build-image minikube-deploy minikube-destroy

dev:
	@echo "Available services: $(SERVICES)"
	@echo "Usage: make SERVICE-dev (e.g., make users-dev)"

generate:
	@for service in $(SERVICES); do \
		echo "Generating for $$service..."; \
		$(MAKE) $$service-generate; \
	done

build-image:
	@for service in $(SERVICES); do \
		echo "Building image for $$service..."; \
		$(MAKE) $$service-build-image; \
	done

minikube-build-image:
	@for service in $(SERVICES); do \
		echo "Building Minikube image for $$service..."; \
		$(MAKE) $$service-minikube-build-image; \
	done

minikube-deploy:
	@for service in $(SERVICES); do \
		echo "Deploying $$service to Minikube..."; \
		$(MAKE) $$service-minikube-deploy; \
	done

minikube-destroy:
	@for service in $(SERVICES); do \
		echo "Destroying $$service from Minikube..."; \
		$(MAKE) $$service-minikube-destroy; \
	done

# Help target
help:
	@echo "Available targets:"
	@echo "  {SERVICE}-dev                 - Run service in development mode"
	@echo "  {SERVICE}-generate            - Run generate scripts for service"
	@echo "  {SERVICE}-build-image         - Build Docker image for service"
	@echo "  {SERVICE}-minikube-build-image - Build Docker image in Minikube environment"
	@echo "  {SERVICE}-minikube-deploy     - Deploy service to Minikube"
	@echo "  {SERVICE}-minikube-destroy    - Remove service from Minikube"
	@echo ""
	@echo "  generate                    - Run generate for all services"
	@echo "  build-image                 - Build images for all services"
	@echo "  minikube-build-image        - Build Minikube images for all services"
	@echo "  minikube-deploy             - Deploy all services to Minikube"
	@echo "  minikube-destroy            - Remove all services from Minikube"
	@echo ""
	@echo "Available services: $(SERVICES)"
	@echo "Examples:"
	@echo "  make users-dev"
	@echo "  make tasks-generate"
	@echo "  make build-image"