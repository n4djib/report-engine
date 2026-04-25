# Makefile
# Run any target with: make <target>


.PHONY: dev-up dev-down server client test lint secure migrate-up migrate-down

dev-up:
	docker compose -f deploy/compose/docker-compose.dev.yml up -d

dev-down:
	docker compose -f deploy/compose/docker-compose.dev.yml down

dev-down-clean:
	docker compose -f deploy/compose/docker-compose.dev.yml down -v

dev-logs:
	docker compose -f deploy/compose/docker-compose.dev.yml logs -f

server: 
	DATABASE_URL="postgres://reporter_dev:devpassword123@localhost:6432/report_dev" \
	NATS_URL="nats://localhost:4222" \
	go run ./cmd/server

client:
	cd cmd/desktop && wails dev -tags webkit2_41

test: 
	go test -race ./...

test-integration:
	go test -tags integration -timeout 120s -race ./...

lint:
	golangci-lint run ./...

secure:
	gosec -severity medium ./...

migrate-up:
	migrate -path ./migrations -database "postgres://reporter_dev:devpassword123@localhost:5432/report_dev?sslmode=disable" up

migrate-down:
	migrate -path ./migrations \
	-database "postgres://reporter_dev:devpassword123@localhost:5432/report_dev?sslmode=disable" \
	down 1

build-server:
	go build -o dist/report-server  ./cmd/server

build-client:
	./scripts/build-client.sh

