# 🤖 Qwen PET — Personal Entertainment Tool

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=fff)
![Qwen](https://img.shields.io/badge/Qwen-3.5_Flash-7C3AED)
![MCP](https://img.shields.io/badge/MCP-Server-blue)
![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?logo=telegram&logoColor=fff)
![Ollama](https://img.shields.io/badge/Ollama-Embeddings-333?logo=ollama)
![Build Status](https://github.com/JuanVilla424/qwen-pet/actions/workflows/ci.yml/badge.svg?branch=main)
![Status](https://img.shields.io/badge/Status-Development-yellow.svg)
![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)

**Qwen PET** is a persistent AI intermediary agent that works as a virtual pet (tamagotchi). It accumulates developer decisions, preferences, and patterns in a semantic KB. When it knows the answer, it responds autonomously with the chosen pet's personality; when it doesn't, it escalates via Telegram.

Integrates with Claude Code and OpenCode via MCP (Model Context Protocol) as a stdio server.

## 📚 Table of Contents

- [Features](#-features)
- [Available Pets](#-available-pets)
- [Personality System](#-personality-system)
- [Architecture](#-architecture)
- [Quickstart](#-quickstart)
- [MCP Configuration](#-mcp-configuration)
- [MCP Tools](#-mcp-tools)
- [Telegram Commands](#-telegram-commands)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [License](#-license)

## 🌟 Features

- **🧠 Knowledge Base** — Persistent vector store (chromem-go) with preferences, decisions, and patterns
- **🔌 MCP Server** — stdio transport, integrates with Claude Code and OpenCode as a tool provider
- **🤖 Decision Engine** — KB → Qwen 3.5 (OpenRouter) → Telegram pipeline with automatic confidence routing
- **🐾 Virtual Pet** — 18 selectable animals, each with unique visual identity (emojis, sounds, moods)
- **🎭 Configurable Personality** — 27 combinable traits (15 virtues + 12 defects) that affect ALL responses
- **💬 Telegram Bot** — Escalation channel when the pet doesn't know the answer
- **💾 Persistence** — Pet state (mood, stats) + semantic KB across sessions

## 🐾 Available Pets

| Animal    | Emoji     | Default | Animal     | Emoji     | Default |
| --------- | --------- | ------- | ---------- | --------- | ------- |
| 🐱 cat    | `U+1F431` | Michi   | 🐬 dolphin | `U+1F42C` | Finn    |
| 🐶 dog    | `U+1F436` | Rex     | 🦎 lizard  | `U+1F98E` | Gecko   |
| 🦊 fox    | `U+1F98A` | Kit     | 🐙 octopus | `U+1F419` | Inky    |
| 🦉 owl    | `U+1F989` | Archie  | 🐰 rabbit  | `U+1F430` | Bun     |
| 🐉 dragon | `U+1F409` | Drakar  | 🐧 penguin | `U+1F427` | Tux     |
| 🦁 lion   | `U+1F981` | Leo     | 🐍 snake   | `U+1F40D` | Slyth   |
| 🐺 wolf   | `U+1F43A` | Fenrir  | 🐼 panda   | `U+1F43C` | Bamboo  |
| 🐻 bear   | `U+1F43B` | Oso     | 🦄 unicorn | `U+1F984` | Sparkle |
| 🦅 eagle  | `U+1F985` | Aquila  | 🐸 frog    | `U+1F438` | Ribbit  |

Each animal has unique mood emojis and sounds for 6 states: 😊 happy, 😐 neutral, 🤔 thinking, 😴 tired, 🤩 excited, 😢 sad.

## 🎭 Personality System

Animals are **visual identity only**. Personality is built by freely combining traits:

**✨ Virtues**: curious, analytical, creative, patient, enthusiastic, loyal, strategic, honest, protective, adaptable, meticulous, humorous, empathetic, pragmatic, assertive

**💀 Defects**: impatient, stubborn, overthinks, sarcastic, forgetful, blunt, perfectionist, anxious, lazy, dramatic, distracted, competitive

Configuration in `configs/pet.yaml`:

```yaml
pet:
  type: "fox"
  name: "Kit"
  traits:
    - curious
    - creative
    - honest
    - sarcastic # defect
    - impatient # defect
```

Defects manifest subtly in responses, not as caricature.

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐
│ Claude Code │     │  OpenCode   │
└──────┬──────┘     └──────┬──────┘
       │    MCP (stdio)    │
       └────────┬──────────┘
                │
       ┌────────▼────────┐
       │   🤖 Qwen PET   │
       │   MCP Server    │
       ├─────────────────┤
       │ 🧠 Decision     │
       │    Engine       │
       │  🎭 Personality │
       ├─────┬─────┬─────┤
       │ 📚  │ 🤖  │ 💬  │
       │ KB  │ AI  │ TG  │
       └──┬──┴──┬──┴──┬──┘
          │     │     │
       Ollama OpenRouter Telegram
       embeds  inference  escalation
```

**📊 Decision pipeline:**

1. 📚 Query KB (chromem-go) — if similarity > 0.82, responds directly with personality
2. 🤖 AI (Qwen 3.5 Flash via OpenRouter) — generates response with personality prompt, stores in KB
3. 💬 Telegram — escalates to user, stores response in KB for future reference

## 🚀 Quickstart

### 📋 Prerequisites

- **Go 1.25+**
- **Docker** with Docker Compose (for Ollama)
- **OpenRouter account** — [openrouter.ai](https://openrouter.ai) for API key
- **Telegram Bot** — Create with [@BotFather](https://t.me/BotFather)

### 🔨 Installation

```bash
# Clone
git clone --recurse-submodules https://github.com/JuanVilla424/qwen-pet.git
cd qwen-pet

# Configure credentials
cp .env.template .env
# Edit .env: OPENROUTER_API_KEY, TELEGRAM_BOT_TOKEN, TELEGRAM_USER_ID

# Full setup (build + Ollama + embeddings model)
make setup

# Run
make run
```

### 🔧 Manual step-by-step setup

```bash
# 1. Build
make build

# 2. Start Ollama and download embeddings model
make ollama-setup

# 3. Run (loads .env automatically)
make run
```

## 🔌 MCP Configuration

### Claude Code

Add to `~/.claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "qwen-pet": {
      "command": "/path/to/qwen-pet/bin/pet",
      "env": {
        "PET_CONFIG": "/path/to/qwen-pet/configs/pet.yaml",
        "OPENROUTER_API_KEY": "sk-or-...",
        "TELEGRAM_BOT_TOKEN": "...",
        "TELEGRAM_USER_ID": "..."
      }
    }
  }
}
```

### OpenCode

Add to MCP servers configuration with stdio transport pointing to the binary.

## 🛠️ MCP Tools

| Tool               | Description                                                                        | Parameters                                      |
| ------------------ | ---------------------------------------------------------------------------------- | ----------------------------------------------- |
| `ask_pet`          | 🤔 Ask the pet a question. Searches KB, reasons with AI, or escalates via Telegram | `question`, `context` (optional)                |
| `check_preference` | 📋 Check a specific preference from the KB                                         | `category`, `key`                               |
| `store_decision`   | 💾 Store a decision in the KB for future reference                                 | `category`, `key`, `value`, `reason` (optional) |
| `approve_artifact` | ✅ Request user approval for an artifact via Telegram                              | `type`, `description`, `file_path` (optional)   |

## 💬 Telegram Commands

| Command   | Description                      |
| --------- | -------------------------------- |
| `/start`  | 👋 Pet greeting with personality |
| `/status` | 📊 Current mood + statistics     |
| `/pet`    | 🐾 Selected animal info          |
| `/help`   | ❓ Command list                  |

When the pet escalates a question, it sends a message with context and waits for the user's response. The response is automatically stored in the KB.

## 📁 Project Structure

```
qwen-pet/
├── cmd/pet/main.go              # 🚀 Wire + entrypoint
├── internal/
│   ├── config/config.go         # ⚙️ YAML + env vars config
│   ├── pet/
│   │   ├── catalog.go           # 🐾 18 animals + 27 traits
│   │   ├── personality.go       # 🎭 Build + WrapPrompt + FormatResponse
│   │   ├── state.go             # 💾 Mood, stats, JSON persistence
│   │   └── images.go            # 🖼️ Twemoji CDN fetcher
│   ├── kb/
│   │   ├── store.go             # 📚 chromem-go vector store wrapper
│   │   └── embeddings.go        # 🧮 Ollama embedding factory
│   ├── ai/
│   │   ├── client.go            # 🌐 OpenRouter HTTP client
│   │   └── decision.go          # 🧠 KB -> AI -> Telegram pipeline
│   ├── mcp/
│   │   ├── server.go            # 🔌 MCP server stdio transport
│   │   └── tools.go             # 🛠️ 4 tools with typed handlers
│   └── telegram/
│       ├── bot.go               # 💬 Polling + escalation
│       └── handlers.go          # 📨 Commands + default handler
├── configs/pet.yaml             # ⚙️ Default configuration
├── Makefile                     # 🔨 Build, test, setup
├── Dockerfile                   # 🐳 Multi-stage build
├── docker-compose.yml           # 🐳 Pet + Ollama
└── data/                        # 💾 Runtime (KB gob + pet state JSON)
```

## 🧑‍💻 Development

```bash
# Tests
make test

# Lint
make lint

# Build
make build

# Docker
make docker-up    # 🟢 Start all services
make docker-down  # 🔴 Stop all services
```

### 🛸 Pre-commit hooks

```bash
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pre-commit install
pre-commit install -t pre-push
```

## 🤝 Contributing

1. Fork the repository
2. Create a branch: `git checkout -b feature/your-feature`
3. Commit using conventional commits: `feat(scope): description`
4. Push and open a PR to `dev`

## 📫 Contact

For inquiries or support, please open an issue or contact [r6ty5r296it6tl4eg5m.constant214@passinbox.com](mailto:r6ty5r296it6tl4eg5m.constant214@passinbox.com).

---

## 📜 License

2026 - This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). Free to use, modify, and distribute under the terms of GPL-3.0. See [LICENSE](LICENSE) for details.
