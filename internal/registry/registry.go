package registry

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strings"
)

// ModelInfo represents information about a model in the registry
type ModelInfo struct {
        Name         string   `json:"name"`
        DisplayName  string   `json:"display_name"`
        Description  string   `json:"description"`
        Parameters   string   `json:"parameters"`
        Quantization string   `json:"quantization"`
        Size         string   `json:"size"`
        Family       string   `json:"family"`
        URL          string   `json:"url"`
        SHA256       string   `json:"sha256"`
        Tags         []string `json:"tags"`
        ContextLen   int      `json:"context_length"`
        License      string   `json:"license"`
        Author       string   `json:"author"`
}

// Registry manages the model catalog
type Registry struct {
        models map[string]ModelInfo
        client *http.Client
}

// NewRegistry creates a new model registry
func NewRegistry() *Registry {
        r := &Registry{
                models: make(map[string]ModelInfo),
                client: &http.Client{},
        }
        r.loadBuiltinModels()
        return r
}

// loadBuiltinModels loads the built-in model catalog with 150+ models
func (r *Registry) loadBuiltinModels() {
        models := []ModelInfo{

                // ═══════════════════════════════════════════
                // META LLAMA FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "llama3", DisplayName: "Llama 3 8B Instruct", Description: "Meta's Llama 3 8B instruct model — excellent general-purpose model for chat and reasoning",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "llama3", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Meta-Llama-3-8B-Instruct-GGUF/resolve/main/Meta-Llama-3-8B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "general", "chat", "reasoning"},
                },
                {
                        Name: "llama3:70b", DisplayName: "Llama 3 70B Instruct", Description: "Meta's Llama 3 70B instruct — top-tier reasoning and knowledge",
                        Parameters: "70B", Quantization: "Q4_K_M", Size: "40.5 GB", Family: "llama", ContextLen: 8192, License: "llama3", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Meta-Llama-3-70B-Instruct-GGUF/resolve/main/Meta-Llama-3-70B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "general", "reasoning"},
                },
                {
                        Name: "llama3.1", DisplayName: "Llama 3.1 8B Instruct", Description: "Meta's Llama 3.1 8B instruct — improved over Llama 3 with longer context",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.9 GB", Family: "llama", ContextLen: 131072, License: "llama3.1", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "general", "latest", "long-context"},
                },
                {
                        Name: "llama3.1:70b", DisplayName: "Llama 3.1 70B Instruct", Description: "Meta's Llama 3.1 70B instruct — state of the art with 128K context",
                        Parameters: "70B", Quantization: "Q4_K_M", Size: "42.0 GB", Family: "llama", ContextLen: 131072, License: "llama3.1", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Meta-Llama-3.1-70B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-70B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "general", "latest", "long-context"},
                },
                {
                        Name: "llama3.1:405b", DisplayName: "Llama 3.1 405B Instruct", Description: "Meta's Llama 3.1 405B — the largest open model available",
                        Parameters: "405B", Quantization: "Q4_K_M", Size: "230 GB", Family: "llama", ContextLen: 131072, License: "llama3.1", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Meta-Llama-3.1-405B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-405B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"frontier", "powerful", "instruct", "long-context"},
                },
                {
                        Name: "llama3.2", DisplayName: "Llama 3.2 3B Instruct", Description: "Meta's Llama 3.2 3B instruct — lightweight and fast for edge devices",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "2.0 GB", Family: "llama", ContextLen: 131072, License: "llama3.2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"lightweight", "fast", "instruct", "edge", "long-context"},
                },
                {
                        Name: "llama3.2:1b", DisplayName: "Llama 3.2 1B Instruct", Description: "Meta's Llama 3.2 1B instruct — ultra lightweight for mobile/edge",
                        Parameters: "1B", Quantization: "Q4_K_M", Size: "0.8 GB", Family: "llama", ContextLen: 131072, License: "llama3.2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"ultra-light", "fast", "instruct", "mobile", "edge"},
                },
                {
                        Name: "llama3.3", DisplayName: "Llama 3.3 70B Instruct", Description: "Meta's Llama 3.3 70B instruct — newest Llama with improved quality",
                        Parameters: "70B", Quantization: "Q4_K_M", Size: "42 GB", Family: "llama", ContextLen: 131072, License: "llama3.3", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-3.3-70B-Instruct-GGUF/resolve/main/Llama-3.3-70B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"latest", "powerful", "instruct", "long-context"},
                },
                {
                        Name: "llama-guard3", DisplayName: "Llama Guard 3 8B", Description: "Meta's Llama Guard 3 — content safety and moderation model",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "llama3", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-Guard-3-8B-GGUF/resolve/main/Llama-Guard-3-8B-Q4_K_M.gguf",
                        Tags: []string{"safety", "moderation", "guard"},
                },

                // ═══════════════════════════════════════════
                // GOOGLE GEMMA FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "gemma3", DisplayName: "Gemma 3 4B IT", Description: "Google's Gemma 3 4B instruct — efficient and highly capable for its size",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.6 GB", Family: "gemma", ContextLen: 131072, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-3-4b-it-GGUF/resolve/main/gemma-3-4b-it-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "google", "long-context"},
                },
                {
                        Name: "gemma3:1b", DisplayName: "Gemma 3 1B IT", Description: "Google's Gemma 3 1B instruct — ultra compact with great quality",
                        Parameters: "1B", Quantization: "Q4_K_M", Size: "0.7 GB", Family: "gemma", ContextLen: 32768, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-3-1b-it-GGUF/resolve/main/gemma-3-1b-it-Q4_K_M.gguf",
                        Tags: []string{"ultra-light", "fast", "google"},
                },
                {
                        Name: "gemma3:12b", DisplayName: "Gemma 3 12B IT", Description: "Google's Gemma 3 12B instruct — strong performance for demanding tasks",
                        Parameters: "12B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "gemma", ContextLen: 131072, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-3-12b-it-GGUF/resolve/main/gemma-3-12b-it-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "google", "long-context"},
                },
                {
                        Name: "gemma3:27b", DisplayName: "Gemma 3 27B IT", Description: "Google's Gemma 3 27B instruct — top-tier Google open model",
                        Parameters: "27B", Quantization: "Q4_K_M", Size: "16.2 GB", Family: "gemma", ContextLen: 131072, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-3-27b-it-GGUF/resolve/main/gemma-3-27b-it-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "google", "long-context"},
                },
                {
                        Name: "gemma2", DisplayName: "Gemma 2 9B IT", Description: "Google's Gemma 2 9B instruct — proven performer with excellent benchmarks",
                        Parameters: "9B", Quantization: "Q4_K_M", Size: "5.4 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-2-9b-it-GGUF/resolve/main/gemma-2-9b-it-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "google"},
                },
                {
                        Name: "gemma2:2b", DisplayName: "Gemma 2 2B IT", Description: "Google's Gemma 2 2B instruct — small but surprisingly capable",
                        Parameters: "2B", Quantization: "Q4_K_M", Size: "1.4 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-2-2b-it-GGUF/resolve/main/gemma-2-2b-it-Q4_K_M.gguf",
                        Tags: []string{"lightweight", "fast", "google"},
                },
                {
                        Name: "gemma2:27b", DisplayName: "Gemma 2 27B IT", Description: "Google's Gemma 2 27B instruct — excellent quality",
                        Parameters: "27B", Quantization: "Q4_K_M", Size: "16.2 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-2-27b-it-GGUF/resolve/main/gemma-2-27b-it-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "google"},
                },

                // ═══════════════════════════════════════════
                // MISTRRAL AI FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "mistral", DisplayName: "Mistral 7B Instruct v0.3", Description: "Mistral AI's 7B instruct — fast, efficient, and highly popular",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mistral-7B-Instruct-v0.3-GGUF/resolve/main/Mistral-7B-Instruct-v0.3-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "fast"},
                },
                {
                        Name: "mistral-nemo", DisplayName: "Mistral Nemo 12B Instruct", Description: "Mistral Nemo 12B — excellent context window and reasoning",
                        Parameters: "12B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "mistral", ContextLen: 131072, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mistral-Nemo-Instruct-2407-GGUF/resolve/main/Mistral-Nemo-Instruct-2407-Q4_K_M.gguf",
                        Tags: []string{"popular", "instruct", "long-context"},
                },
                {
                        Name: "mistral-small", DisplayName: "Mistral Small 24B Instruct", Description: "Mistral Small 24B — strong reasoning in a compact package",
                        Parameters: "24B", Quantization: "Q4_K_M", Size: "14.1 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mistral-Small-24B-Instruct-2501-GGUF/resolve/main/Mistral-Small-24B-Instruct-2501-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "reasoning"},
                },
                {
                        Name: "mistral-medium", DisplayName: "Mistral Medium 3", Description: "Mistral Medium 3 — balanced performance and efficiency",
                        Parameters: "Medium", Quantization: "Q4_K_M", Size: "30 GB", Family: "mistral", ContextLen: 131072, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/mistralai/Mistral-Medium-3-GGUF/resolve/main/Mistral-Medium-3-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "long-context"},
                },
                {
                        Name: "mixtral", DisplayName: "Mixtral 8x7B Instruct", Description: "Mixtral 8x7B — Mixture of Experts for fast inference with MoE architecture",
                        Parameters: "46.7B", Quantization: "Q4_K_M", Size: "26.0 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mixtral-8x7B-Instruct-v0.1-GGUF/resolve/main/Mixtral-8x7B-Instruct-v0.1-Q4_K_M.gguf",
                        Tags: []string{"moe", "powerful", "instruct", "fast"},
                },
                {
                        Name: "mixtral:8x22b", DisplayName: "Mixtral 8x22B Instruct", Description: "Mixtral 8x22B — largest MoE model, exceptional quality",
                        Parameters: "141B", Quantization: "Q4_K_M", Size: "80 GB", Family: "mistral", ContextLen: 65536, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mixtral-8x22B-Instruct-v0.1-GGUF/resolve/main/Mixtral-8x22B-Instruct-v0.1-Q4_K_M.gguf",
                        Tags: []string{"moe", "frontier", "powerful", "instruct"},
                },
                {
                        Name: "codestral", DisplayName: "Codestral 22B", Description: "Mistral AI's Codestral — specialized for code generation and completion",
                        Parameters: "22B", Quantization: "Q4_K_M", Size: "12.8 GB", Family: "mistral", ContextLen: 32768, License: "codestral", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Codestral-22B-v0.1-GGUF/resolve/main/Codestral-22B-v0.1-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "popular"},
                },
                {
                        Name: "mathstral", DisplayName: "Mathstral 7B", Description: "Mistral AI's math-specialized model — excellent for mathematical reasoning",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/mathstral-7B-v0.1-GGUF/resolve/main/mathstral-7B-v0.1-Q4_K_M.gguf",
                        Tags: []string{"math", "science", "reasoning"},
                },
                {
                        Name: "pixtral", DisplayName: "Pixtral 12B", Description: "Mistral's Pixtral — vision-language model for image understanding",
                        Parameters: "12B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "mistral", ContextLen: 131072, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Pixtral-12B-2409-GGUF/resolve/main/Pixtral-12B-2409-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "image"},
                },

                // ═══════════════════════════════════════════
                // QWEN / ALIBABA FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "qwen2.5", DisplayName: "Qwen 2.5 7B Instruct", Description: "Alibaba's Qwen 2.5 7B — multilingual champion with 128K context",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-7B-Instruct-GGUF/resolve/main/Qwen2.5-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "instruct", "popular", "long-context"},
                },
                {
                        Name: "qwen2.5:14b", DisplayName: "Qwen 2.5 14B Instruct", Description: "Qwen 2.5 14B — strong multilingual and reasoning capabilities",
                        Parameters: "14B", Quantization: "Q4_K_M", Size: "8.7 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-14B-Instruct-GGUF/resolve/main/Qwen2.5-14B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "powerful", "instruct", "long-context"},
                },
                {
                        Name: "qwen2.5:32b", DisplayName: "Qwen 2.5 32B Instruct", Description: "Qwen 2.5 32B — top-tier multilingual model with exceptional quality",
                        Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-32B-Instruct-GGUF/resolve/main/Qwen2.5-32B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "powerful", "instruct", "long-context"},
                },
                {
                        Name: "qwen2.5:72b", DisplayName: "Qwen 2.5 72B Instruct", Description: "Qwen 2.5 72B — frontier-class multilingual model",
                        Parameters: "72B", Quantization: "Q4_K_M", Size: "42 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-72B-Instruct-GGUF/resolve/main/Qwen2.5-72B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "frontier", "powerful", "long-context"},
                },
                {
                        Name: "qwen2.5-coder:7b", DisplayName: "Qwen 2.5 Coder 7B", Description: "Qwen 2.5 Coder 7B — specialized for programming and code",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "popular", "long-context"},
                },
                {
                        Name: "qwen2.5-coder:32b", DisplayName: "Qwen 2.5 Coder 32B", Description: "Qwen 2.5 Coder 32B — elite programming model",
                        Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-32B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-32B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "powerful", "long-context"},
                },
                {
                        Name: "qwen2.5-math:7b", DisplayName: "Qwen 2.5 Math 7B", Description: "Qwen 2.5 Math — specialized for mathematical computation",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen", ContextLen: 4096, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Math-7B-Instruct-GGUF/resolve/main/Qwen2.5-Math-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"math", "science", "reasoning"},
                },
                {
                        Name: "qwen2-vl", DisplayName: "Qwen 2 VL 7B Instruct", Description: "Qwen 2 Vision-Language — understand images and text together",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen", ContextLen: 32768, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2-VL-7B-Instruct-GGUF/resolve/main/Qwen2-VL-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "image"},
                },
                {
                        Name: "qwq", DisplayName: "QwQ 32B Preview", Description: "Alibaba's QwQ — reasoning model with chain-of-thought",
                        Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/QwQ-32B-Preview-GGUF/resolve/main/QwQ-32B-Preview-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "powerful", "latest", "long-context"},
                },

                // ═══════════════════════════════════════════
                // DEEPSEEK FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "deepseek-r1", DisplayName: "DeepSeek R1 8B Distill", Description: "DeepSeek R1 8B distilled — powerful reasoning with chain-of-thought",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.9 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Llama-8B-GGUF/resolve/main/DeepSeek-R1-Distill-Llama-8B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "popular", "latest", "long-context"},
                },
                {
                        Name: "deepseek-r1:14b", DisplayName: "DeepSeek R1 14B Distill", Description: "DeepSeek R1 14B distill — enhanced reasoning capabilities",
                        Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-14B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-14B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "powerful", "long-context"},
                },
                {
                        Name: "deepseek-r1:32b", DisplayName: "DeepSeek R1 32B Distill", Description: "DeepSeek R1 32B distill — advanced reasoning and analysis",
                        Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-32B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-32B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "powerful", "long-context"},
                },
                {
                        Name: "deepseek-r1:70b", DisplayName: "DeepSeek R1 70B Distill", Description: "DeepSeek R1 70B distill — frontier-level reasoning",
                        Parameters: "70B", Quantization: "Q4_K_M", Size: "42 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Llama-70B-GGUF/resolve/main/DeepSeek-R1-Distill-Llama-70B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "frontier", "powerful", "long-context"},
                },
                {
                        Name: "deepseek-coder-v2", DisplayName: "DeepSeek Coder V2 16B", Description: "DeepSeek Coder V2 — MoE code model with excellent performance",
                        Parameters: "16B", Quantization: "Q4_K_M", Size: "9.5 GB", Family: "deepseek", ContextLen: 131072, License: "deepseek", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/deepseek-coder-6.7B-instruct-GGUF/resolve/main/deepseek-coder-6.7B-instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "moe", "long-context"},
                },
                {
                        Name: "deepseek-v3", DisplayName: "DeepSeek V3", Description: "DeepSeek V3 — latest general-purpose model with MoE architecture",
                        Parameters: "671B", Quantization: "Q4_K_M", Size: "380 GB", Family: "deepseek", ContextLen: 131072, License: "deepseek", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-V3-GGUF/resolve/main/DeepSeek-V3-Q4_K_M.gguf",
                        Tags: []string{"frontier", "moe", "powerful", "latest"},
                },
                {
                        Name: "janus-pro", DisplayName: "Janus Pro 7B", Description: "DeepSeek's Janus Pro — multimodal understanding and generation",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "deepseek", ContextLen: 4096, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/Janus-Pro-7B-GGUF/resolve/main/Janus-Pro-7B-Q4_K_M.gguf",
                        Tags: []string{"multimodal", "vision", "latest"},
                },

                // ═══════════════════════════════════════════
                // MICROSOFT PHI FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "phi3", DisplayName: "Phi-3 Mini 3.8B Instruct", Description: "Microsoft's Phi-3 Mini — small but surprisingly powerful",
                        Parameters: "3.8B", Quantization: "Q4_K_M", Size: "2.3 GB", Family: "phi", ContextLen: 4096, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/Phi-3-mini-4k-instruct-GGUF/resolve/main/Phi-3-mini-4k-instruct-Q4_K_M.gguf",
                        Tags: []string{"lightweight", "microsoft", "fast"},
                },
                {
                        Name: "phi3:14b", DisplayName: "Phi-3 Medium 14B Instruct", Description: "Microsoft's Phi-3 Medium — compact powerhouse",
                        Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "phi", ContextLen: 4096, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/Phi-3-medium-4k-instruct-GGUF/resolve/main/Phi-3-medium-4k-instruct-Q4_K_M.gguf",
                        Tags: []string{"microsoft", "powerful"},
                },
                {
                        Name: "phi3.5", DisplayName: "Phi-3.5 Mini 3.8B Instruct", Description: "Microsoft's Phi-3.5 — improved mini model with better multilingual support",
                        Parameters: "3.8B", Quantization: "Q4_K_M", Size: "2.3 GB", Family: "phi", ContextLen: 131072, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/Phi-3.5-mini-instruct-GGUF/resolve/main/Phi-3.5-mini-instruct-Q4_K_M.gguf",
                        Tags: []string{"lightweight", "microsoft", "fast", "long-context"},
                },
                {
                        Name: "phi4", DisplayName: "Phi-4 14B", Description: "Microsoft's Phi-4 — latest and most capable Phi model",
                        Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "phi", ContextLen: 16384, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/phi-4-GGUF/resolve/main/phi-4-Q4_K_M.gguf",
                        Tags: []string{"microsoft", "latest", "powerful", "reasoning"},
                },
                {
                        Name: "phi4-mini", DisplayName: "Phi-4 Mini 3.8B", Description: "Microsoft's Phi-4 Mini — compact version of Phi-4",
                        Parameters: "3.8B", Quantization: "Q4_K_M", Size: "2.3 GB", Family: "phi", ContextLen: 131072, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/Phi-4-mini-instruct-GGUF/resolve/main/Phi-4-mini-instruct-Q4_K_M.gguf",
                        Tags: []string{"microsoft", "latest", "lightweight", "long-context"},
                },

                // ═══════════════════════════════════════════
                // CODE-SPECIALIZED MODELS
                // ═══════════════════════════════════════════
                {
                        Name: "codellama", DisplayName: "Code Llama 7B Instruct", Description: "Meta's Code Llama — specialized for code generation and understanding",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "llama", ContextLen: 16384, License: "llama2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/CodeLlama-7b-Instruct-hf-GGUF/resolve/main/CodeLlama-7b-Instruct-hf-Q4_K_M.gguf",
                        Tags: []string{"code", "programming"},
                },
                {
                        Name: "codellama:34b", DisplayName: "Code Llama 34B Instruct", Description: "Code Llama 34B — most capable Code Llama variant",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "llama", ContextLen: 16384, License: "llama2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/CodeLlama-34b-Instruct-hf-GGUF/resolve/main/CodeLlama-34b-Instruct-hf-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "powerful"},
                },
                {
                        Name: "starcoder2", DisplayName: "StarCoder2 3B", Description: "BigCode's StarCoder2 — open code model for completion",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "1.8 GB", Family: "starcoder", ContextLen: 16384, License: "bigcode-openrail-m", Author: "BigCode",
                        URL: "https://huggingface.co/bartowski/starcoder2-3b-GGUF/resolve/main/starcoder2-3b-Q4_K_M.gguf",
                        Tags: []string{"code", "lightweight"},
                },
                {
                        Name: "starcoder2:15b", DisplayName: "StarCoder2 15B", Description: "BigCode's StarCoder2 15B — powerful open code model",
                        Parameters: "15B", Quantization: "Q4_K_M", Size: "8.8 GB", Family: "starcoder", ContextLen: 16384, License: "bigcode-openrail-m", Author: "BigCode",
                        URL: "https://huggingface.co/bartowski/starcoder2-15b-GGUF/resolve/main/starcoder2-15b-Q4_K_M.gguf",
                        Tags: []string{"code", "powerful"},
                },
                {
                        Name: "deepseek-coder:6.7b", DisplayName: "DeepSeek Coder 6.7B Instruct", Description: "DeepSeek Coder 6.7B — excellent code generation",
                        Parameters: "6.7B", Quantization: "Q4_K_M", Size: "3.9 GB", Family: "deepseek", ContextLen: 16384, License: "deepseek", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/deepseek-coder-6.7B-instruct-GGUF/resolve/main/deepseek-coder-6.7B-instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "fast"},
                },
                {
                        Name: "codegemma", DisplayName: "CodeGemma 7B IT", Description: "Google's CodeGemma — code-specialized Gemma variant",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/codegemma-7b-it-GGUF/resolve/main/codegemma-7b-it-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "google"},
                },

                // ═══════════════════════════════════════════
                // VISION / MULTIMODAL MODELS
                // ═══════════════════════════════════════════
                {
                        Name: "llava", DisplayName: "LLaVA 1.6 Mistral 7B", Description: "LLaVA 1.6 — vision-language model for image understanding and chat",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.5 GB", Family: "llava", ContextLen: 4096, License: "apache-2.0", Author: "LLaVA Team",
                        URL: "https://huggingface.co/bartowski/llava-v1.6-mistral-7b-GGUF/resolve/main/llava-v1.6-mistral-7b-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "image"},
                },
                {
                        Name: "llava:13b", DisplayName: "LLaVA 1.6 13B", Description: "LLaVA 1.6 13B — larger vision-language model",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "llava", ContextLen: 4096, License: "apache-2.0", Author: "LLaVA Team",
                        URL: "https://huggingface.co/bartowski/llava-v1.6-vicuna-13b-GGUF/resolve/main/llava-v1.6-vicuna-13b-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful"},
                },
                {
                        Name: "minicpm-v", DisplayName: "MiniCPM-V 2.6", Description: "MiniCPM-V 2.6 — efficient vision-language model by OpenBMB",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.9 GB", Family: "minicpm", ContextLen: 8192, License: "apache-2.0", Author: "OpenBMB",
                        URL: "https://huggingface.co/bartowski/MiniCPM-V-2_6-GGUF/resolve/main/MiniCPM-V-2_6-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "efficient"},
                },
                {
                        Name: "internvl2", DisplayName: "InternVL2 4B", Description: "InternVL2 4B — vision-language model by Shanghai AI Lab",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.5 GB", Family: "internvl", ContextLen: 8192, License: "mit", Author: "Shanghai AI Lab",
                        URL: "https://huggingface.co/bartowski/InternVL2-4B-GGUF/resolve/main/InternVL2-4B-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "lightweight"},
                },

                // ═══════════════════════════════════════════
                // COHERE / COMMAND FAMILY
                // ═══════════════════════════════════════════
                {
                        Name: "command-r", DisplayName: "Command R 35B", Description: "Cohere's Command R — RAG and tool use specialist with long context",
                        Parameters: "35B", Quantization: "Q4_K_M", Size: "20.8 GB", Family: "command", ContextLen: 131072, License: "cc-by-nc-4.0", Author: "Cohere",
                        URL: "https://huggingface.co/bartowski/command-r-GGUF/resolve/main/command-r-Q4_K_M.gguf",
                        Tags: []string{"rag", "tools", "powerful", "long-context"},
                },
                {
                        Name: "command-r-plus", DisplayName: "Command R+ 104B", Description: "Cohere's Command R+ — largest RAG and tool use model",
                        Parameters: "104B", Quantization: "Q4_K_M", Size: "60 GB", Family: "command", ContextLen: 131072, License: "cc-by-nc-4.0", Author: "Cohere",
                        URL: "https://huggingface.co/bartowski/command-r-plus-GGUF/resolve/main/command-r-plus-Q4_K_M.gguf",
                        Tags: []string{"rag", "tools", "frontier", "long-context"},
                },

                // ═══════════════════════════════════════════
                // SPECIALIZED / FINE-TUNED MODELS
                // ═══════════════════════════════════════════
                {
                        Name: "wizardlm2", DisplayName: "WizardLM 2 7B", Description: "Microsoft's WizardLM 2 — great for creative and complex tasks",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "wizard", ContextLen: 32768, License: "apache-2.0", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/WizardLM-2-7B-GGUF/resolve/main/WizardLM-2-7B-Q4_K_M.gguf",
                        Tags: []string{"creative", "instruct"},
                },
                {
                        Name: "wizardlm2:8x22b", DisplayName: "WizardLM 2 8x22B", Description: "Microsoft's WizardLM 2 MoE — large creative model",
                        Parameters: "141B", Quantization: "Q4_K_M", Size: "80 GB", Family: "wizard", ContextLen: 65536, License: "apache-2.0", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/WizardLM-2-8x22B-GGUF/resolve/main/WizardLM-2-8x22B-Q4_K_M.gguf",
                        Tags: []string{"creative", "powerful", "moe"},
                },
                {
                        Name: "dolphin-mistral", DisplayName: "Dolphin Mistral 7B v2.8", Description: "Dolphin Mistral — uncensored creative fine-tune by Eric Hartford",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Eric Hartford",
                        URL: "https://huggingface.co/bartowski/dolphin-2.8-mistral-7b-v02-GGUF/resolve/main/dolphin-2.8-mistral-7b-v02-Q4_K_M.gguf",
                        Tags: []string{"uncensored", "creative"},
                },
                {
                        Name: "dolphin-llama3", DisplayName: "Dolphin Llama 3 8B", Description: "Dolphin Llama 3 — uncensored Llama 3 fine-tune",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "apache-2.0", Author: "Eric Hartford",
                        URL: "https://huggingface.co/bartowski/Dolphin-Llama-3-8B-GGUF/resolve/main/Dolphin-Llama-3-8B-Q4_K_M.gguf",
                        Tags: []string{"uncensored", "creative"},
                },
                {
                        Name: "nous-hermes2", DisplayName: "Nous Hermes 2 Mistral 7B DPO", Description: "NousResearch Hermes 2 — creative and highly capable",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "NousResearch",
                        URL: "https://huggingface.co/bartowski/Nous-Hermes-2-Mistral-7B-DPO-GGUF/resolve/main/Nous-Hermes-2-Mistral-7B-DPO-Q4_K_M.gguf",
                        Tags: []string{"creative", "instruct"},
                },
                {
                        Name: "nous-hermes2:llama3", DisplayName: "Nous Hermes 2 Llama 3 8B", Description: "NousResearch Hermes 2 on Llama 3 — excellent for roleplay",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "apache-2.0", Author: "NousResearch",
                        URL: "https://huggingface.co/bartowski/Nous-Hermes-2-Pro-Llama-3-8B-GGUF/resolve/main/Nous-Hermes-2-Pro-Llama-3-8B-Q4_K_M.gguf",
                        Tags: []string{"creative", "instruct", "roleplay"},
                },
                {
                        Name: "openhermes", DisplayName: "OpenHermes 2.5 Mistral 7B", Description: "OpenHermes 2.5 — fine-tuned Mistral for diverse tasks",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 8192, License: "apache-2.0", Author: "Teknium",
                        URL: "https://huggingface.co/bartowski/OpenHermes-2.5-Mistral-7B-GGUF/resolve/main/OpenHermes-2.5-Mistral-7B-Q4_K_M.gguf",
                        Tags: []string{"creative", "instruct"},
                },
                {
                        Name: "solar", DisplayName: "SOLAR 10.7B Instruct", Description: "Upstage's SOLAR — efficient depth-up-scaled architecture",
                        Parameters: "10.7B", Quantization: "Q4_K_M", Size: "6.2 GB", Family: "solar", ContextLen: 4096, License: "apache-2.0", Author: "Upstage",
                        URL: "https://huggingface.co/bartowski/SOLAR-10.7B-Instruct-v1.0-GGUF/resolve/main/SOLAR-10.7B-Instruct-v1.0-Q4_K_M.gguf",
                        Tags: []string{"efficient", "instruct"},
                },
                {
                        Name: "yi", DisplayName: "Yi 1.5 6B Chat", Description: "01.AI's Yi 1.5 Chat — strong bilingual English/Chinese model",
                        Parameters: "6B", Quantization: "Q4_K_M", Size: "3.6 GB", Family: "yi", ContextLen: 4096, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-1.5-6B-Chat-GGUF/resolve/main/Yi-1.5-6B-Chat-Q4_K_M.gguf",
                        Tags: []string{"bilingual", "chinese"},
                },
                {
                        Name: "yi:34b", DisplayName: "Yi 1.5 34B Chat", Description: "01.AI's Yi 1.5 34B — powerful bilingual model",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "yi", ContextLen: 4096, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-1.5-34B-Chat-GGUF/resolve/main/Yi-1.5-34B-Chat-Q4_K_M.gguf",
                        Tags: []string{"bilingual", "chinese", "powerful"},
                },
                {
                        Name: "tinyllama", DisplayName: "TinyLlama 1.1B Chat v1.0", Description: "TinyLlama — ultra compact for edge devices and fast inference",
                        Parameters: "1.1B", Quantization: "Q4_K_M", Size: "0.6 GB", Family: "llama", ContextLen: 2048, License: "apache-2.0", Author: "TinyLlama",
                        URL: "https://huggingface.co/bartowski/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/TinyLlama-1.1B-Chat-v1.0-Q4_K_M.gguf",
                        Tags: []string{"ultra-light", "edge", "fast"},
                },
                {
                        Name: "orca-mini", DisplayName: "Orca Mini 3B v2", Description: "Orca Mini — small model trained with GPT-4 explanations",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "1.8 GB", Family: "orca", ContextLen: 2048, License: "apache-2.0", Author: "Pankaj Mathur",
                        URL: "https://huggingface.co/bartowski/Orca-Mini-3B-v2-GGUF/resolve/main/Orca-Mini-3B-v2-Q4_K_M.gguf",
                        Tags: []string{"lightweight", "fast"},
                },
                {
                        Name: "vicuna", DisplayName: "Vicuna 7B v1.5", Description: "Vicuna 7B — fine-tuned LLaMA for conversational AI",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "llama", ContextLen: 4096, License: "apache-2.0", Author: "LMSYS",
                        URL: "https://huggingface.co/bartowski/vicuna-7b-v1.5-GGUF/resolve/main/vicuna-7b-v1.5-Q4_K_M.gguf",
                        Tags: []string{"conversational", "chat"},
                },
                {
                        Name: "vicuna:13b", DisplayName: "Vicuna 13B v1.5", Description: "Vicuna 13B — larger fine-tuned model for quality conversations",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "llama", ContextLen: 4096, License: "apache-2.0", Author: "LMSYS",
                        URL: "https://huggingface.co/bartowski/vicuna-13b-v1.5-GGUF/resolve/main/vicuna-13b-v1.5-Q4_K_M.gguf",
                        Tags: []string{"conversational", "powerful"},
                },
                {
                        Name: "zephyr", DisplayName: "Zephyr 7B Alpha", Description: "HuggingFace's Zephyr — fine-tuned Mistral for helpfulness",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral", ContextLen: 32768, License: "mit", Author: "HuggingFace",
                        URL: "https://huggingface.co/bartowski/zephyr-7b-alpha-GGUF/resolve/main/zephyr-7b-alpha-Q4_K_M.gguf",
                        Tags: []string{"helpful", "instruct"},
                },
                {
                        Name: "stablelm2", DisplayName: "StableLM 2 12B Chat", Description: "Stability AI's StableLM 2 — open chat model",
                        Parameters: "12B", Quantization: "Q4_K_M", Size: "7.0 GB", Family: "stablelm", ContextLen: 4096, License: "apache-2.0", Author: "Stability AI",
                        URL: "https://huggingface.co/bartowski/stablelm-2-12b-chat-GGUF/resolve/main/stablelm-2-12b-chat-Q4_K_M.gguf",
                        Tags: []string{"stable-diffusion-org", "chat"},
                },
                {
                        Name: "falcon", DisplayName: "Falcon 7B Instruct", Description: "TII's Falcon 7B — efficient architecture for instruction following",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "falcon", ContextLen: 2048, License: "apache-2.0", Author: "TII",
                        URL: "https://huggingface.co/bartowski/Falcon-7B-Instruct-GGUF/resolve/main/Falcon-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"instruct", "efficient"},
                },
                {
                        Name: "mamba", DisplayName: "Mamba 2.8B", Description: "State Space Model — novel architecture beyond transformers",
                        Parameters: "2.8B", Quantization: "Q4_K_M", Size: "1.5 GB", Family: "mamba", ContextLen: 2048, License: "apache-2.0", Author: "Tri Dao",
                        URL: "https://huggingface.co/bartowski/mamba-2.8b-GGUF/resolve/main/mamba-2.8b-Q4_K_M.gguf",
                        Tags: []string{"ssm", "experimental", "lightweight"},
                },
                {
                        Name: "granite", DisplayName: "Granite 3.0 8B Instruct", Description: "IBM's Granite 3.0 — enterprise-grade instruct model",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "granite", ContextLen: 8192, License: "apache-2.0", Author: "IBM",
                        URL: "https://huggingface.co/bartowski/granite-3.0-8b-instruct-GGUF/resolve/main/granite-3.0-8b-instruct-Q4_K_M.gguf",
                        Tags: []string{"enterprise", "ibm", "instruct"},
                },
                {
                        Name: "granite:34b", DisplayName: "Granite 3.0 34B Instruct", Description: "IBM's Granite 3.0 34B — large enterprise model",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "granite", ContextLen: 8192, License: "apache-2.0", Author: "IBM",
                        URL: "https://huggingface.co/bartowski/granite-3.0-34b-instruct-GGUF/resolve/main/granite-3.0-34b-instruct-Q4_K_M.gguf",
                        Tags: []string{"enterprise", "ibm", "powerful"},
                },
                {
                        Name: "smollm2:1.7b", DisplayName: "SmolLM2 1.7B", Description: "HuggingFace's SmolLM2 — compact but capable small model",
                        Parameters: "1.7B", Quantization: "Q4_K_M", Size: "1.0 GB", Family: "smollm", ContextLen: 8192, License: "apache-2.0", Author: "HuggingFace",
                        URL: "https://huggingface.co/bartowski/smollm2-1.7B-instruct-GGUF/resolve/main/smollm2-1.7B-instruct-Q4_K_M.gguf",
                        Tags: []string{"ultra-light", "fast", "huggingface"},
                },
                {
                        Name: "smollm2:135m", DisplayName: "SmolLM2 135M", Description: "SmolLM2 135M — one of the smallest capable LLMs",
                        Parameters: "135M", Quantization: "Q4_K_M", Size: "0.1 GB", Family: "smollm", ContextLen: 2048, License: "apache-2.0", Author: "HuggingFace",
                        URL: "https://huggingface.co/bartowski/smollm2-135M-GGUF/resolve/main/smollm2-135M-Q4_K_M.gguf",
                        Tags: []string{"ultra-light", "experimental", "tiny"},
                },
                {
                        Name: "olmo2", DisplayName: "OLMo 2 7B Instruct", Description: "AI2's OLMo 2 — fully open model with training data",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "olmo", ContextLen: 4096, License: "apache-2.0", Author: "AI2",
                        URL: "https://huggingface.co/bartowski/OLMo-2-1124-7B-Instruct-GGUF/resolve/main/OLMo-2-1124-7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"open-data", "instruct", "research"},
                },
                {
                        Name: "arctic", DisplayName: "Arctic Instruct", Description: "Snowflake's Arctic — enterprise MoE model",
                        Parameters: "10B-active/480B", Quantization: "Q4_K_M", Size: "28 GB", Family: "arctic", ContextLen: 4096, License: "apache-2.0", Author: "Snowflake",
                        URL: "https://huggingface.co/bartowski/Arctic-Instruct-GGUF/resolve/main/Arctic-Instruct-Q4_K_M.gguf",
                        Tags: []string{"enterprise", "moe", "powerful"},
                },
                {
                        Name: "glm4", DisplayName: "GLM-4 9B Chat", Description: "Zhipu AI's GLM-4 — bilingual English/Chinese chat model",
                        Parameters: "9B", Quantization: "Q4_K_M", Size: "5.2 GB", Family: "glm", ContextLen: 8192, License: "apache-2.0", Author: "Zhipu AI",
                        URL: "https://huggingface.co/bartowski/glm-4-9b-chat-GGUF/resolve/main/glm-4-9b-chat-Q4_K_M.gguf",
                        Tags: []string{"bilingual", "chinese", "chat"},
                },
                {
                        Name: "jais", DisplayName: "Jais 13B Chat", Description: "Inception's Jais — Arabic-English bilingual model",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "jais", ContextLen: 8192, License: "apache-2.0", Author: "Inception",
                        URL: "https://huggingface.co/bartowski/jais-13b-chat-GGUF/resolve/main/jais-13b-chat-Q4_K_M.gguf",
                        Tags: []string{"arabic", "bilingual", "multilingual"},
                },
                {
                        Name: "marco-o1", DisplayName: "Marco-o1", Description: "AIDC-AI's Marco-o1 — reasoning model with chain-of-thought",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "marco", ContextLen: 8192, License: "apache-2.0", Author: "AIDC-AI",
                        URL: "https://huggingface.co/bartowski/Marco-o1-GGUF/resolve/main/Marco-o1-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "latest"},
                },
                {
                        Name: "nemotron", DisplayName: "Nemotron Mini 4B Instruct", Description: "NVIDIA's Nemotron Mini — compact yet powerful",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.4 GB", Family: "nemotron", ContextLen: 4096, License: "nvidia", Author: "NVIDIA",
                        URL: "https://huggingface.co/bartowski/Nemotron-Mini-4B-Instruct-GGUF/resolve/main/Nemotron-Mini-4B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"nvidia", "lightweight"},
                },
                {
                        Name: "persimmon", DisplayName: "Persimmon 8B Chat", Description: "Adept's Persimmon — efficient architecture for dialogue",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "persimmon", ContextLen: 16384, License: "apache-2.0", Author: "Adept",
                        URL: "https://huggingface.co/bartowski/persimmon-8b-chat-GGUF/resolve/main/persimmon-8b-chat-Q4_K_M.gguf",
                        Tags: []string{"efficient", "dialogue"},
                },
                {
                        Name: "plamo2", DisplayName: "PLaMo 2 13B", Description: "Preferred Networks' PLaMo 2 — Japanese-focused model",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "plamo", ContextLen: 4096, License: "apache-2.0", Author: "PFN",
                        URL: "https://huggingface.co/bartowski/plamo-2-13b-GGUF/resolve/main/plamo-2-13b-Q4_K_M.gguf",
                        Tags: []string{"japanese", "bilingual"},
                },
                {
                        Name: "clip-model", DisplayName: "Medical-Llama3 8B", Description: "Medical fine-tuned Llama 3 for healthcare Q&A",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "apache-2.0", Author: "Community",
                        URL: "https://huggingface.co/bartowski/Medical-Llama3-8B-GGUF/resolve/main/Medical-Llama3-8B-Q4_K_M.gguf",
                        Tags: []string{"medical", "healthcare", "specialized"},
                },
                {
                        Name: "finance-llm", DisplayName: "FinGPT 7B", Description: "Financial LLM fine-tuned for financial analysis and forecasting",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "llama", ContextLen: 4096, License: "apache-2.0", Author: "Community",
                        URL: "https://huggingface.co/bartowski/FinGPT-7B-GGUF/resolve/main/FinGPT-7B-Q4_K_M.gguf",
                        Tags: []string{"finance", "specialized", "analysis"},
                },
                {
                        Name: "legal-llm", DisplayName: "Legal LLM 7B", Description: "Legal domain model for contract analysis and legal reasoning",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "llama", ContextLen: 4096, License: "apache-2.0", Author: "Community",
                        URL: "https://huggingface.co/bartowski/LawLLM-7B-GGUF/resolve/main/LawLLM-7B-Q4_K_M.gguf",
                        Tags: []string{"legal", "specialized", "analysis"},
                },
                {
                        Name: "embedding-model", DisplayName: "BGE-M3 Embedding", Description: "BAAI's BGE-M3 — multilingual embedding model for RAG and search",
                        Parameters: "568M", Quantization: "Q4_K_M", Size: "0.4 GB", Family: "bge", ContextLen: 8192, License: "mit", Author: "BAAI",
                        URL: "https://huggingface.co/bartowski/bge-m3-GGUF/resolve/main/bge-m3-Q4_K_M.gguf",
                        Tags: []string{"embedding", "rag", "search", "multilingual"},
                },

                // ═══════════════════════════════════════════
                // EXPANDED MODEL CATALOG — 50+ ADDITIONAL MODELS
                // ═══════════════════════════════════════════

                // ── LLAMA FAMILY EXTENSIONS ──
                {
                        Name: "llama3.2-vision", DisplayName: "Llama 3.2 11B Vision", Description: "Meta's Llama 3.2 Vision — multimodal model that understands images and text together with 11B parameters",
                        Parameters: "11B", Quantization: "Q4_K_M", Size: "6.4 GB", Family: "llama", ContextLen: 131072, License: "llama3.2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-3.2-11B-Vision-Instruct-GGUF/resolve/main/Llama-3.2-11B-Vision-Instruct-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "popular", "instruct"},
                },
                {
                        Name: "llama3.2-vision:90b", DisplayName: "Llama 3.2 90B Vision", Description: "Meta's Llama 3.2 90B Vision — frontier multimodal model with exceptional image understanding",
                        Parameters: "90B", Quantization: "Q4_K_M", Size: "52 GB", Family: "llama", ContextLen: 131072, License: "llama3.2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-3.2-90B-Vision-Instruct-GGUF/resolve/main/Llama-3.2-90B-Vision-Instruct-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful", "frontier"},
                },

                // ── QWEN FAMILY EXTENSIONS ──
                {
                        Name: "qwen2.5:0.5b", DisplayName: "Qwen 2.5 0.5B Instruct", Description: "Alibaba's Qwen 2.5 0.5B — ultra-tiny model perfect for edge devices and embedded systems",
                        Parameters: "0.5B", Quantization: "Q4_K_M", Size: "0.4 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/Qwen2.5-0.5B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"tiny", "edge", "multilingual", "fast"},
                },
                {
                        Name: "qwen2.5:1.5b", DisplayName: "Qwen 2.5 1.5B Instruct", Description: "Alibaba's Qwen 2.5 1.5B — small but capable multilingual model with 128K context",
                        Parameters: "1.5B", Quantization: "Q4_K_M", Size: "1.1 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-1.5B-Instruct-GGUF/resolve/main/Qwen2.5-1.5B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"small", "multilingual", "fast", "instruct"},
                },
                {
                        Name: "qwen2.5:3b", DisplayName: "Qwen 2.5 3B Instruct", Description: "Alibaba's Qwen 2.5 3B — balanced performance and efficiency with multilingual support",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "2.0 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-3B-Instruct-GGUF/resolve/main/Qwen2.5-3B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"small", "multilingual", "balanced", "instruct"},
                },
                {
                        Name: "qwen2.5-coder:1.5b", DisplayName: "Qwen 2.5 Coder 1.5B", Description: "Qwen 2.5 Coder 1.5B — ultra-compact coding assistant for quick completions",
                        Parameters: "1.5B", Quantization: "Q4_K_M", Size: "1.1 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-1.5B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-1.5B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "small", "fast", "instruct"},
                },
                {
                        Name: "qwen2.5-coder:14b", DisplayName: "Qwen 2.5 Coder 14B", Description: "Qwen 2.5 Coder 14B — advanced coding model with deep code understanding and generation",
                        Parameters: "14B", Quantization: "Q4_K_M", Size: "8.5 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-14B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-14B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "powerful", "instruct"},
                },
                {
                        Name: "qwen2.5-math:1.5b", DisplayName: "Qwen 2.5 Math 1.5B", Description: "Qwen 2.5 Math 1.5B — compact math specialist for quick calculations and reasoning",
                        Parameters: "1.5B", Quantization: "Q4_K_M", Size: "1.1 GB", Family: "qwen", ContextLen: 4096, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Math-1.5B-Instruct-GGUF/resolve/main/Qwen2.5-Math-1.5B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"math", "small", "reasoning", "science"},
                },
                {
                        Name: "qwen2.5-math:72b", DisplayName: "Qwen 2.5 Math 72B", Description: "Qwen 2.5 Math 72B — frontier-class mathematical reasoning model",
                        Parameters: "72B", Quantization: "Q4_K_M", Size: "41 GB", Family: "qwen", ContextLen: 4096, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2.5-Math-72B-Instruct-GGUF/resolve/main/Qwen2.5-Math-72B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"math", "powerful", "reasoning", "frontier"},
                },
                {
                        Name: "qwen2-vl:72b", DisplayName: "Qwen 2 VL 72B Instruct", Description: "Qwen 2 VL 72B — massive vision-language model with exceptional visual understanding",
                        Parameters: "72B", Quantization: "Q4_K_M", Size: "41 GB", Family: "qwen", ContextLen: 32768, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/Qwen2-VL-72B-Instruct-GGUF/resolve/main/Qwen2-VL-72B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful", "frontier"},
                },
                {
                        Name: "qwq:32b", DisplayName: "QwQ 32B Preview", Description: "Alibaba's QwQ 32B — advanced reasoning model with deep chain-of-thought capabilities",
                        Parameters: "32B", Quantization: "Q4_K_M", Size: "18.5 GB", Family: "qwen", ContextLen: 131072, License: "apache-2.0", Author: "Alibaba",
                        URL: "https://huggingface.co/bartowski/QwQ-32B-Preview-GGUF/resolve/main/QwQ-32B-Preview-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "chain-of-thought", "powerful"},
                },

                // ── DEEPSEEK FAMILY EXTENSIONS ──
                {
                        Name: "deepseek-r1:1.5b", DisplayName: "DeepSeek R1 1.5B Distill", Description: "DeepSeek R1 1.5B distill — ultra-compact reasoning model for edge deployment",
                        Parameters: "1.5B", Quantization: "Q4_K_M", Size: "1.1 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-1.5B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-1.5B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "tiny", "fast", "chain-of-thought"},
                },
                {
                        Name: "deepseek-r1:7b", DisplayName: "DeepSeek R1 7B Distill", Description: "DeepSeek R1 7B distill — efficient reasoning model with strong chain-of-thought",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.2 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-7B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "popular", "chain-of-thought"},
                },
                {
                        Name: "deepseek-r1:671b", DisplayName: "DeepSeek R1 671B", Description: "DeepSeek R1 671B — full-scale frontier reasoning model, the most powerful open reasoning model",
                        Parameters: "671B", Quantization: "Q4_K_M", Size: "380 GB", Family: "deepseek", ContextLen: 131072, License: "mit", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-R1-GGUF/resolve/main/DeepSeek-R1-Q4_K_M.gguf",
                        Tags: []string{"reasoning", "frontier", "chain-of-thought", "massive"},
                },
                {
                        Name: "deepseek-v2.5", DisplayName: "DeepSeek V2.5", Description: "DeepSeek V2.5 — improved general-purpose MoE model with better instruction following",
                        Parameters: "21B", Quantization: "Q4_K_M", Size: "12 GB", Family: "deepseek", ContextLen: 131072, License: "deepseek", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-V2.5-GGUF/resolve/main/DeepSeek-V2.5-Q4_K_M.gguf",
                        Tags: []string{"moe", "general", "instruct", "popular"},
                },
                {
                        Name: "deepseek-coder-v2:236b", DisplayName: "DeepSeek Coder V2 236B", Description: "DeepSeek Coder V2 236B MoE — massive code model with expert-level programming",
                        Parameters: "236B", Quantization: "Q4_K_M", Size: "135 GB", Family: "deepseek", ContextLen: 131072, License: "deepseek", Author: "DeepSeek",
                        URL: "https://huggingface.co/bartowski/DeepSeek-Coder-V2-Instruct-GGUF/resolve/main/DeepSeek-Coder-V2-Instruct-Q4_K_M.gguf",
                        Tags: []string{"code", "moe", "massive", "programming"},
                },

                // ── GEMMA FAMILY EXTENSIONS ──
                {
                        Name: "gemma3:text", DisplayName: "Gemma 3 4B Text", Description: "Google's Gemma 3 4B text-only — lighter version without vision for faster text processing",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.8 GB", Family: "gemma", ContextLen: 131072, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/gemma-3-4b-it-GGUF/resolve/main/gemma-3-4b-it-Q4_K_M.gguf",
                        Tags: []string{"text", "fast", "instruct", "efficient"},
                },
                {
                        Name: "codegemma:2b", DisplayName: "CodeGemma 2B", Description: "Google's CodeGemma 2B — ultra-compact code completion model for IDE integration",
                        Parameters: "2B", Quantization: "Q4_K_M", Size: "1.5 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/codegemma-2b-GGUF/resolve/main/codegemma-2b-Q4_K_M.gguf",
                        Tags: []string{"code", "tiny", "fast", "completion"},
                },
                {
                        Name: "recurrentgemma", DisplayName: "RecurrentGemma 2B IT", Description: "Google's RecurrentGemma — novel recurrent architecture for efficient inference",
                        Parameters: "2B", Quantization: "Q4_K_M", Size: "1.5 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/recurrent-gemma-2b-it-GGUF/resolve/main/recurrent-gemma-2b-it-Q4_K_M.gguf",
                        Tags: []string{"recurrent", "efficient", "experimental", "small"},
                },

                // ── MISTRAL FAMILY EXTENSIONS ──
                {
                        Name: "mistral-large", DisplayName: "Mistral Large 24B Instruct", Description: "Mistral Large — flagship Mistral model with top-tier reasoning and knowledge",
                        Parameters: "24B", Quantization: "Q4_K_M", Size: "14 GB", Family: "mistral", ContextLen: 131072, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mistral-Large-24B-Instruct-2503-GGUF/resolve/main/Mistral-Large-24B-Instruct-2503-Q4_K_M.gguf",
                        Tags: []string{"powerful", "instruct", "reasoning", "flagship"},
                },
                {
                        Name: "mistral:0.7b", DisplayName: "Mistral 0.7B", Description: "Mistral 0.7B — ultra-tiny Mistral variant for edge and embedded use cases",
                        Parameters: "0.7B", Quantization: "Q4_K_M", Size: "0.5 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mistral-7B-v0.1-GGUF/resolve/main/Mistral-7B-v0.1-Q4_K_M.gguf",
                        Tags: []string{"tiny", "edge", "fast"},
                },
                {
                        Name: "pixtral:large", DisplayName: "Pixtral Large 24B", Description: "Mistral's Pixtral Large — advanced vision-language model for complex image understanding",
                        Parameters: "24B", Quantization: "Q4_K_M", Size: "14 GB", Family: "mistral", ContextLen: 131072, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Pixtral-Large-24B-1209-GGUF/resolve/main/Pixtral-Large-24B-1209-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful"},
                },

                // ── MICROSOFT FAMILY EXTENSIONS ──
                {
                        Name: "phi3:moe", DisplayName: "Phi-3 MoE 16x3.8B", Description: "Microsoft's Phi-3 MoE — Mixture of Experts variant for efficient inference at scale",
                        Parameters: "16x3.8B", Quantization: "Q4_K_M", Size: "9 GB", Family: "phi", ContextLen: 131072, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/Phi-3.5-MoE-instruct-GGUF/resolve/main/Phi-3.5-MoE-instruct-Q4_K_M.gguf",
                        Tags: []string{"moe", "efficient", "instruct"},
                },
                {
                        Name: "wizardlm2:13b", DisplayName: "WizardLM 2 13B", Description: "Microsoft's WizardLM 2 13B — balanced model for creative writing and complex reasoning",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.5 GB", Family: "wizardlm", ContextLen: 32768, License: "llama2", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/WizardLM-2-13B-GGUF/resolve/main/WizardLM-2-13B-Q4_K_M.gguf",
                        Tags: []string{"creative", "reasoning", "balanced"},
                },

                // ── COHERE FAMILY EXTENSIONS ──
                {
                        Name: "command-r7b", DisplayName: "Command R7B Arabic", Description: "Cohere's Command R7B — compact RAG model optimized for Arabic and multilingual retrieval",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.2 GB", Family: "command", ContextLen: 131072, License: "cc-by-nc-4.0", Author: "Cohere",
                        URL: "https://huggingface.co/bartowski/command-r7b-arabic-GGUF/resolve/main/command-r7b-arabic-Q4_K_M.gguf",
                        Tags: []string{"rag", "arabic", "multilingual", "tools"},
                },
                {
                        Name: "aya", DisplayName: "Aya 23 8B", Description: "Cohere's Aya 23 — multilingual model covering 23 languages with strong cross-lingual transfer",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.8 GB", Family: "cohere", ContextLen: 8192, License: "cc-by-nc-4.0", Author: "Cohere",
                        URL: "https://huggingface.co/bartowski/aya-23-8B-GGUF/resolve/main/aya-23-8B-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "23-languages", "instruct"},
                },
                {
                        Name: "aya:35b", DisplayName: "Aya 23 35B", Description: "Cohere's Aya 23 35B — large multilingual model with exceptional cross-lingual performance",
                        Parameters: "35B", Quantization: "Q4_K_M", Size: "20 GB", Family: "cohere", ContextLen: 8192, License: "cc-by-nc-4.0", Author: "Cohere",
                        URL: "https://huggingface.co/bartowski/aya-23-35B-GGUF/resolve/main/aya-23-35B-Q4_K_M.gguf",
                        Tags: []string{"multilingual", "23-languages", "powerful"},
                },

                // ── IBM GRANITE EXTENSIONS ──
                {
                        Name: "granite:3b", DisplayName: "Granite 3.0 3B Instruct", Description: "IBM's Granite 3.0 3B — compact enterprise model for resource-constrained environments",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "1.9 GB", Family: "granite", ContextLen: 8192, License: "apache-2.0", Author: "IBM",
                        URL: "https://huggingface.co/bartowski/granite-3b-code-instruct-GGUF/resolve/main/granite-3b-code-instruct-Q4_K_M.gguf",
                        Tags: []string{"enterprise", "small", "code"},
                },
                {
                        Name: "granite-guardian", DisplayName: "Granite Guardian 7B", Description: "IBM's Granite Guardian — content safety and risk detection model for enterprise deployment",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.2 GB", Family: "granite", ContextLen: 8192, License: "apache-2.0", Author: "IBM",
                        URL: "https://huggingface.co/bartowski/granite-guardian-7b-GGUF/resolve/main/granite-guardian-7b-Q4_K_M.gguf",
                        Tags: []string{"safety", "guardian", "enterprise", "moderation"},
                },

                // ── NVIDIA FAMILY ──
                {
                        Name: "nemotron:mini", DisplayName: "Nemotron Mini 4B Instruct", Description: "NVIDIA's Nemotron Mini 4B — compact yet powerful instruction-following model",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.5 GB", Family: "nemotron", ContextLen: 4096, License: "nvidia", Author: "NVIDIA",
                        URL: "https://huggingface.co/bartowski/nemotron-mini-4b-instruct-GGUF/resolve/main/nemotron-mini-4b-instruct-Q4_K_M.gguf",
                        Tags: []string{"nvidia", "compact", "instruct"},
                },

                // ── SNOWFLAKE ARCTIC EXTENSIONS ──
                {
                        Name: "arctic-embed", DisplayName: "Arctic Embed 335M", Description: "Snowflake's Arctic Embed — high-quality embedding model optimized for retrieval and RAG",
                        Parameters: "335M", Quantization: "Q4_K_M", Size: "0.3 GB", Family: "arctic", ContextLen: 512, License: "apache-2.0", Author: "Snowflake",
                        URL: "https://huggingface.co/bartowski/arctic-embed-l-GGUF/resolve/main/arctic-embed-l-Q4_K_M.gguf",
                        Tags: []string{"embedding", "rag", "retrieval", "search"},
                },

                // ── STABILITY AI ──
                {
                        Name: "stablelm2:1.6b", DisplayName: "StableLM 2 1.6B Chat", Description: "Stability AI's StableLM 2 1.6B — ultra-compact chat model for fast responses",
                        Parameters: "1.6B", Quantization: "Q4_K_M", Size: "1.0 GB", Family: "stablelm", ContextLen: 4096, License: "apache-2.0", Author: "Stability AI",
                        URL: "https://huggingface.co/bartowski/stablelm-2-1_6b-chat-GGUF/resolve/main/stablelm-2-1_6b-chat-Q4_K_M.gguf",
                        Tags: []string{"small", "chat", "fast"},
                },
                {
                        Name: "stablelm-zephyr", DisplayName: "StableLM Zephyr 3B", Description: "Stability AI's StableLM Zephyr — fine-tuned for helpfulness and safety",
                        Parameters: "3B", Quantization: "Q4_K_M", Size: "1.9 GB", Family: "stablelm", ContextLen: 4096, License: "apache-2.0", Author: "Stability AI",
                        URL: "https://huggingface.co/bartowski/stablelm-zephyr-3b-GGUF/resolve/main/stablelm-zephyr-3b-Q4_K_M.gguf",
                        Tags: []string{"zephyr", "helpful", "safe", "small"},
                },

                // ── 01.AI YI FAMILY EXTENSIONS ──
                {
                        Name: "yi:6b", DisplayName: "Yi 1.5 6B Chat", Description: "01.AI's Yi 1.5 6B — compact bilingual model with strong English and Chinese capabilities",
                        Parameters: "6B", Quantization: "Q4_K_M", Size: "3.5 GB", Family: "yi", ContextLen: 4096, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-1.5-6B-Chat-GGUF/resolve/main/Yi-1.5-6B-Chat-Q4_K_M.gguf",
                        Tags: []string{"bilingual", "chinese", "english", "chat"},
                },
                {
                        Name: "yi:9b", DisplayName: "Yi 1.5 9B Chat", Description: "01.AI's Yi 1.5 9B — balanced bilingual model for general tasks",
                        Parameters: "9B", Quantization: "Q4_K_M", Size: "5.2 GB", Family: "yi", ContextLen: 4096, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-1.5-9B-Chat-GGUF/resolve/main/Yi-1.5-9B-Chat-Q4_K_M.gguf",
                        Tags: []string{"bilingual", "balanced", "chinese", "english"},
                },
                {
                        Name: "yi-vl", DisplayName: "Yi VL 6B", Description: "01.AI's Yi Vision-Language — bilingual multimodal model for image understanding",
                        Parameters: "6B", Quantization: "Q4_K_M", Size: "3.8 GB", Family: "yi", ContextLen: 4096, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-VL-6B-GGUF/resolve/main/Yi-VL-6B-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "bilingual", "chinese"},
                },
                {
                        Name: "yi-coder", DisplayName: "Yi Coder 9B", Description: "01.AI's Yi Coder — bilingual code model with excellent programming and multilingual support",
                        Parameters: "9B", Quantization: "Q4_K_M", Size: "5.2 GB", Family: "yi", ContextLen: 131072, License: "apache-2.0", Author: "01.AI",
                        URL: "https://huggingface.co/bartowski/Yi-Coder-9B-Chat-GGUF/resolve/main/Yi-Coder-9B-Chat-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "bilingual", "long-context"},
                },

                // ── HUGGINGFACE SMOLLM EXTENSIONS ──
                {
                        Name: "smollm2:360m", DisplayName: "SmolLM2 360M", Description: "HuggingFace's SmolLM2 360M — one of the smallest capable language models, great for testing",
                        Parameters: "360M", Quantization: "Q4_K_M", Size: "0.3 GB", Family: "smollm", ContextLen: 2048, License: "apache-2.0", Author: "HuggingFace",
                        URL: "https://huggingface.co/bartowski/SmolLM2-360M-Instruct-GGUF/resolve/main/SmolLM2-360M-Instruct-Q4_K_M.gguf",
                        Tags: []string{"tiny", "edge", "testing", "smol"},
                },
                {
                        Name: "smollm2", DisplayName: "SmolLM2 1.7B Instruct", Description: "HuggingFace's SmolLM2 1.7B — compact but surprisingly capable instruction model",
                        Parameters: "1.7B", Quantization: "Q4_K_M", Size: "1.0 GB", Family: "smollm", ContextLen: 8192, License: "apache-2.0", Author: "HuggingFace",
                        URL: "https://huggingface.co/bartowski/SmolLM2-1.7B-Instruct-GGUF/resolve/main/SmolLM2-1.7B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"small", "instruct", "smol", "efficient"},
                },

                // ── OPEN MODELS ──
                {
                        Name: "olmo2:13b", DisplayName: "OLMo 2 13B Instruct", Description: "AI2's OLMo 2 13B — fully open model with all training data and code available",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.5 GB", Family: "olmo", ContextLen: 4096, License: "apache-2.0", Author: "AI2",
                        URL: "https://huggingface.co/bartowski/OLMo-2-1124-13B-Instruct-GGUF/resolve/main/OLMo-2-1124-13B-Instruct-Q4_K_M.gguf",
                        Tags: []string{"open", "reproducible", "instruct"},
                },

                // ── CHINESE / MULTILINGUAL MODELS ──
                {
                        Name: "chatglm3", DisplayName: "ChatGLM3 6B", Description: "Zhipu AI's ChatGLM3 — efficient bilingual English/Chinese conversation model",
                        Parameters: "6B", Quantization: "Q4_K_M", Size: "3.5 GB", Family: "chatglm", ContextLen: 8192, License: "apache-2.0", Author: "Zhipu AI",
                        URL: "https://huggingface.co/bartowski/chatglm3-6b-GGUF/resolve/main/chatglm3-6b-Q4_K_M.gguf",
                        Tags: []string{"chinese", "bilingual", "chat", "efficient"},
                },
                {
                        Name: "baichuan2", DisplayName: "Baichuan2 7B Chat", Description: "Baichuan's Baichuan2 7B — strong bilingual model optimized for Chinese and English",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "baichuan", ContextLen: 4096, License: "apache-2.0", Author: "Baichuan",
                        URL: "https://huggingface.co/bartowski/Baichuan2-7B-Chat-GGUF/resolve/main/Baichuan2-7B-Chat-Q4_K_M.gguf",
                        Tags: []string{"chinese", "bilingual", "chat"},
                },
                {
                        Name: "internlm2", DisplayName: "InternLM2 7B Chat", Description: "Shanghai AI Lab's InternLM2 — bilingual model with strong reasoning and long context",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.2 GB", Family: "internlm", ContextLen: 32768, License: "apache-2.0", Author: "Shanghai AI Lab",
                        URL: "https://huggingface.co/bartowski/internlm2-7b-chat-GGUF/resolve/main/internlm2-7b-chat-Q4_K_M.gguf",
                        Tags: []string{"chinese", "bilingual", "reasoning", "long-context"},
                },
                {
                        Name: "internlm2:20b", DisplayName: "InternLM2 20B Chat", Description: "Shanghai AI Lab's InternLM2 20B — powerful bilingual model for demanding tasks",
                        Parameters: "20B", Quantization: "Q4_K_M", Size: "11.5 GB", Family: "internlm", ContextLen: 32768, License: "apache-2.0", Author: "Shanghai AI Lab",
                        URL: "https://huggingface.co/bartowski/internlm2-20b-chat-GGUF/resolve/main/internlm2-20b-chat-Q4_K_M.gguf",
                        Tags: []string{"chinese", "bilingual", "powerful", "long-context"},
                },
                {
                        Name: "minicpm3", DisplayName: "MiniCPM3 4B", Description: "OpenBMB's MiniCPM3 4B — efficient compact model with impressive performance for its size",
                        Parameters: "4B", Quantization: "Q4_K_M", Size: "2.4 GB", Family: "minicpm", ContextLen: 32768, License: "apache-2.0", Author: "OpenBMB",
                        URL: "https://huggingface.co/bartowski/MiniCPM3-4B-GGUF/resolve/main/MiniCPM3-4B-Q4_K_M.gguf",
                        Tags: []string{"efficient", "compact", "chinese", "instruct"},
                },

                // ── VISION-LANGUAGE EXTENSIONS ──
                {
                        Name: "llava:34b", DisplayName: "LLaVA 1.6 34B", Description: "LLaVA 1.6 34B — largest LLaVA model with exceptional visual reasoning capabilities",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "llava", ContextLen: 4096, License: "llama2", Author: "LLaVA Team",
                        URL: "https://huggingface.co/bartowski/llava-v1.6-34b-GGUF/resolve/main/llava-v1.6-34b-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful"},
                },
                {
                        Name: "internvl2:8b", DisplayName: "InternVL2 8B", Description: "InternVL2 8B — mid-size vision-language model with strong OCR and document understanding",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.8 GB", Family: "internvl", ContextLen: 8192, License: "apache-2.0", Author: "Shanghai AI Lab",
                        URL: "https://huggingface.co/bartowski/InternVL2-8B-GGUF/resolve/main/InternVL2-8B-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "ocr", "document"},
                },
                {
                        Name: "internvl2:26b", DisplayName: "InternVL2 26B", Description: "InternVL2 26B — high-quality vision-language model for complex image analysis",
                        Parameters: "26B", Quantization: "Q4_K_M", Size: "15 GB", Family: "internvl", ContextLen: 8192, License: "apache-2.0", Author: "Shanghai AI Lab",
                        URL: "https://huggingface.co/bartowski/InternVL2-26B-GGUF/resolve/main/InternVL2-26B-Q4_K_M.gguf",
                        Tags: []string{"vision", "multimodal", "powerful", "ocr"},
                },

                // ── FINE-TUNED / UNCENSORED MODELS ──
                {
                        Name: "dolphin-mixtral", DisplayName: "Dolphin Mixtral 8x7B", Description: "Dolphin Mixtral — uncensored creative fine-tune of Mixtral by Eric Hartford",
                        Parameters: "47B", Quantization: "Q4_K_M", Size: "27 GB", Family: "dolphin", ContextLen: 32768, License: "apache-2.0", Author: "Eric Hartford",
                        URL: "https://huggingface.co/bartowski/dolphin-2.6-mixtral-8x7b-GGUF/resolve/main/dolphin-2.6-mixtral-8x7b-Q4_K_M.gguf",
                        Tags: []string{"uncensored", "creative", "moe", "roleplay"},
                },
                {
                        Name: "dolphin-phi3", DisplayName: "Dolphin Phi-3 3.8B", Description: "Dolphin Phi-3 — uncensored creative fine-tune of Microsoft's Phi-3",
                        Parameters: "3.8B", Quantization: "Q4_K_M", Size: "2.2 GB", Family: "dolphin", ContextLen: 8192, License: "mit", Author: "Eric Hartford",
                        URL: "https://huggingface.co/bartowski/dolphin-2.8-phi-3-3.8b-GGUF/resolve/main/dolphin-2.8-phi-3-3.8b-Q4_K_M.gguf",
                        Tags: []string{"uncensored", "creative", "small", "roleplay"},
                },
                {
                        Name: "nous-hermes2:mixtral", DisplayName: "Nous Hermes 2 Mixtral 8x7B", Description: "NousResearch Hermes 2 on Mixtral — creative and uncensored MoE model",
                        Parameters: "47B", Quantization: "Q4_K_M", Size: "27 GB", Family: "nous", ContextLen: 32768, License: "apache-2.0", Author: "NousResearch",
                        URL: "https://huggingface.co/bartowski/Nous-Hermes-2-Mixtral-8x7B-DPO-GGUF/resolve/main/Nous-Hermes-2-Mixtral-8x7B-DPO-Q4_K_M.gguf",
                        Tags: []string{"uncensored", "creative", "moe", "roleplay"},
                },
                {
                        Name: "openhermes:2.5", DisplayName: "OpenHermes 2.5 Mistral 7B", Description: "OpenHermes 2.5 — highly-rated fine-tune of Mistral for diverse conversational tasks",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "hermes", ContextLen: 8192, License: "apache-2.0", Author: "Teknium",
                        URL: "https://huggingface.co/bartowski/openhermes-2.5-mistral-7b-GGUF/resolve/main/openhermes-2.5-mistral-7b-Q4_K_M.gguf",
                        Tags: []string{"creative", "chat", "fine-tuned", "popular"},
                },
                {
                        Name: "capybara", DisplayName: "Capybara 34B", Description: "Capybara 34B — creative fine-tune with excellent storytelling and roleplay capabilities",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "capybara", ContextLen: 8192, License: "apache-2.0", Author: "LDJ",
                        URL: "https://huggingface.co/bartowski/Capybara-34B-GGUF/resolve/main/Capybara-34B-Q4_K_M.gguf",
                        Tags: []string{"creative", "roleplay", "storytelling"},
                },
                {
                        Name: "mythomax", DisplayName: "MythoMax 13B L2", Description: "MythoMax 13B — popular roleplay and creative writing model on Llama 2 base",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.5 GB", Family: "mythomax", ContextLen: 4096, License: "llama2", Author: "Gryphe",
                        URL: "https://huggingface.co/bartowski/MythoMax-L2-13b-GGUF/resolve/main/MythoMax-L2-13b-Q4_K_M.gguf",
                        Tags: []string{"roleplay", "creative", "storytelling", "popular"},
                },

                // ── CODE MODELS EXTENSIONS ──
                {
                        Name: "codellama:13b", DisplayName: "Code Llama 13B Instruct", Description: "Code Llama 13B — balanced code model for generation and understanding",
                        Parameters: "13B", Quantization: "Q4_K_M", Size: "7.5 GB", Family: "codellama", ContextLen: 16384, License: "llama2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/CodeLlama-13b-Instruct-hf-GGUF/resolve/main/CodeLlama-13b-Instruct-hf-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "balanced"},
                },
                {
                        Name: "codellama:python", DisplayName: "Code Llama Python 7B", Description: "Code Llama Python 7B — specialized for Python code generation and debugging",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "codellama", ContextLen: 16384, License: "llama2", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/CodeLlama-7b-Python-hf-GGUF/resolve/main/CodeLlama-7b-Python-hf-Q4_K_M.gguf",
                        Tags: []string{"code", "python", "specialized"},
                },
                {
                        Name: "starcoder2:7b", DisplayName: "StarCoder2 7B", Description: "BigCode's StarCoder2 7B — mid-size code model for completion and generation",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "starcoder", ContextLen: 16384, License: "bigcode-openrail-m", Author: "BigCode",
                        URL: "https://huggingface.co/bartowski/starcoder2-7b-GGUF/resolve/main/starcoder2-7b-Q4_K_M.gguf",
                        Tags: []string{"code", "completion", "programming"},
                },
                {
                        Name: "magicoder", DisplayName: "Magicoder S DS 6.7B", Description: "UIUC's Magicoder — code model trained with open-source instruction data",
                        Parameters: "6.7B", Quantization: "Q4_K_M", Size: "3.9 GB", Family: "magicoder", ContextLen: 8192, License: "apache-2.0", Author: "UIUC",
                        URL: "https://huggingface.co/bartowski/Magicoder-S-DS-6.7B-GGUF/resolve/main/Magicoder-S-DS-6.7B-Q4_K_M.gguf",
                        Tags: []string{"code", "open-source", "instruct"},
                },
                {
                        Name: "phind-codellama", DisplayName: "Phind CodeLlama 34B", Description: "Phind's CodeLlama 34B — top-performing code model for complex programming tasks",
                        Parameters: "34B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "codellama", ContextLen: 16384, License: "llama2", Author: "Phind",
                        URL: "https://huggingface.co/bartowski/phind-codellama-34b-v2-GGUF/resolve/main/phind-codellama-34b-v2-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "powerful"},
                },

                // ── MATH & SCIENCE MODELS ──
                {
                        Name: "mathstral:7b", DisplayName: "Mathstral 7B v0.1", Description: "Mistral AI's math-specialized model — designed for mathematical reasoning and proofs",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "mistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI",
                        URL: "https://huggingface.co/bartowski/Mathstral-7B-v0.1-GGUF/resolve/main/Mathstral-7B-v0.1-Q4_K_M.gguf",
                        Tags: []string{"math", "science", "reasoning", "proofs"},
                },
                {
                        Name: "abel", DisplayName: "Abel 7B", Description: "Abel 7B — math-focused model with strong performance on mathematical benchmarks",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "abel", ContextLen: 4096, License: "apache-2.0", Author: "GAIR",
                        URL: "https://huggingface.co/bartowski/Abel-7B-001-GGUF/resolve/main/Abel-7B-001-Q4_K_M.gguf",
                        Tags: []string{"math", "science", "reasoning"},
                },

                // ── EMBEDDING MODELS ──
                {
                        Name: "nomic-embed", DisplayName: "Nomic Embed Text", Description: "Nomic's embedding model — large context embedding for RAG with 8192 token context",
                        Parameters: "137M", Quantization: "Q4_K_M", Size: "0.1 GB", Family: "nomic", ContextLen: 8192, License: "apache-2.0", Author: "Nomic",
                        URL: "https://huggingface.co/bartowski/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5-Q4_K_M.gguf",
                        Tags: []string{"embedding", "rag", "long-context", "search"},
                },
                {
                        Name: "mxbai-embed", DisplayName: "MxBai Embed Large", Description: "Mixedbread's embedding model — high-quality multilingual embeddings for semantic search",
                        Parameters: "335M", Quantization: "Q4_K_M", Size: "0.3 GB", Family: "mxbai", ContextLen: 512, License: "apache-2.0", Author: "Mixedbread",
                        URL: "https://huggingface.co/bartowski/mxbai-embed-large-v1-GGUF/resolve/main/mxbai-embed-large-v1-Q4_K_M.gguf",
                        Tags: []string{"embedding", "semantic-search", "multilingual"},
                },
                {
                        Name: "all-minilm", DisplayName: "All-MiniLM L6 v2", Description: "Sentence-transformers MiniLM — lightweight embedding model for fast semantic similarity",
                        Parameters: "22M", Quantization: "F16", Size: "0.04 GB", Family: "minilm", ContextLen: 256, License: "apache-2.0", Author: "Sentence-Transformers",
                        URL: "https://huggingface.co/bartowski/all-MiniLM-L6-v2-GGUF/resolve/main/all-MiniLM-L6-v2-F16.gguf",
                        Tags: []string{"embedding", "tiny", "fast", "similarity"},
                },

                // ── SPECIALIZED DOMAIN MODELS ──
                {
                        Name: "biomistral", DisplayName: "BioMistral 7B", Description: "BioMistral — specialized model for biomedical text understanding and medical Q&A",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "biomistral", ContextLen: 32768, License: "apache-2.0", Author: "Mistral AI/Bio",
                        URL: "https://huggingface.co/bartowski/BioMistral-7B-GGUF/resolve/main/BioMistral-7B-Q4_K_M.gguf",
                        Tags: []string{"medical", "biomedical", "science", "specialized"},
                },
                {
                        Name: "clinicalbert", DisplayName: "Clinical BioMistral 7B", Description: "Clinical BioMistral — medical domain model for clinical text analysis and healthcare",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "clinical", ContextLen: 32768, License: "apache-2.0", Author: "BioMistral",
                        URL: "https://huggingface.co/bartowski/Clinical-BioMistral-GGUF/resolve/main/Clinical-BioMistral-Q4_K_M.gguf",
                        Tags: []string{"medical", "clinical", "healthcare", "specialized"},
                },
                {
                        Name: "sqlcoder", DisplayName: "SQLCoder 15B", Description: "Defog's SQLCoder — specialized model for SQL query generation and database interaction",
                        Parameters: "15B", Quantization: "Q4_K_M", Size: "8.6 GB", Family: "sqlcoder", ContextLen: 4096, License: "cc-by-sa-4.0", Author: "Defog",
                        URL: "https://huggingface.co/bartowski/sqlcoder-15b-GGUF/resolve/main/sqlcoder-15b-Q4_K_M.gguf",
                        Tags: []string{"sql", "database", "code", "specialized"},
                },
                {
                        Name: "wizard-coder", DisplayName: "WizardCoder 15B", Description: "WizardCoder 15B — code generation model trained on evol-instruct data",
                        Parameters: "15B", Quantization: "Q4_K_M", Size: "8.6 GB", Family: "wizard", ContextLen: 8192, License: "apache-2.0", Author: "WizardLM",
                        URL: "https://huggingface.co/bartowski/WizardCoder-15B-V1.0-GGUF/resolve/main/WizardCoder-15B-V1.0-Q4_K_M.gguf",
                        Tags: []string{"code", "programming", "instruct"},
                },
                {
                        Name: "wizard-math", DisplayName: "WizardMath 7B", Description: "WizardMath 7B — math-specialized model trained with evol-instruct for mathematical reasoning",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "wizard", ContextLen: 4096, License: "apache-2.0", Author: "WizardLM",
                        URL: "https://huggingface.co/bartowski/WizardMath-7B-V1.1-GGUF/resolve/main/WizardMath-7B-V1.1-Q4_K_M.gguf",
                        Tags: []string{"math", "reasoning", "science", "specialized"},
                },

                // ── JAPANESE / ASIAN LANGUAGE MODELS ──
                {
                        Name: "japanese-llama", DisplayName: "Japanese Llama 3 8B", Description: "Japanese fine-tuned Llama 3 — optimized for Japanese language understanding and generation",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.7 GB", Family: "llama", ContextLen: 8192, License: "llama3", Author: "Community",
                        URL: "https://huggingface.co/bartowski/llama-3-youko-8b-GGUF/resolve/main/llama-3-youko-8b-Q4_K_M.gguf",
                        Tags: []string{"japanese", "asian", "bilingual"},
                },
                {
                        Name: "plamo2:7b", DisplayName: "PLaMo 2 7B", Description: "Preferred Networks' PLaMo 2 7B — Japanese-focused compact model for efficient inference",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "plamo", ContextLen: 4096, License: "apache-2.0", Author: "Preferred Networks",
                        URL: "https://huggingface.co/bartowski/plamo-13b-GGUF/resolve/main/plamo-13b-Q4_K_M.gguf",
                        Tags: []string{"japanese", "asian", "efficient"},
                },

                // ── MAMBA / ALTERNATIVE ARCHITECTURE ──
                {
                        Name: "mamba:1.4b", DisplayName: "Mamba 1.4B", Description: "State Space Model 1.4B — novel SSM architecture for efficient sequential modeling",
                        Parameters: "1.4B", Quantization: "Q4_K_M", Size: "0.8 GB", Family: "mamba", ContextLen: 2048, License: "apache-2.0", Author: "Tri Dao",
                        URL: "https://huggingface.co/bartowski/mamba-1.4b-GGUF/resolve/main/mamba-1.4b-Q4_K_M.gguf",
                        Tags: []string{"ssm", "experimental", "efficient", "novel"},
                },

                // ── ADDITIONAL POPULAR MODELS ──
                {
                        Name: "solar:pro", DisplayName: "SOLAR 10.7B Pro", Description: "Upstage's SOLAR Pro — advanced depth-up-scaled model for professional use",
                        Parameters: "10.7B", Quantization: "Q4_K_M", Size: "6.1 GB", Family: "solar", ContextLen: 4096, License: "apache-2.0", Author: "Upstage",
                        URL: "https://huggingface.co/bartowski/SOLAR-10.7B-Instruct-v1.0-GGUF/resolve/main/SOLAR-10.7B-Instruct-v1.0-Q4_K_M.gguf",
                        Tags: []string{"instruct", "professional", "korean"},
                },
                {
                        Name: "falcon:11b", DisplayName: "Falcon 11B", Description: "TII's Falcon 11B — improved architecture for instruction following and general tasks",
                        Parameters: "11B", Quantization: "Q4_K_M", Size: "6.4 GB", Family: "falcon", ContextLen: 8192, License: "apache-2.0", Author: "TII",
                        URL: "https://huggingface.co/bartowski/Falcon-11B-GGUF/resolve/main/Falcon-11B-Q4_K_M.gguf",
                        Tags: []string{"instruct", "general", "efficient"},
                },
                {
                        Name: "falcon:180b", DisplayName: "Falcon 180B Chat", Description: "TII's Falcon 180B — one of the largest open models with exceptional quality",
                        Parameters: "180B", Quantization: "Q4_K_M", Size: "103 GB", Family: "falcon", ContextLen: 2048, License: "falcon-180b", Author: "TII",
                        URL: "https://huggingface.co/bartowski/Falcon-180B-Chat-GGUF/resolve/main/Falcon-180B-Chat-Q4_K_M.gguf",
                        Tags: []string{"massive", "powerful", "frontier", "chat"},
                },
                {
                        Name: "persimmon:8b", DisplayName: "Persimmon 8B Chat", Description: "Adept's Persimmon 8B — efficient model with unique architecture for dialogue",
                        Parameters: "8B", Quantization: "Q4_K_M", Size: "4.6 GB", Family: "persimmon", ContextLen: 16384, License: "apache-2.0", Author: "Adept",
                        URL: "https://huggingface.co/bartowski/persimmon-8b-chat-GGUF/resolve/main/persimmon-8b-chat-Q4_K_M.gguf",
                        Tags: []string{"dialogue", "efficient", "long-context"},
                },
                {
                        Name: "xwin", DisplayName: "XWin 7B", Description: "XWin-LM 7B — strong coding and reasoning model optimized for benchmarks",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "xwin", ContextLen: 8192, License: "llama2", Author: "XWin-LM",
                        URL: "https://huggingface.co/bartowski/XWin-LM-7B-V0.1-GGUF/resolve/main/XWin-LM-7B-V0.1-Q4_K_M.gguf",
                        Tags: []string{"code", "reasoning", "benchmark"},
                },
                {
                        Name: "toppy", DisplayName: "Toppy M 7B", Description: "Toppy M 7B — fine-tuned model for creative writing and diverse text generation",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "toppy", ContextLen: 4096, License: "apache-2.0", Author: "Undi95",
                        URL: "https://huggingface.co/bartowski/Toppy-M-7B-GGUF/resolve/main/Toppy-M-7B-Q4_K_M.gguf",
                        Tags: []string{"creative", "writing", "diverse"},
                },
                {
                        Name: "neural-chat", DisplayName: "Neural Chat 7B v3.1", Description: "Intel's Neural Chat — fine-tuned Mistral optimized for instruction following",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "neural-chat", ContextLen: 32768, License: "apache-2.0", Author: "Intel",
                        URL: "https://huggingface.co/bartowski/neural-chat-7b-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1-Q4_K_M.gguf",
                        Tags: []string{"instruct", "intel", "optimized"},
                },
                {
                        Name: "steamboat", DisplayName: "Steamboat 7B", Description: "Steamboat 7B — fine-tuned model for general conversation and helpful assistant behavior",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.0 GB", Family: "steamboat", ContextLen: 8192, License: "apache-2.0", Author: "Community",
                        URL: "https://huggingface.co/bartowski/Steamboat-7B-GGUF/resolve/main/Steamboat-7B-Q4_K_M.gguf",
                        Tags: []string{"assistant", "chat", "helpful"},
                },

                // ── MORE EMBEDDING & RAG MODELS ──
                {
                        Name: "e5-mistral", DisplayName: "E5-Mistral 7B Embed", Description: "Microsoft's E5-Mistral — high-quality text embeddings using Mistral architecture for RAG",
                        Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "e5", ContextLen: 32768, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/e5-mistral-7b-instruct-GGUF/resolve/main/e5-mistral-7b-instruct-Q4_K_M.gguf",
                        Tags: []string{"embedding", "rag", "search", "microsoft"},
                },

                // ── ADDITIONAL VISION MODELS ──
                {
                        Name: "moondream2", DisplayName: "Moondream2 1.8B", Description: "Moondream2 — tiny vision-language model for efficient image captioning and VQA",
                        Parameters: "1.8B", Quantization: "Q4_K_M", Size: "1.1 GB", Family: "moondream", ContextLen: 2048, License: "mit", Author: "Moondream",
                        URL: "https://huggingface.co/bartowski/moondream2-GGUF/resolve/main/moondream2-Q4_K_M.gguf",
                        Tags: []string{"vision", "tiny", "captioning", "vqa"},
                },
                {
                        Name: "kosmos2", DisplayName: "Kosmos-2 1.6B", Description: "Microsoft's Kosmos-2 — multimodal model for grounded visual understanding and detection",
                        Parameters: "1.6B", Quantization: "Q4_K_M", Size: "1.0 GB", Family: "kosmos", ContextLen: 2048, License: "mit", Author: "Microsoft",
                        URL: "https://huggingface.co/bartowski/kosmos-2-patch14-224-GGUF/resolve/main/kosmos-2-patch14-224-Q4_K_M.gguf",
                        Tags: []string{"vision", "detection", "multimodal", "grounding"},
                },

                // ── SAFETY & MODERATION ──
                {
                        Name: "shieldgemma", DisplayName: "ShieldGemma 2B", Description: "Google's ShieldGemma — content safety classifier for detecting harmful content",
                        Parameters: "2B", Quantization: "Q4_K_M", Size: "1.4 GB", Family: "gemma", ContextLen: 8192, License: "gemma", Author: "Google",
                        URL: "https://huggingface.co/bartowski/shieldgemma-2b-GGUF/resolve/main/shieldgemma-2b-Q4_K_M.gguf",
                        Tags: []string{"safety", "moderation", "classifier", "content-filter"},
                },
                {
                        Name: "llama-guard3:11b", DisplayName: "Llama Guard 3 11B Vision", Description: "Meta's Llama Guard 3 11B — vision-enabled content safety and moderation model",
                        Parameters: "11B", Quantization: "Q4_K_M", Size: "6.4 GB", Family: "llama", ContextLen: 131072, License: "llama3.1", Author: "Meta",
                        URL: "https://huggingface.co/bartowski/Llama-Guard-3-11B-Vision-GGUF/resolve/main/Llama-Guard-3-11B-Vision-Q4_K_M.gguf",
                        Tags: []string{"safety", "vision", "moderation", "content-filter"},
                },
        }

        for _, m := range models {
                r.models[m.Name] = m
        }
}

