package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go.yozora.eu/nethack-mcp/internal/tmux"
)

func main() {
	app := &cli.App{
		Name:  "nethack-ctl",
		Usage: "Control NetHack running in tmux for AI agents",
		Commands: []*cli.Command{
			{
				Name:  "screen",
				Usage: "Capture and display the current NetHack screen",
				Action: func(c *cli.Context) error {
					target := tmux.GetTarget()

					output, err := tmux.CapturePane(target, false)
					if err != nil {
						return err
					}

					fmt.Print(output)
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
					target := tmux.GetTarget()

					keysStr := c.String("keys")
					if keysStr == "" {
						return fmt.Errorf("no keys to send")
					}

					keys := strings.Split(keysStr, "|")

					output, err := tmux.SendKeys(target, keys)
					if err != nil {
						return err
					}

					fmt.Print(output)
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
