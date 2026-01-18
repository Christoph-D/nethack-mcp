---
id: nh-45v9
title: Build and test nethack-ctl
type: task
status: fixed
created: "2026-01-18T15:48:56+01:00"
changed: "2026-01-18T16:03:43+01:00"
blocked-by:
    - nh-ulka
---
# Build and test nethack-ctl

## Build

- Run `go build ./cmd/nethack-ctl`
- Verify binary is created

## Manual Testing

1. Start a tmux session with nethack:
   ```bash
   tmux new-session -d -s nethack
   tmux send-keys -t nethack "nethack" Enter
   ```

2. Test screen capture:
   ```bash
   NETHACK_TMUX_SESSION=nethack ./nethack-ctl screen
   ```

3. Test sending keys:
   ```bash
   NETHACK_TMUX_SESSION=nethack ./nethack-ctl send h j k l
   ```

## Error Cases to Verify

- Missing env var: should show clear error
- Invalid session: should show tmux error
- No keys to send: should show usage error