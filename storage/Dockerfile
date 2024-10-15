FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# Устанавливаем временную зону
RUN apt-get update && apt-get install -y tzdata

# install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz -o migrate.tar.gz
RUN tar -xvzf migrate.tar.gz
RUN mv migrate /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o telegram-reminder-bot ./cmd/main.go

CMD ["./telegram-reminder-bot"]