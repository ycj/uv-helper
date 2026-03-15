# uv-helper

`uv-helper` is a single-file Go helper that prints shell snippets to activate or deactivate a Python virtual environment when you enter or leave a project directory containing a virtualenv (e.g. `.venv`). It supports macOS/Linux (bash, zsh, fish) and Windows (PowerShell, cmd).

Build

```bash
go build -o uv-helper uv-helper.go
# optional: move to a directory in your PATH
sudo mv uv-helper /usr/local/bin/
```

Quick usage

Basic command:

```
uv-helper shell <bash|zsh|fish|powershell|cmd> <cwd>
```

Example (eval the output to activate in your shell):

```bash
eval "$(uv-helper shell bash "$PWD")"
```

Install / Uninstall hooks

`uv-helper` includes `install` and `uninstall` subcommands that append recommended hook snippets to your shell config or remove them. They support interactive confirmation, backups, force and dry-run modes.

Examples:

- Preview install (no write):
  - `uv-helper install zsh --target=/path/to/file --dry-run`
- Interactive install (prompt):
  - `uv-helper install zsh --target=/path/to/file`
- Non-interactive install (create timestamped backup):
  - `uv-helper install zsh --target=/path/to/file --yes`
- Force install (append even if snippet exists):
  - `uv-helper install zsh --target=/path/to/file --yes --force`

- Preview uninstall:
  - `uv-helper uninstall zsh --target=/path/to/file --dry-run`
- Interactive uninstall:
  - `uv-helper uninstall zsh --target=/path/to/file`
- Non-interactive uninstall (create timestamped backup):
  - `uv-helper uninstall zsh --target=/path/to/file --yes`

Flags summary:

- `--target=PATH` : specify target config file (useful for testing).
- `--yes` / `-y` : proceed without interactive prompt.
- `--no-backup` : do not create backup (default is to create timestamped backup when writing). Backups are stored under `~/.uv-helper/backups/` by default.
- `--force` : force append even if snippet is already present.
- `--dry-run` : preview changes only, do not write.

Behavior

- `uv-helper shell` searches from `<cwd>` up to root for `.venv` (preferred), then `venv`, then a directory containing `pyvenv.cfg`. If found, it prints a shell snippet that sets `VIRTUAL_ENV`, prefixes `PATH` with the virtualenv `bin` (or `Scripts` on Windows), and sets `PYTHONNOUSERSITE=1`. It also preserves previous `PATH` in temporary vars to support deactivation.
- If no env is found, it prints a deactivation snippet to restore the previous `PATH` and clear env variables.

Manual test

1. Create a `.venv` directory in a test project.
2. Run:

```bash
eval "$(./uv-helper shell bash /path/to/project)"
echo $VIRTUAL_ENV
```

Install/uninstall test example

```bash
uv-helper install zsh --target=/tmp/test_zshrc --dry-run
uv-helper install zsh --target=/tmp/test_zshrc --yes
uv-helper uninstall zsh --target=/tmp/test_zshrc --dry-run
uv-helper uninstall zsh --target=/tmp/test_zshrc --yes
```

Security notes

- The tool operates with a convenience-first default (auto-activation when a recognized virtualenv is found). For untrusted or multi-user environments, prefer previewing changes with `--dry-run` and avoid `--yes`. Consider only installing hooks in accounts or shells you control.

Files

- Source: `uv-helper.go`
- Shell examples: `examples/shells.md`
- Chinese README: `README.md`

Managing backups

Backups are stored under `~/.uv-helper/backups/`. Useful commands:

- List all backups (names):

```bash
ls -1 ~/.uv-helper/backups/
```

- Show the 10 most recent backups:

```bash
ls -1t ~/.uv-helper/backups/ | head -n 10
```

- Show details for the latest backup:

```bash
ls -l ~/.uv-helper/backups/$(ls -1t ~/.uv-helper/backups | head -n1)
```

- Remove a single backup:

```bash
rm ~/.uv-helper/backups/<backup-filename>
```

- Remove backups older than 30 days (preview first):

```bash
find ~/.uv-helper/backups -type f -mtime +30 -print
# then delete if satisfied:
find ~/.uv-helper/backups -type f -mtime +30 -delete
```

- Interactive removal (confirm each file):

```bash
find ~/.uv-helper/backups -type f -mtime +30 -ok rm {} \;
```

