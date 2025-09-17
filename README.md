# consol - Conflict Resolver

Git merge conflict resolution TUI tool that provides an interactive interface
for resolving conflicts line by line.

## Features

- Auto-discovery of Git conflict files
- Interactive file selection
- Conflict-by-conflict navigation
- Conflict resolution

## Installation

### Manual Installation

**macOS (Apple Silicon):**

```bash
curl -L -o consol https://github.com/blackzarifa/consol/releases/latest/download/consol-darwin-arm64
chmod +x consol
sudo mv consol /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L -o consol https://github.com/blackzarifa/consol/releases/latest/download/consol-darwin-amd64
chmod +x consol
sudo mv consol /usr/local/bin/
```

**Linux (x64):**

```bash
curl -L -o consol https://github.com/blackzarifa/consol/releases/latest/download/consol-linux-amd64
chmod +x consol
sudo mv consol /usr/local/bin/
```

**Linux (ARM64):**

```bash
curl -L -o consol https://github.com/blackzarifa/consol/releases/latest/download/consol-linux-arm64
chmod +x consol
sudo mv consol /usr/local/bin/
```

**Windows:**

1. Download [consol-windows-amd64.exe](https://github.com/blackzarifa/consol/releases/latest/download/consol-windows-amd64.exe)
2. Rename to `consol.exe`
3. Add to your PATH

## Usage

### Basic Usage

Run `consol` in a Git repository with merge conflicts:

```bash
consol
```

Or specify a specific conflict file:

```bash
consol filename.txt
```

### Interface

**Navigation:**

- `j` - Move down a line
- `k` - Move up a line
- `n` - Next conflict
- `p` - Previous conflict
- `o` - Choose "ours" (current branch)
- `t` - Choose "theirs" (incoming changes)
- `w` - Save file
- `b` - Return to file selector
- `q` - Quit

## Screenshots

