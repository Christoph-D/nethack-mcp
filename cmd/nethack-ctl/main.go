package main

import (
	"fmt"
	"os"

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

					fmt.Print(output)
					return nil
				},
			},
			{
				Name:      "send",
				Usage:     "Send keystrokes to NetHack",
				ArgsUsage: "<keys...>",
				Action: func(c *cli.Context) error {
					target, err := getTarget()
					if err != nil {
						return err
					}

					keys := c.Args().Slice()
					if len(keys) == 0 {
						return fmt.Errorf("no keys to send")
					}

					return tmux.SendKeys(target, keys)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
