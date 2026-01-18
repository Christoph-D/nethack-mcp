---
id: nh-tqb8
title: Implement internal/tmux package
type: task
status: new
created: "2026-01-18T15:48:46+01:00"
changed: "2026-01-18T15:50:25+01:00"
blocked-by:
    - nh-juqz
---
# Implement internal/tmux package

## File: `internal/tmux/tmux.go`

Implement two functions for tmux interaction:

### `CapturePane(target string) (string, error)`

- Execute `tmux capture-pane -p -t <target>`
- Return captured pane content as string
- Handle errors (tmux not found, session not found)

### `SendKeys(target string, keys []string) error`

- Execute `tmux send-keys -t <target> <keys...>`
- Pass all keys as arguments to send-keys
- Handle errors appropriately

## Error Handling

- Wrap exec.ExitError to extract stderr message
- Return clear error messages for common failures