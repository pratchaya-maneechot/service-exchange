# Makefile - Simple & Maintainable
SERVICES := users
SERVICE ?= users
NAME ?= default_migration
TOOLS_SCRIPTS_DIR := tools/scripts

# env
export DB_URL = postgres://root:Hl7FudwaSNzOhhioo0GxlmmMD0LM+I8StQIqJCZ1TPg@localhost:5432

.PHONY: help $(SERVICES)

# Service-specific targets
%-dev:
	@echo "Starting $* service..."
	@go run ./apps/$*/cmd/api/main.go

%-generate:
	@echo "Generating $* code..."
	@[ -f $(TOOLS_SCRIPTS_DIR)/generate-proto.sh ] && sh $(TOOLS_SCRIPTS_DIR)/generate-proto.sh $* || true
	@[ -f $(TOOLS_SCRIPTS_DIR)/generate-wire.sh ] && sh $(TOOLS_SCRIPTS_DIR)/generate-wire.sh $* || true
	@[ -f $(TOOLS_SCRIPTS_DIR)/generate-sqlc.sh ] && sh $(TOOLS_SCRIPTS_DIR)/generate-sqlc.sh $* || true

%-sync-proto-gateway:
	@echo "Running sync $* proto to api-gateway..."
	@cp apps/$*/api/proto/**/*.proto apps/api-gateway/proto
	@echo "Completed to sync $* proto to api-gateway."
	
%-build-image:
	@echo "Building $* image..."
	@docker build -t $*-service:latest -f apps/$*/deploy/docker/Dockerfile .

%-migrate-create:
	@echo "Creating migration $(NAME) for $*..."
	@migrate create -ext sql -dir ./apps/$*/internal/infra/persistence/postgres/migrations -seq $(NAME)

%-migrate-up:
	@echo "Running $* migrations..."
	@[ -f $(TOOLS_SCRIPTS_DIR)/migrate-db-up.sh ] && sh $(TOOLS_SCRIPTS_DIR)/migrate-db-up.sh $* || true

%-k8s-deploy:
	@echo "Deploying $* to k8s..."
	@kubectl apply -k apps/$*/deploy/k8s/overlays/dev/

%-k8s-destroy:
	@echo "Removing $* from k8s..."
	@kubectl delete -k apps/$*/deploy/k8s/overlays/dev/

# Bulk operations
generate:
	@for s in $(SERVICES); do $(MAKE) $$s-generate; done

sync-proto-gateway:
	@for s in $(SERVICES); do $(MAKE) $$s-sync-proto-gateway; done

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