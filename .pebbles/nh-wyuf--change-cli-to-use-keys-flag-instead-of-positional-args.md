---
id: nh-wyuf
title: Change CLI to use --keys flag instead of positional args
type: feature
status: fixed
created: "2026-01-18T17:29:07+01:00"
changed: "2026-01-18T17:30:01+01:00"
---
Modify the `send` command to accept keys via a `--keys` flag with comma-separated values instead of positional arguments.

Current: `nethack-ctl send h`
New: `nethack-ctl send --keys=a,b,c`

This avoids conflict with the built-in `-h` help flag.