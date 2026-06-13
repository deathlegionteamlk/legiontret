package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/deathlegionteam/legiontret/internal/api"
	"github.com/deathlegionteam/legiontret/internal/config"
	"github.com/deathlegionteam/legiontret/internal/download"
	"github.com/deathlegionteam/legiontret/internal/llama"
	"github.com/deathlegionteam/legiontret/internal/model"
	"github.com/deathlegionteam/legiontret/internal/progress"
	"github.com/deathlegionteam/legiontret/internal/registry"
	"github.com/deathlegionteam/legiontret/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cfg := config.DefaultConfig()
	if err := cfg.EnsureDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	reg := registry.NewRegistry()
	mgr := model.NewManager(cfg, reg)
	llamaServer := llama.NewServer(cfg)

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "run":
		cmdRun(cfg, mgr, reg, llamaServer, args)
	case "pull":
		cmdPull(cfg, mgr, reg, args)
	case "serve":
		cmdServe(cfg, mgr, reg, llamaServer, args)
	case "list", "ls":
		cmdList(cfg, mgr, reg, args)
	case "rm", "delete":
		cmdDelete(mgr, args)
	case "show":
		cmdShow(mgr, reg, args)
	case "search":
		cmdSearch(reg, args)
	case "info":
		cmdInfo(cfg)
	case "api":
		cmdServe(cfg, mgr, reg, llamaServer, args)
	case "create":
		cmdCreate(args)
	case "update":
		cmdUpdate(cfg)
	case "help", "--help", "-h":
		printUsage()
	case "version", "--version", "-v":
		fmt.Printf("legiontret version %s\nby %s\n", config.Version, config.OrgName)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	progress.PrintBanner()
	fmt.Println("Usage:")
	fmt.Println("  legiontret <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <model>          Download (if needed) and run a model interactively")
	fmt.Println("  pull <model>         Download a model")
	fmt.Println("  serve                Start the API server (without running a model)")
	fmt.Println("  list, ls             List local and available models")
	fmt.Println("  rm <model>           Delete a local model")
	fmt.Println("  show <model>         Show model details")
	fmt.Println("  search <query>       Search for models")
	fmt.Println("  info                 Show system information")
	fmt.Println("  api                  Start the API server")
	fmt.Println("  create <name>        Create a custom model from a Modelfile")
	fmt.Println("  update               Update LegionTret to the latest version")
	fmt.Println("  version              Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  legiontret run gemma3              # Run Gemma 3")
	fmt.Println("  legiontret run llama3              # Run Llama 3")
	fmt.Println("  legiontret run mistral             # Run Mistral")
	fmt.Println("  legiontret pull deepseek-r1        # Download DeepSeek R1")
	fmt.Println("  legiontret list                    # List all models")
	fmt.Println("  legiontret search code             # Search for code models")
	fmt.Println()
	fmt.Printf("By Death Legion Team | v%s\n", config.Version)
}

