---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: []
date: 2026-03-07
author: Core
status: complete
---

# 🤖 Product Brief: qwen-pet

> 🧠 Personal Entertainment Tool — Agente intermediario con contexto persistente

🟢 Go 1.25 | 🔮 Qwen 3.5 Flash | 🔌 MCP Server | 📬 Telegram Bot | 🖥️ Vue 3 + TS | 📋 Planning Phase | 🔓 Public

---

## 📄 Executive Summary

**qwen-pet** es un agente intermediario con contexto persistente que elimina la perdida de conocimiento entre sesiones de desarrollo AI. Actua como oraculo de contexto entre las herramientas de desarrollo (Claude Code / OpenCode / BMAD agents) y el desarrollador, acumulando decisiones, preferencias y patrones para responder autonomamente cuando tiene la respuesta, o escalar via Telegram cuando necesita input humano.

El problema fundamental: cuando un desarrollador trabaja con AI, ninguno de los dos tiene conocimiento profundo acumulado — son "dos desconocidos hablando cada vez". qwen-pet resuelve esto siendo el tercero que **SI** sabe, que acumula, que recuerda, y que decide.

---

## 🎯 Core Vision

### ❗ Problem Statement

Los desarrolladores que trabajan con herramientas AI enfrentan tres problemas criticos:

1. **🔄 Perdida de contexto entre sesiones** — cada nueva sesion empieza casi desde cero. MEMORY.md y CLAUDE.md ayudan parcialmente pero son manuales, limitados, y no capturan el full picture
2. **🚫 Bloqueo por preguntas repetitivas** — los agentes AI necesitan decisiones del usuario que podrian resolverse automaticamente si tuvieran acceso a decisiones previas similares
3. **⏸️ Ausencia del desarrollador** — cuando el usuario no esta disponible, todo el pipeline de desarrollo se detiene completamente
4. **🤷 Falta de intermediario inteligente** — no existe una capa de conocimiento acumulado entre el AI y el humano que pueda tomar decisiones informadas

### 💥 Problem Impact

| Impacto               | Descripcion                                                               |
| --------------------- | ------------------------------------------------------------------------- |
| ⏱️ **Productividad**  | 15-30 min perdidos por sesion re-explicando contexto y preferencias       |
| ⚖️ **Calidad**        | Decisiones inconsistentes porque el AI no recuerda decisiones previas     |
| 🔴 **Disponibilidad** | Pipeline completamente bloqueado cuando el usuario no esta online         |
| 📉 **Conocimiento**   | Decisiones valiosas se pierden al cerrar sesion — no hay acumulacion      |
| 😤 **Frustracion**    | El usuario repite las mismas instrucciones y correcciones multiples veces |

### 🔍 Why Existing Solutions Fall Short

| Solucion Actual                | Limitacion                                                       |
| ------------------------------ | ---------------------------------------------------------------- |
| 📝 **CLAUDE.md / MEMORY.md**   | Manual, limitado a 200 lineas, sin busqueda semantica, no decide |
| 💬 **Chat history**            | Se pierde al compactar, no consultable, sin estructura           |
| 📖 **Docs del proyecto**       | Estatica, no se actualiza, no responde preguntas                 |
| 🤖 **Copilot / AI assistants** | Sin memoria persistente, sin contexto de usuario                 |
| 📋 **Notion / wikis**          | No integrado con el flujo de desarrollo, mantenimiento manual    |

Ninguna solucion combina: **memoria semantica** + **decision autonoma** + **escalacion inteligente** + **integracion nativa con herramientas de desarrollo**.

### 💡 Proposed Solution

Un agente Go que corre 24/7 en la maquina de desarrollo, expuesto como MCP Server nativo para Claude Code y OpenCode.

```
🖥️ [Claude Code / OpenCode / BMAD Agents]
              |
         (MCP stdio)
              |
     🤖 [qwen-pet]
         ├── 🧠 Knowledge Base (chromem-go) ─── decisiones, preferencias, patrones
         ├── 🔮 Decision Engine (Qwen 3.5 Flash via OpenRouter) ─── razonamiento
         ├── 📬 Telegram Bot (go-telegram/bot) ─── escalacion al usuario
         └── 🖥️ Web UI (Vue 3) ─── configuracion y visualizacion
```

**🔄 Flujo de decision:**

1. Claude Code invoca tool `ask_pet("como prefiere Core manejar auth?")`
2. PET busca en KB por similitud semantica
3. ✅ Si encuentra respuesta con alta confianza → responde inmediatamente
4. 🔮 Si no → invoca Qwen 3.5 Flash con el contexto del KB para razonar
5. 📬 Si aun no tiene confianza → escala por Telegram
6. 👤 Usuario responde → PET almacena en KB y retorna
7. 💾 La respuesta se guarda para futuras consultas similares

