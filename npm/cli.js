#!/usr/bin/env node

import { spawn } from "child_process";
import { existsSync } from "fs";
import { getBinaryPath } from "./platform.js";

function main() {
  const binaryPath = getBinaryPath();

  if (!existsSync(binaryPath)) {
    console.error(`Binary not found: ${binaryPath}`);
    console.error("Try reinstalling the package: npm install -g git-consol");
    process.exit(1);
  }

  // Pass arguments to bin
  const args = process.argv.slice(2);

  const child = spawn(binaryPath, args, {
    stdio: "inherit",
    cwd: process.cwd(),
  });

  child.on("error", (error) => {
    console.error(`Failed to start consol: ${error.message}`);
    process.exit(1);
  });

  child.on("exit", (code, signal) => {
    if (signal) process.kill(process.pid, signal);
    else process.exit(code || 0);
  });
}

main();
