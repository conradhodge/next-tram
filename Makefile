AWS_ACCOUNT_ID=431809455209
AwS_REGION=us-east-1

.DEFAULT_GOAL := explain
.PHONY: explain
explain:
	#### Next tram
	#   _  _                     _                _
	#  | \| |    ___    __ __   | |_      o O O  | |_      _ _   __ _    _ __
	#  | .` |   / -_)   \ \ /   |  _|    o       |  _|    | '_| / _` |  | '  \
	#  |_|\_|   \___|   /_\_\   _\__|   TS__[O]  _\__|   _|_|_  \__,_|  |_|_|_|
	# _|"""""|_|"""""|_|"""""|_|"""""| {======|_|"""""|_|"""""|_|"""""|_|"""""|
	# "`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'./o--000'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'
	#
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
# Setup targets
##

.PHONY: setup
setup: clean install vet build test ## Set up for development

.PHONY: install
install: ## Install any specific tooling
	npm install
	go get golang.org/x/lint/golint
	go get github.com/securego/gosec/cmd/gosec

.PHONY: clean
clean: ## Clean the local filesystem
	rm -fr node_modules
	rm -fr cdk.out
	git clean -fdx


##
## Vet targets
##

.PHONY: vet
vet: vet-go ## Vet the code

.PHONY: vet-go
vet-go: ## Vet the Go code
	@echo "Vet the code..."
	go vet -v ./...

	@echo "Lint the code..."
	$$GOPATH/bin/golint -set_exit_status $(shell go list ./...)

	@echo "Inspect code for security vulnerabilities..."
	$$GOPATH/bin/gosec -exclude-dir build ./...


##
# Build targets
##

.PHONY: build
build: builders build-cdk ## Build everything

DIRS=$(shell find src/lambda/* -type d)

.PHONY: builders $(DIRS)
builders: $(DIRS) ## Build all the underlying lambdas

$(DIRS): ## Build each lambda and zip up
	cd $@ && GOOS=linux go build -o main ./...
	cd $@ && zip handler.zip ./main

.PHONY: build-cdk
build-cdk: ## Build the CDK stacks
	npm run build


##
# Test targets
##

.PHONY: test
test: test-cdk ## Run all the tests

.PHONY: test-cdk
test-cdk: ## Run the CDK tests
	npm run test


##
# Deployment targets
##

.PHONY: bootstrap
bootstrap: ## Bootstrap the CDK
	npx cdk bootstrap aws://${AWS_ACCOUNT_ID}/${AwS_REGION}

.PHONY: deploy
deploy: build bootstrap ## Create or update the infrastructure on AWS
	npx cdk --app "npx ts-node ./infrastructure/bin/next-tram.ts" deploy next-tram-stack

.PHONY: destroy
destroy: build ## Destroy the infrastructure in AWS
	npx cdk --app "npx ts-node ./infrastructure/bin/next-tram.ts" destroy next-tram-stack
