#!/usr/bin/env node

import { createWriteStream, chmodSync, readFileSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import https from "https";
import { getPlatformBinary, getBinaryPath } from "./platform.js";

function downloadFile(url, destination) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${url}...`);

    https
      .get(url, (response) => {
        if (response.statusCode === 302 || response.statusCode === 301) {
          return downloadFile(response.headers.location, destination)
            .then(resolve)
            .catch(reject);
        }

        if (response.statusCode !== 200) {
          reject(
            new Error(`Download failed with status ${response.statusCode}`),
          );
          return;
        }

        const file = createWriteStream(destination);
        response.pipe(file);

        file.on("finish", () => {
          file.close();
          try {
            chmodSync(destination, "755");
          } catch (err) {
            // Ignore chmod errors on Windows
          }
          console.log(`Downloaded and installed ${destination}`);
          resolve();
        });

        file.on("error", reject);
      })
      .on("error", reject);
  });
}

async function main() {
  try {
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = dirname(__filename);
    const packageJson = JSON.parse(
      readFileSync(join(__dirname, "package.json"), "utf8"),
    );
    const version = packageJson.version;
    const binaryName = getPlatformBinary();
    const downloadUrl = `https://github.com/blackzarifa/consol/releases/download/v${version}/${binaryName}`;

    await downloadFile(downloadUrl, getBinaryPath());
    console.log("Binary installed successfully!");
  } catch (error) {
    console.error("Failed to download binary:", error.message);
    process.exit(1);
  }
}

main();
