package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		usageAndExit()
	}

	cmd := os.Args[1]
	switch cmd {
	case "shell":
		if len(os.Args) < 4 {
			usageAndExit()
		}
		shell := os.Args[2]
		cwd := os.Args[3]
		venv := findVenv(cwd)
		script := generateScript(shell, venv)
		fmt.Print(script)
	case "install":
		if len(os.Args) < 3 {
			usageAndExit()
		}
		shell := os.Args[2]
		yes := false
		target := ""
		noBackup := false
		force := false
		dryRun := false
		for _, a := range os.Args[3:] {
			switch {
			case a == "--yes" || a == "-y":
				yes = true
			case a == "--no-backup":
				noBackup = true
			case a == "--force":
				force = true
			case a == "--dry-run":
				dryRun = true
			case strings.HasPrefix(a, "--target="):
				target = strings.TrimPrefix(a, "--target=")
			}
		}
		exe, _ := os.Executable()
		installForShell(shell, exe, target, yes, !noBackup, force, dryRun)
	case "uninstall":
		if len(os.Args) < 3 {
			usageAndExit()
		}
		shell := os.Args[2]
		yes := false
		target := ""
		noBackup := false
		force := false
		dryRun := false
		for _, a := range os.Args[3:] {
			switch {
			case a == "--yes" || a == "-y":
				yes = true
			case a == "--no-backup":
				noBackup = true
			case a == "--force":
				force = true
			case a == "--dry-run":
				dryRun = true
			case strings.HasPrefix(a, "--target="):
				target = strings.TrimPrefix(a, "--target=")
			}
		}
		uninstallForShell(shell, target, yes, !noBackup, force, dryRun)
	default:
		usageAndExit()
	}
}

func usageAndExit() {
	fmt.Fprintf(os.Stderr, "usage:\n  %s shell <bash|zsh|fish|powershell|cmd> <cwd>\n  %s install <shell> [--yes] [--target=/path/to/config] [--no-backup] [--force] [--dry-run]\n  %s uninstall <shell> [--yes] [--target=/path/to/config] [--no-backup] [--force] [--dry-run]\n", os.Args[0], os.Args[0], os.Args[0])
	os.Exit(2)
}

