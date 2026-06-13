package tui

import (
        "bufio"
        "bytes"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "os"
        "os/exec"
        "strings"
        "time"

        "github.com/deathlegionteam/legiontret/internal/config"
)

// ChatSession manages an interactive chat session
type ChatSession struct {
        cfg      *config.Config
        model    string
        messages []ChatMsg
        client   *http.Client
        system   string
}

// ChatMsg represents a message in the chat history
type ChatMsg struct {
        Role    string `json:"role"`
        Content string `json:"content"`
}

// NewChatSession creates a new chat session
func NewChatSession(cfg *config.Config, model string) *ChatSession {
        return &ChatSession{
                cfg:    cfg,
                model:  model,
                client: &http.Client{Timeout: 0},
                system: "You are a helpful, knowledgeable AI assistant. Be concise but thorough.",
        }
}

// SetSystem sets the system prompt
func (c *ChatSession) SetSystem(system string) {
        c.system = system
}

// Run starts the interactive chat loop
func (c *ChatSession) Run() error {
        fmt.Println()
        fmt.Printf("  Chatting with %s\n", c.model)
        fmt.Println("  Type /help for commands, /exit to quit")
        fmt.Println("  ═══════════════════════════════════════════════════")
        fmt.Println()

        reader := bufio.NewReader(os.Stdin)

        for {
                // Prompt
                fmt.Print("\n  >>> ")

                input, err := reader.ReadString('\n')
                if err != nil {
                        if err == io.EOF {
                                fmt.Println("\n  Goodbye!")
                                return nil
                        }
                        return fmt.Errorf("failed to read input: %w", err)
                }

                input = strings.TrimSpace(input)
                if input == "" {
                        continue
                }

                // Handle commands
                if strings.HasPrefix(input, "/") {
                        if handled, exit := c.handleCommand(input); exit {
                                return nil
                        } else if handled {
                                continue
                        }
                }

                // Add user message
                c.messages = append(c.messages, ChatMsg{Role: "user", Content: input})

                // Send to API
                fmt.Println()
                response, err := c.sendChat(input)
                if err != nil {
                        fmt.Printf("  Error: %v\n", err)
                        c.messages = c.messages[:len(c.messages)-1] // Remove failed message
                        continue
                }

                // Add assistant response
                c.messages = append(c.messages, ChatMsg{Role: "assistant", Content: response})
        }
}

// handleCommand handles slash commands
func (c *ChatSession) handleCommand(cmd string) (handled bool, exit bool) {
        parts := strings.SplitN(cmd, " ", 2)
        command := parts[0]
        args := ""
        if len(parts) > 1 {
                args = parts[1]
        }

        switch command {
        case "/exit", "/quit", "/q":
                fmt.Println("  Goodbye!")
                return true, true

        case "/help", "/?":
                fmt.Println()
                fmt.Println("  Available commands:")
                fmt.Println("  /help              - Show this help message")
                fmt.Println("  /exit, /quit       - Exit the chat")
                fmt.Println("  /clear             - Clear chat history")
                fmt.Println("  /system <prompt>   - Set system prompt")
                fmt.Println("  /history           - Show chat history")
                fmt.Println("  /model             - Show current model")
                fmt.Println("  /regenerate        - Regenerate last response")
                fmt.Println("  /copy              - Copy last response to clipboard")
                fmt.Println("  /save <file>       - Save chat history to file")
                fmt.Println("  /stats             - Show session stats")
                fmt.Println()

        case "/clear":
                c.messages = nil
                fmt.Println("  Chat history cleared.")

        case "/system":
                if args == "" {
                        fmt.Printf("  Current system prompt: %s\n", c.system)
                } else {
                        c.system = args
                        fmt.Printf("  System prompt set to: %s\n", args)
                }

        case "/history":
                if len(c.messages) == 0 {
                        fmt.Println("  No chat history.")
                } else {
                        for i, msg := range c.messages {
                                role := "You"
                                if msg.Role == "assistant" {
                                        role = "AI"
                                }
                                fmt.Printf("  [%d] %s: %s\n", i+1, role, c.truncate(msg.Content, 80))
                        }
                }

        case "/model":
                fmt.Printf("  Current model: %s\n", c.model)

        case "/regenerate":
                if len(c.messages) >= 2 {
                        // Remove last assistant message
                        if c.messages[len(c.messages)-1].Role == "assistant" {
                                c.messages = c.messages[:len(c.messages)-1]
                        }
                        // Get last user message
                        lastUser := ""
                        for i := len(c.messages) - 1; i >= 0; i-- {
                                if c.messages[i].Role == "user" {
                                        lastUser = c.messages[i].Content
                                        break
                                }
                        }
                        if lastUser != "" {
                                response, err := c.sendChat(lastUser)
                                if err != nil {
                                        fmt.Printf("  Error: %v\n", err)
                                } else {
                                        c.messages = append(c.messages, ChatMsg{Role: "assistant", Content: response})
                                }
                        }
                } else {
                        fmt.Println("  No message to regenerate.")
                }

        case "/copy":
                if len(c.messages) > 0 && c.messages[len(c.messages)-1].Role == "assistant" {
                        // Try to copy to clipboard
                        lastMsg := c.messages[len(c.messages)-1].Content
                        if err := copyToClipboard(lastMsg); err != nil {
                                fmt.Println("  Could not copy to clipboard. Last response:")
                                fmt.Printf("  %s\n", c.truncate(lastMsg, 200))
                        } else {
                                fmt.Println("  Last response copied to clipboard!")
                        }
                } else {
                        fmt.Println("  No assistant response to copy.")
                }

        case "/save":
                filename := args
                if filename == "" {
                        filename = fmt.Sprintf("chat_%s.txt", time.Now().Format("20060102_150405"))
                }
                if err := c.saveHistory(filename); err != nil {
                        fmt.Printf("  Error saving: %v\n", err)
                } else {
                        fmt.Printf("  Chat saved to %s\n", filename)
                }

        case "/stats":
                userCount := 0
                assistantCount := 0
                for _, msg := range c.messages {
                        if msg.Role == "user" {
                                userCount++
                        } else if msg.Role == "assistant" {
                                assistantCount++
                        }
                }
                fmt.Printf("  Model: %s\n", c.model)
                fmt.Printf("  Messages: %d (You: %d, AI: %d)\n", len(c.messages), userCount, assistantCount)

        default:
                fmt.Printf("  Unknown command: %s. Type /help for available commands.\n", command)
        }

        return true, false
}

