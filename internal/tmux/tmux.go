package tmux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func CapturePane(target string) (*Output, error) {
	args := []string{"capture-pane", "-p", "-t", target}
	cmd := exec.Command("tmux", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("tmux capture-pane failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("tmux capture-pane failed: %w", err)
	}

	output := Parse(stdout.String())
	return &output, nil
}

func SendKeys(target string, keys []string) (*Output, error) {
	if len(keys) > 1 {
		return nil, fmt.Errorf("Please send one key at a time")
	}

	args := []string{"send-keys", "-t", target}

	// Translate literal characters for tmux compatibility
	for i, key := range keys {
		if key == " " {
			keys[i] = "Space"
		}
		if key == "\n" {
			keys[i] = "Enter"
		}
	}

	args = append(args, keys...)

	cmd := exec.Command("tmux", args...)

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("tmux send-keys failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("tmux send-keys failed: %w", err)
	}

	time.Sleep(500 * time.Millisecond)
	return CapturePane(target)
}
