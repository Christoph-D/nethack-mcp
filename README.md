# NetHack MCP

A plugin for opencode to play nethack.

## Build

To build the project, use the provided Makefile:

```bash
make install
```

This commands builds the binary and places it into your `PATH`, assuming you
[installed Go](https://go.dev/doc/install).

Other useful make targets:

- `make test`: Run tests
- `make deps`: Download dependencies
- `make help`: Show all targets

## How to Run

Follow these steps to start the system and the agent:

1. **Start NetHack**: Run the helper script to launch NetHack in a tmux session.

   ```bash
   ./run.sh
   ```

   This creates a tmux session named `nethack`.

2. **Navigate to Harness**:

   ```bash
   cd harness
   ```

3. **Start the Agent Environment**: Run `opencode` with the environment variable
   pointing to the tmux session.

   ```bash
   NETHACK_TMUX_SESSION=nethack opencode
   ```

4. **Play**: Once inside the `opencode` session, you can instruct the agent:
   > Play nethack

## Project Structure

- `bin/`: Compiled binaries.
- `cmd/`: Entry points for Go applications.
- `harness/`: working directory for the agent execution.
- `internal/`: Internal Go packages and logic.
- `sh/`: Shell scripts for running NetHack.
