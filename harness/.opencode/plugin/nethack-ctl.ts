import type { Plugin } from "@opencode-ai/plugin";
import { tool } from "@opencode-ai/plugin";
import { spawn } from "bun";

export const NethackPlugin: Plugin = async () => {
  return {
    tool: {
      nethack_screen: tool({
        description: "Capture and display the current NetHack screen (╋number╋ is a column counter and is not part of the game)",
        args: {},
        async execute() {
          const proc = spawn(['nethack-ctl', 'screen'], {
            stdout: 'pipe',
            stderr: 'pipe',
          });
          const [stdout, stderr] = await Promise.all([
            new Response(proc.stdout).text(),
            new Response(proc.stderr).text(),
          ]);
          await proc.exited;
          const stderrText = stderr.trim();
          const result = stdout.trim();
          return stderrText ? `${result}\n\nError: ${stderrText}` : result;
        },
      }),
      nethack_send: tool({
        description: "Send tmux keystrokes to NetHack and get the new screen (╋number╋ is a column counter and is not part of the game)",
        args: {
          keys: tool.schema.array(tool.schema.string()).describe("Array of tmux keystrokes to send (e.g., ['h', 'y', 'Escape']). You must use tmux key names, e.g. Space, Escape."),
        },
        async execute(args) {
          const proc = spawn(['nethack-ctl', 'send', '--keys=' + args.keys.join('|')], {
            stdout: 'pipe',
            stderr: 'pipe',
          });
          const [stdout, stderr] = await Promise.all([
            new Response(proc.stdout).text(),
            new Response(proc.stderr).text(),
          ]);
          await proc.exited;
          const stderrText = stderr.trim();
          const result = stdout.trim();
          return stderrText ? `${result}\n\nError: ${stderrText}` : result;
        },
      }),
    },
  };
};

export default NethackPlugin;
