mock:
	mockery --all --keeptree

migrate:
	migrate -source file://internal/migrations \
			-database postgres://postgres:bert@127.0.0.1:5432/graphql_chat?sslmode=disable up

rollback:
	migrate -source file://internal/migrations \
            			 -database postgres://postgres:bert@127.0.0.1:5432/graphql_chat?sslmode=disable down 1


drop:
	migrate -source file://internal/migrations \
            			-database postgres://postgres:bert@127.0.0.1:5432/graphql_chat?sslmode=disable drop

#migration:
#	@if [ -z "$(name)" ]; then \
#		echo "Migration name is required. Usage: make migration name=<migration_name>"; \
#		exit 1; \
#	fi; \
#	migrate create -ext sql -dir internal/migrations $(name)

migration:
	@read -p "Enter migration name: " name; \
  	migrate create -ext sql -dir internal/migrations $$name

generate:
	go get github.com/99designs/gqlgen
#    go run github.com/99designs/gqlgen generate
	go generate ./...

run:
	go run ./cmd/main.go

