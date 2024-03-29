GO_CODE_PATH=./lambda/...

.DEFAULT_GOAL:=help
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-17s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup

.PHONY: setup
setup: clean install build ## Set up for development

.PHONY: install
install: ## Install any specific tooling
ifeq ($(CI),true)
	npm ci
else
	npm install
endif
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	go install github.com/golang/mock/mockgen@latest
	go generate $(GO_CODE_PATH)

.PHONY: clean
clean: ## Clean the local filesystem
	rm -fr node_modules
	rm -fr cdk.out
	git clean -fdX


##@ Vet

.PHONY: vet
vet: vet-go lint-cdk prettier ## Vet the code

.PHONY: vet-go
vet-go: ## Vet the Go code
	@echo "Vet the Go code..."
	go vet -v $(GO_CODE_PATH)

.PHONY: lint-go
lint-go: ## Lint the Go code
	@echo "Lint the Go code..."
	./bin/golangci-lint run -v

.PHONY: lint-cdk
lint-cdk: ## Lint the CDK code
	@echo "Lint the CDK code..."
	npm run lint

.PHONY: prettier
prettier: ## Run Prettier
	@echo "Run Prettier"
	npx prettier --check .

##@ Build

.PHONY: build
build: builders build-cdk ## Build everything

DIRS=$(shell find lambda/* -type d)

.PHONY: builders $(DIRS)
builders: $(DIRS) ## Build all the underlying lambdas

$(DIRS): ## Build each lambda - https://www.wolfe.id.au/2023/08/09/rip-aws-go-lambda-runtime/#why-is-this-hard
	cd $@ && GOOS=linux CGO_ENABLED=0 go build -o bootstrap ./...

.PHONY: build-cdk
build-cdk: ## Build the CDK stacks
	npm run build


##@ Test

.PHONY: test
test: test-go test-cdk ## Run all the tests

.PHONY: test-go
test-go: ## Run the Go tests
	go test $(GO_CODE_PATH) -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: test-cdk
test-cdk: build-cdk ## Run the CDK tests
	npm run test


##@ Deployment

.PHONY: check-aws-details
check-aws-details: ## Check that the AWS details have been given
ifeq ($(AWS_ACCOUNT_ID),)
	@echo "[Error] Please specify an AWS_ACCOUNT_ID"
	@exit 1;
endif
ifeq ($(AWS_REGION),)
	@echo "[Error] Please specify an AWS_REGION"
	@exit 1;
endif

.PHONY: check-api-creds
check-api-creds: ## Check that the API credentials have been given
ifeq ($(USERNAME),)
	@echo "[Error] Please specify a USERNAME for the Traveline API"
	@exit 1;
endif
ifeq ($(PASSWORD),)
	@echo "[Error] Please specify a PASSWORD for the Traveline API"
	@exit 1;
endif

.PHONY: bootstrap
bootstrap: check-aws-details ## Bootstrap the CDK
	npx cdk bootstrap aws://${AWS_ACCOUNT_ID}/${AWS_REGION}

.PHONY: deploy
deploy: check-api-creds build bootstrap ## Create or update the infrastructure on AWS
	npx cdk deploy next-tram-stack \
		-c travelineApiUsername=${USERNAME} \
		-c travelineApiPassword=${PASSWORD} \
		-c naptanCode=${NAPTAN_CODE}
	./scripts/add-alexa-permission.sh

.PHONY: diff
diff: build ## Compare the infrastructure with stack on AWS
	npx cdk diff next-tram-stack

.PHONY: synth
synth: build ## Synthasise the infrastructure stack
	npx cdk synth next-tram-stack

.PHONY: destroy
destroy: ## Destroy the infrastructure in AWS
	npx cdk destroy next-tram-stack


##@ AWS SAM

.PHONY: sam-synth-cdk
sam-synth-cdk: check-api-creds build ## Synthasise the infrastructure stack for AWS SAM
	npx cdk synth next-tram-stack --no-staging \
		-c travelineApiUsername=${USERNAME} \
		-c travelineApiPassword=${PASSWORD} \
		-c naptanCode=${NAPTAN_CODE}

.PHONY: sam-local
sam-local: ## Run Lambda locally using AWS SAM
	sam local invoke GetNextTramLambda \
		-e sam-event.json \
		-t ./cdk.out/next-tram-stack.template.json
