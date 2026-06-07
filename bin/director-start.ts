#!/usr/bin/env -S /home/tunnel49/.deno/bin/deno run --allow-run=git,openspec --allow-read=. --allow-write=. --allow-env=PATH
// bin/director-start.ts

import { parseArgs } from "https://deno.land/std@0.224.0/cli/parse_args.ts";

const args = parseArgs(Deno.args);
const changeName = args._[0];

if (!changeName || typeof changeName !== "string") {
  console.error("%cError: Missing change name argument.", "color: red; font-weight: bold");
  console.log("Usage: director-start <change-name>");
  Deno.exit(1);
}

// Validate kebab-case
const kebabRegex = /^[a-z0-9]+(-[a-z0-9]+)*$/;
if (!kebabRegex.test(changeName)) {
  console.error(`%cError: Change name '${changeName}' must be kebab-case (e.g. 'my-new-feature').`, "color: red");
  Deno.exit(1);
}

const branchName = `change/${changeName}`;
console.log(`%c=== Starting OpenSpec Change: ${changeName} ===`, "color: blue; font-weight: bold");

// 1. Verify Git workspace
try {
  const gitCheck = new Deno.Command("git", { args: ["rev-parse", "--is-inside-work-tree"] });
  const { success } = await gitCheck.output();
  if (!success) throw new Error();
} catch {
  console.error("%cError: Git command failed. Ensure git is installed and you are in a repo.", "color: red");
  Deno.exit(1);
}

// 2. Checkout or switch branch
console.log(`Checking out branch: ${branchName}`);
const checkoutCmd = new Deno.Command("git", { args: ["checkout", "-b", branchName] });
const { success: checkoutSuccess } = await checkoutCmd.output();

if (!checkoutSuccess) {
  console.log(`Branch already exists. Switching to existing branch...`);
  const switchCmd = new Deno.Command("git", { args: ["checkout", branchName] });
  const { success: switchSuccess } = await switchCmd.output();
  if (!switchSuccess) {
    console.error("%cError: Failed to checkout or switch to branch.", "color: red");
    Deno.exit(1);
  }
}

// 3. Instantiate OpenSpec change
console.log(`Initializing OpenSpec change...`);
try {
  const openspecCmd = new Deno.Command("openspec", { args: ["new", "change", changeName] });
  const { success, stderr } = await openspecCmd.output();
  if (!success) {
    const errorString = new TextDecoder().decode(stderr);
    console.error(`%cError: openspec CLI failed:\n${errorString}`, "color: red");
    Deno.exit(1);
  }
} catch {
  console.error("%cError: openspec CLI not found in path.", "color: red");
  Deno.exit(1);
}

console.log(`%c=== Successfully initialized branch and change! ===`, "color: green; font-weight: bold");
console.log(`Active branch: ${branchName}`);
console.log(`You can begin editing planning artifacts at: openspec/changes/${changeName}/`);
