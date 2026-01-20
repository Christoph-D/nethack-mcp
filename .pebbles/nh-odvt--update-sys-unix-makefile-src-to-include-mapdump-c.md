---
id: nh-odvt
title: Update sys/unix/Makefile.src to include mapdump.c
type: task
status: new
created: "2026-01-20T10:26:01+01:00"
changed: "2026-01-20T10:33:18+01:00"
---
Update sys/unix/Makefile.src to include the new mapdump.c source file in the build process.

## Implementation

**File to modify:** /workspaces/nethack-mcp/nethack-llm/sys/unix/Makefile.src

**Location to modify:** In the HACKCSRC variable definition, around line 138

**Current structure:**
```makefile
HACKCSRC = allmain.c alloc.c apply.c artifact.c attrib.c ball.c bones.c \
           botl.c cmd.c dbridge.c decl.c detect.c dig.c display.c dlb.c do.c \
           do_name.c do_wear.c dog.c dogmove.c dokick.c dothrow.c drawing.c \
           dungeon.c eat.c end.c engrave.c exper.c explode.c extralev.c \
           files.c fountain.c hack.c hacklib.c invent.c isaac64.c light.c \
           lock.c mail.c makemon.c mapglyph.c mcastu.c mhitm.c mhitu.c \
           minion.c mklev.c mkmap.c mkmaze.c mkobj.c mkroom.c mon.c mondata.c \
           monmove.c monst.c mplayer.c mthrowu.c muse.c music.c o_init.c \
           objects.c objnam.c options.c pager.c pickup.c pline.c polyself.c potion.c \
           pray.c priest.c quest.c questpgr.c read.c rect.c region.c restore.c \
           rip.c rnd.c role.c rumors.c save.c shk.c shknam.c sit.c sounds.c \
           sp_lev.c spell.c steal.c steed.c sys.c teleport.c timeout.c \
           topten.c track.c trap.c u_init.c uhitm.c vault.c version.c vision.c \
           weapon.c were.c wield.c windows.c wizard.c worm.c worn.c write.c zap.c
```

**Change required:**
```makefile
HACKCSRC = allmain.c alloc.c apply.c artifact.c attrib.c ball.c bones.c \
           botl.c cmd.c dbridge.c decl.c detect.c dig.c display.c dlb.c do.c \
           do_name.c do_wear.c dog.c dogmove.c dokick.c dothrow.c drawing.c \
           dungeon.c eat.c end.c engrave.c exper.c explode.c extralev.c \
           files.c fountain.c hack.c hacklib.c invent.c isaac64.c light.c \
           lock.c mail.c makemon.c mapglyph.c mcastu.c mhitm.c mhitu.c \
           minion.c mklev.c mkmap.c mkmaze.c mkobj.c mkroom.c mon.c mondata.c \
           monmove.c monst.c mplayer.c mthrowu.c muse.c music.c o_init.c \
           objects.c objnam.c options.c pager.c pickup.c pline.c polyself.c potion.c \
           pray.c priest.c quest.c questpgr.c read.c rect.c region.c restore.c \
           rip.c rnd.c role.c rumors.c save.c shk.c shknam.c sit.c sounds.c \
           sp_lev.c spell.c steal.c steed.c sys.c teleport.c timeout.c \
           topten.c track.c trap.c u_init.c uhitm.c vault.c version.c vision.c \
           weapon.c were.c wield.c windows.c wizard.c worm.c worn.c write.c zap.c mapdump.c
```

**Key change:** Append `mapdump.c` after `zap.c`

## Purpose

HACKCSRC defines all C source files that are compiled to build NetHack:
- Each .c file is compiled to a .o object file
- All object files are linked together to create the NetHack binary
- Adding mapdump.c ensures it's compiled and included in the final binary

## Makefile Variable Structure

```makefile
HACKCSRC = allmain.c alloc.c ... zap.c mapdump.c
```

- Uses backslash (\) continuation lines for readability
- All source files are in src/ directory (parent of sys/unix/)
- Compiled with $(CC) using $(CFLAGS) from build configuration
- Output objects with same basename: mapdump.c → mapdump.o

## Build Process

When you run `make` in sys/unix directory:

1. Reads Makefile.src
2. For each file in HACKCSRC:
   - Calls compiler: $(CC) $(CFLAGS) -I../include -c ../src/filename.c -o filename.o
3. For mapdump.c:
   - $(CC) $(CFLAGS) -I../include -c ../src/mapdump.c -o mapdump.o
4. Links all .o files into NetHack binary

## Alternative Variables

Makefile.src also defines other source lists:
- SYSSRC - System-specific source files
- WIN*SRC - Window system source files (TTY, X11, curses, etc.)
- HACKCSRC - Core NetHack game source files (where mapdump.c belongs)

We only modify HACKCSRC since mapdump.c is a core game module.

## Acceptance Criteria

Task is complete when:
1. mapdump.c is added to end of HACKCSRC list
2. Makefile.src has valid syntax (no parse errors)
3. File follows existing formatting conventions (backslash continuation)
4. mapdump.c is placed after zap.c (alphabetically at end)
5. make runs successfully and creates mapdump.o
6. NetHack binary links successfully

## Dependencies

- Requires: nh-x6ym (mapdump.c must exist)
- Enables: nh-v91z (can build once Makefile is updated)
- Independent of: nh-29n0, nh-psbs (can be done in parallel)

## Verification Steps

After updating Makefile.src:

1. **Check syntax:**
   ```bash
   cd /workspaces/nethack-mcp/nethack-llm/sys/unix
   make -n  # Dry-run to see what make would do
   ```

2. **Attempt build:**
   ```bash
   ./setup.sh  # Run setup if needed
   make
   ```

3. **Verify compilation:**
   - Check that mapdump.o is created
   - Check that make completes without errors
   - Check that NetHack binary is created

4. **Check for warnings:**
   - Look for implicit function declaration warnings (means extern.h task incomplete)
   - Look for unused function warnings (means hook in allmain.c task incomplete)

## Potential Issues

1. **Whitespace/formatting:** Makefiles are whitespace-sensitive. Ensure backslash continuation has no trailing spaces after the backslash.

2. **File order:** Not strictly required, but alphabetical order (zap.c before mapdump.c) follows convention.

3. **Missing prerequisites:** If mapdump.c doesn't exist yet, make will fail. This is expected until nh-x6ym task is complete.

4. **Object name collision:** If another file named mapdump.c existed, it would cause issues. This is a new file name, so unlikely.

## References

- sys/unix/Makefile.src - Build configuration for Unix systems
- sys/unix/Makefile.top - Top-level makefile that includes Makefile.src
- Files listing - Shows mapdump.c should be in src/ directory
- src/mapdump.c - Source file to be included in build