### ⚡ Key Differentiators

| Diferenciador                 | Detalle                                                               |
| ----------------------------- | --------------------------------------------------------------------- |
| 🔌 **MCP nativo**             | MCP server que Claude Code y OpenCode consumen directamente via stdio |
| 🧠 **Memoria semantica**      | Busqueda por significado, no keywords                                 |
| 🤖 **Decision autonoma**      | Si tiene confianza suficiente, responde sin molestar al usuario       |
| 📬 **Escalacion inteligente** | Escala por Telegram con contexto completo                             |
| 📈 **Aprendizaje continuo**   | Cada interaccion enriquece el KB                                      |
| 🏠 **Corre local**            | El conocimiento vive en la maquina del desarrollador                  |

---

## 👥 Target Users

### 👤 Primary User: Core (El Desarrollador)

**Perfil:**

- 🛠️ Desarrollador intermedio-avanzado en cloud e infraestructura
- 🔧 Stack: Go, Rust/Axum, Python/FastAPI, Vue 3/TS/Tailwind
- ☁️ Trabaja con AWS y OCI
- 🤖 Usa Claude Code y OpenCode con BMAD system (73 workflows)
- 🗣️ Respuestas directas, codigo en ingles, comunicacion en espanol

**😤 Pain Points:**

- Re-explica preferencias en cada sesion
- Pierde decisiones arquitectonicas al cerrar sesion
- Se bloquea cuando AI hace preguntas ya respondidas
- No puede alejarse sin que todo se detenga

**🎉 Success Moment:** Abrir nueva sesion y que PET ya sepa todo sin explicar nada.

### 🤖 Secondary Users: BMAD Agents / Claude Code / OpenCode

**Perfil:** Agentes AI ejecutando workflows, necesitan decisiones para avanzar.

**🎉 Success Moment:** Invocar `ask_pet` y obtener respuesta inmediata del historial.

### 🗺️ User Journey

| Fase               | Accion                                               | Componente      |
| ------------------ | ---------------------------------------------------- | --------------- |
| 🔧 **Setup**       | Instalar PET, configurar Telegram, registrar MCP     | CLI + config    |
| 🌱 **Seeding**     | Cargar preferencias iniciales                        | Web UI + import |
| 📅 **Uso diario**  | Claude/OpenCode invoca PET via MCP                   | MCP Server      |
| 📬 **Escalacion**  | PET pregunta por Telegram cuando no sabe             | Telegram Bot    |
| ✅ **Aprobacion**  | Aprobar wireframes/arquitectura via inline keyboards | Telegram Bot    |
| 📈 **Crecimiento** | KB se enriquece con cada interaccion                 | chromem-go      |
| ⚙️ **Config**      | Ajustar thresholds, ver historial                    | Web UI          |

---

## 📊 Success Metrics

### 👤 User Success

| Metrica                       | Target                     | Medicion                          |
| ----------------------------- | -------------------------- | --------------------------------- |
| ✅ % preguntas sin escalar    | >60% (2 sem), >80% (2 mes) | Ratio KB / total                  |
| ⚡ Latencia KB                | <500ms p95                 | Tiempo de `ask_pet`               |
| 📬 Respuesta Telegram         | <60s mediana               | Envio → respuesta                 |
| 🔄 Reduccion re-explicaciones | >70%                       | Preguntas repetidas antes/despues |
| 🎯 Precision respuestas       | >90%                       | Tasa de correcciones              |

### 📋 Business Objectives

| Objetivo              | Target                  | Horizonte |
| --------------------- | ----------------------- | --------- |
| 🤖 Autonomia pipeline | >30% avance sin usuario | 3 meses   |
| 📈 Crecimiento KB     | >500 entradas           | 2 meses   |
| 🟢 Uptime             | 99.5% (systemd)         | 1 mes     |
| 💰 Costo API          | <$5/mes OpenRouter      | Continuo  |

### 📉 KPIs

| KPI                    | Formula                           | Frecuencia |
| ---------------------- | --------------------------------- | ---------- |
| 📊 **Resolution Rate** | consultas_kb / total \* 100       | Diario     |
| 📬 **Escalation Rate** | consultas_telegram / total \* 100 | Diario     |
| 📈 **KB Growth**       | nuevas_entradas / dia             | Semanal    |
| 🎯 **Accuracy**        | 1 - (correcciones / respuestas)   | Semanal    |
| 💰 **API Cost**        | tokens \* precio                  | Mensual    |

---

## 🏗️ MVP Scope

### ✅ Core Features

#### 🔌 F1 — MCP Server (Claude Code + OpenCode Integration)

