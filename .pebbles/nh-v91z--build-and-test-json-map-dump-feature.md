---
id: nh-v91z
title: Build and test JSON map dump feature
type: task
status: new
created: "2026-01-20T10:26:01+01:00"
changed: "2026-01-20T10:33:18+01:00"
---
Build NetHack with the new JSON map dump functionality and verify it works correctly through comprehensive testing.

## Build Process

### Step 1: Setup Build Environment

```bash
cd /workspaces/nethack-mcp/nethack-llm/sys/unix
./setup.sh
```

**What setup.sh does:**
- Runs ./sys/unix/hints/linux (or appropriate hint file)
- Creates sys/unix/Makefile from templates
- Creates include/config.h from configuration options
- Prepares build system for compilation

**Expected output:** No errors, messages about creating Makefile and config.h

### Step 2: Compile NetHack

```bash
make
```

**What make does:**
- Reads Makefile (created by setup.sh)
- For each file in HACKCSRC:
  - Compiles src/*.c to *.o
  - Compiles sys/**/*.c to *.o
  - Compiles win/**/*.c to *.o
- Links all .o files into NetHack binary

**Expected output:**
- Compilation messages for each .c file
- Final linking message
- No errors
- NetHack binary created in src/ directory

**Specific files to verify:**
```
cc -c ../src/mapdump.c -o mapdump.o
```
Should appear and complete without errors.

### Step 3: Check for Warnings

Review make output for warnings:
- **Implicit function declaration:** dump_map_json not in extern.h (nh-29n0 issue)
- **Unused variable/parameter:** May indicate code issues
- **Type mismatch:** Function prototype doesn't match implementation

Warnings don't prevent compilation, but should be investigated.

## Verification Checklist

### Build Verification

- [ ] mapdump.o is created
- [ ] NetHack binary is created
- [ ] No compilation errors
- [ ] No critical warnings
- [ ] File size reasonable (indicating successful linking)

### Runtime Verification

#### Test 1: Environment Variable Control

**Setup:**
```bash
# Without environment variable (should not create file)
unset NETHACK_DUMP_FILENAME
./nethack
# Play a few turns
# Check: /tmp/nethack_dump.json should NOT exist
```

**Setup:**
```bash
# With environment variable (should create file)
export NETHACK_DUMP_FILENAME=/tmp/nethack_dump.json
./nethack
# Make a move (press a direction key)
# Check: /tmp/nethack_dump.json should exist
```

**Expected results:**
- File created when env var is set
- No file when env var is not set
- No console output either way (silent operation)

#### Test 2: File Overwrite (Not Append)

```bash
# First turn
echo "Turn 1" > /tmp/test_turn.json
export NETHACK_DUMP_FILENAME=/tmp/test_turn.json
# Make move in NetHack
# Check file: Should contain JSON, not "Turn 1"

# Second turn
# Make another move in NetHack
# Check file: Should contain new JSON, not appended
```

