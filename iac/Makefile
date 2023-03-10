SHELL := /bin/bash
.DEFAULT_GOAL := plan

# Use this makefile to ensure the lambda and dashboard have been built before provisioning

.PHONY: plan
plan: build-lambda build-dashboard
	./ter.sh plan

.PHONY: apply
apply: build-lambda build-dashboard
	./ter.sh apply

.PHONY: build-lambda
build-lambda:
	@echo Building the lambda project...
	cd lambda; make build

.PHONY: build-dashboard
build-dashboard:
	@echo Building the dashboard project...
	cd ../dashboard; make build

.PHONY: clean
clean:
	cd lambda; make clean
	cd ../dashboard; make clean
