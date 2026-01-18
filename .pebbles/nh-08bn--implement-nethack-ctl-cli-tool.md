---
id: nh-08bn
title: Implement nethack-ctl CLI tool
type: epic
status: new
created: "2026-01-18T15:48:37+01:00"
changed: "2026-01-18T15:50:26+01:00"
blocked-by:
    - nh-45v9
---
# nethack-ctl

A Go CLI tool for AI agents to interact with NetHack running in tmux.

## Overview

The CLI provides two core subcommands:
- `screen` - Capture and display the current NetHack screen
- `send` - Send keystrokes to NetHack

## Configuration

- Environment variable: `NETHACK_TMUX_SESSION`
- Format: `session:window.pane` (e.g., `nethack:0.0`)

## Project Structure

```
nethack-mcp/
├── cmd/
│   └── nethack-ctl/
│       └── main.go           # Entry point, urfave/cli setup
├── internal/
│   └── tmux/
│       └── tmux.go           # tmux interaction (capture, send-keys)
├── go.mod
└── go.sum
```

## Dependencies

- `github.com/urfave/cli/v2` — CLI framework

## Usage Examples

```bash
# Read the current game state
$ nethack-ctl screen

# Send movement keys
$ nethack-ctl send h j k l

# Send special keys
$ nethack-ctl send Escape
$ nethack-ctl send Enter
```
