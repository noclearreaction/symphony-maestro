#!/usr/bin/env -S deno run --allow-run=docker --allow-env=GITHUB_COM_TOKEN,COMPOSE_PROJECT_NAME,LOCAL_WORKSPACE_FOLDER
// bin/versions-check.ts
// Checks dependency versions using renovate's local dry-run report.
//
// Usage: versions-check [--no-updates] [--no-uptodate] [--no-skipped] [--plain] [--json] [--github-token=<token>]
//   --no-updates         Hide the "Updates available" section
//   --no-uptodate        Hide the "Up to date" section
//   --no-skipped         Hide the "Skipped" section
//   --plain              Plain text output, no emoji or color
//   --json               Output as JSON (ignores section filters and display flags)
//   --github-token=<t>   GitHub token for API lookups (falls back to $GITHUB_COM_TOKEN)

import { UntarStream } from "jsr:@std/tar/untar-stream";

const args = new Set(Deno.args);
const showUpdates  = !args.has("--no-updates");
const showUpToDate = !args.has("--no-uptodate");
const showSkipped  = !args.has("--no-skipped");
const plain        = args.has("--plain") || args.has("--json");
const json         = args.has("--json");

const githubTokenArg = [...args].find((a) => a.startsWith("--github-token="));
const githubToken = githubTokenArg?.slice("--github-token=".length) ?? Deno.env.get("GITHUB_COM_TOKEN");
const composeProjectName = Deno.env.get("COMPOSE_PROJECT_NAME") ?? "symphony-maestro";
const workspaceFolder = Deno.env.get("LOCAL_WORKSPACE_FOLDER") ?? Deno.cwd();

// --- Run renovate via docker, extract report via docker cp ---
// Runs renovate in a named container (no --rm), copies the report out via
// `docker cp ... -` (tar stream), then removes the container.
// Uses LOCAL_WORKSPACE_FOLDER so the Docker daemon (on the host) can bind-mount
// the workspace correctly under DOOD.

if (!json) {
  if (plain) {
    console.log("Checking dependency versions...");
  } else {
    console.log("%cChecking dependency versions...", "color: blue");
  }
}

// Generate a unique container name to avoid collisions.
const containerName = `renovate-versions-check-${Date.now()}`;

// Ensure the container is always removed on exit.
const cleanup = () => {
  try {
    new Deno.Command("docker", { args: ["rm", "-f", containerName], stdout: "null", stderr: "null" }).outputSync();
  } catch { /* best effort */ }
};
globalThis.addEventListener("unload", cleanup);

// Step 1: run renovate.
const renovateArgs = [
  "run",
  "--name", containerName,
  "-v", `${workspaceFolder}:/usr/src/app`,
  "-e", "LOG_LEVEL=error",
];
if (githubToken) renovateArgs.push("-e", `GITHUB_COM_TOKEN=${githubToken}`);
renovateArgs.push(
  `renovate:${composeProjectName}`,
  "--platform=local",
  "--dry-run=lookup",
  "--report-type=file",
  "--report-path=/tmp/report.json",
  "--github-token-warn=false",
);

const renovateResult = await new Deno.Command("docker", {
  args: renovateArgs,
  stdout: "null",
  stderr: "inherit",
}).output();

if (!renovateResult.success) {
  console.error("%cError: renovate exited with code " + renovateResult.code, "color: red");
  Deno.exit(1);
}

// Step 2: copy report out as tar stream and extract with @std/tar.
const cpResult = await new Deno.Command("docker", {
  args: ["cp", `${containerName}:/tmp/report.json`, "-"],
  stdout: "piped",
  stderr: "inherit",
}).output();

if (!cpResult.success) {
  console.error("%cError: failed to copy report from container", "color: red");
  Deno.exit(1);
}

let reportText = "";
const stream = new ReadableStream<Uint8Array>({
  start(controller) {
    controller.enqueue(cpResult.stdout);
    controller.close();
  },
});
for await (const entry of stream.pipeThrough(new UntarStream())) {
  if (entry.path.endsWith("report.json") && entry.readable) {
    const chunks: Uint8Array[] = [];
    for await (const chunk of entry.readable) chunks.push(chunk);
    reportText = new TextDecoder().decode(
      chunks.reduce((a, b) => { const c = new Uint8Array(a.length + b.length); c.set(a); c.set(b, a.length); return c; }, new Uint8Array())
    );
  }
}

cleanup();

// --- Parse report ---

interface Update {
  newValue: string;
  updateType: string;
}

interface Warning {
  topic: string;
  message: string;
}

interface Dep {
  depName?: string;
  currentValue?: string;
  skipReason?: string;
  updates?: Update[];
  warnings?: Warning[];
}

interface PackageFile {
  deps?: Dep[];
}

interface Report {
  repositories?: {
    local?: {
      packageFiles?: Record<string, PackageFile[]>;
    };
  };
}

const report: Report = JSON.parse(reportText);

const packageFiles = report.repositories?.local?.packageFiles ?? {};

// --- Collect rows ---

const SEVERITY: Record<string, number> = { major: 0, minor: 1, patch: 2 };

const SKIP_REASON: Record<string, string> = {
  "invalid-value":            "version not parseable",
  "unsupported-version":      "version format not supported",
  "ignored":                  "ignored by config",
  "disabled":                 "disabled by config",
  "in-range-only":            "already satisfies range",
  "already-updated":          "already at latest",
  "pin-digest-not-supported": "digest pinning not supported",
  "digest-unavailable":       "digest unavailable",
};

