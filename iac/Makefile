SHELL:=/bin/bash

AWS_PROFILE:=saml
TERRAFORM_CMD:=./ter.sh

LAMBDA_DIR:=./lambda
DASHBOARD_DIR:=../dashboard

# Use this makefile to ensure the lambda and dashboard have been built before provisioning

.PHONY: help
# Run "make" or "make help" to get a list of user targets
# Adapted from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?##.*$$' $(MAKEFILE_LIST) | awk 'BEGIN { \
	 FS = ":.*?## "; \
	 printf "\033[1m%-30s\033[0m %s\n", "TARGET", "DESCRIPTION" \
	} \
	{ printf "\033[32m%-30s\033[0m %s\n", $$1, $$2 }'

.PHONY: plan
plan: build ## Preview the Terraform plan
	$(TERRAFORM_CMD) plan

.PHONY: apply
apply: build ## Apply the Terraform plan
	$(TERRAFORM_CMD) apply

.PHONY: refresh
refresh: ## Refresh Terraform outputs
	@echo Refreshing terraform outputs...
	$(TERRAFORM_CMD) refresh

.PHONY: build
build: build-lambda build-dashboard ## Build all dependencies

.PHONY: clean
clean: clean-dashboard clean-lambda clean-connection-profile ## Clean dependenies of all artifacts

# Dashboard

DASHBOARD_DEVENV_FILE:=$(DASHBOARD_DIR)/.env.development
DASHBOARD_PRODENV_FILE:=$(DASHBOARD_DIR)/.env.prod

.PHONY: build-dashboard
build-dashboard: $(DASHBOARD_PRODENV_FILE) ## Build the dashboard
	@echo Building the dashboard project...
	$(MAKE) -C "$(DASHBOARD_DIR)" build AWS_PROFILE="$(AWS_PROFILE)"

.PHONY: environment-dashboard
environment-dashboard: $(DASHBOARD_DEVENV_FILE) $(DASHBOARD_PRODENV_FILE) ## Produce dashboard environment files

$(DASHBOARD_DEVENV_FILE):
	$(TERRAFORM_CMD) output -raw vite_dev_env > "$(DASHBOARD_DEVENV_FILE)"

$(DASHBOARD_PRODENV_FILE):
	$(TERRAFORM_CMD) output -raw vite_prod_env > "$(DASHBOARD_PRODENV_FILE)"

.PHONY: clean-dashboard
clean-dashboard: ## Clean dashboard build artifacts and environment files
	$(MAKE) -C "$(DASHBOARD_DIR)" clean
	rm -f "$(DASHBOARD_DEVENV_FILE)"
	rm -f "$(DASHBOARD_PRODENV_FILE)"

# Lambda

.PHONY: build-lambda
build-lambda: ## Build the lambda source
	@echo Building the lambda project...
	$(MAKE) -C "$(LAMBDA_DIR)" build AWS_PROFILE="$(AWS_PROFILE)"

.PHONY: clean-lambda
clean-lambda: ## Clean lambda build artifacts
	$(MAKE) -C "$(LAMBDA_DIR)" clean

# Connection Profile

CONNECTION_PROFILE_FILE:=$(shell $(TERRAFORM_CMD) output -raw connection_profile_file)

.PHONY: connection-profile
connection-profile: $(CONNECTION_PROFILE_FILE) ## Generate the connection profile for the given workspace

UTIL_DIR:=../util
VIRTUAL_ENV_ACTIVATION:=$(UTIL_DIR)/venv/bin/activate

$(VIRTUAL_ENV_ACTIVATION):
	python3 -m venv "$(UTIL_DIR)/venv"
	source "$(VIRTUAL_ENV_ACTIVATION)" && \
		python3 -m pip install -r "$(UTIL_DIR)/requirements.txt"

$(CONNECTION_PROFILE_FILE): $(VIRTUAL_ENV_ACTIVATION)
	source $(VIRTUAL_ENV_ACTIVATION) && \
		AWS_PROFILE=$(AWS_PROFILE) $(UTIL_DIR)/gen-connection-profile.py \
			--network_id "$(shell $(TERRAFORM_CMD) output -raw network_id)" \
			--channels "$(shell $(TERRAFORM_CMD) output -raw channel_name)" \
			--tlsCaCertPath "./managedblockchain-tls-chain.pem" \
		> $(CONNECTION_PROFILE_FILE)

.PHONY: clean-connection-profile
clean-connection-profile: ## Clean the connection profile
	rm -f $(CONNECTION_PROFILE_FILE)
