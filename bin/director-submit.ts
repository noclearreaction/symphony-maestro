#!/usr/bin/env -S /home/tunnel49/.deno/bin/deno run --allow-run=git,openspec,gh --allow-read=. --allow-write=. --allow-env=PATH,GITHUB_TOKEN,GH_TOKEN
// bin/director-submit.ts

console.log("%c=== Submitting Active OpenSpec Change ===", "color: blue; font-weight: bold");

// 1. Get current branch
const gitBranchCmd = new Deno.Command("git", { args: ["rev-parse", "--abbrev-ref", "HEAD"] });
const { success: branchSuccess, stdout } = await gitBranchCmd.output();
if (!branchSuccess) {
  console.error("%cError: Git branch check failed.", "color: red");
  Deno.exit(1);
}

const currentBranch = new TextDecoder().decode(stdout).trim();
const branchRegex = /^change\/(.+)$/;
const match = currentBranch.match(branchRegex);

if (!match) {
  console.error(`%cError: You are on branch '${currentBranch}'. You must be on a 'change/<name>' branch to submit.`, "color: red");
  Deno.exit(1);
}

const changeName = match[1];
console.log(`Active change identified: ${changeName}`);

// 2. Sync delta specifications to main specs directory
console.log("Synchronizing delta specifications to main specs directory...");
try {
  const deltaSpecsDir = `openspec/changes/${changeName}/specs`;
  const mainSpecsDir = `openspec/specs`;

  // Verify the delta specs directory exists
  try {
    await Deno.stat(deltaSpecsDir);
    // Read all capability directories inside changes/.../specs/
    for await (const entry of Deno.readDir(deltaSpecsDir)) {
      if (entry.isDirectory) {
        const capability = entry.name;
        const targetDir = `${mainSpecsDir}/${capability}`;
        const sourceFile = `${deltaSpecsDir}/${capability}/spec.md`;
        const targetFile = `${targetDir}/spec.md`;

        console.log(`Syncing capability: ${capability}`);
        await Deno.mkdir(targetDir, { recursive: true });

        // Read delta spec
        const deltaContent = await Deno.readTextFile(sourceFile);

        // Extract the ADDED Requirements block
        const addedMatch = deltaContent.match(/## ADDED Requirements([\s\S]*?)(## MODIFIED|## REMOVED|## RENAMED|$)/);
        const addedRequirements = addedMatch ? addedMatch[1].trim() : "";

        if (addedRequirements) {
          // If main spec doesn't exist, create it
          try {
            const existingContent = await Deno.readTextFile(targetFile);
            console.log(`Main spec already exists at ${targetFile}. Appending new requirements...`);
            if (!existingContent.includes(addedRequirements)) {
              await Deno.writeTextFile(targetFile, `${existingContent.trim()}\n\n${addedRequirements}\n`);
            }
          } catch {
            console.log(`Creating new main spec at ${targetFile}...`);
            const mainSpecTemplate = `# Capability: ${capability}

## Purpose
Structured Git, GitHub PR, and custom builder agent local workflow for Director development.

## Requirements

${addedRequirements}
`;
            await Deno.writeTextFile(targetFile, mainSpecTemplate);
          }
        }
      }
    }
  } catch {
    console.log("No delta specifications directory found to sync.");
  }
} catch (error) {
  console.error("%cError: Failed to synchronize specifications dynamically:", "color: red", error);
  Deno.exit(1);
}

// 3. Stage changes safely
console.log("Staging files...");
const filesToStage = [".gitignore", "opencode.json", "openspec/", ".opencode/", "bin/"];
const addCmd = new Deno.Command("git", { args: ["add", ...filesToStage] });
await addCmd.output();

// 4. Check for staged modifications before committing
const diffCmd = new Deno.Command("git", { args: ["diff", "--cached", "--quiet"] });
const { success: noChangesToCommit } = await diffCmd.output();

if (noChangesToCommit) {
  console.log("No modifications staged to commit. Working tree is clean.");
} else {
  const commitMsg = `docs(change): complete ${changeName} and promote specs`;
  console.log(`Committing synced changes: "${commitMsg}"`);
  const commitCmd = new Deno.Command("git", { args: ["commit", "-m", commitMsg] });
  const { success: commitSuccess } = await commitCmd.output();
  if (!commitSuccess) {
    console.error("%cError: Git commit failed.", "color: red");
    Deno.exit(1);
  }
}

// 5. Detect target remote dynamically
const remoteListCmd = new Deno.Command("git", { args: ["remote"] });
const { success: remoteSuccess, stdout: remoteStdout } = await remoteListCmd.output();
const remotes = new TextDecoder().decode(remoteStdout).trim().split("\n").filter(Boolean);

if (remotes.length === 0) {
  console.log("\n%c=== ⚠️ Upstream Remote Required ===", "color: orange; font-weight: bold");
  console.log("No git remote (like 'upstream' or 'origin') is configured in this repository yet.");
  console.log("To set up your remote repository, run:");
  console.log("  git remote add upstream <your-github-repository-url>");
  console.log("=========================================\n");
  Deno.exit(0);
}

const remoteToPush = remotes.includes("upstream") ? "upstream" : remotes[0];
console.log(`Target remote identified: ${remoteToPush}`);

// 6. Verify and Use GitHub CLI
let hasGh = false;
try {
  const ghCheck = new Deno.Command("gh", { args: ["auth", "status"] });
  const { success } = await ghCheck.output();
  hasGh = success;
} catch {
  // gh CLI not available or not logged in
}

if (hasGh) {
  console.log(`Pushing branch '${currentBranch}' to remote '${remoteToPush}'...`);
  const pushCmd = new Deno.Command("git", { args: ["push", "-u", remoteToPush, currentBranch] });
  const { success: pushSuccess } = await pushCmd.output();
  if (!pushSuccess) {
    console.error("%cError: git push failed.", "color: red");
    Deno.exit(1);
  }

  console.log("Creating GitHub Pull Request...");
  const prCmd = new Deno.Command("gh", { args: ["pr", "create", "--fill"] });
  const { success: prSuccess } = await prCmd.output();
  if (!prSuccess) {
    console.warn("%cWarning: 'gh pr create' failed. You may need to create the PR manually.", "color: orange");
  }
} else {
  console.log("\n%c=== ℹ️ Manual Remote Submission ===", "color: blue; font-weight: bold");
  console.log("To push changes to GitHub and create a Pull Request manually, run:");
  console.log(`  1. Push:       git push -u ${remoteToPush} ${currentBranch}`);
  console.log("  2. PR Link:    Go to your repository on GitHub and click 'Compare & Pull Request'.");
  console.log("====================================");
}
