---
id: nh-9v9c
title: Update plugin to use new --keys flag format
type: task
status: fixed
created: "2026-01-18T17:30:26+01:00"
changed: "2026-01-18T17:30:46+01:00"
---
Update the harness plugin to use the new --keys flag format for nethack-ctl send command.

Current: ['nethack-ctl', 'send', ...args.keys]
New: ['nethack-ctl', 'send', '--keys=' + args.keys.join(',')]