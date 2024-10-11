build:
	docker-compose build telegram_reminder_bot

run:
	docker-compose up telegram_reminder_bot

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:secret@0.0.0.0:5436/postgres?sslmode=disable' up
