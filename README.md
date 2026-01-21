# NetHack MCP

A plugin for opencode to play nethack.

## Build

Build the project:

```bash
make install
```

This commands builds the binary and places it into your `PATH`, assuming you
[installed Go](https://go.dev/doc/install).

Other useful make targets:

- `make test`: Run tests
- `make deps`: Download dependencies
- `make help`: Show all targets

## Build NetHack

Before running the system, build NetHack:

```bash
cd nethack-llm/sys/unix
sh setup.sh hints/linux
cd ../../
make
```

## How to Run

Run these scripts in separate terminals:

1. **Terminal 1 - Start NetHack**:

   ```bash
   ./run.sh
   ```

   This launches NetHack in a tmux session named `nethack`.

2. **Terminal 2 - Start the Agent**:

   ```bash
   ./harness/opencode-nethack.sh
   ```

   This sets up the required environment variables and launches the opencode
   agent.

3. Once inside the `opencode` session, instruct the agent:
   > Play nethack

## Environment Variables

- **`NETHACK_TMUX_SESSION`** (default: `nethack`): The tmux session name where NetHack is running. Used by `run.sh` to name the tmux session and by `nethack-ctl` and the harness to target the correct session.

- **`NETHACK_DUMP_FILENAME`** (optional): Path to the JSON file where NetHack writes the annotated map after each turn. If unset, defaults to `/tmp/${NETHACK_TMUX_SESSION}-map.json`.

## Running Multiple Agents in Parallel

To run multiple NetHack-playing agents simultaneously, use different session
names by setting the `NETHACK_TMUX_SESSION` environment variable:

**Agent 1 (default):**

```bash
# Terminal 1
./run.sh

# Terminal 2
./harness/opencode-nethack.sh
```

**Agent 2:**

```bash
# Terminal 3
NETHACK_TMUX_SESSION=nethack2 ./run.sh

# Terminal 4
cp -r harness harness-copy
NETHACK_TMUX_SESSION=nethack2 ./harness-copy/opencode-nethack.sh
```

## Architecture

Components:

- `nethack-ctl`: A Go CLI tool to capture NetHack output and send keys presses
  to it.
- Instrumented NetHack: Modified to output an annotated map in JSON format after
  every turn.
- Opencode MCP plugin: Wraps `nethack-cli` commands as MCP tools.
- Opencode harness: AGENTS.md instructs the agent to play NetHack.

Example map output:

<!-- prettier-ignore -->
```json
{
  "turn": 2,
  "dungeon_level": 1,
  "tiles": {
    "walls": [[69,10], [70,10], [72,10], [73,10], [74,10], [75,10], [76,10], [69,11], [76,11], [69,12], [76,12], [69,13], [76,13], [69,14], [70,14], [71,14], [72,14], [73,14], [74,14], [75,14], [76,14]],
    "open_spaces": [[71,10], [70,11], [71,11], [72,11], [73,11], [74,11], [75,11], [70,12], [72,12], [73,12], [74,12], [75,12], [70,13], [71,13], [72,13], [73,13], [74,13], [75,13]],
    "special": [[[71,12],"stairs up"]]
  },
  "monsters": [[[70,11],"the kitten","peaceful"], [[73,13],"the grid bug","hostile"]],
  "items": [[[73,11],"dwarf corpse"], [[73,11],"blue gem"], [[73,11],"10 darts"]],
  "hero": "70,13"
}
```

## Project Structure

- `bin/`: Compiled binaries
- `cmd/`: Entry points of `nethack-ctl`, the CLI backend of the MCP tools
- `internal/`: Implementation of `nethack-ctl`
- `harness/`: Agent execution workspace
  - `opencode-nethack.sh`: Launch opencode agent with NetHack MCP tools
  - `.opencode/plugin/nethack-ctl.ts`: opencode plugin that wraps `nethack-ctl`
    CLI in MCP tools
- `nethack-llm/`: NetHack game fork (git submodule)
- `docs/`: Documentation files
- `AGENTS.md`: Development guide for agents
- `run.sh`: Launch NetHack in tmux
- `Makefile`: Build system configuration
