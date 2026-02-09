# shui

A terminal UI for live shell command previewing. Type a command, see its output instantly.

![demo](docs/demo.gif)

## Name

**shui** = **sh**ell **ui** — also happens to be the Chinese word for water (水, shuǐ), evoking the fluid, flowing nature of the interface.

## Features

- **Live preview**: Output updates automatically as you type (with 300ms debounce)
- **Error resilience**: Failed commands show errors in the status bar while preserving the last successful output
- **Freeze mode**: Pause updates to reference output while editing your command
- **Stdin piping**: Pipe data into shui and build processing pipelines interactively
- **Output capture**: Write results to a file while editing

## Installation

```bash
go install github.com/edofic/shui@latest
```

Or build from source:

```bash
git clone https://github.com/edofic/shui
cd shui
go build -o shui .
```

## Shell Integration

Add to your shell config to enable command pre-population after exiting:

```bash
# ~/.zshrc
eval "$(shui init zsh)"

# ~/.bashrc
eval "$(shui init bash)"
```

With shell integration, when you exit shui, the command you were editing will be placed on your command line ready to run or continue editing.

## Usage

```bash
shui
```

Then start typing shell commands. The output panel updates automatically after you stop typing.

### Piping Input

Pipe data into shui to build processing pipelines interactively:

```bash
curl https://api.example.com/data | shui
```

The piped data:
- Displays immediately in the output pane so you can see what you're working with
- Gets fed as stdin to every command you type
- Reappears when you clear the editor (`Ctrl+D`)

This lets you iteratively build a pipeline like `jq '.items[]' | sort | head -5` while seeing results live.

### Output Capture

Write results to a file with the `-o` flag:

```bash
curl https://api.example.com/data | shui -o results.json
```

Or set the output file interactively with `Ctrl+O` while shui is running.

Every successful command execution writes to the output file. When the editor is empty and stdin data is present, the raw stdin is written instead — so you can use shui as a pass-through:

```bash
# Preview and optionally process before saving
curl https://api.example.com/data | shui -o data.json
# See the data, maybe type 'jq .items' to extract, or just exit to save raw
```

## Key Bindings

| Key | Action |
|-----|--------|
| `Ctrl+C` / `Esc` | Quit |
| `Ctrl+F` | Toggle freeze (pause/resume live updates) |
| `Ctrl+D` | Clear editor (shows stdin if present) |
| `Ctrl+L` | Clear output |
| `Ctrl+O` | Set output file |
| `PgUp` / `PgDown` | Scroll output |

## License

MIT