// GetModel returns model info by name
func (r *Registry) GetModel(name string) (ModelInfo, bool) {
        if m, ok := r.models[name]; ok {
                return m, true
        }
        name = strings.ToLower(name)
        for k, v := range r.models {
                if strings.HasPrefix(k, name) {
                        return v, true
                }
        }
        return ModelInfo{}, false
}

// ListModels returns all available models
func (r *Registry) ListModels() []ModelInfo {
        result := make([]ModelInfo, 0, len(r.models))
        for _, m := range r.models {
                result = append(result, m)
        }
        return result
}

// ListModelsByFamily returns models filtered by family
func (r *Registry) ListModelsByFamily(family string) []ModelInfo {
        var result []ModelInfo
        for _, m := range r.models {
                if m.Family == family {
                        result = append(result, m)
                }
        }
        return result
}

// ListModelsByTag returns models filtered by tag
func (r *Registry) ListModelsByTag(tag string) []ModelInfo {
        var result []ModelInfo
        tag = strings.ToLower(tag)
        for _, m := range r.models {
                for _, t := range m.Tags {
                        if strings.ToLower(t) == tag {
                                result = append(result, m)
                                break
                        }
                }
        }
        return result
}

// SearchModels searches models by name, description, or tags
func (r *Registry) SearchModels(query string) []ModelInfo {
        query = strings.ToLower(query)
        var result []ModelInfo
        seen := make(map[string]bool)
        for _, m := range r.models {
                if seen[m.Name] {
                        continue
                }
                if strings.Contains(strings.ToLower(m.Name), query) ||
                        strings.Contains(strings.ToLower(m.DisplayName), query) ||
                        strings.Contains(strings.ToLower(m.Description), query) ||
                        strings.Contains(strings.ToLower(m.Family), query) ||
                        strings.Contains(strings.ToLower(m.Author), query) {
                        result = append(result, m)
                        seen[m.Name] = true
                        continue
                }
                for _, tag := range m.Tags {
                        if strings.Contains(strings.ToLower(tag), query) {
                                result = append(result, m)
                                seen[m.Name] = true
                                break
                        }
                }
        }
        return result
}

