package tmux

import (
	"bytes"
	"fmt"
	"os"
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

	output := stdout.String()

	if dumpFile := os.Getenv("NETHACK_DUMP_FILENAME"); dumpFile != "" {
		if content, err := os.ReadFile(dumpFile); err == nil {
			output += "\n" + string(content)
		}
	}

	return output, nil
}

func SendKeys(target string, keys []string) (string, error) {
	if len(keys) > 5 {
		return "", fmt.Errorf("Please send at most 5 keys at a time")
	}

	for _, key := range keys {
		isValid := false

		if len(key) == 1 {
			isValid = true
		} else if strings.HasPrefix(key, "C-") && len(key) == 3 {
			isValid = true
		} else if key == "Enter" || key == "Space" || key == "Escape" {
			isValid = true
		}

		if !isValid {
			return "", fmt.Errorf("invalid key '%s': must be a single character, C-<char> (e.g. C-x), or Enter/Space/Escape", key)
		}
	}

	// Translate literal characters for tmux compatibility
	for i, key := range keys {
		if key == " " {
			keys[i] = "Space"
		}
		if key == "\n" {
			keys[i] = "Enter"
		}
	}

	for _, key := range keys {
		cmd := exec.Command("tmux", "send-keys", "-t", target, key)

		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return "", fmt.Errorf("tmux send-keys failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
			}
			return "", fmt.Errorf("tmux send-keys failed: %w", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	return CapturePane(target)
}
