SHELL:=/bin/bash

AMB_CERT_FILE:=dist/managedblockchain-tls-chain.pem
# The AWS CLI profile to use
AWS_PROFILE:=saml

SRC_DIR:=./src
# Used to rebuild the project
SOURCE_FILES:=$(shell find "$(SRC_DIR)" -name '*.ts')

.PHONY: help
# Run "make" or "make help" to get a list of user targets
# Adapted from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?##.*$$' $(MAKEFILE_LIST) | awk 'BEGIN { \
	 FS = ":.*?## "; \
	 printf "\033[1m%-30s\033[0m %s\n", "TARGET", "DESCRIPTION" \
	} \
	{ printf "\033[32m%-30s\033[0m %s\n", $$1, $$2 }'

# Built in a way that does not update modification times unless necessary (prevents spuriously triggered re-uploads)
.PHONY: build
build: $(SOURCE_FILES) node_modules dist/node_modules $(AMB_CERT_FILE) ## Build the lambda project
	npx tsc --build

AMB_CERT_S3_URL:=s3://us-east-1.managedblockchain/etc/managedblockchain-tls-chain.pem

$(AMB_CERT_FILE):
	@echo Downloading Amazon Managed Blockchain TLS Certificate
	./../../util/test_reauthenticate_aws_saml.sh
	mkdir -p "$(dir $(AMB_CERT_FILE))"
	AWS_PROFILE="$(AWS_PROFILE)" aws s3 cp "$(AMB_CERT_S3_URL)" "$(AMB_CERT_FILE)"

.PHONY: dependencies
dependencies: node_modules ## Download project dependencies

node_modules: package.json package-lock.json
	@echo Downloading dependencies...
	npm ci

# Download runtime dependencies to the node_modules folder of the output artifact
dist/node_modules: package.json package-lock.json
	@echo Downloading lambda runtime dependencies
	mkdir -p dist/node_modules
	cp package.json package-lock.json dist/
	npm ci --prefix dist/ --omit=dev
	rm dist/{package,package-lock}.json

.PHONY: clean
clean: ## Clean the project of build and dependency artifacts
	rm -fr dist node_modules
