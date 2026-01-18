---
id: nh-5l6d
title: Use pipe (|) as key separator instead of comma
type: feature
status: fixed
created: "2026-01-18T17:31:47+01:00"
changed: "2026-01-18T17:32:18+01:00"
---
Change the --keys flag separator from comma to pipe (|).

Current: --keys=a,b,c
New: --keys=a|b|c

Update both the CLI tool and the harness plugin.