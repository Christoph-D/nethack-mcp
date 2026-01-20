---
id: nh-nwzi
title: Implement JSON map dump feature for NetHack
type: epic
status: new
created: "2026-01-20T10:25:32+01:00"
changed: "2026-01-20T10:31:13+01:00"
blocked-by:
    - nh-x6ym
    - nh-29n0
    - nh-psbs
    - nh-odvt
    - nh-v91z
---
Implement a feature to dump NetHack map data as JSON after each turn controlled by NETHACK_DUMP_FILENAME environment variable.

## Feature Overview

This feature enables external tools to access NetHack's game state in a structured JSON format after every turn. It provides tiered information based on tile visibility (currently visible, explored but not visible, unexplored).

## Architecture

**Components:**
1. **src/mapdump.c** - New module containing all JSON dumping logic
2. **include/extern.h** - Function declaration to make module accessible
3. **src/allmain.c** - Hook in main game loop to trigger dump
4. **sys/unix/Makefile.src** - Build system integration

**Control Flow:**
```
moveloop() in src/allmain.c
  → After per-turn processing completes (line 340+)
  → Check NETHACK_DUMP_FILENAME environment variable
  → If set: call dump_map_json()
  → dump_map_json() in src/mapdump.c
    → Iterate through all tiles (80x21 grid)
    → Build tiered JSON based on visibility
    → Write to file specified by env var
```

## JSON Structure

```json
{
  "turn": 123,
  "dungeon_level": 5,
  "hero": {"x": 40, "y": 10},
  "tiles": [
    {
      "x": 0, "y": 0,
      "visible": true,
      "terrain_type": "ROOM",
      "glyph": ".",
      "object": {
        "class": "POTION",
        "name": "egg",
        "count": 1
      },
      "monster": {
        "name": "goblin",
        "peaceful": false
      }
    },
    {
      "x": 10, "y": 5,
      "visible": false,
      "terrain_type": "CORR"
    }
  ],
  "unexplored_tiles": [
    {"x": 60, "y": 10}
  ]
}
```

## Tiered Information Strategy

**Visible tiles (cansee(x,y) == true):**
- terrain_type: Full terrain type name (36 possible types)
- glyph: Displayed character from glyph_at()
- object: class, name (from xname()), count (quan)
- monster: name (from mon_nam()), peaceful (mpeaceful)

**Explored but not visible (seenv > 0 && cansee == false):**
- terrain_type only
- Represents fog of war memory

**Unexplored tiles (seenv == 0):**
- Listed in unexplored_tiles array
- Adjacency filter: Only include if adjacent (8 directions) to explored tile
- Includes all terrain types (no STONE filtering)

## Terrain Types

All 36 terrain types from `enum levl_typ_types` in include/rm.h:
- Walls: STONE, VWALL, HWALL, TLCORNER, TRCORNER, BLCORNER, BRCORNER, CROSSWALL, TUWALL, TDWALL, TLWALL, TRWALL, DBWALL
- Features: TREE, SDOOR, SCORR, POOL, MOAT, WATER, DRAWBRIDGE_UP, LAVAPOOL, IRONBARS, DOOR, CORR
- Rooms: ROOM, STAIRS, LADDER, FOUNTAIN, THRONE, SINK, GRAVE, ALTAR, ICE, DRAWBRIDGE_DOWN
- Special: AIR, CLOUD

## Key Design Decisions

1. **Silent activation:** Only activates when NETHACK_DUMP_FILENAME is set; zero overhead otherwise
2. **File overwrite:** Overwrites file each turn rather than appending
3. **Silent failure:** No console output if file cannot be opened/written
4. **Manual JSON generation:** Uses fprintf() instead of external JSON library
5. **Placement in allmain.c:** After all per-turn updates but before player input

## Dependencies

Blocked by (all must be fixed before epic can be completed):
- nh-x6ym: Create src/mapdump.c with JSON dumping logic
- nh-29n0: Add dump_map_json declaration to include/extern.h
- nh-psbs: Add dump_map_json hook to src/allmain.c
- nh-odvt: Update sys/unix/Makefile.src to include mapdump.c
- nh-v91z: Build and test JSON map dump feature

## Success Criteria

Epic is complete when:
1. All 5 blocking subtasks are marked as fixed
2. NetHack compiles successfully with new code
3. JSON dump is generated after each turn when NETHACK_DUMP_FILENAME is set
4. JSON contains correct tiered information based on visibility
5. Unexplored tiles are properly filtered by adjacency
6. No performance impact when environment variable is not set
7. No crashes or instability during gameplay