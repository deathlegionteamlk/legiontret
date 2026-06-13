# LegionTret 🏴

**By Death Legion Team**

**Run open-source large language models locally. Simple. Fast. Free.**

LegionTret lets you download and run 92+ AI models — including Llama 3, Gemma 3, Mistral, DeepSeek R1, Qwen 2.5, Phi-4, Codestral, and many more — right on your own computer. No cloud. No API keys. No tracking. Complete privacy.

---

## ✨ Features

- 🤖 **92+ Models** — Llama, Gemma, Mistral, Qwen, DeepSeek, Phi, and many more
- 🚀 **GPU Acceleration** — NVIDIA CUDA, Apple Metal, and CPU fallback
- 🔄 **Ollama-Compatible API** — Drop-in replacement for Ollama
- 🤖 **OpenAI-Compatible API** — Works with any OpenAI SDK/tool
- 🐍 **Python SDK** — `pip install legiontret`
- 📦 **JavaScript/TypeScript SDK** — `npm install legiontret`
- 🐳 **Docker Support** — One-command deploy with GPU passthrough
- 💬 **Interactive Chat** — Beautiful TUI with slash commands
- 📥 **Smart Downloads** — Resume support, progress bars, checksum verification
- 🏷️ **Model Discovery** — Search, tags, families, authors
- 📊 **Export** — JSON/CSV model catalog export
- 🔒 **100% Local & Private** — No data leaves your machine
- 🌍 **Multilingual** — English, Chinese, Arabic, Japanese, Korean, and more
- 🔧 **Cross-Platform** — macOS, Linux, Windows (5 architectures)
- ⚡ **Lightweight** — Single ~6MB binary, no dependencies
- 🎯 **One Command** — `legiontret run gemma3` and you're chatting

---

## 🚀 Quick Start

### Install

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/deathlegionteamlk/legiontret/main/scripts/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/deathlegionteamlk/legiontret/main/scripts/install.ps1 | iex
```

**Docker:**
```bash
docker pull deathlegion/legiontret
docker run -p 11434:11434 deathlegion/legiontret
```

### Run a Model

```bash
legiontret run gemma3         # Google Gemma 3 4B
legiontret run llama3         # Meta Llama 3 8B
legiontret run mistral        # Mistral 7B
legiontret run deepseek-r1    # DeepSeek R1 8B (reasoning)
legiontret run qwen2.5        # Qwen 2.5 7B (multilingual)
legiontret run codestral      # Codestral 22B (code)
legiontret run phi4           # Phi-4 14B (Microsoft)
legiontret run mixtral        # Mixtral 8x7B MoE
legiontret run qwq            # QwQ 32B (reasoning)
legiontret run pixtral        # Pixtral 12B (vision)
```

### Download Models

```bash
legiontret pull gemma3        # Download without running
legiontret pull llama3.1      # Llama 3.1 with 128K context
legiontret pull qwen2.5-coder:7b  # Code specialist
```

### Discover Models

```bash
legiontret list               # List downloaded models
legiontret list --all         # List all 92+ available models
legiontret search code        # Find coding models
legiontret search reasoning   # Find reasoning models
legiontret tags               # Browse by category
legiontret families           # Browse by model family
legiontret authors            # Browse by organization
legiontret count              # Model statistics
```

---

## 🤖 Supported Models (92+)

| Family | Models | Sizes | Best For |
|---|---|---|---|
| **Llama 3/3.1/3.2/3.3** | 9 models | 1B — 405B | General, chat, reasoning |
| **Gemma 2/3** | 7 models | 1B — 27B | Google ecosystem, efficiency |
| **Mistral/Nemo/Mixtral** | 9 models | 7B — 141B | Fast inference, MoE |
| **Qwen 2.5/Qwen Coder** | 9 models | 7B — 72B | Multilingual, code |
| **DeepSeek R1/V3/Coder** | 7 models | 6.7B — 671B | Reasoning, code, frontier |
| **Phi-3/3.5/4** | 5 models | 3.8B — 14B | Compact, Microsoft |
| **Code Llama/CodeGemma** | 6 models | 3B — 34B | Code generation |
| **Codestral/StarCoder2** | 4 models | 3B — 22B | Programming specialist |
| **LLaVA/Pixtral/InternVL** | 5 models | 4B — 13B | Vision, multimodal |
| **Command R/R+** | 2 models | 35B — 104B | RAG, tool use |
| **Dolphin/Nous Hermes** | 4 models | 7B — 8B | Creative, uncensored |
| **Granite (IBM)** | 2 models | 8B — 34B | Enterprise |
| **Yi/GLM/Jais** | 4 models | 6B — 34B | Chinese, Arabic, bilingual |
| **TinyLlama/SmolLM** | 3 models | 135M — 1.7B | Ultra-light, edge |
| **Specialized** | 8 models | Various | Medical, legal, finance, math, embedding |
| **And more...** | 9+ models | Various | Vicuna, Zephyr, Falcon, Mamba, etc. |

---

## 🔌 REST API

### Ollama-Compatible Endpoints

```bash
# List models
curl http://localhost:11434/api/tags

# Generate text
curl http://localhost:11434/api/generate -d '{
  "model": "gemma3",
  "prompt": "What is Python?"
}'

