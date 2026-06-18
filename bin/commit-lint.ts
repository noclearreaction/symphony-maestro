#!/usr/bin/env -S deno run --allow-read --allow-env
// bin/commit-lint.ts

const commitMsgFile = Deno.args[0];

if (!commitMsgFile) {
  console.error("%cError: Missing commit message or file path argument.", "color: red; font-weight: bold");
  console.log("Usage: commit-lint <commit-msg-or-file-path>");
  Deno.exit(1);
}

let commitMsg = commitMsgFile;

try {
  // Check if argument is a file path and read it
  const fileInfo = await Deno.stat(commitMsgFile);
  if (fileInfo.isFile) {
    commitMsg = await Deno.readTextFile(commitMsgFile);
  }
} catch {
  // Not a file or cannot be read, treat commitMsgFile as the raw commit message
}

commitMsg = commitMsg.trim();

// Strip comment lines starting with # (Git automatically removes these)
const lines = commitMsg.split("\n")
  .map(line => line.trim())
  .filter(line => !line.startsWith("#") && line.length > 0);

if (lines.length === 0) {
  console.error("%cError: Empty commit message.", "color: red; font-weight: bold");
  Deno.exit(1);
}

const subjectLine = lines[0];

// Regex for Conventional Commits
// type(scope): description
const conventionalRegex = /^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(?:\([a-z0-9/_-]+\))?:\s+.+/;

if (!conventionalRegex.test(subjectLine)) {
  console.error("%cError: Commit message does not conform to Conventional Commits standard.", "color: red; font-weight: bold");
  console.error(`Invalid message: "${subjectLine}"\n`);
  console.log("%cExpected format:", "font-weight: bold; color: blue");
  console.log("  <type>(<scope>): <subject>\n");
  console.log("%cAllowed types:", "font-weight: bold");
  console.log("  feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert\n");
  console.log("%cExamples:", "font-weight: bold");
  console.log("  feat(git): add commit linter script");
  console.log("  docs(readme): update build instructions");
  console.log("  fix(auth): resolve memory leak in login flow");
  Deno.exit(1);
}

console.log("%c✓ Commit message conforms to Conventional Commits.", "color: green");
Deno.exit(0);
