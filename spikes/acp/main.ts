#!/usr/bin/env -S deno run --allow-run --allow-read --allow-env
// spikes/acp/main.ts
//
// Stage 0 of the Ensemble build order (docs/ensemble/README.md §8): drive a Vibe
// player over the Agent Client Protocol and report what the contract actually
// gives back. Throwaway — the point is the findings, not the code.
//
// It answers four questions:
//
//   1. Does an ACP session run a real turn against the installed vibe-acp?
//   2. Do custom instructions land, i.e. can we choose which agent profile —
//      and therefore which system prompt — a player runs under?
//   3. Is the session free of the 385-char headless fragment that `-p` forces?
//   4. Does the traffic route through rubato?
//
// (4) is checked outside this program, by watching the proxy's log while it runs.
// The rest are answered from the session record Vibe writes to $VIBE_HOME.

import * as acp from "@agentclientprotocol/sdk";

// ---------------------------------------------------------------------------
// Arguments
// ---------------------------------------------------------------------------

interface Options {
  agent: string | undefined;
  prompt: string;
  workdir: string;
}

function parseArgs(argv: string[]): Options {
  const opts: Options = {
    agent: undefined,
    prompt: "Reply with exactly the word MANGO and nothing else.",
    workdir: "/workspace",
  };

  for (let i = 0; i < argv.length; i++) {
    const next = () => {
      const v = argv[++i];
      if (v === undefined) {
        console.error(`Missing value for ${argv[i - 1]}`);
        Deno.exit(2);
      }
      return v;
    };
    switch (argv[i]) {
      case "--agent":
        opts.agent = next();
        break;
      case "--prompt":
        opts.prompt = next();
        break;
      case "--workdir":
        opts.workdir = next();
        break;
      case "-h":
      case "--help":
        console.log(
          [
            "Usage: deno task spike [options]",
            "",
            "  --agent NAME     Vibe agent profile to run the player as.",
            "                   Passed as VIBE_DEFAULT_AGENT, since vibe/acp/agent.py",
            "                   builds SessionOptions without an agent parameter —",
            "                   config is the only lever, and VIBE_ env vars override it.",
            "  --prompt TEXT    Prompt to send (default: a one-word reply).",
            "  --workdir DIR    Working directory for the session (default: /workspace).",
          ].join("\n"),
        );
        Deno.exit(0);
        break;
      default:
        console.error(`Unknown argument: ${argv[i]}`);
        Deno.exit(2);
    }
  }
  return opts;
}

// ---------------------------------------------------------------------------
// Session record — what Vibe wrote to disk for the session we just ran
// ---------------------------------------------------------------------------

interface SessionRecord {
  dir: string;
  systemPromptChars: number;
  headlessFragment: boolean;
  agentProfile: string;
  firstPromptLine: string;
}