# Chat completion
curl http://localhost:11434/api/chat -d '{
  "model": "gemma3",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}'

# Get embeddings
curl http://localhost:11434/api/embeddings -d '{
  "model": "gemma3",
  "prompt": "embed this text"
}'
```

### OpenAI-Compatible Endpoints

```bash
# List models (OpenAI format)
curl http://localhost:11434/api/v1/models

# Chat completion (OpenAI format)
curl http://localhost:11434/api/v1/chat/completions -d '{
  "model": "gemma3",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}'
```

---

## 🐍 Python SDK

[![PyPI](https://img.shields.io/pypi/v/legiontret)](https://pypi.org/project/legiontret/)

```bash
pip install legiontret
```

```python
from legiontret import Client

client = Client()

# Generate text
response = client.generate("gemma3", "Why is the sky blue?")
print(response["response"])

# Chat
response = client.chat("llama3", [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Explain quantum computing"}
])

# Stream responses
for chunk in client.generate_stream("mistral", "Tell me a story"):
    print(chunk, end="")

# List models
models = client.list_models()

# OpenAI-compatible
response = client.openai_chat("gemma3", [
    {"role": "user", "content": "Hello!"}
])

# Get embeddings
embedding = client.embeddings("gemma3", "embed this text")

# Search models
results = client.search("code")

# Check server status
if client.is_running():
    print(f"Server version: {client.version()}")
```

---

## 📦 JavaScript/TypeScript SDK

[![npm](https://img.shields.io/npm/v/legiontret)](https://www.npmjs.com/package/legiontret)

```bash
npm install legiontret
```

```javascript
const { Client } = require('legiontret');

const client = new Client();

// Generate text
const response = await client.generate('gemma3', 'Why is the sky blue?');

// Chat
const chat = await client.chat('llama3', [
  { role: 'user', content: 'Hello!' }
]);

// Stream responses
for await (const chunk of client.generateStream('mistral', 'Tell me a story')) {
  process.stdout.write(chunk);
}

// List models
const models = await client.listModels();

// OpenAI-compatible
const result = await client.openaiChat('gemma3', [
  { role: 'user', content: 'Hello!' }
]);
```

TypeScript types included automatically.

---

## 🐳 Docker

```bash
# Pull from Docker Hub
docker pull deathlegion/legiontret:latest

# Run with GPU
docker compose up -d

# Run manually
docker run -p 11434:11434 -v legiontret-models:/root/.legiontret/models deathlegion/legiontret
```

---

## 💬 Chat Interface

```
  >>> Hello!
  Hi there! How can I help you today?

  >>> /help
  /help              - Show help
  /exit              - Exit chat
  /clear             - Clear history
  /system <prompt>   - Set system prompt
  /history           - Show chat history
  /regenerate        - Regenerate last response
  /save <file>       - Save chat history
  /stats             - Show session stats
```

---

## 🏗️ Architecture

```
┌──────────────────────────────────────────────┐
│                  LegionTret                   │
├──────────────────────────────────────────────┤
│  CLI (run, pull, list, serve, search, etc.)  │
├──────────────────────────────────────────────┤
│  REST API (Ollama + OpenAI compatible)       │
├──────────────────────────────────────────────┤
│  Model Registry (92+ GGUF models)            │
│  Download Manager (resume, progress, SHA256) │
│  llama.cpp Integration (GPU + CPU)           │
├──────────────────────────────────────────────┤
│  Python SDK  │  JavaScript SDK  │  Docker    │
└──────────────────────────────────────────────┘
```

---

## 🔧 Building from Source

```bash
git clone https://github.com/deathlegionteamlk/legiontret.git
cd legiontret
make build

# Cross-compile all platforms
make build-all
```

---

## 📋 Requirements

- **OS:** macOS (Intel/Apple Silicon), Linux (x64/ARM64), Windows (x64)
- **RAM:** 8 GB minimum (varies by model size)
- **Storage:** 0.5 GB — 400 GB depending on models
- **GPU:** Optional — NVIDIA, Apple Metal, or CPU-only
- **llama.cpp:** Auto-downloaded or use existing install

---

## 🏷️ Keywords

`llm` `ai` `local-ai` `local-llm` `language-model` `chatbot` `inference` `ollama` `ollama-alternative`
`gemma` `llama` `mistral` `deepseek` `qwen` `phi` `codestral` `mixtral` `llama3` `gemma3`
`code-generation` `gguf` `llama-cpp` `text-generation` `ai-assistant` `machine-learning` `nlp`
`natural-language-processing` `transformer` `open-source-ai` `on-device-ai` `privacy-ai` `offline-ai`
`huggingface` `model-runner` `local-inference` `chat-completion` `openai-compatible` `ollama-compatible`
`rest-api` `python-sdk` `javascript-sdk` `docker` `gpu-acceleration` `cuda` `metal`
`quantization` `ggml` `cpu-inference` `edge-ai` `embedded-ai` `self-hosted-ai`
`reasoning-model` `code-model` `vision-model` `multimodal` `rag` `embedding` `multilingual`

---

## 📄 License

MIT License — By Death Legion Team

---

**LegionTret** — Run LLMs locally. Simple. Fast. Free. 🏴
