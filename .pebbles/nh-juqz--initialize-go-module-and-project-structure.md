---
id: nh-juqz
title: Initialize Go module and project structure
type: task
status: fixed
created: "2026-01-18T15:48:42+01:00"
changed: "2026-01-18T16:00:00+01:00"
---
# Initialize Go module and project structure

## Tasks

1. Initialize Go module: `go mod init nethack-ctl` (full project name: https://github.com/Christoph-D/nethack-mcp)
2. Create directory structure:
   - `cmd/nethack-ctl/`
   - `internal/tmux/`
3. Add urfave/cli dependency: `go get github.com/urfave/cli/v2`
