# KDex Tech Docs App Makefile

NPM_REGISTRY ?= http://npm.test/

.PHONY: all
all: build deploy

##@ Build

.PHONY: build
build: build-app

.PHONY: build-app
build-app: ## Increment version and publish the docs-app
	@echo "==> Incrementing version and publishing docs-app..."
	npm version patch --no-git-tag-version
	@APP_VERSION=$$(jq -r .version package.json); \
	echo "==> New version: $$APP_VERSION"; \
	sed -i "s/version: [0-9.]*/version: $$APP_VERSION/" ../../k8s/100_app_docs.yaml
	npm publish --registry $(NPM_REGISTRY)

##@ Deploy

.PHONY: deploy
deploy: ## Delete the docs-app pod
	@echo "==> Applying Kubernetes resources..."
	kubectl apply -f ../../k8s/100_app_docs.yaml
	@echo "==> Deleting docs-app pod..."
	kubectl -n kdex-site delete pod -l kdex.dev/backend=app-docs

##@ Helper

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
