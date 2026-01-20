---
id: nh-x6ym
title: Create src/mapdump.c with JSON dumping logic
type: task
status: new
created: "2026-01-20T10:26:01+01:00"
changed: "2026-01-20T10:31:13+01:00"
---
Create src/mapdump.c implementing JSON map dumping functionality.

## Implementation Overview

Create new source file containing all JSON dumping logic with helper functions for terrain mapping, object class mapping, and adjacency checking.

## File Structure

**Location:** /workspaces/nethack-mcp/nethack-llm/src/mapdump.c

**Includes:**
```c
#include "hack.h"
#include <stdio.h>
#include <stdlib.h>
```

## Functions to Implement

### Main Function

```c
void dump_map_json(void);
```

**Logic:**
1. Check NETHACK_DUMP_FILENAME environment variable
2. If not set or empty, return silently
3. Open file for writing (overwrite mode "w")
4. If file cannot be opened, return silently
5. Write JSON header with metadata:
   - turn (from global `moves`)
   - dungeon_level (from `depth(&u.uz)`)
   - hero position (from `u.ux`, `u.uy`)
6. Iterate through all tiles (COLNO=80, ROWNO=21):
   - For explored tiles (`seenv > 0`):
     - Write base info: x, y, visible, terrain_type
     - If visible (`cansee(x,y)`): add glyph, objects, monsters
   - Skip unexplored tiles in main loop
7. Write unexplored tiles array:
   - Filter: `seenv == 0` AND adjacent to explored tile
   - Adjacency: 8-direction check (horizontal, vertical, diagonal)
   - Include ALL terrain types (no STONE filtering)
8. Close JSON structure and file

**JSON Building Pattern:**
```c
fprintf(fp, "{\n");
fprintf(fp, "  \"turn\": %ld,\n", moves);
// ... metadata ...
fprintf(fp, "  \"tiles\": [\n");

boolean first = TRUE;
for (int y = 0; y < ROWNO; y++) {
    for (int x = 0; x < COLNO; x++) {
        // ... tile processing ...
        if (!first) fprintf(fp, ",\n");
        first = FALSE;
        fprintf(fp, "    {\"x\": %d, \"y\": %d, ...", x, y);
    }
}

fprintf(fp, "\n  ],\n  \"unexplored_tiles\": [");
// ... unexplored tiles ...
fprintf(fp, "]\n}\n");
```

### Helper Functions

#### terrain_name(schar typ)

```c
static const char *terrain_name(schar typ);
```

**Purpose:** Map terrain type enum values to string names

**Implementation:** Switch statement covering all 36 terrain types:
- STONE, VWALL, HWALL, TLCORNER, TRCORNER, BLCORNER, BRCORNER
- CROSSWALL, TUWALL, TDWALL, TLWALL, TRWALL, DBWALL, TREE
- SDOOR, SCORR, POOL, MOAT, WATER, DRAWBRIDGE_UP, LAVAPOOL
- IRONBARS, DOOR, CORR, ROOM, STAIRS, LADDER, FOUNTAIN
- THRONE, SINK, GRAVE, ALTAR, ICE, DRAWBRIDGE_DOWN, AIR, CLOUD
- Default: "UNKNOWN"

**Note:** IRONBARS should be "IRON_BARS" for consistency with NetHack's type_names array

#### object_class(int oclass)

```c
static const char *object_class(int oclass);
```

**Purpose:** Map object class character to string name

**Implementation:** Switch statement covering all object classes:
- COIN_CLASS → "COIN"
- POTION_CLASS → "POTION"
- SCROLL_CLASS → "SCROLL"
- WAND_CLASS → "WAND"
- RING_CLASS → "RING"
- AMULET_CLASS → "AMULET"
- FOOD_CLASS → "FOOD"
- TOOL_CLASS → "TOOL"
- WEAPON_CLASS → "WEAPON"
- ARMOR_CLASS → "ARMOR"
- GEM_CLASS → "GEM"
- ROCK_CLASS → "ROCK"
- BALL_CLASS → "BALL"
- CHAIN_CLASS → "CHAIN"
- SPBOOK_CLASS → "SPELLBOOK"
- VENOM_CLASS → "VENOM"
- Default: "UNKNOWN"

**Source:** Object class definitions from include/obj.h

#### adjacent_to_explored(int x, int y)

```c
static boolean adjacent_to_explored(int x, int y);
```

