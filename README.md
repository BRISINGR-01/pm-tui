# pm-tui

A terminal UI for managing packages across multiple package managers, built with [Bubble Tea](https://charm.land/bubbletea) and [Bubbles](https://charm.land/bubbles).

## Features

- Browse and search installed packages
- Install, update, and remove packages
- Fuzzy search across available packages
- Switch between package managers at runtime
- Update system
- Recent packages history

## Supported Package Managers

| Manager | Platform |
|---------|----------|
| `pacman` | Arch Linux |
| `yay` | Arch Linux (AUR) |
| `apt` | Ubuntu / Debian |
| `npm` | Node.js |
| `pip` | Python |

## Keybindings

| Key | Action |
|-----|--------|
| `↑ / ↓` | Navigate list |
| `enter` | Select |
| `/` | Filter list |
| `esc` | Go back |
| `i` | Install selected |
| `u` | Update selected |
| `d` | Remove selected |
| `q` | Quit |

## Installation

### Download
```bash
git clone https://github.com/BRISINGR-01/pm-tui
cd pm-tui
```

### Run
```bash
go build -o pm-tui .
./pm-tui
```

or 
```bash
go run .
```

## Test Environments

Test environments are provided as Docker containers to verify behaviour across different package managers without affecting your host system.

```
tests/
├── ubuntu/     apt-based environment
├── npm/        Node.js / npm environment
└── pip/        Python / pip environment
```

### Running a test environment

From the repo root:

```bash
# Ubuntu / apt
docker build -t pm-tui-ubuntu . -f ./tests/ubuntu/Dockerfile
docker run -it --rm pm-tui-ubuntu
```

## Debug Logging

A `debug.log` file is optioanlly written to the project root at runtime. Tail it while running to trace messages:

```bash
tail -f debug.log
```

To enable it, run with:

```bash
go run -tags dev .
```