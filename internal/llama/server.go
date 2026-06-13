package llama

import (
        "context"
        "fmt"
        "io"
        "net/http"
        "os"
        "os/exec"
        "path/filepath"
        "runtime"
        "strings"
        "sync"
        "time"

        "github.com/deathlegionteam/legiontret/internal/config"
)

// Server manages a llama.cpp server process
type Server struct {
        cfg     *config.Config
        cmd     *exec.Cmd
        host    string
        port    int
        mu      sync.Mutex
        running bool
}

// NewServer creates a new llama.cpp server manager
func NewServer(cfg *config.Config) *Server {
        return &Server{
                cfg:  cfg,
                host: cfg.Host,
                port: cfg.Port,
        }
}

// Start starts the llama.cpp server with the given model
func (s *Server) Start(ctx context.Context, modelPath string, opts ...ServerOption) error {
        s.mu.Lock()
        defer s.mu.Unlock()

        if s.running {
                return fmt.Errorf("server is already running")
        }

        // Apply options
        params := &ServerParams{
                ContextSize:    4096,
                NumGPULayers:   -1, // auto-detect
                Threads:        0,  // auto-detect
                BatchSize:      512,
                FlashAttention: true,
                Mlock:          false,
                MMProjPath:     "",
        }
        for _, opt := range opts {
                opt(params)
        }

        // Build command args
        binaryPath := s.cfg.LlamaCppBinaryPath()
        args := []string{
                "--model", modelPath,
                "--host", s.host,
                "--port", fmt.Sprintf("%d", s.port),
                "--ctx-size", fmt.Sprintf("%d", params.ContextSize),
                "--batch-size", fmt.Sprintf("%d", params.BatchSize),
        }

        if params.NumGPULayers != 0 {
                args = append(args, "--n-gpu-layers", fmt.Sprintf("%d", params.NumGPULayers))
        }
        if params.Threads > 0 {
                args = append(args, "--threads", fmt.Sprintf("%d", params.Threads))
        }
        if params.FlashAttention {
                args = append(args, "--flash-attn")
        }
        if params.Mlock {
                args = append(args, "--mlock")
        }
        if params.MMProjPath != "" {
                args = append(args, "--mmproj", params.MMProjPath)
        }

        s.cmd = exec.CommandContext(ctx, binaryPath, args...)
        s.cmd.Stdout = os.Stdout
        s.cmd.Stderr = os.Stderr

        // Set process group for clean shutdown (platform-specific)
        setProcessGroupAttr(s.cmd)

        if err := s.cmd.Start(); err != nil {
                return fmt.Errorf("failed to start llama.cpp server: %w", err)
        }

        s.running = true

        // Wait for server to be ready
        go s.waitForExit()

        if err := s.waitForReady(ctx, 120*time.Second); err != nil {
                s.Stop()
                return fmt.Errorf("server failed to become ready: %w", err)
        }

        return nil
}

// Stop stops the llama.cpp server
func (s *Server) Stop() error {
        s.mu.Lock()
        defer s.mu.Unlock()

        if !s.running || s.cmd == nil || s.cmd.Process == nil {
                return nil
        }

        // Stop the process (platform-specific)
        stopProcess(s.cmd)

        // Wait with timeout
        done := make(chan error, 1)
        go func() {
                done <- s.cmd.Wait()
        }()

        select {
        case <-done:
        case <-time.After(10 * time.Second):
                s.cmd.Process.Kill()
        }

        s.running = false
        return nil
}

// IsRunning returns whether the server is currently running
func (s *Server) IsRunning() bool {
        s.mu.Lock()
        defer s.mu.Unlock()
        return s.running
}

// WaitForReady waits for the server to accept connections
func (s *Server) waitForReady(ctx context.Context, timeout time.Duration) error {
        deadline := time.Now().Add(timeout)
        url := fmt.Sprintf("http://%s:%d/health", s.host, s.port)

        for time.Now().Before(deadline) {
                select {
                case <-ctx.Done():
                        return ctx.Err()
                default:
                }

                resp, err := http.Get(url)
                if err == nil {
                        resp.Body.Close()
                        if resp.StatusCode == http.StatusOK {
                                return nil
                        }
                }
                time.Sleep(500 * time.Millisecond)
        }

        return fmt.Errorf("server did not become ready within %v", timeout)
}

// waitForExit monitors the server process
func (s *Server) waitForExit() {
        err := s.cmd.Wait()
        s.mu.Lock()
        s.running = false
        s.mu.Unlock()
        if err != nil {
                fmt.Fprintf(os.Stderr, "llama.cpp server exited: %v\n", err)
        }
}