// findVenv walks from start up to root searching for .venv, venv or pyvenv.cfg
func findVenv(start string) string {
	p, err := filepath.Abs(start)
	if err != nil {
		p = start
	}
	for {
		// check .venv
		check := filepath.Join(p, ".venv")
		if existsDir(check) {
			return check
		}
		// check venv
		check = filepath.Join(p, "venv")
		if existsDir(check) {
			return check
		}
		// check pyvenv.cfg in this dir
		cfg := filepath.Join(p, "pyvenv.cfg")
		if existsFile(cfg) {
			return p
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return ""
}

func existsDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func existsFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func singleQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func generateScript(shell, venv string) string {
	switch shell {
	case "bash", "zsh":
		return generatePosix(venv)
	case "fish":
		return generateFish(venv)
	case "powershell":
		return generatePowerShell(venv)
	case "cmd":
		return generateCmd(venv)
	default:
		return ""
	}
}

func generatePosix(venv string) string {
	if venv != "" {
		q := singleQuote(venv)
		return fmt.Sprintf("# uv-helper activation (POSIX)\nif [ \"${VIRTUAL_ENV:-}\" != %s ]; then\n  export __UV_HELPER_PREV_PATH=\"${PATH:-}\"\n  export __UV_HELPER_PREV_VENV=\"${VIRTUAL_ENV:-}\"\n  export VIRTUAL_ENV=%s\n  export PATH=\"$VIRTUAL_ENV/bin:$PATH\"\n  export PYTHONNOUSERSITE=1\nfi\n", q, q)
	}
	return "# uv-helper deactivation (POSIX)\nif [ -n \"${VIRTUAL_ENV:-}\" ]; then\n  export PATH=\"${__UV_HELPER_PREV_PATH:-$PATH}\"\n  export VIRTUAL_ENV=\"${__UV_HELPER_PREV_VENV:-}\"\n  unset __UV_HELPER_PREV_PATH\n  unset __UV_HELPER_PREV_VENV\n  if [ -z \"$VIRTUAL_ENV\" ]; then\n    unset VIRTUAL_ENV\n  fi\n  unset PYTHONNOUSERSITE\nfi\n"
}

func generateFish(venv string) string {
	if venv != "" {
		esc := escapeFish(venv)
		return fmt.Sprintf("## uv-helper activation (fish)\nif not test \"$VIRTUAL_ENV\" = \"%s\"\n  set -gx __UV_HELPER_PREV_PATH $PATH\n  set -gx __UV_HELPER_PREV_VENV $VIRTUAL_ENV\n  set -gx VIRTUAL_ENV %s\n  set -gx PATH \"$VIRTUAL_ENV/bin\" $PATH\n  set -gx PYTHONNOUSERSITE 1\nend\n", esc, esc)
	}
	return "## uv-helper deactivation (fish)\nif test -n \"$VIRTUAL_ENV\"\n  set -gx PATH $__UV_HELPER_PREV_PATH\n  set -gx VIRTUAL_ENV $__UV_HELPER_PREV_VENV\n  set -e __UV_HELPER_PREV_PATH\n  set -e __UV_HELPER_PREV_VENV\n  set -e PYTHONNOUSERSITE\nend\n"
}

func escapeFish(s string) string {
	if strings.ContainsAny(s, " '\\") {
		return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
	}
	return s
}

func generatePowerShell(venv string) string {
	if venv != "" {
		scripts := filepath.Join(venv, "Scripts")
		return fmt.Sprintf("# uv-helper activation (PowerShell)\nif ($env:VIRTUAL_ENV -ne \"%s\") {\n  $env:__UV_HELPER_PREV_PATH = $env:PATH\n  $env:__UV_HELPER_PREV_VENV = $env:VIRTUAL_ENV\n  $env:VIRTUAL_ENV = \"%s\"\n  $env:PATH = \"%s;\" + $env:PATH\n  $env:PYTHONNOUSERSITE = \"1\"\n}\n", escapePowerShell(venv), escapePowerShell(venv), escapePowerShell(scripts))
	}
	return "# uv-helper deactivation (PowerShell)\nif ($env:VIRTUAL_ENV) {\n  $env:PATH = $env:__UV_HELPER_PREV_PATH\n  $env:VIRTUAL_ENV = $env:__UV_HELPER_PREV_VENV\n  Remove-Item Env:\\__UV_HELPER_PREV_PATH -ErrorAction SilentlyContinue\n  Remove-Item Env:\\__UV_HELPER_PREV_VENV -ErrorAction SilentlyContinue\n  Remove-Item Env:\\PYTHONNOUSERSITE -ErrorAction SilentlyContinue\n}\n"
}

func escapePowerShell(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
}

func generateCmd(venv string) string {
	if venv != "" {
		return fmt.Sprintf("@echo off\nrem uv-helper (cmd) - to activate run: call %s\\Scripts\\activate.bat\n", venv)
	}
	return "@echo off\nrem uv-helper (cmd) - deactivation: restore PATH manually or reopen shell\n"
}

// installForShell appends or prints the hook snippet for a given shell.
func installForShell(shell, exePath, target string, yes, backup, force, dryRun bool) {
	home, _ := os.UserHomeDir()
	var cfg string
	switch shell {
	case "zsh":
		cfg = filepath.Join(home, ".zshrc")
	case "bash":
		cfg = filepath.Join(home, ".bashrc")
	case "fish":
		cfg = filepath.Join(home, ".config", "fish", "config.fish")
	case "powershell":
		cfg = filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	case "cmd":
		cfg = filepath.Join(home, "_uv_helper_cmd.bat")
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
		os.Exit(2)
	}
	if target != "" {
		cfg = target
	}

	snippet := generateInstallSnippet(shell, exePath)

	// ensure directory exists
	os.MkdirAll(filepath.Dir(cfg), 0o755)

	// check if already installed
	if existsFile(cfg) {
		data, _ := os.ReadFile(cfg)
		if strings.Contains(string(data), "# >>> uv-helper init >>>") && !force {
			fmt.Printf("uv-helper already installed in %s (use --force to append anyway)\n", cfg)
			return
		}
	}

	if dryRun || !yes {
		action := "append"
		if dryRun {
			action = "(dry-run) append"
		}
		fmt.Printf("Would %s the following snippet to: %s\n---BEGIN SNIPPET---\n%s---END SNIPPET---\n", action, cfg, snippet)
		if dryRun {
			return
		}
		if !yes {
			// interactive prompt
			if !confirmPrompt(fmt.Sprintf("Append snippet to %s? [y/N]: ", cfg)) {
				fmt.Println("Aborted")
				return
			}
		}
	}

	// create centralized backup if requested
	if backup && existsFile(cfg) {
		home, _ := os.UserHomeDir()
		bdir := filepath.Join(home, ".uv-helper", "backups")
		_ = os.MkdirAll(bdir, 0o755)
		ts := time.Now().Format("20060102T150405")
		base := filepath.Base(cfg)
		safe := strings.ReplaceAll(base, string(filepath.Separator), "_")
		bak := filepath.Join(bdir, safe+"."+ts+".bak")
		_ = os.WriteFile(bak, mustReadFile(cfg), 0o644)
		fmt.Printf("Backup created: %s\n", bak)
	}

	f, err := os.OpenFile(cfg, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", cfg, err)
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.WriteString("\n" + snippet + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to %s: %v\n", cfg, err)
		os.Exit(1)
	}
	fmt.Printf("Appended uv-helper snippet to %s\n", cfg)
}

func mustReadFile(p string) []byte {
	b, _ := os.ReadFile(p)
	return b
}

func confirmPrompt(prompt string) bool {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "y" || s == "yes"
}

func generateInstallSnippet(shell, exePath string) string {
	// make exe path safe for shells
	exe := exePath
	switch shell {
	case "zsh", "bash":
		// Use a robust snippet that handles both zsh and bash. Use double quotes around exe for safety.
		return fmt.Sprintf("# >>> uv-helper init >>>\n# added by uv-helper\nuv_auto() { eval \"$(\"%s\" shell %s \"$PWD\")\"; }\nif type add-zsh-hook >/dev/null 2>&1; then\n  autoload -U add-zsh-hook\n  add-zsh-hook chpwd uv_auto\nfi\n# run once\nuv_auto\n# <<< uv-helper init <<<", exe, shell)
	case "fish":
		return fmt.Sprintf("# >>> uv-helper init >>>\n# added by uv-helper\nfunction chpwd\n  eval ( %s shell fish $PWD )\nend\n# run once\neval ( %s shell fish $PWD )\n# <<< uv-helper init <<<", exe, exe)
	case "powershell":
		return fmt.Sprintf("# >>> uv-helper init >>>\n# added by uv-helper\nfunction Set-Location { param($Path) Microsoft.PowerShell.Management\\Set-Location $Path; Invoke-Expression (& \"%s\" shell powershell $PWD) }\n# run once\nInvoke-Expression (& \"%s\" shell powershell $PWD)\n# <<< uv-helper init <<<", exe, exe)
	case "cmd":
		return fmt.Sprintf("rem >>> uv-helper init >>>\nrem added by uv-helper\nrem To use uv-helper in cmd, run: call \"%s\" shell cmd %%CD%%\nrem <<< uv-helper init <<<", exe)
	}
	return ""
}

func uninstallForShell(shell, target string, yes, backup, force, dryRun bool) {
	home, _ := os.UserHomeDir()
	var cfg string
	switch shell {
	case "zsh":
		cfg = filepath.Join(home, ".zshrc")
	case "bash":
		cfg = filepath.Join(home, ".bashrc")
	case "fish":
		cfg = filepath.Join(home, ".config", "fish", "config.fish")
	case "powershell":
		cfg = filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	case "cmd":
		cfg = filepath.Join(home, "_uv_helper_cmd.bat")
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
		os.Exit(2)
	}
	if target != "" {
		cfg = target
	}

	if !existsFile(cfg) {
		fmt.Printf("config file not found: %s\n", cfg)
		return
	}

	data, err := os.ReadFile(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", cfg, err)
		os.Exit(1)
	}
	s := string(data)
	startMarker := "# >>> uv-helper init >>>"
	endMarker := "# <<< uv-helper init <<<"
	si := strings.Index(s, startMarker)
	if si == -1 {
		fmt.Printf("no uv-helper snippet found in %s\n", cfg)
		return
	}
	ei := strings.Index(s[si:], endMarker)
	var newS string
	if ei == -1 {
		newS = s[:si]
	} else {
		newS = s[:si] + s[si+ei+len(endMarker):]
	}

	if dryRun || !yes {
		action := "remove"
		if dryRun {
			action = "(dry-run) remove"
		}
		fmt.Printf("Would %s uv-helper snippet from: %s\n---BEGIN REMAINING---\n%s\n---END REMAINING---\n", action, cfg, newS)
		if dryRun {
			return
		}
		if !yes {
			if !confirmPrompt(fmt.Sprintf("Remove uv-helper snippet from %s? [y/N]: ", cfg)) {
				fmt.Println("Aborted")
				return
			}
		}
	}

	// centralized backup unless disabled
	if backup {
		home, _ := os.UserHomeDir()
		bdir := filepath.Join(home, ".uv-helper", "backups")
		_ = os.MkdirAll(bdir, 0o755)
		ts := time.Now().Format("20060102T150405")
		base := filepath.Base(cfg)
		safe := strings.ReplaceAll(base, string(filepath.Separator), "_")
		bak := filepath.Join(bdir, safe+"."+ts+".bak")
		_ = os.WriteFile(bak, data, 0o644)
		fmt.Printf("Backup created: %s\n", bak)
	}
	err = os.WriteFile(cfg, []byte(newS), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", cfg, err)
		os.Exit(1)
	}
	fmt.Printf("Removed uv-helper snippet from %s\n", cfg)
}
