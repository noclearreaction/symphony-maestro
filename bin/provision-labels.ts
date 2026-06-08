#!/usr/bin/env -S /home/tunnel49/.deno/bin/deno run --allow-run=gh
// bin/provision-labels.ts

console.log("%c=== Provisioning Standardized GitHub Labels ===", "color: blue; font-weight: bold");

interface Label {
  name: string;
  color: string;
  description: string;
}

const labels: Label[] = [
  { name: "type:feature", color: "0E8A16", description: "New functionality or intent." },
  { name: "type:bug", color: "D93F0B", description: "Unexpected behavior or failure." },
  { name: "type:chore", color: "EDEDED", description: "Internal maintenance, CI, or configuration." },
  { name: "type:spike", color: "FBCA04", description: "Investigation, technical spikes, or documentation (preferred over research)." },
  { name: "status:backlog", color: "EDEDED", description: "Default for new/untriaged work." },
  { name: "status:accepted", color: "FEF2C0", description: "Approved for implementation." },
  { name: "status:in-progress", color: "FEF2C0", description: "Active execution phase." },
  { name: "status:completed", color: "0E8A16", description: "Work finished, PR submitted." },
  { name: "status:blocked", color: "000000", description: "Process halted by external factor." },
  { name: "priority:high", color: "B60205", description: "Critical blocker for milestones." },
  { name: "priority:medium", color: "FBCA04", description: "Standard prioritized work." },
  { name: "priority:low", color: "EDEDED", description: "Elective polish or minor backlog items." },
];

let successCount = 0;
let failureCount = 0;

for (const label of labels) {
  // Strip any leading '#' from the hex color
  const cleanColor = label.color.replace(/^#/, "");

  console.log(`Provisioning label: %c${label.name}%c (Color: ${cleanColor}, Description: "${label.description}")`, "font-weight: bold", "color: inherit");

  const command = new Deno.Command("gh", {
    args: [
      "label",
      "create",
      label.name,
      "--color",
      cleanColor,
      "--description",
      label.description,
      "--force",
    ],
  });

  try {
    const { success, stdout, stderr } = await command.output();
    const decoder = new TextDecoder();
    const errText = decoder.decode(stderr).trim();
    const outText = decoder.decode(stdout).trim();

    if (success) {
      console.log(`%c✓ Successfully provisioned/updated label '${label.name}'`, "color: green");
      if (outText) console.log(outText);
      successCount++;
    } else {
      console.error(`%c✗ Failed to provision label '${label.name}':`, "color: red");
      if (errText) console.error(errText);
      failureCount++;
    }
  } catch (error) {
    console.error(`%c✗ Unexpected error running 'gh' command for '${label.name}':`, "color: red", error);
    failureCount++;
  }
}

console.log("\n%c=== Provisioning Summary ===", "color: blue; font-weight: bold");
console.log(`Total: ${labels.length}`);
console.log(`%cSuccess: ${successCount}`, "color: green");
if (failureCount > 0) {
  console.log(`%cFailure: ${failureCount}`, "color: red");
  Deno.exit(1);
} else {
  console.log("%cAll labels provisioned successfully! 🎉", "color: green");
}
