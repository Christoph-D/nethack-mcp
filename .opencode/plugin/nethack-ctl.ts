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
        description: "Send keystrokes to NetHack",
        args: {
          keys: tool.schema.array(tool.schema.string()).describe("Array of keystrokes to send (e.g., ['h', 'y', 'e'])"),
        },
        async execute(args) {
          const proc = spawn(['nethack-ctl', 'send', ...args.keys], {
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
