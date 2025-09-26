import { join, dirname } from "path";
import { fileURLToPath } from "url";

export function getPlatformBinary() {
  const platform = process.platform;
  const arch = process.arch;

  if (platform === "darwin") {
    return arch === "arm64" ? "consol-darwin-arm64" : "consol-darwin-amd64";
  } else if (platform === "linux") {
    return arch === "arm64" ? "consol-linux-arm64" : "consol-linux-amd64";
  } else if (platform === "win32") {
    return "consol-windows-amd64.exe";
  } else {
    throw new Error(`Unsupported platform: ${platform}-${arch}`);
  }
}

export function getBinaryPath() {
  const __filename = fileURLToPath(import.meta.url);
  const __dirname = dirname(__filename);
  return join(__dirname, getPlatformBinary());
}

