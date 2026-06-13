package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Downloader handles file downloads with progress tracking
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a new downloader
func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 0, // No timeout for large files
		},
	}
}

// ProgressCallback is called with download progress info
type ProgressCallback func(downloaded, total int64, speed float64, eta time.Duration)

// Download downloads a file with progress reporting
func (d *Downloader) Download(url, destPath string, progressFn ProgressCallback) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check for partial download
	partialPath := destPath + ".part"
	var offset int64
	if info, err := os.Stat(partialPath); err == nil {
		offset = info.Size()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Resume support
	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	// Handle range response
	var total int64
	var startOffset int64
	if resp.StatusCode == http.StatusPartialContent {
		total = resp.ContentLength + offset
		startOffset = offset
	} else if resp.StatusCode == http.StatusOK {
		total = resp.ContentLength
		startOffset = 0
		offset = 0
	} else {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Open file for writing
	var outFile *os.File
	if startOffset > 0 {
		outFile, err = os.OpenFile(partialPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		outFile, err = os.Create(partialPath)
	}
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Download with progress tracking
	buf := make([]byte, 32*1024) // 32KB buffer
	var downloaded int64 = startOffset
	lastUpdate := time.Now()
	var lastDownloaded int64 = startOffset

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			written, writeErr := outFile.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write data: %w", writeErr)
			}
			downloaded += int64(written)
		}

		now := time.Now()
		elapsed := now.Sub(lastUpdate)
		if elapsed >= 500*time.Millisecond && progressFn != nil {
			speed := float64(downloaded-lastDownloaded) / elapsed.Seconds()
			var eta time.Duration
			if speed > 0 {
				remaining := total - downloaded
				eta = time.Duration(float64(remaining)/speed) * time.Second
			}
			progressFn(downloaded, total, speed, eta)
			lastUpdate = now
			lastDownloaded = downloaded
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("download error: %w", err)
		}
	}

	// Final progress update
	if progressFn != nil {
		progressFn(downloaded, total, 0, 0)
	}

	outFile.Close()

	// Rename partial to final
	if err := os.Rename(partialPath, destPath); err != nil {
		return fmt.Errorf("failed to finalize download: %w", err)
	}

	return nil
}

// FormatSize formats bytes into a human-readable string
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatSpeed formats bytes per second into a human-readable string
func FormatSpeed(bytesPerSec float64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)
	switch {
	case bytesPerSec >= MB:
		return fmt.Sprintf("%.1f MB/s", bytesPerSec/MB)
	case bytesPerSec >= KB:
		return fmt.Sprintf("%.1f KB/s", bytesPerSec/KB)
	default:
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	}
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

// ValidateURL checks if a URL is reachable
func (d *Downloader) ValidateURL(url string) (int64, error) {
	resp, err := d.client.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to check URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return resp.ContentLength, nil
}

// IsURL checks if a string is a URL
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
