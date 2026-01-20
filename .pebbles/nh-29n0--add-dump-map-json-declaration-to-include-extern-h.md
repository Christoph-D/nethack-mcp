---
id: nh-29n0
title: Add dump_map_json declaration to include/extern.h
type: task
status: new
created: "2026-01-20T10:26:01+01:00"
changed: "2026-01-20T10:33:18+01:00"
---
Add dump_map_json function declaration to include/extern.h to make it accessible to all NetHack source files.

## Implementation

**File to modify:** /workspaces/nethack-mcp/nethack-llm/include/extern.h

**Location to add:** Around line 384, in the display and map functions section

**Add declaration:**
```c
void FDECL(dump_map_json, (void));
```

**Context:** This section contains similar function declarations like:
```c
int FDECL(glyph_at, (XCHAR_P x, XCHAR_P y));
int FDECL(back_to_glyph, (XCHAR_P, XCHAR_P));
int FDECL(mapglyph, (int glyph, int *ochar, int *ocolor, unsigned *ospecial, int x, int y, unsigned mgflags));
void FDECL(newsym, (int, int));
void FDECL(show_glyph, (int, int, int));
void FDECL(vision_recalc, (int));
```

**Place after:** `vision_recalc` declaration or in nearby grouping

## Purpose

Makes the dump_map_json() function (defined in src/mapdump.c) callable from:
- src/allmain.c - Main game loop hook
- Any other NetHack source files if needed in future

## FDECL Macro

The FDECL macro handles function prototype declaration syntax:
- On some platforms/configurations, it adds calling convention modifiers
- Ensures consistent function declarations across NetHack codebase
- Standard pattern used throughout extern.h

## Acceptance Criteria

Task is complete when:
1. extern.h contains: `void FDECL(dump_map_json, (void));`
2. Declaration is in appropriate section (near display/map functions)
3. Follows extern.h formatting conventions
4. No syntax errors in extern.h
5. Function is callable from other source files

## Dependencies

- Requires: nh-x6ym (src/mapdump.c must exist with dump_map_json implementation)
- Enables: nh-psbs (allmain.c can call dump_map_json after declaration exists)

## Testing

After making this change:
1. Verify extern.h is syntactically correct
2. Build should succeed when all other tasks are complete
3. No compilation errors about undefined dump_map_json in allmain.c
4. No warnings about function prototype mismatches

## References

- include/extern.h - Function declarations for display and map functions
- src/mapdump.c - Function definition that will be declared here