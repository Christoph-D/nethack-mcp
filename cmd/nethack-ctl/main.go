package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Christoph-D/nethack-mcp/internal/tmux"
	"github.com/urfave/cli/v2"
)

func getTarget() (string, error) {
	target := os.Getenv("NETHACK_TMUX_SESSION")
	if target == "" {
		return "", fmt.Errorf("NETHACK_TMUX_SESSION environment variable not set")
	}
	return target, nil
}

func formatOutput(output *tmux.Output) (string, error) {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal output: %w", err)
	}
	var result strings.Builder
	fmt.Fprintf(&result, "%s\nParsed data:\n", output.RawScreen)
	result.WriteString(string(data))
	return result.String(), nil
}

func main() {
	app := &cli.App{
		Name:  "nethack-ctl",
		Usage: "Control NetHack running in tmux for AI agents",
		Commands: []*cli.Command{
			{
				Name:  "screen",
				Usage: "Capture and display the current NetHack screen",
				Action: func(c *cli.Context) error {
					target, err := getTarget()
					if err != nil {
						return err
					}

					output, err := tmux.CapturePane(target)
					if err != nil {
						return err
					}

					formatted, err := formatOutput(output)
					if err != nil {
						return err
					}

					fmt.Print(formatted)
					return nil
				},
			},
			{
				Name:  "send",
				Usage: "Send keystrokes to NetHack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "keys",
						Usage:    "Pipe-separated keys to send (e.g., a|b|c)",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					target, err := getTarget()
					if err != nil {
						return err
					}

					keysStr := c.String("keys")
					if keysStr == "" {
						return fmt.Errorf("no keys to send")
					}

					keys := strings.Split(keysStr, "|")

					output, err := tmux.SendKeys(target, keys)
					if err != nil {
						return err
					}

					formatted, err := formatOutput(output)
					if err != nil {
						return err
					}

					fmt.Print(formatted)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
