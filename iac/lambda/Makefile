SHELL := /bin/bash
.DEFAULT_GOAL := build

amb_cert_file := dist/managedblockchain-tls-chain.pem
aws_profile := saml

.PHONY: build
build: node_modules prebuild
	npx tsc; cp -r node_modules dist

.PHONY: prebuild
prebuild: $(amb_cert_file) dist/node_modules
	rm -f dist/*.js

$(amb_cert_file):
	@echo Downloading Amazon Managed Blockchain TLS Certificate
	./../../util/test_reauthenticate_aws_saml.sh
	mkdir dist
	AWS_PROFILE=saml aws s3 cp s3://us-east-1.managedblockchain/etc/managedblockchain-tls-chain.pem dist/

node_modules: package.json package-lock.json
	@echo Downloading dependencies...
	npm ci

dist/node_modules: package.json package-lock.json
	@echo Downloading lambda runtime dependencies
	mkdir -p dist/node_modules
	cp package.json package-lock.json dist/node_modules
	npm install --prefix dist/node_modules --omit=dev
	rm dist/node_modules/{package,package-lock}.json

.PHONY: clean
clean:
	rm -fr dist node_modules