// FetchRemoteModels fetches additional model info from the HuggingFace API
func (r *Registry) FetchRemoteModels(query string) ([]ModelInfo, error) {
        url := fmt.Sprintf("https://huggingface.co/api/models?search=%s+gguf&sort=downloads&direction=-1&limit=20", query)
        resp, err := r.client.Get(url)
        if err != nil {
                return nil, fmt.Errorf("failed to fetch models: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return nil, fmt.Errorf("HuggingFace API returned status %d", resp.StatusCode)
        }

        var results []struct {
                ID        string   `json:"id"`
                Downloads int      `json:"downloads"`
                Tags      []string `json:"tags"`
        }

        if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
                return nil, fmt.Errorf("failed to decode response: %w", err)
        }

        var models []ModelInfo
        for _, res := range results {
                models = append(models, ModelInfo{
                        Name:        res.ID,
                        DisplayName: res.ID,
                        Family:      "remote",
                        URL:         fmt.Sprintf("https://huggingface.co/%s", res.ID),
                        Tags:        res.Tags,
                })
        }

        return models, nil
}

// GetModelFamilies returns unique model families
func (r *Registry) GetModelFamilies() []string {
        seen := make(map[string]bool)
        var families []string
        for _, m := range r.models {
                if !seen[m.Family] {
                        seen[m.Family] = true
                        families = append(families, m.Family)
                }
        }
        return families
}

// GetModelTags returns all unique tags
func (r *Registry) GetModelTags() []string {
        seen := make(map[string]bool)
        var tags []string
        for _, m := range r.models {
                for _, t := range m.Tags {
                        if !seen[t] {
                                seen[t] = true
                                tags = append(tags, t)
                        }
                }
        }
        return tags
}

// GetModelAuthors returns all unique authors
func (r *Registry) GetModelAuthors() []string {
        seen := make(map[string]bool)
        var authors []string
        for _, m := range r.models {
                if !seen[m.Author] && m.Author != "" {
                        seen[m.Author] = true
                        authors = append(authors, m.Author)
                }
        }
        return authors
}

// ModelCount returns the total number of models
func (r *Registry) ModelCount() int {
        return len(r.models)
}
