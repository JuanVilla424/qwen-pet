# 🤖 Qwen PET — Personal Entertainment Tool

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=fff)
![Qwen](https://img.shields.io/badge/Qwen-3.5-7C3AED)
![MCP](https://img.shields.io/badge/MCP-Server-blue)
![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?logo=telegram&logoColor=fff)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=fff)
![Build Status](https://github.com/JuanVilla424/qwen-pet/actions/workflows/ci.yml/badge.svg?branch=main)
![Status](https://img.shields.io/badge/Status-Development-yellow.svg)
![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)

**Qwen PET** is a persistent AI intermediary agent that integrates with Claude Code and OpenCode via MCP (Model Context Protocol). It maintains a Knowledge Base of developer preferences, architectural decisions, and project patterns — responding autonomously when confident, escalating to the developer via Telegram when it needs input.

## 📚 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
  - [Prerequisites](#-prerequisites)
  - [Installation](#-installation)
  - [Environment Setup](#-environment-setup)
  - [Pre-Commit Hooks](#-pre-commit-hooks)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## 🌟 Features

- **🧠 Knowledge Base** — Persistent vector store (chromem-go) with developer preferences, decisions, and patterns
- **🔌 MCP Server** — stdio transport, integrates with Claude Code and OpenCode as a tool provider
- **🤖 AI Decision Engine** — Qwen 3.5 via OpenRouter API for intelligent responses when KB doesn't have enough context
- **💬 Telegram Bot** — Escalation channel with inline keyboards for approval workflows
- **🌐 Web UI** — Vue 3 + TypeScript + Tailwind CSS configuration panel
- **📊 Confidence Scoring** — Automatic routing: KB direct → AI-assisted → Telegram escalation
- **🔄 Context Persistence** — No more losing decisions between sessions

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐
│ Claude Code │     │  OpenCode   │
└──────┬──────┘     └──────┬──────┘
       │    MCP (stdio)    │
       └────────┬──────────┘
                │
       ┌────────▼────────┐
       │    Qwen PET     │
       │   MCP Server    │
       ├─────────────────┤
       │ Decision Engine │
       ├─────┬─────┬─────┤
       │ KB  │ AI  │ TG  │
       └─────┴─────┴─────┘
```

## 🚀 Getting Started

### 📋 Prerequisites

- **Go 1.25+** — Backend runtime
- **Python 3.12+** — CICD tooling (pre-commit hooks, bump2version, scripts)
- **Git** — With submodule support
- **Telegram Bot Token** — From [@BotFather](https://t.me/BotFather)
- **OpenRouter API Key** — For Qwen 3.5 inference

### 🔨 Installation

1. **Clone the Repository**

   ```bash
   git clone --recurse-submodules https://github.com/JuanVilla424/qwen-pet.git
   cd qwen-pet
   ```

2. **Initialize Submodules** (if cloned without `--recurse-submodules`)

   ```bash
   git submodule update --init --recursive
   ```

### 🔧 Environment Setup

1. **Create Python Virtual Environment** (for CICD tooling)

   ```bash
   python -m venv venv
   source venv/bin/activate
   pip install --upgrade pip
   pip install poetry
   poetry lock
   poetry install
   ```

2. **Configure Environment Variables**

   ```bash
   cp .env.template .env
   # Edit .env with your credentials
   ```

3. **Build the Project**

   ```bash
   go build ./cmd/...
   ```

### 🛸 Pre-Commit Hooks

```bash
source venv/bin/activate
pre-commit install
pre-commit install -t pre-commit
pre-commit install -t pre-push
pre-commit run --all-files
```

## 🛠️ Usage

1. **Configure MCP** — Add qwen-pet as an MCP server in Claude Code or OpenCode configuration
2. **Set Up Secrets** — Add `ACCESS_TOKEN` (PAT) in GitHub repository settings for workflow triggers
3. **Version Bumping** — Add `[patch candidate]`, `[minor candidate]`, or `[major candidate]` to commit messages to trigger version bumps
4. **Branch Workflow** — Push to `dev` → auto PR to `test` → `prod` → `main` with version tags and releases

## 📁 Project Structure

```
qwen-pet/
├── cmd/                    # Application entrypoints
│   └── pet/                # Main binary
├── internal/               # Private application code
│   ├── mcp/                # MCP server implementation
│   ├── kb/                 # Knowledge Base (chromem-go)
│   ├── ai/                 # AI inference (OpenRouter/Qwen)
│   ├── telegram/           # Telegram bot integration
│   └── config/             # Configuration management
├── web/                    # Vue 3 frontend
├── configs/                # Configuration files
├── data/                   # Runtime data (KB storage)
├── scripts/                # CICD tooling submodule
├── docs/                   # Documentation
├── .github/                # GitHub Actions workflows
├── go.mod                  # Go dependencies
├── pyproject.toml          # Python CICD tooling config
└── docker-compose.yml      # Container orchestration
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) and follow the [Code of Conduct](CODE_OF_CONDUCT.md).

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit using conventional commits: `feat(scope): description`
4. Push and open a PR into `dev` branch

## 📫 Contact

For inquiries or support, please open an issue or contact [r6ty5r296it6tl4eg5m.constant214@passinbox.com](mailto:r6ty5r296it6tl4eg5m.constant214@passinbox.com).

---

## 📜 License

2026 - This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). You are free to use, modify, and distribute this software under the terms of the GPL-3.0 license. For more details, please refer to the [LICENSE](LICENSE) file included in this repository.
