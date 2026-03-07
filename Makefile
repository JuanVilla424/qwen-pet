.PHONY: build run test lint clean ollama-setup docker-up docker-down

BINARY := bin/pet
SRC := ./cmd/pet

build:
	go build -o $(BINARY) $(SRC)

run: build
	@test -f .env && set -a && . ./.env && set +a; $(BINARY)

test:
	go test ./internal/...

lint:
	go vet ./...
	gofmt -l .

clean:
	rm -rf bin/ data/

ollama-setup:
	ollama pull nomic-embed-text

docker-up:
	docker compose up -d

docker-down:
	docker compose down
