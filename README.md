# shui

A terminal UI for live shell command previewing. Type a command, see its output instantly.

## Name

**shui** = **sh**ell **ui** — also happens to be the Chinese word for water (水, shuǐ), evoking the fluid, flowing nature of the interface.

## Features

- **Live preview**: Output updates automatically as you type (with 300ms debounce)
- **Error resilience**: Failed commands show errors in the status bar while preserving the last successful output
- **Freeze mode**: Pause updates to reference output while editing your command

## Installation

```bash
go install github.com/yourusername/shui@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/shui
cd shui
go build -o shui .
```

## Usage

```bash
./shui
```

Then start typing shell commands. The output panel updates automatically after you stop typing.

## Key Bindings

| Key | Action |
|-----|--------|
| `Ctrl+C` / `Esc` | Quit |
| `Ctrl+F` | Toggle freeze (pause/resume live updates) |
| `Ctrl+D` | Clear editor |
| `Ctrl+L` | Clear output |
| `PgUp` / `PgDown` | Scroll output |

## License

MIT
