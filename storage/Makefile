build:
	docker-compose build telegram-reminder-bot

run:
	@echo "Running Docker containers..."
	docker-compose up --build

test:
	@echo "Running tests..."
	go test -v ./...

db_init:
	@echo "Initializing database..."
	docker run --name tg-bot-db -e POSTGRES_PASSWORD='secret' -p 5436:5432 -d --rm postgres

db_stop:
	@echo "Stopping database..."
	docker stop tg-bot-db

migrate:
	@echo "Running database migrations..."
	migrate -path ./schema -database 'postgres://postgres:secret@0.0.0.0:5436/postgres?sslmode=disable' up

un_migrate:
	@echo "Rolling back database migrations..."
	migrate -path ./schema -database 'postgres://postgres:secret@0.0.0.0:5436/postgres?sslmode=disable' down

