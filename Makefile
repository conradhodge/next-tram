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

.PHONY: install
install: ## Install any specific tooling
	npm install

.PHONY: clean
clean: ## Clean the local filesystem
	rm -fr node_modules
	rm -fr cdk.out
	git clean -fdx


##
# Build targets
##

.PHONY: build
build: build-cdk ## Build everything

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

.PHONY: deploy
deploy: build ## Create or update the infrastructure on AWS
	npx cdk --app "npx ts-node ./infrastructure/bin/next-tram.ts" deploy next-tram-stack

.PHONY: destroy
destroy: build ## Destroy the infrastructure in AWS
	npx cdk --app "npx ts-node ./infrastructure/bin/next-tram.ts" destroy next-tram-stack