// EnsureBinary ensures the llama.cpp binary exists, downloading if necessary
func (s *Server) EnsureBinary(ctx context.Context) error {
        binaryPath := s.cfg.LlamaCppBinaryPath()
        if _, err := os.Stat(binaryPath); err == nil {
                return nil
        }

        fmt.Println("Downloading llama.cpp binary...")

        // Determine platform-specific download URL
        var arch, ext string
        switch runtime.GOARCH {
        case "amd64":
                arch = "x64"
        case "arm64":
                arch = "arm64"
        default:
                return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
        }

        var osName string
        switch runtime.GOOS {
        case "darwin":
                osName = "macos"
        case "linux":
                osName = "linux"
        case "windows":
                osName = "windows"
                ext = ".exe"
        default:
                return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
        }

        // For now, download from llama.cpp releases
        releaseURL := fmt.Sprintf(
                "https://github.com/ggerganov/llama.cpp/releases/latest/download/llama-server-%s-%s%s",
                osName, arch, ext,
        )

        if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
                return fmt.Errorf("failed to create binary directory: %w", err)
        }

        // Download the binary
        req, err := http.NewRequestWithContext(ctx, "GET", releaseURL, nil)
        if err != nil {
                return fmt.Errorf("failed to create download request: %w", err)
        }

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                // Fallback: try building from source or alternative download
                return s.fallbackBinarySetup(ctx, binaryPath)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return s.fallbackBinarySetup(ctx, binaryPath)
        }

        outFile, err := os.Create(binaryPath)
        if err != nil {
                return fmt.Errorf("failed to create binary file: %w", err)
        }
        defer outFile.Close()

        if _, err := io.Copy(outFile, resp.Body); err != nil {
                os.Remove(binaryPath)
                return fmt.Errorf("failed to download binary: %w", err)
        }

        // Make executable
        if runtime.GOOS != "windows" {
                if err := os.Chmod(binaryPath, 0755); err != nil {
                        return fmt.Errorf("failed to make binary executable: %w", err)
                }
        }

        fmt.Println("llama.cpp binary downloaded successfully!")
        return nil
}

// fallbackBinarySetup handles binary setup when direct download fails
func (s *Server) fallbackBinarySetup(ctx context.Context, binaryPath string) error {
        // Check if llama-server is available in PATH
        if path, err := exec.LookPath("llama-server"); err == nil {
                // Create a symlink or copy
                if err := os.Symlink(path, binaryPath); err != nil {
                        // Copy instead
                        src, err := os.Open(path)
                        if err != nil {
                                return fmt.Errorf("failed to open llama-server: %w", err)
                        }
                        defer src.Close()

                        dst, err := os.Create(binaryPath)
                        if err != nil {
                                return fmt.Errorf("failed to create binary: %w", err)
                        }
                        defer dst.Close()

                        if _, err := io.Copy(dst, src); err != nil {
                                return fmt.Errorf("failed to copy binary: %w", err)
                        }
                        os.Chmod(binaryPath, 0755)
                }
                return nil
        }

        // Check for ollama's llama server as a last resort
        ollamaPaths := []string{
                "/usr/local/bin/ollama",
                "/usr/bin/ollama",
                filepath.Join(os.Getenv("HOME"), ".local/bin/ollama"),
        }
        for _, p := range ollamaPaths {
                if _, err := os.Stat(p); err == nil {
                        fmt.Println("Found Ollama installation. LegionTret can use Ollama's backend.")
                        // Create a wrapper script
                        wrapper := fmt.Sprintf(`#!/bin/bash
# LegionTret wrapper using Ollama's llama.cpp
exec ollama serve "$@"
`, binaryPath)
                        return os.WriteFile(binaryPath, []byte(wrapper), 0755)
                }
        }

        return fmt.Errorf("could not find or download llama.cpp binary. Please install llama.cpp manually:\n" +
                "  1. Visit https://github.com/ggerganov/llama.cpp\n" +
                "  2. Build or download the llama-server binary\n" +
                fmt.Sprintf("  3. Place it at: %s", binaryPath))
}

// ServerParams holds parameters for the llama.cpp server
type ServerParams struct {
        ContextSize    int
        NumGPULayers   int
        Threads        int
        BatchSize      int
        FlashAttention bool
        Mlock          bool
        MMProjPath     string
}

// ServerOption is a function that configures server parameters
type ServerOption func(*ServerParams)

// WithContextSize sets the context size
func WithContextSize(size int) ServerOption {
        return func(p *ServerParams) { p.ContextSize = size }
}

// WithGPULayers sets the number of GPU layers
func WithGPULayers(n int) ServerOption {
        return func(p *ServerParams) { p.NumGPULayers = n }
}

// WithThreads sets the number of threads
func WithThreads(n int) ServerOption {
        return func(p *ServerParams) { p.Threads = n }
}

// WithBatchSize sets the batch size
func WithBatchSize(n int) ServerOption {
        return func(p *ServerParams) { p.BatchSize = n }
}

// WithFlashAttention enables flash attention
func WithFlashAttention(enable bool) ServerOption {
        return func(p *ServerParams) { p.FlashAttention = enable }
}

// WithMlock enables mlock
func WithMlock(enable bool) ServerOption {
        return func(p *ServerParams) { p.Mlock = enable }
}

// WithMMProj sets the multimodal projection path (for vision models)
func WithMMProj(path string) ServerOption {
        return func(p *ServerParams) { p.MMProjPath = path }
}

// GetSystemInfo returns system information for display
func GetSystemInfo() string {
        var info strings.Builder
        info.WriteString(fmt.Sprintf("OS: %s/%s\n", runtime.GOOS, runtime.GOARCH))
        info.WriteString(fmt.Sprintf("CPU Threads: %d\n", runtime.NumCPU()))

        // Check for GPU
        if hasNvidiaGPU() {
                info.WriteString("GPU: NVIDIA detected\n")
        } else if hasAppleSilicon() {
                info.WriteString("GPU: Apple Silicon (Metal)\n")
        } else {
                info.WriteString("GPU: Not detected (CPU mode)\n")
        }

        return info.String()
}

func hasNvidiaGPU() bool {
        _, err := exec.LookPath("nvidia-smi")
        return err == nil
}

func hasAppleSilicon() bool {
        if runtime.GOOS != "darwin" {
                return false
        }
        out, err := exec.Command("sysctl", "-n", "hw.optional.arm64").Output()
        return err == nil && strings.TrimSpace(string(out)) == "1"
}
