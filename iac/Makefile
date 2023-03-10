SHELL := /bin/bash
.DEFAULT_GOAL := plan

# Use this makefile to ensure the lambda and dashboard have been built before provisioning

.PHONY: plan
plan: build-lambda
	./ter.sh plan

.PHONY: apply
apply: build-lambda
	./ter.sh apply

.PHONY: build-lambda
build-lambda:
	@echo Building the lambda project...
	cd lambda; make build
