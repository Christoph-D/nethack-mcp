---
id: nh-ulka
title: Implement CLI entry point and subcommands
type: task
status: new
created: "2026-01-18T15:48:51+01:00"
changed: "2026-01-18T15:50:25+01:00"
blocked-by:
    - nh-tqb8
---
# Implement CLI entry point and subcommands

## File: `cmd/nethack-ctl/main.go`

### CLI Setup

- Use urfave/cli/v2 for CLI framework
- App name: `nethack-ctl`
- Usage: "Control NetHack running in tmux for AI agents"

### Helper Function

`getTarget() (string, error)` - Read `NETHACK_TMUX_SESSION` env var

### Subcommand: `screen`

- Usage: "Capture and display the current NetHack screen"
- Action: Call `tmux.CapturePane()`, print output to stdout

### Subcommand: `send`

- Usage: "Send keystrokes to NetHack"
- ArgsUsage: "<keys...>"
- Action: Validate args provided, call `tmux.SendKeys()`

### Error Handling

- Print errors to stderr with "error: " prefix
- Exit with code 1 on error