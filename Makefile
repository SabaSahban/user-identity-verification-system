#@IgnoreInspection BashAddShebang

export POSTGRES_ADDRESS=suleiman.db.elephantsql.com:5432
export POSTGRES_DATABASE=vnexexsq
export POSTGRES_USER=vnexexsq
export POSTGRES_PASSWORD=mtQWMkoyIZIHmJtZvKRoDTkL8ctSDBx9
export POSTGRES_DSN="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable"

migrate-create:
	migrate create -ext sql -dir ./migrations $(NAME)

migrate-up:
	migrate -verbose  -path ./migrations -database $(POSTGRES_DSN) up

migrate-down:
	 migrate -path ./migrations -database $(POSTGRES_DSN) down

migrate-reset:
	 migrate -path ./migrations -database $(POSTGRES_DSN) drop

run-server:
	go run . server

run-consumer:
	go run . consumer

up:
	docker-compose up