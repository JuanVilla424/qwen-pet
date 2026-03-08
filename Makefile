.PHONY: build run test lint clean setup ollama-setup docker-up docker-down
.PHONY: install uninstall start stop restart status logs

BINARY := bin/pet
SRC := ./cmd/pet

# Install paths (XDG standard)
PREFIX       := $(HOME)/.local
BINDIR       := $(PREFIX)/bin
CONFIG_DIR   := $(HOME)/.config/qwen-pet
DATA_DIR     := $(PREFIX)/share/qwen-pet
SERVICE_NAME := qwen-pet
SERVICE_DIR  := $(HOME)/.config/systemd/user
SERVICE_TPL  := deployments/$(SERVICE_NAME).service.tpl
SERVICE_FILE := $(SERVICE_DIR)/$(SERVICE_NAME).service
ENV_FILE     := $(CONFIG_DIR)/.env

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

# --- Service management ---

install: build
	@echo "Installing qwen-pet..."
	install -d $(BINDIR)
	install -m 755 $(BINARY) $(BINDIR)/qwen-pet
	install -d $(CONFIG_DIR)
	test -f $(CONFIG_DIR)/pet.yaml || cp configs/pet.yaml $(CONFIG_DIR)/pet.yaml
	test -f $(ENV_FILE) || cp .env.template $(ENV_FILE)
	@test -s $(ENV_FILE) || echo "WARNING: Fill credentials in $(ENV_FILE)"
	install -d $(DATA_DIR)
	test -f data/pet-profile.json && cp -n data/pet-profile.json $(DATA_DIR)/ || true
	test -f data/pet.json && cp -n data/pet.json $(DATA_DIR)/ || true
	test -d data/chromem && cp -rn data/chromem $(DATA_DIR)/ || true
	install -d $(SERVICE_DIR)
	sed -e 's|{{BINARY}}|$(BINDIR)/qwen-pet|g' \
	    -e 's|{{CONFIG}}|$(CONFIG_DIR)/pet.yaml|g' \
	    -e 's|{{DATA_DIR}}|$(DATA_DIR)|g' \
	    -e 's|{{ENV_FILE}}|$(ENV_FILE)|g' \
	    $(SERVICE_TPL) > $(SERVICE_FILE)
	systemctl --user daemon-reload
	systemctl --user enable $(SERVICE_NAME)
	loginctl enable-linger $(USER)
	@echo ""
	@echo "=== Installation complete ==="
	@echo "1. Edit credentials: $(ENV_FILE)"
	@echo "2. Start service:    make start"
	@echo "3. MCP config for Claude Code:"
	@echo '   {'
	@echo '     "mcpServers": {'
	@echo '       "qwen-pet": {'
	@echo '         "command": "$(BINDIR)/qwen-pet",'
	@echo '         "env": {'
	@echo '           "PET_CONFIG": "$(CONFIG_DIR)/pet.yaml",'
	@echo '           "PET_DATA_DIR": "$(DATA_DIR)"'
	@echo '         }'
	@echo '       }'
	@echo '     }'
	@echo '   }'

uninstall:
	systemctl --user stop $(SERVICE_NAME) || true
	systemctl --user disable $(SERVICE_NAME) || true
	rm -f $(SERVICE_FILE)
	rm -f $(BINDIR)/qwen-pet
	systemctl --user daemon-reload
	@echo "Uninstalled. Config/data preserved in $(CONFIG_DIR) and $(DATA_DIR)"

start:
	systemctl --user start $(SERVICE_NAME)

stop:
	systemctl --user stop $(SERVICE_NAME)

restart: build
	install -m 755 $(BINARY) $(BINDIR)/qwen-pet
	systemctl --user restart $(SERVICE_NAME)

status:
	systemctl --user status $(SERVICE_NAME)

logs:
	journalctl --user -u $(SERVICE_NAME) -f