**Purpose:** Check if a tile is adjacent (8 directions) to any explored tile

**Implementation:**
```c
static boolean adjacent_to_explored(int x, int y) {
    for (int dy = -1; dy <= 1; dy++) {
        for (int dx = -1; dx <= 1; dx++) {
            if (dx == 0 && dy == 0) continue; // Skip center tile
            
            int nx = x + dx;
            int ny = y + dy;
            
            if (isok(nx, ny) && levl[nx][ny].seenv > 0) {
                return TRUE;
            }
        }
    }
    return FALSE;
}
```

**Algorithm:**
- Iterate through dx: -1, 0, 1
- Iterate through dy: -1, 0, 1
- Skip (0,0) - the tile itself
- Check if neighbor coordinates are valid with `isok(nx, ny)`
- Check if neighbor is explored: `levl[nx][ny].seenv > 0`
- Return TRUE if ANY neighbor is explored, FALSE otherwise

**Purpose:** Filter unexplored tiles to only frontier adjacent to known area

## Visible Tile Data Collection

For visible tiles, collect additional data:

**Glyph:**
```c
int glyph = glyph_at(x, y);
int ochar, ocolor;
unsigned ospecial;
mapglyph(glyph, &ochar, &ocolor, &ospecial, x, y, 0);
fprintf(fp, "\"glyph\": \"%c\", ", (char)ochar);
```

**Objects:**
```c
if (OBJ_AT(x, y)) {
    struct obj *obj = vobj_at(x, y);
    fprintf(fp, "\"object\": {");
    fprintf(fp, "\"class\": \"%s\", ", object_class(obj->oclass));
    fprintf(fp, "\"name\": \"%s\", ", xname(obj));
    fprintf(fp, "\"count\": %ld", obj->quan);
    fprintf(fp, "}");
}
```

**Monsters:**
```c
if (MON_AT(x, y)) {
    struct monst *mon = m_at(x, y);
    fprintf(fp, "\"monster\": {");
    fprintf(fp, "\"name\": \"%s\", ", mon_nam(mon));
    fprintf(fp, "\"peaceful\": %s", mon->mpeaceful ? "true" : "false");
    fprintf(fp, "}");
}
```

## Code Style

- Follow NetHack coding conventions
- Use static functions for helpers (file-local scope)
- Use fprintf() for JSON generation
- Use boolean variables for clarity
- Proper indentation in JSON output (2 spaces per level)
- No trailing commas in JSON arrays

## Dependencies on Other Files

- include/hack.h - Core headers
- include/rm.h - Map tile structure (levl[x][y], seenv, typ)
- include/extern.h - Function declarations (glyph_at, mapglyph, vobj_at, m_at, xname, mon_nam, isok, cansee)
- include/obj.h - Object structure (oclass, quan)
- include/monst.h - Monster structure (mpeaceful)
- include/you.h - Hero position (u.ux, u.uy) and moves counter
- include/dungeon.h - depth() function

## Acceptance Criteria

Task is complete when:
1. src/mapdump.c file is created with all required functions
2. All 36 terrain types are mapped to strings
3. All 15+ object classes are mapped to strings
4. Adjacency checking correctly identifies frontier tiles
5. JSON output follows specified structure
6. Visible tiles include: glyph, objects, monsters
7. Explored not-visible tiles include only terrain_type
8. Unexplored tiles are filtered by adjacency
9. No STONE filtering (all terrain types included)
10. File compiles without errors or warnings

## Potential Issues to Watch

1. **JSON escaping:** Object/monster names may contain special characters that need escaping (quotes, backslashes). xname() and mon_nam() typically return safe strings, but verify.

2. **Performance:** Iterating through 1680 tiles (80x21) per turn may have overhead. Only runs when NETHACK_DUMP_FILENAME is set.

3. **File permissions:** Silent failure means user won't know if dump fails. This is intentional per requirements.

4. **Object/monster stacking:** vobj_at() and m_at() return first object/monster at location. NetHack typically doesn't stack multiple monsters, but objects can stack. xname() handles quantity.

## References

- docs/mapdata.md - Map data structure and visibility system
- include/rm.h - Tile structures and terrain types (lines 38-78)
- include/obj.h - Object structure and class definitions
- include/monst.h - Monster structure
- include/extern.h - Function declarations for mapglyph, glyph_at, etc.
- src/display.c - type_names array showing terrain name strings (lines 1967-1972)