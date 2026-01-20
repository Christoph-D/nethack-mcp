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

   This sets up the required environment variables and launches the opencode agent.

3. Once inside the `opencode` session, instruct the agent:
   > Play nethack

## Running Multiple Agents in Parallel

To run multiple NetHack-playing agents simultaneously, use different session names by setting the `NETHACK_TMUX_SESSION` environment variable:

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

## Project Structure

- `bin/`: Compiled binaries
- `cmd/`: Entry points of `nethack-ctl`, the CLI backend of the MCP tools
- `internal/`: Implementation of `nethack-ctl`
- `harness/`: Agent execution workspace
  - `opencode-nethack.sh`: Launch opencode agent with NetHack MCP tools
  - `.opencode/plugin/nethack-ctl.ts`: opencode plugin that wraps `nethack-ctl` CLI in MCP tools
- `nethack-llm/`: NetHack game fork (git submodule)
- `docs/`: Documentation files
- `AGENTS.md`: Development guide for agents
- `run.sh`: Launch NetHack in tmux
- `Makefile`: Build system configuration
