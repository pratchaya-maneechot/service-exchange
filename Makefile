# Makefile - Simple & Maintainable
SERVICES := users tasks
SERVICE ?= users
NAME ?= default_migration
TOOLS_SCRIPTS_DIR := tools/scripts

# Database config
DB_URL = postgres://root:Hl7FudwaSNzOhhioo0GxlmmMD0LM+I8StQIqJCZ1TPg@localhost:5422 # Changed port to 5422 to avoid conflict with default


.PHONY: help $(SERVICES)

# Service-specific targets
%-dev:
	@echo "Starting $* service..."
	@go run ./apps/$*/cmd/api/main.go

%-generate:
	@echo "Generating $* code..."
	@[ -f $(TOOLS_SCRIPTS_DIR)/generate-proto.sh ] && sh $(TOOLS_SCRIPTS_DIR)/generate-proto.sh $* || true
	@[ -f $(TOOLS_SCRIPTS_DIR)/generate-wire.sh ] && sh $(TOOLS_SCRIPTS_DIR)/generate-wire.sh $* || true
	@[ -d ./apps/$*/internal/infra/persistence/postgres ] && cd ./apps/$*/internal/infra/persistence/postgres && sqlc generate || true

%-build-image:
	@echo "Building $* image..."
	@docker build -t $*-service:latest -f apps/$*/deploy/docker/Dockerfile .

%-migrate-create:
	@echo "Creating migration $(NAME) for $*..."
	@migrate create -ext sql -dir ./apps/$*/internal/infra/persistence/postgres/migrations -seq $(NAME)

%-migrate-up:
	@echo "Running $* migrations..."
	@migrate -path ./apps/$*/internal/infra/persistence/postgres/migrations -database "$(DB_URL)/$*?sslmode=disable" up

%-k8s-deploy:
	@echo "Deploying $* to k8s..."
	@kubectl apply -k apps/$*/deploy/k8s/overlays/dev/

%-k8s-destroy:
	@echo "Removing $* from k8s..."
	@kubectl delete -k apps/$*/deploy/k8s/overlays/dev/

# Bulk operations
generate:
	@for s in $(SERVICES); do $(MAKE) $$s-generate; done

build-images:
	@for s in $(SERVICES); do $(MAKE) $$s-build-image; done

migrate-up:
	@for s in $(SERVICES); do $(MAKE) $$s-migrate-up; done

k8s-deploy:
	@for s in $(SERVICES); do $(MAKE) $$s-k8s-deploy; done

k8s-destroy:
	@for s in $(SERVICES); do $(MAKE) $$s-k8s-destroy; done

# Utility
dev:
	@echo "Available: $(SERVICES)"
	@echo "Usage: make <service>-dev"

help:
	@echo "Services: $(SERVICES)"
	@echo ""
	@echo "Commands:"
	@echo "  <service>-dev           - Run service"
	@echo "  <service>-generate      - Generate code"
	@echo "  <service>-build-image   - Build Docker image"
	@echo "  <service>-migrate-create NAME=<name> - Create migration"
	@echo "  <service>-migrate-up    - Run migrations (make sure the database \"<service>\" already exists)"
	@echo "  <service>-k8s-deploy    - Deploy to k8s"
	@echo "  <service>-k8s-destroy   - Remove from k8s"
	@echo ""
	@echo "Bulk:"
	@echo "  generate, build-images, migrate-up, k8s-deploy, k8s-destroy"
	@echo ""
	@echo "Examples: make users-dev, make generate"