.PHONY: build

build:
	sam build

dev: build
	sam local start-api --skip-pull-image --env-vars env.dev.json

deploy: build
	sam package /
	sam deploy