/** Vibe names session directories `session_<date>_<first 8 of uuid>`. */
async function readSessionRecord(sessionId: string): Promise<SessionRecord | undefined> {
  const home = Deno.env.get("VIBE_HOME") ?? `${Deno.env.get("HOME")}/.vibe`;
  const root = `${home}/logs/session`;
  const short = sessionId.slice(0, 8);

  let dir: string | undefined;
  try {
    for await (const entry of Deno.readDir(root)) {
      if (entry.isDirectory && entry.name.endsWith(`_${short}`)) dir = `${root}/${entry.name}`;
    }
  } catch {
    return undefined;
  }
  if (!dir) return undefined;

  const meta = JSON.parse(await Deno.readTextFile(`${dir}/meta.json`));
  const prompt: string = meta.system_prompt?.content ?? "";

  return {
    dir,
    systemPromptChars: prompt.length,
    headlessFragment: prompt.includes("# Headless Mode"),
    agentProfile: meta.agent_profile?.name ?? meta.config?.default_agent ?? "(unrecorded)",
    firstPromptLine: prompt.split("\n").find((l: string) => l.trim() !== "") ?? "",
  };
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

async function main() {
  const opts = parseArgs(Deno.args);

  // vibe-acp inherits the environment, and VIBE_-prefixed variables override
  // config.toml. That is the only way to choose a profile per player, because
  // the ACP front-end does not accept one over the protocol.
  const env: Record<string, string> = {};
  if (opts.agent) env["VIBE_DEFAULT_AGENT"] = opts.agent;

  const child = new Deno.Command("vibe-acp", {
    stdin: "piped",
    stdout: "piped",
    stderr: "piped",
    cwd: opts.workdir,
    env,
  }).spawn();

  // Surface agent diagnostics rather than swallowing them — a spike that fails
  // silently teaches nothing.
  child.stderr.pipeTo(
    new WritableStream({
      write(chunk) {
        Deno.stderr.write(chunk);
      },
    }),
  ).catch(() => {});

  const stream = acp.ndJsonStream(child.stdin, child.stdout);

  let sessionId = "";
  let toolCalls = 0;
  let permissionRequests = 0;

  try {
    const result = await acp
      .client({ name: "ensemble-acp-spike" })
      .onRequest(acp.methods.client.session.requestPermission, (ctx) => {
        permissionRequests++;
        const opt = ctx.params.options[0];
        console.log(`  permission asked: ${ctx.params.toolCall.title} -> auto "${opt.name}"`);
        return { outcome: { outcome: "selected" as const, optionId: opt.optionId } };
      })
      .onRequest(acp.methods.client.fs.readTextFile, () => ({ content: "" }))
      .onRequest(acp.methods.client.fs.writeTextFile, () => ({}))
      .connectWith(stream, async (ctx) => {
        const init = await ctx.request(acp.methods.agent.initialize, {
          protocolVersion: acp.PROTOCOL_VERSION,
          clientCapabilities: { fs: { readTextFile: true, writeTextFile: true } },
        });
        console.log(`initialize      protocol v${init.protocolVersion}`);
        console.log(`                agent ${init.agentInfo?.name} ${init.agentInfo?.version}`);

        return ctx.buildSession(opts.workdir).withSession(async (session) => {
          sessionId = session.sessionId;
          console.log(`session/new     ${sessionId}`);
          console.log(`prompt          ${JSON.stringify(opts.prompt)}\n`);

          session.prompt(opts.prompt);

          for (;;) {
            const message = await session.nextUpdate();
            if (message.kind === "stop") return message.response;

            const update = message.notification.update;
            switch (update.sessionUpdate) {
              case "agent_message_chunk":
                if (update.content.type === "text") {
                  Deno.stdout.write(
                    new TextEncoder().encode(update.content.text),
                  );
                }
                break;
              case "tool_call":
                toolCalls++;
                console.log(`\n  tool_call: ${update.title} (${update.status})`);
                break;
              default:
                break;
            }
          }
        });
      });

    console.log(`\n\nstopReason      ${result.stopReason}`);
    if (result.usage) {
      console.log(
        `usage           ${result.usage.totalTokens} total ` +
          `(${result.usage.inputTokens} in / ${result.usage.outputTokens} out)`,
      );
    }
    console.log(`tool calls      ${toolCalls}`);
    console.log(`permissions     ${permissionRequests}`);

    const record = await readSessionRecord(sessionId);
    if (!record) {
      console.log(`\nNo session record found for ${sessionId}.`);
    } else {
      console.log(`\nsession record  ${record.dir}`);
      console.log(`  agent         ${record.agentProfile}`);
      console.log(`  system prompt ${record.systemPromptChars} chars`);
      console.log(`  first line    ${JSON.stringify(record.firstPromptLine)}`);
      console.log(
        `  headless      ${record.headlessFragment ? "PRESENT — unexpected" : "absent"}`,
      );
    }
  } finally {
    try {
      child.kill();
    } catch {
      // already gone
    }
    await child.status;
  }
}

if (import.meta.main) await main();
