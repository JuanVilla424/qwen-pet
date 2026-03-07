.PHONY: build run test lint clean setup ollama-setup docker-up docker-down

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

setup: build ollama-setup
	@echo "qwen-pet ready. Configure MCP in your editor and run: make run"

ollama-setup:
	docker compose up -d ollama
	@echo "Waiting for Ollama to start..."
	@sleep 3
	docker compose exec ollama ollama pull nomic-embed-text

docker-up:
	docker compose up -d

docker-down:
	docker compose down
