# Shell hook examples for uv-helper

zsh (add to `~/.zshrc`):

```sh
autoload -U add-zsh-hook
uv_auto() { eval "$(/path/to/uv-helper shell zsh "$PWD")"; }
add-zsh-hook chpwd uv_auto
# run once for current dir
uv_auto
```

bash (add to `~/.bashrc`):

```sh
cd() { builtin cd "$@" && eval "$(/path/to/uv-helper shell bash "$PWD")"; }
# or use PROMPT_COMMAND to run every prompt
# PROMPT_COMMAND='eval "$(/path/to/uv-helper shell bash "$PWD")"'
```

fish (add to `~/.config/fish/config.fish`):

```fish
function chpwd
  eval ( /path/to/uv-helper shell fish $PWD )
end
# run once at startup
eval ( /path/to/uv-helper shell fish $PWD )
```

PowerShell (profile):

```powershell
function Set-Location { param($Path) Microsoft.PowerShell.Management\Set-Location $Path; Invoke-Expression (& "C:\path\to\uv-helper.exe" shell powershell $PWD) }
# run once for current dir
Invoke-Expression (& "C:\path\to\uv-helper.exe" shell powershell $PWD)
```

cmd.exe: limited support — recommend using PowerShell on Windows.
