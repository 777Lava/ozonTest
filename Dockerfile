FROM golang AS builder

WORKDIR /app

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download
RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . ./
RUN go build -o ./bin/main ./cmd

# wait 2 seconds for init db
CMD sleep 2 && go test ./internal/database/tests && migrate -database ${POSTGRESQL_URL} -path sql/migrations up && ./bin/main 

EXPOSE 8000
