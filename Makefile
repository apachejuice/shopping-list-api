SHELL = /bin/bash

_:
	@echo "Specify a target to execute, or prod-all/test-all to make everything"

prod-all: oapi-generate prod-db
test-all: oapi-generate test-db

test-create-db:
	(source ".test_env" && sqlboiler mysql)

test-create-model:
	(source ".test_env" && sql-migrate up -env=development)

prod-create-db:
	(source ".env" && sqlboiler mysql)

prod-create-model:
	(source ".env" && sql-migrate up -env=production)

prod-db: prod-create-db prod-create-model

test-db: test-create-db test-create-model

test-schema:
	@(read -p "Enter schema name: " SCHEMA && sql-migrate new -env=development $$SCHEMA)

prod-schema:
	@(read -p "Enter schema name: " SCHEMA && sql-migrate new -env=production $$SCHEMA)

oapi-generate:
	oapi-codegen -generate "types,gin,spec" -o internal/apispec/spec.gen.go -package apispec spec/swagger.yaml 

swaggerui:
	./scripts/swagger-ui.sh

.PHONY: test-db prod-db test-schema prod-schema oapi-generate swaggerui
