build:
	docker-compose build telegram-reminder-bot


run:
	docker-compose up telegram-reminder-bot


test:
	go test -v ./...

db_init:
	docker run --name tg-bot-db -e POSTGRES_PASSWORD='secret' -p 5436:5432 -d --rm postgres

db_stop:
	docker stop tg-bot-db

migrate:
	migrate -path ./schema -database 'postgres://postgres:secret@0.0.0.0:5436/postgres?sslmode=disable' up

un_migrate:
	migrate -path ./schema -database 'postgres://postgres:secret@0.0.0.0:5436/postgres?sslmode=disable' down
