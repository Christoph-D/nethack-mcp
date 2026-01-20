---
id: nh-q210
title: Fix undefined reference to dump_map_json
type: bug
status: fixed
created: "2026-01-20T18:57:06+01:00"
changed: "2026-01-20T18:57:36+01:00"
---
The nethack-llm build fails with an undefined reference error for `dump_map_json` in allmain.c:343.

The function is declared in include/extern.h and implemented in src/mapdump.c, but mapdump.c is not being compiled into the binary.

The file mapdump.c is already listed in HACKCSRC but mapdump.o is missing from the HOBJ variable in the Makefile.

Fix: Add mapdump.o to the HOBJ variable in src/Makefile after mapglyph.o (around line 552).