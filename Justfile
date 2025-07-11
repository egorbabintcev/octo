goose-create name:
	docker run --rm \
	  -v ./internal/migrations:/migrations \
	  -e GOOSE_COMMAND="create" \
	  -e GOOSE_COMMAND_ARG="{{name}} sql" \
	  ghcr.io/kukymbr/goose-docker:latest

sqlc-generate:
	docker run --rm \
		-v $(pwd):/src \
		-w /src \
		sqlc/sqlc generate
