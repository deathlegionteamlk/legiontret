package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ModelInfo represents information about a model in the registry
type ModelInfo struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Description string            `json:"description"`
	Parameters  string            `json:"parameters"` // e.g., "7B", "13B", "70B"
	Quantization string           `json:"quantization"` // e.g., "Q4_0", "Q5_1"
	Size        string            `json:"size"`        // e.g., "4.1 GB"
	Family      string            `json:"family"`      // e.g., "llama", "gemma", "mistral"
	URL         string            `json:"url"`         // Download URL
	SHA256      string            `json:"sha256"`
	Tags        []string          `json:"tags"`
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

// loadBuiltinModels loads the built-in model catalog
func (r *Registry) loadBuiltinModels() {
	models := []ModelInfo{
		// Llama 3 family
		{
			Name: "llama3", DisplayName: "Llama 3 8B", Description: "Meta's Llama 3 8B instruct model - excellent general-purpose model",
			Parameters: "8B", Quantization: "Q4_0", Size: "4.7 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Meta-Llama-3-8B-Instruct-GGUF/resolve/main/Meta-Llama-3-8B-Instruct-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "general"},
		},
		{
			Name: "llama3:70b", DisplayName: "Llama 3 70B", Description: "Meta's Llama 3 70B instruct model - top-tier reasoning",
			Parameters: "70B", Quantization: "Q4_0", Size: "40.5 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Meta-Llama-3-70B-Instruct-GGUF/resolve/main/Meta-Llama-3-70B-Instruct-Q4_K_M.gguf",
			Tags: []string{"powerful", "instruct", "general"},
		},
		{
			Name: "llama3.1", DisplayName: "Llama 3.1 8B", Description: "Meta's Llama 3.1 8B instruct - improved over Llama 3",
			Parameters: "8B", Quantization: "Q4_K_M", Size: "4.9 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "general", "latest"},
		},
		{
			Name: "llama3.1:70b", DisplayName: "Llama 3.1 70B", Description: "Meta's Llama 3.1 70B instruct - state of the art",
			Parameters: "70B", Quantization: "Q4_K_M", Size: "42.0 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Meta-Llama-3.1-70B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-70B-Instruct-Q4_K_M.gguf",
			Tags: []string{"powerful", "instruct", "general", "latest"},
		},
		{
			Name: "llama3.2", DisplayName: "Llama 3.2 3B", Description: "Meta's Llama 3.2 3B instruct - lightweight and fast",
			Parameters: "3B", Quantization: "Q4_K_M", Size: "2.0 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf",
			Tags: []string{"lightweight", "fast", "instruct"},
		},
		{
			Name: "llama3.2:1b", DisplayName: "Llama 3.2 1B", Description: "Meta's Llama 3.2 1B instruct - ultra lightweight",
			Parameters: "1B", Quantization: "Q4_K_M", Size: "0.8 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q4_K_M.gguf",
			Tags: []string{"ultra-light", "fast", "instruct"},
		},

		// Gemma family
		{
			Name: "gemma3", DisplayName: "Gemma 3 4B", Description: "Google's Gemma 3 4B instruct - efficient and capable",
			Parameters: "4B", Quantization: "Q4_K_M", Size: "2.6 GB", Family: "gemma",
			URL: "https://huggingface.co/bartowski/gemma-3-4b-it-GGUF/resolve/main/gemma-3-4b-it-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "google"},
		},
		{
			Name: "gemma3:12b", DisplayName: "Gemma 3 12B", Description: "Google's Gemma 3 12B instruct - strong performance",
			Parameters: "12B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "gemma",
			URL: "https://huggingface.co/bartowski/gemma-3-12b-it-GGUF/resolve/main/gemma-3-12b-it-Q4_K_M.gguf",
			Tags: []string{"powerful", "instruct", "google"},
		},
		{
			Name: "gemma3:27b", DisplayName: "Gemma 3 27B", Description: "Google's Gemma 3 27B instruct - top Google model",
			Parameters: "27B", Quantization: "Q4_K_M", Size: "16.2 GB", Family: "gemma",
			URL: "https://huggingface.co/bartowski/gemma-3-27b-it-GGUF/resolve/main/gemma-3-27b-it-Q4_K_M.gguf",
			Tags: []string{"powerful", "instruct", "google"},
		},
		{
			Name: "gemma3:1b", DisplayName: "Gemma 3 1B", Description: "Google's Gemma 3 1B - ultra compact",
			Parameters: "1B", Quantization: "Q4_K_M", Size: "0.7 GB", Family: "gemma",
			URL: "https://huggingface.co/bartowski/gemma-3-1b-it-GGUF/resolve/main/gemma-3-1b-it-Q4_K_M.gguf",
			Tags: []string{"ultra-light", "fast", "google"},
		},
		{
			Name: "gemma2", DisplayName: "Gemma 2 9B", Description: "Google's Gemma 2 9B instruct - proven performer",
			Parameters: "9B", Quantization: "Q4_K_M", Size: "5.4 GB", Family: "gemma",
			URL: "https://huggingface.co/bartowski/gemma-2-9b-it-GGUF/resolve/main/gemma-2-9b-it-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "google"},
		},

		// Mistral family
		{
			Name: "mistral", DisplayName: "Mistral 7B", Description: "Mistral AI's 7B instruct model - fast and efficient",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/Mistral-7B-Instruct-v0.3-GGUF/resolve/main/Mistral-7B-Instruct-v0.3-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "fast"},
		},
		{
			Name: "mistral-nemo", DisplayName: "Mistral Nemo 12B", Description: "Mistral Nemo 12B instruct - excellent context window",
			Parameters: "12B", Quantization: "Q4_K_M", Size: "7.4 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/Mistral-Nemo-Instruct-2407-GGUF/resolve/main/Mistral-Nemo-Instruct-2407-Q4_K_M.gguf",
			Tags: []string{"popular", "instruct", "long-context"},
		},
		{
			Name: "mistral-small", DisplayName: "Mistral Small 24B", Description: "Mistral Small 24B instruct - strong reasoning",
			Parameters: "24B", Quantization: "Q4_K_M", Size: "14.1 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/Mistral-Small-24B-Instruct-2501-GGUF/resolve/main/Mistral-Small-24B-Instruct-2501-Q4_K_M.gguf",
			Tags: []string{"powerful", "instruct"},
		},
		{
			Name: "mixtral", DisplayName: "Mixtral 8x7B", Description: "Mixtral 8x7B instruct - MoE architecture for speed",
			Parameters: "46.7B", Quantization: "Q4_K_M", Size: "26.0 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/Mixtral-8x7B-Instruct-v0.1-GGUF/resolve/main/Mixtral-8x7B-Instruct-v0.1-Q4_K_M.gguf",
			Tags: []string{"moe", "powerful", "instruct"},
		},

		// Qwen family
		{
			Name: "qwen2.5", DisplayName: "Qwen 2.5 7B", Description: "Alibaba's Qwen 2.5 7B instruct - multilingual champion",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen",
			URL: "https://huggingface.co/bartowski/Qwen2.5-7B-Instruct-GGUF/resolve/main/Qwen2.5-7B-Instruct-Q4_K_M.gguf",
			Tags: []string{"multilingual", "instruct", "popular"},
		},
		{
			Name: "qwen2.5:14b", DisplayName: "Qwen 2.5 14B", Description: "Qwen 2.5 14B instruct - strong multilingual model",
			Parameters: "14B", Quantization: "Q4_K_M", Size: "8.7 GB", Family: "qwen",
			URL: "https://huggingface.co/bartowski/Qwen2.5-14B-Instruct-GGUF/resolve/main/Qwen2.5-14B-Instruct-Q4_K_M.gguf",
			Tags: []string{"multilingual", "powerful", "instruct"},
		},
		{
			Name: "qwen2.5:32b", DisplayName: "Qwen 2.5 32B", Description: "Qwen 2.5 32B instruct - top-tier multilingual",
			Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "qwen",
			URL: "https://huggingface.co/bartowski/Qwen2.5-32B-Instruct-GGUF/resolve/main/Qwen2.5-32B-Instruct-Q4_K_M.gguf",
			Tags: []string{"multilingual", "powerful", "instruct"},
		},
		{
			Name: "qwen2.5-coder:7b", DisplayName: "Qwen 2.5 Coder 7B", Description: "Qwen 2.5 Coder - specialized for programming",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "qwen",
			URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-7B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-7B-Instruct-Q4_K_M.gguf",
			Tags: []string{"code", "programming", "popular"},
		},
		{
			Name: "qwen2.5-coder:32b", DisplayName: "Qwen 2.5 Coder 32B", Description: "Qwen 2.5 Coder 32B - elite programming model",
			Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "qwen",
			URL: "https://huggingface.co/bartowski/Qwen2.5-Coder-32B-Instruct-GGUF/resolve/main/Qwen2.5-Coder-32B-Instruct-Q4_K_M.gguf",
			Tags: []string{"code", "programming", "powerful"},
		},

		// DeepSeek family
		{
			Name: "deepseek-r1", DisplayName: "DeepSeek R1 8B", Description: "DeepSeek R1 8B - powerful reasoning model",
			Parameters: "8B", Quantization: "Q4_K_M", Size: "4.9 GB", Family: "deepseek",
			URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Llama-8B-GGUF/resolve/main/DeepSeek-R1-Distill-Llama-8B-Q4_K_M.gguf",
			Tags: []string{"reasoning", "popular", "latest"},
		},
		{
			Name: "deepseek-r1:14b", DisplayName: "DeepSeek R1 14B", Description: "DeepSeek R1 14B distill - enhanced reasoning",
			Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "deepseek",
			URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-14B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-14B-Q4_K_M.gguf",
			Tags: []string{"reasoning", "powerful"},
		},
		{
			Name: "deepseek-r1:32b", DisplayName: "DeepSeek R1 32B", Description: "DeepSeek R1 32B distill - advanced reasoning",
			Parameters: "32B", Quantization: "Q4_K_M", Size: "19.5 GB", Family: "deepseek",
			URL: "https://huggingface.co/bartowski/DeepSeek-R1-Distill-Qwen-32B-GGUF/resolve/main/DeepSeek-R1-Distill-Qwen-32B-Q4_K_M.gguf",
			Tags: []string{"reasoning", "powerful"},
		},

		// Phi family
		{
			Name: "phi3", DisplayName: "Phi-3 Mini 3.8B", Description: "Microsoft's Phi-3 Mini - small but mighty",
			Parameters: "3.8B", Quantization: "Q4_K_M", Size: "2.3 GB", Family: "phi",
			URL: "https://huggingface.co/bartowski/Phi-3-mini-4k-instruct-GGUF/resolve/main/Phi-3-mini-4k-instruct-Q4_K_M.gguf",
			Tags: []string{"lightweight", "microsoft", "fast"},
		},
		{
			Name: "phi3:14b", DisplayName: "Phi-3 Medium 14B", Description: "Microsoft's Phi-3 Medium - compact powerhouse",
			Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "phi",
			URL: "https://huggingface.co/bartowski/Phi-3-medium-4k-instruct-GGUF/resolve/main/Phi-3-medium-4k-instruct-Q4_K_M.gguf",
			Tags: []string{"microsoft", "powerful"},
		},
		{
			Name: "phi4", DisplayName: "Phi-4 14B", Description: "Microsoft's Phi-4 - latest Phi model",
			Parameters: "14B", Quantization: "Q4_K_M", Size: "8.4 GB", Family: "phi",
			URL: "https://huggingface.co/bartowski/phi-4-GGUF/resolve/main/phi-4-Q4_K_M.gguf",
			Tags: []string{"microsoft", "latest", "powerful"},
		},

		// Code-specific models
		{
			Name: "codellama", DisplayName: "Code Llama 7B", Description: "Meta's Code Llama - specialized for code generation",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/CodeLlama-7b-Instruct-hf-GGUF/resolve/main/CodeLlama-7b-Instruct-hf-Q4_K_M.gguf",
			Tags: []string{"code", "programming"},
		},
		{
			Name: "starcoder2", DisplayName: "StarCoder2 3B", Description: "BigCode's StarCoder2 - open code model",
			Parameters: "3B", Quantization: "Q4_K_M", Size: "1.8 GB", Family: "starcoder",
			URL: "https://huggingface.co/bartowski/starcoder2-3b-GGUF/resolve/main/starcoder2-3b-Q4_K_M.gguf",
			Tags: []string{"code", "lightweight"},
		},

		// Specialized models
		{
			Name: "wizardlm2", DisplayName: "WizardLM 2 7B", Description: "Microsoft's WizardLM 2 - great for creative tasks",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.1 GB", Family: "wizard",
			URL: "https://huggingface.co/bartowski/WizardLM-2-7B-GGUF/resolve/main/WizardLM-2-7B-Q4_K_M.gguf",
			Tags: []string{"creative", "instruct"},
		},
		{
			Name: "yi", DisplayName: "Yi 1.5 6B", Description: "01.AI's Yi 1.5 Chat - strong bilingual model",
			Parameters: "6B", Quantization: "Q4_K_M", Size: "3.6 GB", Family: "yi",
			URL: "https://huggingface.co/bartowski/Yi-1.5-6B-Chat-GGUF/resolve/main/Yi-1.5-6B-Chat-Q4_K_M.gguf",
			Tags: []string{"bilingual", "chinese"},
		},
		{
			Name: "solar", DisplayName: "Solar 10.7B", Description: "Upstage's Solar - efficient architecture",
			Parameters: "10.7B", Quantization: "Q4_K_M", Size: "6.2 GB", Family: "solar",
			URL: "https://huggingface.co/bartowski/SOLAR-10.7B-Instruct-v1.0-GGUF/resolve/main/SOLAR-10.7B-Instruct-v1.0-Q4_K_M.gguf",
			Tags: []string{"efficient", "instruct"},
		},
		{
			Name: "command-r", DisplayName: "Command R 35B", Description: "Cohere's Command R - RAG and tool use specialist",
			Parameters: "35B", Quantization: "Q4_K_M", Size: "20.8 GB", Family: "command",
			URL: "https://huggingface.co/bartowski/command-r-GGUF/resolve/main/command-r-Q4_K_M.gguf",
			Tags: []string{"rag", "tools", "powerful"},
		},
		{
			Name: "tinyllama", DisplayName: "TinyLlama 1.1B", Description: "TinyLlama - ultra compact for edge devices",
			Parameters: "1.1B", Quantization: "Q4_K_M", Size: "0.6 GB", Family: "llama",
			URL: "https://huggingface.co/bartowski/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/TinyLlama-1.1B-Chat-v1.0-Q4_K_M.gguf",
			Tags: []string{"ultra-light", "edge", "fast"},
		},
		{
			Name: "dolphin-mistral", DisplayName: "Dolphin Mistral 7B", Description: "Uncensored Mistral fine-tune by Eric Hartford",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/dolphin-2.8-mistral-7b-v02-GGUF/resolve/main/dolphin-2.8-mistral-7b-v02-Q4_K_M.gguf",
			Tags: []string{"uncensored", "creative"},
		},
		{
			Name: "nous-hermes2", DisplayName: "Nous Hermes 2 7B", Description: "NousResearch Hermes 2 - creative and capable",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/Nous-Hermes-2-Mistral-7B-DPO-GGUF/resolve/main/Nous-Hermes-2-Mistral-7B-DPO-Q4_K_M.gguf",
			Tags: []string{"creative", "instruct"},
		},
		{
			Name: "llava", DisplayName: "LLaVA 1.6 7B", Description: "Vision-language model - understand images and text",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.5 GB", Family: "llava",
			URL: "https://huggingface.co/bartowski/llava-v1.6-mistral-7b-GGUF/resolve/main/llava-v1.6-mistral-7b-Q4_K_M.gguf",
			Tags: []string{"vision", "multimodal"},
		},
		{
			Name: "mathstral", DisplayName: "Mathstral 7B", Description: "Mistral's math-specialized model",
			Parameters: "7B", Quantization: "Q4_K_M", Size: "4.4 GB", Family: "mistral",
			URL: "https://huggingface.co/bartowski/mathstral-7B-v0.1-GGUF/resolve/main/mathstral-7B-v0.1-Q4_K_M.gguf",
			Tags: []string{"math", "science"},
		},
	}

	for _, m := range models {
		r.models[m.Name] = m
	}
}

// GetModel returns model info by name
func (r *Registry) GetModel(name string) (ModelInfo, bool) {
	// Direct lookup
	if m, ok := r.models[name]; ok {
		return m, true
	}

	// Try with :latest suffix
	if m, ok := r.models[name+":latest"]; ok {
		return m, true
	}

	// Try fuzzy matching - find first model that starts with the name
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

// SearchModels searches models by name, description, or tags
func (r *Registry) SearchModels(query string) []ModelInfo {
	query = strings.ToLower(query)
	var result []ModelInfo
	for _, m := range r.models {
		if strings.Contains(strings.ToLower(m.Name), query) ||
			strings.Contains(strings.ToLower(m.DisplayName), query) ||
			strings.Contains(strings.ToLower(m.Description), query) ||
			strings.Contains(strings.ToLower(m.Family), query) {
			result = append(result, m)
			continue
		}
		for _, tag := range m.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				result = append(result, m)
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
		ID      string `json:"id"`
		Downloads int  `json:"downloads"`
		Tags    []string `json:"tags"`
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
