import type { Plugin } from "@opencode-ai/plugin";
import { tool } from "@opencode-ai/plugin";
import { spawn } from "bun";

export const NethackPlugin: Plugin = async () => {
  return {
    tool: {
      nethack_screen: tool({
        description: "Capture and display the current NetHack screen",
        args: {},
        async execute() {
          const proc = spawn(['nethack-ctl', 'screen'], {
            stdout: 'pipe',
          });
          const result = await new Response(proc.stdout).text();
          await proc.exited;
          return result.trim();
        },
      }),
      nethack_send: tool({
        description: "Send keystrokes to NetHack and get the new screen",
        args: {
          keys: tool.schema.array(tool.schema.string()).describe("Array of tmux keystrokes to send (e.g., ['h', 'y', 'Escape'])"),
        },
        async execute(args) {
          const proc = spawn(['nethack-ctl', 'send', '--keys=' + args.keys.join('|')], {
            stdout: 'pipe',
          });
          const result = await new Response(proc.stdout).text();
          await proc.exited;
          return result.trim();
        },
      }),
    },
  };
};

export default NethackPlugin;
