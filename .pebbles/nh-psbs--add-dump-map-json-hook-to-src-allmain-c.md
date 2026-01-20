---
id: nh-psbs
title: Add dump_map_json hook to src/allmain.c
type: task
status: new
created: "2026-01-20T10:26:01+01:00"
changed: "2026-01-20T10:33:18+01:00"
---
Add hook call to src/allmain.c to trigger JSON map dump after each turn completes.

## Implementation

**File to modify:** /workspaces/nethack-mcp/nethack-llm/src/allmain.c

**Location to add:** After line 340, immediately after the closing comment block for per-turn processing

**Find this pattern:**
```c
        } /* actual time passed */

        /****************************************/
        /* once-per-player-input things go here */
        /****************************************/
```

**Add hook between comment blocks:**
```c
        } /* actual time passed */

        /* Dump map data if NETHACK_DUMP_FILENAME is set */
        dump_map_json();

        /****************************************/
        /* once-per-player-input things go here */
        /****************************************/
```

## Placement Rationale

**Context of moveloop() in allmain.c:**

The moveloop function structure (simplified):
```c
for (;;) {  // Main game loop
    // ... get events ...
    
    if (context.move) {  // Actual time passed block
        // ... monster moves ...
        // ... hero movement ...
        // ... status updates ...
        
        /********************************/
        /* once-per-turn things go here */
        /********************************/
        // ... timeout, effects, regen ...
        
        if (u.utotype) deferred_goto();
        
    } /* actual time passed */
    
    // ⬅️ ADD HOOK HERE ⬅️
    
    /****************************************/
    /* once-per-player-input things go here */
    /****************************************/
    
    // ... display updates ...
    // ... player input ...
}
```

**Why this location is optimal:**

1. **After all per-turn state updates:**
   - Monsters have moved (monmov() loop)
   - Hero has moved (domove())
   - Time has advanced (moves++ at line 173)
   - All status effects applied

2. **After vision system updated:**
   - vision_recalc() called at line 362
   - viz_array reflects current visibility
   - cansee(x,y) will return accurate values

3. **After display buffer updated:**
   - bot() called at line 365 (status bar)
   - curs_on_u() called (cursor position)
   - Screen state is consistent

4. **Before player input:**
   - mapglyph() buffer is current
   - gbuf[][] reflects what's displayed
   - File dump won't race with display updates

**What happens at this point:**
- `moves` counter is current
- `u.ux`, `u.uy` are current hero position
- `levl[][]` has current map state
- `cansee(x,y)` returns accurate visibility
- `level.objects[][]`, `level.monsters[][]` are current

## Call Conventions

The hook uses:
```c
dump_map_json();
```

No parameters, no return value. The function:
- Checks NETHACK_DUMP_FILENAME environment variable
- Opens file if env var is set
- Reads global state (moves, level, u.ux, u.uy, etc.)
- Writes JSON to file
- Returns silently if env var is not set or file cannot be opened

## Dependencies on Allmain State

dump_map_json() relies on these allmain/global state variables:
- `moves` - Turn counter (line 173 increment)
- `u.ux`, `u.uy` - Hero position (updated by movement)
- `levl[][]` - Map tile data (updated by movement/events)
- `level.objects[][]` - Object positions
- `level.monsters[][]` - Monster positions
- `viz_array[][]` - Visibility data (updated by vision_recalc)

All of these are current and stable at the hook location.

## Error Handling

The hook itself doesn't need error handling because dump_map_json():
- Returns silently if env var is not set
- Returns silently if file cannot be opened
- Gracefully handles any internal errors
- Won't crash or corrupt game state

## Performance Impact

**When NETHACK_DUMP_FILENAME is NOT set:**
- Environment variable check is a single getenv() call
- Returns immediately with zero overhead
- No performance impact on normal gameplay

**When NETHACK_DUMP_FILENAME IS set:**
- Function runs after every turn
- Iterates through all tiles (80x21 = 1680 tiles)
- Writes JSON to file
- Acceptable for debugging/analysis use case

## Acceptance Criteria

Task is complete when:
1. Hook is added after line 340 in allmain.c
2. Hook calls dump_map_json()
3. Comment is added explaining the hook's purpose
4. Hook is between "actual time passed" and "player input" comment blocks
5. No syntax errors in allmain.c
6. Hook is reachable during gameplay (not inside #ifdef or condition)
7. Code compiles successfully

## Dependencies

- Requires: nh-29n0 (extern.h must declare dump_map_json)
- Requires: nh-x6ym (mapdump.c must implement dump_map_json)
- Blocks: nh-v91z (cannot build/test until hook is in place)

## Testing

After making this change:
1. Start NetHack and set NETHACK_DUMP_FILENAME
2. Make moves (walk around)
3. Verify JSON file is created after each move
4. Verify JSON contains current game state
5. Verify file is overwritten (not appended)
6. Unset NETHACK_DUMP_FILENAME
7. Verify no more files are created
8. Verify gameplay is not affected

## Potential Issues

1. **Hook placement:** If placed incorrectly (e.g., inside other logic blocks), may not be called or may be called at wrong time. Must be after per-turn block.

2. **Compiler warnings:** Ensure proper function prototype exists before use (handled by nh-29n0 task).

3. **Race conditions:** Dump happens before player input, so state should be stable. If other code runs between hook and input, ensure it doesn't modify state incompatibly.

## References

- src/allmain.c - Main game loop structure (lines 22-500)
- docs/mapdata.md - Map data access patterns
- nh-29n0 - Function declaration in extern.h
- nh-x6ym - Function implementation in mapdump.c