# KDex Tech Main Site Holistic Makefile

REPOSITORY ?= k3d-registry:5000
PLATFORMS ?= linux/amd64
NAMESPACE ?= kdex-site
NPM_REGISTRY ?= http://npm.test/

.PHONY: all
all: build deploy

##@ Build

.PHONY: build
build: build-app build-content build-theme ## Build all components (NPM app and OCI images)

.PHONY: build-app
build-app: ## Increment version, update CR and publish the docs-app
	@echo "==> Patching version and publishing docs-app..."
	cd docs-app/app && npm version patch --no-git-tag-version
	@APP_VERSION=$$(jq -r .version docs-app/app/package.json); \
	echo "==> New version: $$APP_VERSION"; \
	sed -i "s/version: [0-9.]*/version: $$APP_VERSION/" k8s/100_app_docs.yaml
	cd docs-app/app && npm publish --registry $(NPM_REGISTRY)

.PHONY: build-content
build-content: ## Build and push the documentation content image
	@echo "==> Building content image..."
	cd docs-app/content && $(MAKE) docker-buildx REPOSITORY=$(REPOSITORY) PLATFORMS=$(PLATFORMS)
	@echo "==> Capturing content digest..."
	@DIGEST=$$(docker buildx imagetools inspect $(REPOSITORY)/kdex-tech/kdex-docs-content:latest --format '{{json .}}' | jq -r '.manifest.digest // .container_config.digest // .digest'); \
	echo "==> Content digest: $$DIGEST"; \

##	sed -i "s|staticImage: .*/kdex-docs-content:latest.*|staticImage: $(REPOSITORY)/kdex-tech/kdex-docs-content:latest@$$DIGEST|" k8s/100_app_docs.yaml

.PHONY: build-theme
build-theme: ## Build and push the theme image
	@echo "==> Building theme image..."
	cd kdex-main-theme && $(MAKE) docker-buildx REPOSITORY=$(REPOSITORY) PLATFORMS=$(PLATFORMS)
	@echo "==> Capturing theme digest..."
	@DIGEST=$$(docker buildx imagetools inspect $(REPOSITORY)/kdex-tech/kdex-main-theme:latest --format '{{json .}}' | jq -r '.manifest.digest // .container_config.digest // .digest'); \
	echo "==> Theme digest: $$DIGEST"; \

##	sed -i "s|staticImage: .*/kdex-main-theme:latest.*|staticImage: $(REPOSITORY)/kdex-tech/kdex-main-theme:latest@$$DIGEST|" k8s/030_theme.yaml

##@ Deploy

.PHONY: deploy
deploy: ## Apply Kubernetes resources to the cluster
	@echo "==> Applying Kubernetes resources..."
	kubectl apply -f k8s/
	kubectl -n kdex-site delete pod -l kdex.dev/backend=app-docs
	kubectl -n kdex-site delete pod -l kdex.dev/backend=kdex-main-theme

##@ Helper

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
