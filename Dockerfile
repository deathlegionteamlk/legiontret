# ─── Build Stage ────────────────────────────────────────────────────
FROM golang:1.22-bookworm AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /legiontret ./cmd/legiontret

# ─── Runtime Stage ──────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install llama.cpp server binary
RUN if [ "$(uname -m)" = "x86_64" ]; then \
      LLAMA_ARCH="x64"; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
      LLAMA_ARCH="arm64"; \
    fi && \
    curl -L -o /usr/local/bin/llama-server \
      "https://github.com/ggerganov/llama.cpp/releases/latest/download/llama-server-linux-${LLAMA_ARCH}" && \
    chmod +x /usr/local/bin/llama-server || echo "llama.cpp download failed - install manually"

# Copy LegionTret binary
COPY --from=builder /legiontret /usr/local/bin/legiontret

# Create data directories
RUN mkdir -p /root/.legiontret/models /root/.legiontret/bin

# Environment
ENV LEGIONTRET_HOST=0.0.0.0
ENV LEGIONTRET_PORT=11434

# Expose API port
EXPOSE 11434

# Volume for models
VOLUME ["/root/.legiontret/models"]

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:11434/health || exit 1

# Entry point
ENTRYPOINT ["legiontret"]
CMD ["serve", "--host", "0.0.0.0"]