MCP server Go via stdio (`modelcontextprotocol/go-sdk`).

| Tool                  | Descripcion                      | Input                        | Output                     |
| --------------------- | -------------------------------- | ---------------------------- | -------------------------- |
| 🔍 `ask_pet`          | Pregunta al KB + decision engine | question, context            | answer, confidence, source |
| ⚙️ `check_preference` | Consulta preferencia especifica  | category, key                | value, last_updated        |
| 💾 `store_decision`   | Almacena decision/preferencia    | category, key, value, reason | success, id                |
| ✅ `approve_artifact` | Pide aprobacion de artefacto     | type, description, file_path | approved, comments         |

#### 🧠 F2 — Knowledge Base (chromem-go)

| Categoria                    | Ejemplos                                     |
| ---------------------------- | -------------------------------------------- |
| ⚙️ `preferences.stack`       | "Frontend: Vue 3 + TS + Tailwind"            |
| 📝 `preferences.conventions` | "Commits: type(scope): desc"                 |
| 🎨 `preferences.ui`          | "Dark mode default, toast notifications"     |
| 🏗️ `decisions.architecture`  | "BadgerDB para Kursaal, chromem-go para PET" |
| 🔧 `decisions.patterns`      | "Circuit breaker, token bucket, retry"       |
| 📁 `projects.context`        | Contexto de cada proyecto activo             |
| ✅ `approvals.history`       | Historial aprobaciones/rechazos              |

#### 🔮 F3 — Decision Engine (Qwen 3.5 Flash)

| Paso | Accion                               | Condicion              |
| ---- | ------------------------------------ | ---------------------- |
| 1️⃣   | Recibe pregunta + contexto de MCP    | Siempre                |
| 2️⃣   | Busca KB por similitud (top 5)       | Siempre                |
| 3️⃣   | Responde directo del KB              | match > threshold alto |
| 4️⃣   | Invoca Qwen 3.5 Flash con KB context | match parcial          |
| 5️⃣   | Retorna respuesta Qwen               | confianza alta         |
| 6️⃣   | Escala a Telegram                    | confianza baja         |

**Config:** OpenRouter, `qwen/qwen3.5-flash` ($0.10/$0.40 per M), fallback `qwen/qwen3-4b-free`, context 1M, timeout 30s.

#### 📬 F4 — Telegram Bot (Escalacion)

| Tipo                  | Formato                 | Respuesta     |
| --------------------- | ----------------------- | ------------- |
| ❓ Pregunta texto     | Mensaje + contexto      | Texto libre   |
| ✅ Aprobacion binaria | Descripcion + [Si / No] | Click boton   |
| 🖼️ Aprobacion imagen  | Imagen + desc + botones | Click boton   |
| 📋 Seleccion multiple | Opciones + teclado      | Click opcion  |
| ⏰ Timeout            | Recordatorio a N min    | Texto o click |

### 🚫 Out of Scope (MVP)

| Feature                         | Fase    |
| ------------------------------- | ------- |
| 🖥️ Web UI completa              | v2      |
| 🏠 Modelo local (Ollama)        | v2      |
| 👥 Multi-usuario                | v3+     |
| 💬 Slack/Discord                | v3+     |
| ⚙️ Ejecucion autonoma de codigo | Evaluar |
| 🎙️ Voice interface              | v3+     |
| 📥 Import masivo CLAUDE.md      | v2      |

### 🏁 MVP Success Criteria

| Criterio                        | Validacion                     |
| ------------------------------- | ------------------------------ |
| 🔌 `ask_pet` retorna respuestas | Test E2E: MCP tool → respuesta |
| 🧠 >50% resuelto sin escalar    | Metrica resolution rate        |
| 📬 Escalacion Telegram funciona | Test request-response          |
| 💾 KB persiste entre reinicios  | Test seed → restart → query    |
| 🖼️ Aprobacion imagen            | Test PNG → botones → respuesta |
| 💰 Costo <$5/mes                | Monitoreo tokens               |

---

## 🔮 Future Vision

### 📦 v2 — Web UI + Local Model

- 🖥️ Web UI (Vue 3): dashboard KB, historial, metricas
- 🏠 Modelo local: Qwen3.5 4B via ollama
- 📥 Import automatico: CLAUDE.md, MEMORY.md, CONTEXT.md
- 💡 Sugerencias proactivas por patrones

### 📦 v3 — Autonomia Avanzada

- 💬 Multi-canal (Slack, Discord, email)
- ⚙️ Workflow automation (branches, PRs)
- 📁 Multi-proyecto con cross-references
- 📊 Analytics dashboard

### 📦 v4 — Ecosistema