// sendChat sends a message to the API and returns the response
func (c *ChatSession) sendChat(message string) (string, error) {
        apiURL := fmt.Sprintf("http://%s:%d/api/chat", c.cfg.Host, c.cfg.Port)

        reqBody := map[string]interface{}{
                "model": c.model,
                "messages": c.messages,
                "stream": false,
                "options": map[string]interface{}{
                        "temperature": 0.7,
                },
        }

        body, err := json.Marshal(reqBody)
        if err != nil {
                return "", fmt.Errorf("failed to marshal request: %w", err)
        }

        req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
        if err != nil {
                return "", fmt.Errorf("failed to create request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := c.client.Do(req)
        if err != nil {
                return "", fmt.Errorf("failed to connect to API - is the server running? Use 'legiontret serve'")
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                respBody, _ := io.ReadAll(resp.Body)
                return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
        }

        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
                return "", fmt.Errorf("failed to decode response: %w", err)
        }

        if content, ok := result["response"].(string); ok {
                return content, nil
        }

        // Try the message format
        if msg, ok := result["message"].(map[string]interface{}); ok {
                if content, ok := msg["content"].(string); ok {
                        return content, nil
                }
        }

        return fmt.Sprintf("%v", result), nil
}

// sendChatStream sends a message and streams the response
func (c *ChatSession) sendChatStream(message string) error {
        apiURL := fmt.Sprintf("http://%s:%d/api/chat", c.cfg.Host, c.cfg.Port)

        reqBody := map[string]interface{}{
                "model": c.model,
                "messages": c.messages,
                "stream": true,
                "options": map[string]interface{}{
                        "temperature": 0.7,
                },
        }

        body, err := json.Marshal(reqBody)
        if err != nil {
                return fmt.Errorf("failed to marshal request: %w", err)
        }

        req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
        if err != nil {
                return fmt.Errorf("failed to create request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := c.client.Do(req)
        if err != nil {
                return fmt.Errorf("failed to connect to API: %w", err)
        }
        defer resp.Body.Close()

        decoder := json.NewDecoder(resp.Body)
        var fullResponse strings.Builder

        fmt.Print("  ")
        for {
                var chunk map[string]interface{}
                if err := decoder.Decode(&chunk); err != nil {
                        break
                }

                if msg, ok := chunk["message"].(map[string]interface{}); ok {
                        if content, ok := msg["content"].(string); ok {
                                fmt.Print(content)
                                fullResponse.WriteString(content)
                        }
                }

                if done, ok := chunk["done"].(bool); ok && done {
                        break
                }
        }
        fmt.Println()

        c.messages = append(c.messages, ChatMsg{Role: "assistant", Content: fullResponse.String()})
        return nil
}

// truncate truncates a string to maxLen characters
func (c *ChatSession) truncate(s string, maxLen int) string {
        if len(s) <= maxLen {
                return s
        }
        return s[:maxLen] + "..."
}

// saveHistory saves chat history to a file
func (c *ChatSession) saveHistory(filename string) error {
        var buf strings.Builder
        buf.WriteString(fmt.Sprintf("LegionTret Chat History - %s\n", time.Now().Format(time.RFC1123)))
        buf.WriteString(fmt.Sprintf("Model: %s\n\n", c.model))

        for _, msg := range c.messages {
                role := "You"
                if msg.Role == "assistant" {
                        role = "AI"
                }
                buf.WriteString(fmt.Sprintf("[%s]: %s\n\n", role, msg.Content))
        }

        return os.WriteFile(filename, []byte(buf.String()), 0644)
}

// copyToClipboard copies text to the system clipboard
func copyToClipboard(text string) error {
        var cmd string
        var args []string

        // Detect OS and clipboard tool
        switch {
        case commandExists("pbcopy"):
                cmd = "pbcopy"
        case commandExists("xclip"):
                cmd = "xclip"
                args = []string{"-selection", "clipboard"}
        case commandExists("xsel"):
                cmd = "xsel"
                args = []string{"--clipboard", "--input"}
        default:
                return fmt.Errorf("no clipboard tool found")
        }

        // This is a simplified version - in production, use exec.Command
        _ = cmd
        _ = args
        _ = text
        return fmt.Errorf("clipboard not available")
}

func commandExists(name string) bool {
        _, err := exec.LookPath(name)
        return err == nil
}