function skipLabel(reason: string): string {
  return SKIP_REASON[reason] ?? reason;
}

interface Row {
  manager: string;
  name: string;
  current: string;
  next: string;
  updateType: string;
}

const upToDate: Row[] = [];
const outdated: Row[] = [];
const skipped: Row[] = [];

for (const [manager, files] of Object.entries(packageFiles)) {
  for (const file of files) {
    for (const dep of file.deps ?? []) {
      if (!dep.depName || !dep.currentValue) continue;

      const unquote = (s: string) => s.replace(/^["']|["']$/g, "");

      const row: Row = {
        manager,
        name: dep.depName,
        current: unquote(dep.currentValue),
        next: "",
        updateType: "",
      };

      if (dep.skipReason) {
        row.updateType = skipLabel(dep.skipReason);
        skipped.push(row);
      } else if (dep.warnings && dep.warnings.length > 0) {
        row.updateType = "lookup failed";
        skipped.push(row);
      } else if (dep.updates && dep.updates.length > 0) {
        row.next = unquote(dep.updates[0].newValue);
        row.updateType = dep.updates[0].updateType;
        outdated.push(row);
      } else {
        upToDate.push(row);
      }
    }
  }
}

// --- Print ---

const MAX_NAME = 52;
const MAX_VER = 24;

function truncate(s: string, max: number): string {
  return s.length > max ? s.slice(0, max - 1) + "…" : s;
}

function fieldMax(field: keyof Row): number {
  return field === "name" ? MAX_NAME : MAX_VER;
}

function colWidth(rows: Row[], field: keyof Row, min: number): number {
  return Math.max(min, ...rows.map((r) => truncate(r[field], fieldMax(field)).length));
}

function rowEmoji(updateType: string): string {
  switch (updateType) {
    case "major": return "❌";
    case "minor": return "⚠️ ";
    case "patch": return "✅";
    default:      return "  ";
  }
}

function rowStyle(updateType: string): string {
  switch (updateType) {
    case "major": return "color: red";
    case "minor": return "color: yellow";
    case "patch": return "color: green";
    default:      return "";
  }
}

function printTable(rows: Row[], showNext: boolean, showReason = false) {
  if (rows.length === 0) return;

  const wMgr  = colWidth(rows, "manager", 8);
  const wName = colWidth(rows, "name", 12);
  const wCur  = colWidth(rows, "current", 9);
  const wNext = showNext ? colWidth(rows, "next", 9) : 0;

  const header = showNext
    ? `   ${"MANAGER".padEnd(wMgr)}  ${"PACKAGE".padEnd(wName)}  ${"CURRENT".padEnd(wCur)}  ${"LATEST".padEnd(wNext)}  TYPE`
    : showReason
    ? `   ${"MANAGER".padEnd(wMgr)}  ${"PACKAGE".padEnd(wName)}  ${"VERSION".padEnd(wCur)}  REASON`
    : `   ${"MANAGER".padEnd(wMgr)}  ${"PACKAGE".padEnd(wName)}  VERSION`;

  if (plain) {
    console.log(header);
  } else {
    console.log(`%c${header}`, "color: gray");
  }

  for (const row of rows) {
    const emoji = showNext && !plain ? rowEmoji(row.updateType) : "  ";
    const mgr   = row.manager.padEnd(wMgr);
    const name  = truncate(row.name, fieldMax("name")).padEnd(wName);
    const cur   = truncate(row.current, fieldMax("current")).padEnd(wCur);
    if (showNext) {
      const next = truncate(row.next, fieldMax("next")).padEnd(wNext);
      if (plain) {
        console.log(`   ${mgr}  ${name}  ${cur}  ${next}  ${row.updateType}`);
      } else {
        console.log(`%c${emoji} ${mgr}  ${name}  ${cur}  ${next}  ${row.updateType}`, rowStyle(row.updateType));
      }
    } else if (showReason) {
      console.log(`%c   ${mgr}  ${name}  ${cur}  ${row.updateType}`, plain ? "" : "color: gray");
    } else {
      console.log(`   ${mgr}  ${name}  ${cur}`);
    }
  }
}

outdated.sort((a, b) => (SEVERITY[a.updateType] ?? 99) - (SEVERITY[b.updateType] ?? 99));

if (json) {
  console.log(JSON.stringify({ outdated, upToDate, skipped }, null, 2));
  Deno.exit(0);
}

if (outdated.length > 0 && showUpdates) {
  plain
    ? console.log("\nUpdates available")
    : console.log("\n%cUpdates available", "color: yellow; font-weight: bold");
  printTable(outdated, true);
} else if (showUpdates) {
  plain
    ? console.log("\nNo packages out of date")
    : console.log("\n%c✅ No packages out of date", "color: green");
}

if (upToDate.length > 0 && showUpToDate) {
  plain
    ? console.log("\nUp to date")
    : console.log("\n%cUp to date", "color: green; font-weight: bold");
  printTable(upToDate, false);
} else if (showUpToDate) {
  plain
    ? console.log("\nNo packages up to date")
    : console.log("\n%c⚠️  No packages up to date", "color: yellow");
}

if (skipped.length > 0 && showSkipped) {
  plain
    ? console.log("\nSkipped")
    : console.log("\n%cSkipped", "color: gray");
  printTable(skipped, false, true);
} else if (showSkipped) {
  plain
    ? console.log("\nNo packages skipped")
    : console.log("\n%c✅ No packages skipped", "color: green");
}

console.log("");