- 👥 Multi-usuario con permisos
- 🔌 Plugin system
- 🌐 API publica
- 🎓 Fine-tuning con decisiones del usuario

---

## 🔧 Technical Stack

| Componente            | Tecnologia                    | Justificacion                          |
| --------------------- | ----------------------------- | -------------------------------------- |
| 🟢 **Runtime**        | Go 1.25+                      | Consistente con stack, MCP SDK oficial |
| 🔌 **MCP Server**     | modelcontextprotocol/go-sdk   | SDK oficial, stdio, Google             |
| 🔮 **AI Inference**   | OpenRouter (Qwen3.5-Flash)    | $0.10/$0.40 per M, 1M context          |
| 🧠 **Knowledge Base** | chromem-go                    | Pure Go, embedded, zero deps           |
| 🔗 **Embeddings**     | OpenRouter / nomic-embed-text | Busqueda semantica                     |
| 📬 **Telegram Bot**   | go-telegram/bot               | API moderna, inline keyboards          |
| 🖥️ **Web UI**         | Vue 3 + TS + Tailwind CSS 4   | Consistente con stack frontend         |
| ⚙️ **Config**         | YAML + env vars               | Patron establecido                     |
| 💾 **Persistencia**   | gob files + JSON              | Ligero, sin DB server                  |

### 💻 Hardware Constraints

| Recurso  | Disponible   | Presupuesto PET    |
| -------- | ------------ | ------------------ |
| 🧮 RAM   | 7.6 GB total | Max 1 GB idle      |
| ⚙️ CPU   | 2 vCPU Xeon  | Compartido         |
| 💾 Disco | 34 GB libres | Max 2 GB datos     |
| 🎮 GPU   | Ninguna      | Inferencia via API |

### 📂 Project Structure

```
qwen-pet/
├── cmd/pet/main.go                 # Entrypoint
├── internal/
│   ├── config/config.go            # YAML + env config
│   ├── mcp/server.go               # MCP server + transport
│   ├── mcp/tools.go                # Tool definitions + handlers
│   ├── kb/store.go                 # chromem-go wrapper
│   ├── kb/embeddings.go            # Embedding provider
│   ├── ai/inference.go             # OpenRouter client
│   ├── ai/decision.go              # Decision engine
│   ├── telegram/bot.go             # Telegram bot lifecycle
│   ├── telegram/handlers.go        # Message/callback handlers
│   └── web/server.go               # Web UI API (v2)
├── web/                            # Vue 3 frontend (v2)
├── data/                           # KB persistence (gitignored)
├── configs/pet.yaml                # Default config
├── docs/PRODUCT-BRIEF.md           # Este documento
├── scripts/                        # CICD submodule
├── .env.template
├── .gitignore
├── go.mod
├── go.sum
├── docker-compose.yml
└── README.md
```

---

## ⚠️ Risks & Mitigations

| Riesgo                                | Prob.    | Impacto | Mitigacion                             |
| ------------------------------------- | -------- | ------- | -------------------------------------- |
| 🔮 Qwen3.5-Flash calidad insuficiente | Media    | Alto    | Threshold conservador, fallback a Plus |
| 🌐 OpenRouter downtime                | Baja     | Alto    | Cache, retry, fallback gratis          |
| 🧠 chromem-go no escala               | Muy baja | Medio   | 100K docs en 40ms sobra                |
| 💰 Costo API crece                    | Baja     | Medio   | Rate limiting, modelo barato           |
| 🔌 Complejidad MCP                    | Media    | Medio   | SDK oficial, empezar 2 tools           |
| 💾 Corrupcion KB                      | Baja     | Alto    | Backup, gzip, checksums                |

---

## 📅 Implementation Phases

| Fase | Features                                                  | Deps      |
| ---- | --------------------------------------------------------- | --------- |
| 1️⃣   | ⚙️ Config + 🧠 KB + 🔗 Embeddings                         | —         |
| 2️⃣   | 🔌 MCP Server (`ask_pet`, `check_preference`)             | Fase 1    |
| 3️⃣   | 🔮 Decision Engine (OpenRouter/Qwen)                      | Fase 1, 2 |
| 4️⃣   | 📬 Telegram Bot + escalacion                              | Fase 3    |
| 5️⃣   | 💾 Tools completos (`store_decision`, `approve_artifact`) | Fase 2, 4 |
| 6️⃣   | 🌱 Seeding KB + testing E2E                               | Fase 1-5  |
| 7️⃣   | 🚀 Systemd service + deployment                           | Fase 6    |

---

_📅 Creado: 2026-03-07_
_📁 Proyecto: qwen-pet_
_📍 Ubicacion: ~/Projects/qwen-pet/_
_📋 Status: Planning Phase — Product Brief Complete_
