# NetHack Map Data Analysis

This document describes how the map is stored in NetHack and how to access the currently visible map data.

## Map Storage Structure

### Primary Map Data

Location: `nethack-llm/include/rm.h:613`

The entire level map is stored in the `level` variable of type `dlevel_t`:

```c
typedef struct {
    struct rm locations[COLNO][ROWNO];     // Primary tile data (terrain, etc.)
    struct obj *objects[COLNO][ROWNO];    // Objects at each location
    struct monst *monsters[COLNO][ROWNO]; // Monsters at each location
    struct obj *objlist;                  // Linked list of all objects
    struct obj *buriedobjlist;            // Linked list of buried objects
    struct monst *monlist;                // Linked list of all monsters
    struct damage *damagelist;
    struct cemetery *bonesinfo;
    struct levelflags flags;
} dlevel_t;

extern dlevel_t level;
#define levl level.locations  // Common macro alias
```

### Map Dimensions

Location: `nethack-llm/include/global.h:327-328`

- `COLNO = 80` (columns, x-axis)
- `ROWNO = 21` (rows, y-axis)

## Individual Tile Structure

### Tile Data

Location: `nethack-llm/include/rm.h:418-430`

Each map location contains:

```c
struct rm {
    int glyph;               // What the hero thinks is there (memory)
    schar typ;               // What is really there (terrain type)
    uchar seenv;             // Seen vector (8 directions for wall rendering)
    Bitfield(flags, 5);      // Extra info (door state, wall mode, etc)
    Bitfield(horizontal, 1); // Wall/door is horizontal
    Bitfield(lit, 1);        // Currently lit
    Bitfield(waslit, 1);     // Was previously lit (for memory)
    Bitfield(roomno, 6);     // Room number (for special rooms)
    Bitfield(edge, 1);       // Marks boundaries for special rooms
    Bitfield(candig, 1);     // Exception to Can_dig_down; was a trapdoor
};
```

### Terrain Types

Location: `nethack-llm/include/rm.h:38-78`

Common terrain types include:
- `STONE`, `VWALL`, `HWALL` - Walls
- `DOOR`, `SDOOR` - Doors
- `CORR` - Corridor
- `ROOM` - Room floor
- `STAIRS`, `LADDER` - Level connections
- `POOL`, `WATER`, `LAVAPOOL` - Liquids
- `TREE` - Trees
- And more...

### Key Macros

Location: `nethack-llm/include/rm.h:85-106`

```c
#define IS_WALL(typ)    ((typ) && (typ) <= DBWALL)
#define IS_DOOR(typ)    ((typ) == DOOR)
#define IS_ROOM(typ)    ((typ) >= ROOM)
#define ACCESSIBLE(typ) ((typ) >= DOOR)
#define IS_POOL(typ)    ((typ) >= POOL && (typ) <= DRAWBRIDGE_UP)
```

## Visibility System

### Visibility Array

Location: `nethack-llm/src/decl.c:275`

```c
NEARDATA char **viz_array;  // Row pointers to visibility data
```

### Visibility Flags

Location: `nethack-llm/include/vision.h:14-16`

```c
#define COULD_SEE 0x1  // Location could be seen if lit (line of sight)
#define IN_SIGHT 0x2   // Location is currently visible
#define TEMP_LIT 0x4   // Location is temporarily lit
```

### Visibility Check Macros

Location: `nethack-llm/include/vision.h:30-32`

```c
#define cansee(x, y)    (viz_array[y][x] & IN_SIGHT)  // Is currently visible?
#define couldsee(x, y)   (viz_array[y][x] & COULD_SEE) // Has line of sight?
#define templit(x, y)    (viz_array[y][x] & TEMP_LIT)   // Is temporarily lit?
```

### Seen Vector

Location: `nethack-llm/include/rm.h:517-525`

The `seenv` field tracks which 8 directions a tile has been seen from (for proper wall rendering):

```c
#define SV0   0x01  // North
#define SV1   0x02  // Northeast
#define SV2   0x04  // East
#define SV3   0x08  // Southeast
#define SV4   0x10  // South
#define SV5   0x20  // Southwest
#define SV6   0x40  // West
#define SV7   0x80  // Northwest
```

## Display Buffer

### Glyph Display Buffer

Location: `nethack-llm/src/display.c`

```c
typedef struct {
    xchar new;  // Flag for display update
    int glyph;  // Displayed glyph
} gbuf_entry;

static gbuf_entry gbuf[ROWNO][COLNO];  // Display buffer
```

### Getting Display Glyph

Location: `nethack-llm/src/display.c:1877`

```c
int glyph_at(x, y) {
    if (x < 0 || y < 0 || x >= COLNO || y >= ROWNO)
        return cmap_to_glyph(S_room); /* XXX */
    return gbuf[y][x].glyph;
}
```

## Key Functions

### Map Data Access

Location: `nethack-llm/include/extern.h:384`

