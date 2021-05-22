GO_CODE_PATH=./src/...

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
setup: clean install build ## Set up for development

.PHONY: install
install: ## Install any specific tooling
ifeq ($(CI),true)
	npm ci
else
	npm install
endif
	go install github.com/golang/mock/mockgen@v1.5.0
	go generate $(GO_CODE_PATH)

.PHONY: clean
clean: ## Clean the local filesystem
	rm -fr node_modules
	rm -fr cdk.out
	git clean -fdX


##
## Vet targets
##

.PHONY: vet
vet: vet-go lint-cdk prettier ## Vet the code

.PHONY: vet-go
vet-go: ## Vet the Go code
	@echo "Vet the Go code..."
	go vet -v $(GO_CODE_PATH)

.PHONY: lint-go
lint-go: ## Lint the Go code
	@echo "Lint the Go code..."
	golangci-lint run -v

.PHONY: lint-cdk
lint-cdk: ## Lint the CDK code
	@echo "Lint the CDK code..."
	npm run lint

.PHONY: prettier
prettier: ## Run Prettier
	@echo "Run Prettier"
	npx prettier --check .

##
# Build targets
##

.PHONY: build
build: builders build-cdk ## Build everything

DIRS=$(shell find src/lambda/* -type d)

.PHONY: builders $(DIRS)
builders: $(DIRS) ## Build all the underlying lambdas

$(DIRS): ## Build each lambda
	cd $@ && GOOS=linux go build -o main ./...

.PHONY: build-cdk
build-cdk: ## Build the CDK stacks
	npm run build


##
# Test targets
##

.PHONY: test
test: test-go test-cdk ## Run all the tests

.PHONY: test-go
test-go: ## Run the Go tests
	go test $(GO_CODE_PATH) -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: test-cdk
test-cdk: build-cdk ## Run the CDK tests
	npm run test


##
# Deployment targets
##

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

.PHONY: check-api-credentials
check-api-credentials: ## Check that the API credentials have been given
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
deploy: check-api-credentials build bootstrap ## Create or update the infrastructure on AWS
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
