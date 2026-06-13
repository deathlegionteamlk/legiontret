package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deathlegionteam/legiontret/internal/download"
)

// Bar renders a progress bar in the terminal
type Bar struct {
	mu       sync.Mutex
	total    int64
	width    int
	desc     string
	started  time.Time
	lastLine string
}

// NewBar creates a new progress bar
func NewBar(desc string, total int64) *Bar {
	return &Bar{
		total:   total,
		width:   40,
		desc:    desc,
		started: time.Now(),
	}
}

// Update updates the progress bar
func (b *Bar) Update(downloaded, total int64, speed float64, eta time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if total > 0 {
		b.total = total
	}

	pct := float64(0)
	if b.total > 0 {
		pct = float64(downloaded) / float64(b.total)
	}

	// Build progress bar
	filled := int(pct * float64(b.width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", b.width-filled)

	line := fmt.Sprintf("\r%s [%s] %s/%s %.1f%% %s ETA: %s",
		b.desc,
		bar,
		download.FormatSize(downloaded),
		download.FormatSize(b.total),
		pct*100,
		download.FormatSpeed(speed),
		download.FormatDuration(eta),
	)

	b.lastLine = line
	fmt.Print(line)
}

// Finish completes the progress bar
func (b *Bar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()

	elapsed := time.Since(b.started)
	line := fmt.Sprintf("\r%s [%s] %s 100%% completed in %s    \n",
		b.desc,
		strings.Repeat("█", b.width),
		download.FormatSize(b.total),
		download.FormatDuration(elapsed),
	)

	fmt.Print(line)
}

// Spinner shows a simple spinner for indeterminate progress
type Spinner struct {
	mu      sync.Mutex
	desc    string
	chars   []string
	current int
	stop    chan struct{}
}

// NewSpinner creates a new spinner
func NewSpinner(desc string) *Spinner {
	return &Spinner{
		desc:  desc,
		chars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:  make(chan struct{}),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				s.mu.Lock()
				char := s.chars[s.current%len(s.chars)]
				s.current++
				fmt.Printf("\r%s %s", s.desc, char)
				s.mu.Unlock()
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop(msg string) {
	s.stop <- struct{}{}
	fmt.Printf("\r%s %s    \n", s.desc, msg)
}

// UpdateDesc updates the spinner description
func (s *Spinner) UpdateDesc(desc string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.desc = desc
}

// PrintBanner prints the application banner
func PrintBanner() {
	banner := `
  ╔═══════════════════════════════════════════════════════╗
  ║                                                       ║
  ║   ██╗     ██╗███████╗██████╗  ██████╗ ███████╗       ║
  ║   ██║     ██║██╔════╝██╔══██╗██╔═══██╗██╔════╝       ║
  ║   ██║     ██║█████╗  ██████╔╝██║   ██║███████╗       ║
  ║   ██║     ██║██╔══╝  ██╔══██╗██║   ██║╚════██║       ║
  ║   ███████╗██║███████╗██████╔╝╚██████╔╝███████║       ║
  ║   ╚══════╝╚═╝╚══════╝╚═════╝  ╚═════╝ ╚══════╝       ║
  ║                                                       ║
  ║         LegionTret by Death Legion Team               ║
  ║         Run LLMs locally. Simple. Fast. Free.         ║
  ║                                                       ║
  ╚═══════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintModelTable prints a formatted table of models
func PrintModelTable(models []interface {
	GetName() string
	GetDisplayName() string
	GetParameters() string
	GetSize() string
	GetFamily() string
}, downloaded map[string]bool) {
	if len(models) == 0 {
		fmt.Println("No models found.")
		return
	}

	// Header
	fmt.Println()
	fmt.Printf("  %-20s %-25s %-10s %-10s %-12s %-10s\n", "NAME", "DISPLAY NAME", "PARAMS", "SIZE", "FAMILY", "STATUS")
	fmt.Println("  " + strings.Repeat("─", 90))

	for _, m := range models {
		status := "remote"
		if downloaded[m.GetName()] {
			status = "✓ local"
		}
		fmt.Printf("  %-20s %-25s %-10s %-10s %-12s %-10s\n",
			m.GetName(), m.GetDisplayName(), m.GetParameters(), m.GetSize(), m.GetFamily(), status)
	}
	fmt.Println()
}