```c
int FDECL(glyph_at, (XCHAR_P x, XCHAR_P y));      // Get displayed glyph at location
int FDECL(back_to_glyph, (XCHAR_P, XCHAR_P));     // Get background terrain glyph
int FDECL(mapglyph, (int glyph, int *ochar, int *ocolor, unsigned *ospecial, int x, int y, unsigned mgflags));
void FDECL(newsym, (int, int));                   // Update display at location
void FDECL(show_glyph, (int, int, int));          // Show specific glyph at location
void FDECL(vision_recalc, (int));                 // Recalculate visibility
```

### Object/Monster Access

Location: `nethack-llm/include/rm.h:630-663`

```c
#define OBJ_AT(x, y)      (level.objects[x][y] != (struct obj *) 0)
#define MON_AT(x, y)      (level.monsters[x][y] != (struct monst *) 0)
#define m_at(x, y)        (MON_AT(x, y) ? level.monsters[x][y] : (struct monst *) 0)
#define vobj_at(x, y)     (level.objects[x][y])
```

## Getting the Visible Map

To iterate through the currently visible map:

```c
for (int y = 0; y < ROWNO; y++) {
    for (int x = 0; x < COLNO; x++) {
        // Check if tile is currently visible
        if (cansee(x, y)) {
            // Get the remembered glyph (what the hero thinks is there)
            int glyph = levl[x][y].glyph;

            // Get the actual terrain type
            int terrain_type = levl[x][y].typ;

            // Get objects at this location
            struct obj *obj = level.objects[x][y];

            // Get monsters at this location
            struct monst *mon = level.monsters[x][y];

            // Get seen vector (which directions tile was seen from)
            uchar seen_from = levl[x][y].seenv;

            // Get lighting info
            boolean is_lit = levl[x][y].lit;
            boolean was_lit = levl[x][y].waslit;

            // Process visible tile...
        }
    }
}
```

### Alternative: Using Display Buffer

For getting what's actually displayed on screen:

```c
for (int y = 0; y < ROWNO; y++) {
    for (int x = 0; x < COLNO; x++) {
        if (cansee(x, y)) {
            int displayed_glyph = glyph_at(x, y);
            // Convert glyph to character, color, and special flags
            int char, color;
            unsigned special;
            mapglyph(displayed_glyph, &char, &color, &special, x, y, 0);
            // Use char and color to render...
        }
    }
}
```

## Map Data Layers

NetHack uses three layers for map data:

1. **Physical Layer** (`levl[][]`) - What's actually on the map (terrain, objects, monsters)
   - `levl[x][y].typ` - Actual terrain type
   - `level.objects[x][y]` - Objects at location
   - `level.monsters[x][y]` - Monsters at location

2. **Mental Layer** (`levl[][].glyph`) - What the hero remembers/thinks is there
   - `levl[x][y].glyph` - Remembered glyph (includes memory of visited areas)
   - `levl[x][y].seenv` - Which directions tile was seen from
   - `levl[x][y].waslit` - Whether tile was lit when seen

3. **Display Layer** (`gbuf[][]`) - What's currently shown on screen
   - `gbuf[y][x].glyph` - Currently displayed glyph
   - `glyph_at(x, y)` - Get displayed glyph

4. **Vision Layer** (`viz_array[][]`) - What's currently visible (fog of war)
   - `viz_array[y][x] & IN_SIGHT` - Currently visible
   - `viz_array[y][x] & COULD_SEE` - Has line of sight

## Example: Dump Visible Map

```c
void dump_visible_map() {
    for (int y = 0; y < ROWNO; y++) {
        for (int x = 0; x < COLNO; x++) {
            if (cansee(x, y)) {
                // Get what's displayed
                int glyph = glyph_at(x, y);
                int char, color;
                unsigned special;
                mapglyph(glyph, &char, &color, &special, x, y, 0);
                printf("%c", (char)char);
            } else {
                // Not visible
                printf(" ");
            }
        }
        printf("\n");
    }
}
```

## Key Files

- `nethack-llm/include/rm.h` - Map and tile structures
- `nethack-llm/include/vision.h` - Visibility system
- `nethack-llm/include/decl.h` - External variable declarations
- `nethack-llm/include/extern.h` - External function declarations
- `nethack-llm/src/vision.c` - Vision calculation logic
- `nethack-llm/src/display.c` - Display functions
- `nethack-llm/src/mapglyph.c` - Convert glyphs to characters/colors
- `nethack-llm/src/decl.c` - Global variable definitions

## Notes

- Map coordinates are 0-indexed
- The hero position is in `u.ux` and `u.uy`
- The `cansee(x, y)` macro already implements fog of war - only returns true for currently visible tiles
- Unexplored tiles will have `levl[x][y].seenv == 0`
- Explored but currently not visible tiles will have `seenv > 0` but `cansee(x, y) == 0`
- The `glyph` field in `levl` stores the player's memory of what's there (not necessarily what's currently visible)
- For getting the actual displayed glyph (including currently visible monsters), use `glyph_at(x, y)` or access `gbuf[y][x].glyph`
