# LegionTret 🏴

**By Death Legion Team**

Run open-source large language models locally. Simple. Fast. Free.

LegionTret lets you download and run models like Llama 3, Gemma 3, Mistral, DeepSeek R1, Qwen 2.5, and 30+ more — right on your own computer. No cloud. No API keys. No tracking.

## Quick Start

### Install

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/deathlegionteam/legiontret/main/scripts/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/deathlegionteam/legiontret/main/scripts/install.ps1 | iex
```

**Docker:**
```bash
docker compose up -d
```

### Run a Model

```bash
# Download and chat with Gemma 3
legiontret run gemma3

# Run Llama 3
legiontret run llama3

# Run Mistral
legiontret run mistral

# Run DeepSeek R1
legiontret run deepseek-r1
```

### Download Models

```bash
# Download without running
legiontret pull gemma3
legiontret pull llama3.1
legiontret pull qwen2.5
```

### List & Search

```bash
# List downloaded models
legiontret list

# List all available models
legiontret list --all

# Search for models
legiontret search code
legiontret search reasoning
```

## Supported Models

| Model Family | Models | Sizes |
|---|---|---|
| **Llama 3/3.1/3.2** | Meta's flagship models | 1B — 70B |
| **Gemma 2/3** | Google's open models | 1B — 27B |
| **Mistral/Nemo/Mixtral** | Mistral AI models | 7B — 46.7B |
| **Qwen 2.5** | Alibaba's multilingual models | 7B — 32B |
| **DeepSeek R1** | Reasoning specialist | 8B — 32B |
| **Phi-3/4** | Microsoft's compact models | 3.8B — 14B |
| **Code Llama** | Code generation | 7B |
| **Qwen 2.5 Coder** | Programming specialist | 7B — 32B |
| **LLaVA** | Vision-language model | 7B |
| **TinyLlama** | Ultra-compact | 1.1B |
| **Command R** | RAG & tool use | 35B |
| **And more...** | Solar, Yi, Dolphin, etc. | Various |

## REST API

When running, LegionTret exposes a local REST API on `http://127.0.0.1:11434`:

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
```

### OpenAI-Compatible Endpoints

```bash
# List models
curl http://localhost:11434/api/v1/models

# Chat completion (OpenAI format)
curl http://localhost:11434/api/v1/chat/completions -d '{
  "model": "gemma3",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}'
```

## Python Client

```python
from legiontret import Client

client = Client()

# Generate
response = client.generate("gemma3", "Why is the sky blue?")
print(response["response"])

# Chat
response = client.chat("gemma3", [
    {"role": "user", "content": "Explain quantum computing"}
])

# List models
models = client.list_models()

# OpenAI-compatible
response = client.openai_chat("gemma3", [
    {"role": "user", "content": "Hello!"}
])
```

Install: `pip install legiontret`

## JavaScript Client

```javascript
const { Client } = require('legiontret');

const client = new Client();

// Generate
const response = await client.generate('gemma3', 'Why is the sky blue?');

// Chat
const chat = await client.chat('gemma3', [
  { role: 'user', content: 'Hello!' }
]);

// List models
const models = await client.listModels();
```

Install: `npm install legiontret`

## Docker

```bash
# Build and run
docker compose up -d

# Pull an image
docker pull ghcr.io/deathlegionteam/legiontret:latest

# Run with GPU support
docker compose up -d
```

## Chat Interface

When you run `legiontret run <model>`, you get an interactive chat with slash commands:

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

## Architecture

LegionTret works by:

1. **Model Registry** — Built-in catalog of 30+ GGUF models from HuggingFace
2. **Download Manager** — Downloads GGUF model files with resume support and progress bars
3. **llama.cpp Integration** — Runs models via llama.cpp server (auto-downloaded or uses system install)
4. **REST API** — Ollama-compatible + OpenAI-compatible endpoints
5. **Client Libraries** — Python and JavaScript wrappers for the API

## Building from Source

```bash
# Prerequisites: Go 1.22+

git clone https://github.com/deathlegionteam/legiontret.git
cd legiontret
go build ./cmd/legiontret

# Cross-compile
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o legiontret.exe ./cmd/legiontret
```

## Requirements

- **OS:** macOS, Linux, or Windows
- **RAM:** 8 GB minimum (varies by model)
- **Storage:** 1-50 GB depending on models
- **GPU:** Optional (CPU inference supported, GPU much faster)
- **llama.cpp:** Auto-downloaded or use existing install

## License

MIT License — By Death Legion Team

---

**LegionTret** — Run LLMs locally. Simple. Fast. Free. 🏴