**Expected results:**
- File is overwritten each turn
- No appending (file doesn't grow with each turn)
- File always contains current state only

#### Test 3: JSON Structure Validity

After dumping, examine /tmp/nethack_dump.json:

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
    }
  ],
  "unexplored_tiles": [
    {"x": 60, "y": 10}
  ]
}
```

**Validation tools:**
```bash
python3 -m json.tool /tmp/nethack_dump.json
# or
jq '.' /tmp/nethack_dump.json
```

**Expected results:**
- Valid JSON (no parse errors)
- All required fields present
- No trailing commas
- No unescaped special characters

#### Test 4: Tiered Information

**Test visible tile:**
- Stand in room with object and monster
- Verify tile in JSON has: visible: true, terrain_type, glyph, object, monster
- Verify object has: class, name, count
- Verify monster has: name, peaceful

**Test explored not-visible tile:**
- Explore a room, then walk away (fog of war)
- Verify tile in JSON has: visible: false, terrain_type
- Verify tile does NOT have: glyph, object, monster

**Test unexplored tile:**
- Stand at frontier of explored area
- Verify unexplored_tiles array contains only adjacent tiles
- Verify unexplored tiles at map edge are not included (no explored neighbors)

**Expected results:**
- Visible tiles: Full information
- Explored not-visible: Terrain only
- Unexplored: Adjacency-filtered list

#### Test 5: All Terrain Types

Play NetHack and visit different terrain types:
- Stand on ROOM, CORR, STAIRS, LADDER
- Stand near VWALL, HWALL, DOOR, FOUNTAIN, ALTAR
- Stand on POOL, LAVAPOOL, ICE, TREE

For each terrain type, verify:
- terrain_type field matches expected name
- No "UNKNOWN" values (except for invalid tiles)

**Expected results:**
- All 36 terrain types produce correct string names
- No crashes on unusual terrain types
- IRON_BARS (not IRONBARS) for consistency

#### Test 6: Object Details

Place or encounter various objects:
- **Potions:** Verify class: "POTION", name from xname()
- **Weapons:** Verify class: "WEAPON", name includes enchantment
- **Food:** Verify class: "FOOD", name like "apple"
- **Corpses:** Verify name like "Kobold corpse" or "newt corpse"
- **Stacks:** Pick up multiple items, drop them; verify count field
- **Eggs:** Verify name "egg" (not detailed type)

**Expected results:**
- All object classes map correctly
- Names match what hero sees
- Count is correct for stacks
- No crashes on unusual objects

#### Test 7: Monster Details

Encounter various monsters:
- **Peaceful monsters:** Verify peaceful: true
- **Hostile monsters:** Verify peaceful: false
- **Invisible monsters:** If visible, name shows correctly
- **Pets:** Verify peaceful: true
- **Monsters with items:** Monster name is still correct

**Expected results:**
- Monster names match what hero sees (from mon_nam())
- Peaceful status matches game state
- No crashes on unusual monsters

#### Test 8: Level Transitions

- Move between dungeon levels (stairs, ladder, magic trap)
- After each level change, check JSON:
  - dungeon_level updates correctly
  - tiles array reflects new level
  - unexplored_tiles reflects new level frontier

**Expected results:**
- JSON dumps work correctly on all levels
- Level number is accurate
- Map data resets correctly between levels

#### Test 9: Game State Updates

Perform actions that change state, verify JSON reflects changes:
- **Move around:** hero.x, hero.y update
- **Wait multiple turns:** turn counter increments
- **Pick up items:** object disappears from tile
- **Drop items:** object appears on tile
- **Monster moves:** monster changes positions in tiles array

**Expected results:**
- JSON always reflects current state
- No stale data from previous turns
- Updates happen immediately after action

#### Test 10: Error Handling

Test failure conditions:
- **File permissions:**
  ```bash
  export NETHACK_DUMP_FILENAME=/root/nethack_dump.json  # No write permission
  ./nethack
  ```
  Should continue without errors, no crash

- **Directory doesn't exist:**
  ```bash
  export NETHACK_DUMP_FILENAME=/nonexistent/path/dump.json
  ./nethack
  ```
  Should continue without errors, no crash

- **Disk full:** (If possible to simulate)
  Should continue without errors, no crash

**Expected results:**
- Silent failure (no console output)
- No crash or game interruption
- Gameplay continues normally

## Performance Testing

### Test: Performance Impact Without Dump

```bash
# Without NETHACK_DUMP_FILENAME
unset NETHACK_DUMP_FILENAME
./nethack
# Play for 100 turns
# Note: Should feel normal, no lag
```

### Test: Performance Impact With Dump

```bash
# With NETHACK_DUMP_FILENAME
export NETHACK_DUMP_FILENAME=/tmp/nethack_dump.json
./nethack
# Play for 100 turns
# Note: May have slight delay on each move (acceptable)
```

**Comparison:**
- Without dump: Standard NetHack performance
- With dump: Slight overhead per turn (acceptable for debugging/analysis)
- Overhead should be minimal (just file I/O)

## Edge Cases

### Case 1: Completely Unexplored Level

Start a new game, immediately check JSON after first turn:
- tiles array should have only tiles around starting position
- unexplored_tiles should have adjacent frontier tiles
- No crashes on sparse explored area

### Case 2: Completely Explored Level

Explore entire level:
- All tiles should be in tiles array
- unexplored_tiles should be empty []
- No crashes on fully explored level

### Case 3: Single Explored Tile

Teleport or level load where only one tile is explored:
- tiles array has one tile
- unexplored_tiles has 8 neighbors
- Adjacency algorithm works correctly

### Case 4: Boundary Conditions

Move to map edges (corners, edges):
- No out-of-bounds crashes
- isok() checks work correctly
- Adjacency works at boundaries

## Validation Checklist

All tasks (nh-x6ym, nh-29n0, nh-psbs, nh-odvt) must be fixed before this task can be tested.

Task is complete when:
1. Build succeeds with no errors
2. All 10 runtime tests pass
3. JSON structure is valid and parseable
4. Performance impact is acceptable
5. No crashes or stability issues
6. All edge cases handled correctly
7. File permissions handled gracefully (silent failure)
8. Level transitions work correctly
9. All terrain types produce correct names
10. Object and monster details are accurate

## Dependencies

- Requires: nh-x6ym, nh-29n0, nh-psbs, nh-odvt (all blocking tasks)
- Enables: nh-nwzi (epic complete when this task is fixed)

## Troubleshooting

**If build fails:**
1. Check that all 4 blocking tasks are complete
2. Check Makefile.src syntax (spaces vs tabs)
3. Verify mapdump.c exists in src/ directory
4. Check for missing includes in mapdump.c

**If runtime tests fail:**
1. Verify NETHACK_DUMP_FILENAME is set correctly
2. Check file permissions on output file/directory
3. Review JSON for syntax errors
4. Check NetHack console for any messages

**If JSON is invalid:**
1. Check for missing commas between array elements
2. Check for trailing commas in arrays
3. Check for unescaped quotes in names
4. Verify first/last element handling

## References

- sys/unix/Makefile.src - Build configuration
- sys/unix/setup.sh - Build setup script
- src/mapdump.c - Implementation being tested
- docs/mapdata.md - Map data specifications
- include/rm.h - Terrain type definitions
- include/extern.h - Function declarations