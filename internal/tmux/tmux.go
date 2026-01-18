package tmux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func CapturePane(target string) (string, error) {
	args := []string{"capture-pane", "-p", "-t", target}
	cmd := exec.Command("tmux", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("tmux capture-pane failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("tmux capture-pane failed: %w", err)
	}

	return stdout.String(), nil
}

func SendKeys(target string, keys []string) error {
	args := []string{"send-keys", "-t", target}
	args = append(args, keys...)

	cmd := exec.Command("tmux", args...)

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("tmux send-keys failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return fmt.Errorf("tmux send-keys failed: %w", err)
	}

	time.Sleep(500 * time.Millisecond)
	pane, err := CapturePane(target)
	if err != nil {
		return fmt.Errorf("failed to capture pane after sending keys: %w")
	}

	fmt.Printf("Screen after sending keys:\n%s", pane)

	return nil
}