func cmdRun(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, llamaServer *llama.Server, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret run <model>")
		fmt.Fprintln(os.Stderr, "Example: legiontret run gemma3")
		os.Exit(1)
	}

	modelName := model.ResolveModelName(args[0])
	progress.PrintBanner()

	if !mgr.IsDownloaded(modelName) {
		fmt.Printf("  Model %q not found locally. Downloading...\n\n", modelName)
		if err := doDownload(cfg, mgr, reg, modelName); err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading model: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("  Model %q found locally.\n", modelName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ensure llama.cpp binary exists
	if err := llamaServer.EnsureBinary(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	modelPath := mgr.GetModelPath(modelName)
	fmt.Printf("  Loading model from %s\n", modelPath)
	fmt.Println("  Starting inference server...")

	if err := llamaServer.Start(ctx, modelPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting model server: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure llama.cpp is installed.")
		os.Exit(1)
	}
	defer llamaServer.Stop()

	// Start API server in background
	apiServer := api.NewServer(cfg, mgr, reg, llamaServer)
	go apiServer.Start()
	defer apiServer.Stop()

	fmt.Println()
	fmt.Println("  ═══════════════════════════════════════════════════")
	fmt.Printf("  Model %s is ready!\n", modelName)
	fmt.Println("  ═══════════════════════════════════════════════════")

	chat := tui.NewChatSession(cfg, modelName)
	if err := chat.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Chat error: %v\n", err)
		os.Exit(1)
	}
}

func cmdPull(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret pull <model>")
		os.Exit(1)
	}

	modelName := model.ResolveModelName(args[0])
	progress.PrintBanner()

	if mgr.IsDownloaded(modelName) {
		fmt.Printf("  Model %q already downloaded.\n", modelName)
		return
	}

	if err := doDownload(cfg, mgr, reg, modelName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func doDownload(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, modelName string) error {
	regModel, found := reg.GetModel(modelName)
	if !found {
		return fmt.Errorf("model %q not found in registry. Use 'legiontret search' to find models", modelName)
	}

	fmt.Printf("  Downloading %s (%s, %s)...\n", regModel.DisplayName, regModel.Parameters, regModel.Size)
	fmt.Println()

	downloader := download.NewDownloader()
	bar := progress.NewBar("  Pulling", 0)

	err := downloader.Download(regModel.URL, cfg.ModelPath(modelName), func(downloaded, total int64, speed float64, eta time.Duration) {
		bar.Update(downloaded, total, speed, eta)
	})
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	bar.Finish()

	// Save metadata
	meta := &model.ModelMetadata{
		Name:         modelName,
		URL:          regModel.URL,
		DownloadedAt: time.Now(),
		Family:       regModel.Family,
	}
	if info, err := os.Stat(cfg.ModelPath(modelName)); err == nil {
		meta.Size = info.Size()
	}
	mgr.SaveMetadata(modelName, meta)

	fmt.Printf("\n  Success! Model %q downloaded.\n", modelName)
	fmt.Printf("  Run it with: legiontret run %s\n", modelName)
	return nil
}

func cmdServe(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, llamaServer *llama.Server, args []string) {
	progress.PrintBanner()

	serveFlags := flag.NewFlagSet("serve", flag.ContinueOnError)
	host := serveFlags.String("host", cfg.Host, "Host to bind to")
	port := serveFlags.Int("port", cfg.Port, "Port to bind to")
	serveFlags.Parse(args)

	cfg.Host = *host
	cfg.Port = *port

	fmt.Printf("  Starting LegionTret API server on %s:%d...\n", cfg.Host, cfg.Port)

	apiServer := api.NewServer(cfg, mgr, reg, llamaServer)
	if err := apiServer.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cmdList(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, args []string) {
	progress.PrintBanner()

	localModels, err := mgr.ListLocal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	downloaded := make(map[string]bool)
	for _, m := range localModels {
		downloaded[m.Name] = true
	}

	showAll := len(args) > 0 && args[0] == "--all"

	if showAll {
		fmt.Println("  All Available Models:")
		fmt.Println("  ─────────────────────────────────────────────────")

		allModels := reg.ListModels()
		sort.Slice(allModels, func(i, j int) bool {
			return allModels[i].Name < allModels[j].Name
		})

		families := make(map[string][]registry.ModelInfo)
		for _, m := range allModels {
			families[m.Family] = append(families[m.Family], m)
		}

		for _, family := range sortedKeys(families) {
			fmt.Printf("\n  [%s]\n", strings.ToUpper(family))
			for _, m := range families[family] {
				status := "remote"
				if downloaded[m.Name] {
					status = "local"
				}
				fmt.Printf("    %-25s %-25s %8s %10s %s\n",
					m.Name, m.DisplayName, m.Parameters, m.Size, status)
			}
		}
	} else {
		if len(localModels) == 0 {
			fmt.Println("  No models downloaded yet.")
			fmt.Println()
			fmt.Println("  Download a model with: legiontret pull <model>")
			fmt.Println("  Or run directly with:  legiontret run <model>")
			fmt.Println()
			fmt.Println("  Popular models to try:")
			popular := []struct{ name, desc string }{
				{"gemma3", "Google's Gemma 3 4B - lightweight and capable"},
				{"llama3", "Meta's Llama 3 8B - excellent all-rounder"},
				{"mistral", "Mistral 7B - fast and efficient"},
				{"deepseek-r1", "DeepSeek R1 8B - powerful reasoning"},
				{"qwen2.5", "Qwen 2.5 7B - multilingual champion"},
				{"tinyllama", "TinyLlama 1.1B - ultra compact"},
			}
			for _, m := range popular {
				fmt.Printf("    %-20s %s\n", m.name, m.desc)
			}
			return
		}

		fmt.Println("  Locally Available Models:")
		fmt.Println("  ─────────────────────────────────────────────────")
		for _, m := range localModels {
			fmt.Printf("    %-25s %8s %10s  %s\n",
				m.Name, m.Parameters, download.FormatSize(m.Size), m.ModifiedAt.Format("2006-01-02"))
		}
	}

	fmt.Println()
	fmt.Println("  Use 'legiontret list --all' to see all available models.")
	fmt.Printf("  Models directory: %s\n", cfg.ModelsDir)
}

func cmdDelete(mgr *model.Manager, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret rm <model>")
		os.Exit(1)
	}
	if err := mgr.Delete(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Deleted model %q\n", args[0])
}

func cmdShow(mgr *model.Manager, reg *registry.Registry, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret show <model>")
		os.Exit(1)
	}
	modelName := model.ResolveModelName(args[0])
	info, err := mgr.GetModelInfo(modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("  Model: %s\n", info.DisplayName)
	fmt.Printf("  Name:  %s\n", info.Name)
	fmt.Printf("  Family: %s\n", info.Family)
	fmt.Printf("  Parameters: %s\n", info.Parameters)
	fmt.Printf("  Downloaded: %v\n", info.IsDownloaded)
	if info.IsDownloaded {
		fmt.Printf("  Size: %s\n", download.FormatSize(info.Size))
	}
	fmt.Printf("  Description: %s\n", info.Description)
	if len(info.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(info.Tags, ", "))
	}
	fmt.Println()
}

func cmdSearch(reg *registry.Registry, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret search <query>")
		os.Exit(1)
	}

	query := strings.Join(args, " ")
	results := reg.SearchModels(query)

	progress.PrintBanner()
	fmt.Printf("  Search results for %q:\n\n", query)

	if len(results) == 0 {
		fmt.Println("  No models found. Try a different search term.")
		return
	}

	for _, m := range results {
		fmt.Printf("  %-25s %s\n", m.Name, m.DisplayName)
		fmt.Printf("  %-25s %s, %s\n", "", m.Parameters, m.Size)
		fmt.Printf("  %-25s %s\n", "", m.Description)
		if len(m.Tags) > 0 {
			fmt.Printf("  %-25s Tags: %s\n", "", strings.Join(m.Tags, ", "))
		}
		fmt.Println()
	}
}

func cmdInfo(cfg *config.Config) {
	progress.PrintBanner()
	fmt.Println("  System Information:")
	fmt.Println("  ─────────────────────────────────────────────────")
	fmt.Printf("  LegionTret Version: %s\n", config.Version)
	fmt.Printf("  Organization: %s\n", config.OrgName)
	fmt.Printf("  Models Directory: %s\n", cfg.ModelsDir)
	fmt.Printf("  Binaries Directory: %s\n", cfg.BinariesDir)
	fmt.Printf("  API URL: %s\n", cfg.APIBaseURL())
	fmt.Println()
	fmt.Print(llama.GetSystemInfo())
}

func cmdCreate(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: legiontret create <name> -f <modelfile>")
		os.Exit(1)
	}
	fmt.Println("  Custom model creation is coming soon!")
}

func cmdUpdate(cfg *config.Config) {
	fmt.Println("  Checking for updates...")
	fmt.Printf("  Current version: %s\n", config.Version)
	fmt.Println("  Update functionality will be available after the first release.")
}

func sortedKeys(m map[string][]registry.ModelInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Ensure signal package is used
var _ = syscall.SIGINT
var _ = signal.Notify
var _ = context.Background
var _ = flag.NewFlagSet
