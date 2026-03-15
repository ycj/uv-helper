

# uv-helper

`uv-helper` 是一个单文件 Go 工具，用于在进入或离开包含 Python 虚拟环境（例如 `.venv`）的目录时，自动产生命令片段由 shell `eval`/`Invoke-Expression` 执行，从而激活或撤销虚拟环境。支持 macOS/Linux（bash、zsh、fish）和 Windows（PowerShell、cmd）。

构建

```bash
go build -o uv-helper uv-helper.go
# 可选：安装到 PATH
sudo mv uv-helper /usr/local/bin/
```

快速用法

基础命令：

```
uv-helper shell <bash|zsh|fish|powershell|cmd> <cwd>
```

示例（在 shell 中 eval 输出以激活/撤销）：

```bash
eval "$(uv-helper shell bash "$PWD")"
```

安装 / 卸载 hook（交互与选项）

`uv-helper` 提供内建的 `install` 与 `uninstall` 子命令，用于把推荐的 hook 片段追加到或从 shell 配置文件移除：

- 预览安装（不写入）：
	- `uv-helper install zsh --target=/path/to/file --dry-run`
- 交互式安装：
	- `uv-helper install zsh --target=/path/to/file`（会提示确认）
- 非交互直接安装（自动备份）：
	- `uv-helper install zsh --target=/path/to/file --yes`
- 强制追加（即使已存在）：
	- `uv-helper install zsh --target=/path/to/file --yes --force`

- 预览卸载：
	- `uv-helper uninstall zsh --target=/path/to/file --dry-run`
- 交互式卸载：
	- `uv-helper uninstall zsh --target=/path/to/file`
- 非交互卸载（创建时间戳备份）：
	- `uv-helper uninstall zsh --target=/path/to/file --yes`

安装 / 卸载常用标志说明：

- `--target=PATH` : 指定安装/卸载的目标配置文件（可用于测试，例如 `/tmp/test_zshrc`）。
- `--yes` / `-y` : 跳过交互提示，直接执行。
- `--no-backup` : 安装/卸载时不创建备份（默认会在写入前创建时间戳备份，备份存放到 `~/.uv-helper/backups/`）。
- `--force` : 强制追加安装片段，即使文件已检测到同样的片段。
- `--dry-run` : 仅预览要执行的变更，不写入。

行为说明

- 当在 `<cwd>` 或其父目录中发现 `.venv`（优先）、其次 `venv`、或 `pyvenv.cfg` 时，`uv-helper shell` 会输出适用于指定 shell 的激活脚本：设置 `VIRTUAL_ENV`、将 `VIRTUAL_ENV/bin`（或 Windows 的 `Scripts`）放到 `PATH` 前端，并设置 `PYTHONNOUSERSITE=1`。同时会保存原 `PATH` 到临时变量以便撤销。
- 当未发现虚拟环境时，`uv-helper` 输出撤销脚本以恢复之前的 `PATH` 与清理环境变量。

示例 - 手动测试

1. 在项目中创建 `.venv`（可以为空或包含 `pyvenv.cfg`）。
2. 运行：

```bash
eval "$(./uv-helper shell bash /path/to/project)"
echo $VIRTUAL_ENV
```

示例 - 安装到测试文件并回退

```bash
# 预览并安装到临时文件
uv-helper install zsh --target=/tmp/test_zshrc --dry-run
uv-helper install zsh --target=/tmp/test_zshrc --yes

# 卸载并创建备份
uv-helper uninstall zsh --target=/tmp/test_zshrc --dry-run
uv-helper uninstall zsh --target=/tmp/test_zshrc --yes
```

安全建议

- 当前实现默认按用户指示使用“自动允许”的工作流（检测到 `.venv` 就激活）。在多人协作或受信任度低的目录中，建议使用 `--dry-run` 先预览或在 `install` 时选择不自动应用（不使用 `--yes`），并尽量启用显式授权/白名单策略（将来版本可能支持）。

备份管理

`uv-helper` 的备份默认存放在 `~/.uv-helper/backups/`。下面是一些常用的管理命令：

- 列出所有备份（按名称）：

```bash
ls -1 ~/.uv-helper/backups/
```

- 按时间排序显示最近的 10 个备份：

```bash
ls -1t ~/.uv-helper/backups/ | head -n 10
```

- 查看某个备份的详细信息（例如最近的一个）：

```bash
ls -l ~/.uv-helper/backups/$(ls -1t ~/.uv-helper/backups | head -n1)
```

- 删除单个备份文件：

```bash
rm ~/.uv-helper/backups/<backup-filename>
```

- 删除早于 30 天的备份（慎用，先用 `-print` 预览）：

```bash
find ~/.uv-helper/backups -type f -mtime +30 -print
# 如果满意：
find ~/.uv-helper/backups -type f -mtime +30 -delete
```

- 交互式删除（逐个确认）：

```bash
find ~/.uv-helper/backups -type f -mtime +30 -ok rm {} \;
```

更多

- 示例 shell 片段：`examples/shells.md`
- 源码（单文件）：`uv-helper.go`

如果需要，我可以把当前 `install` 默认行为改为总是先创建备份到 `~/.uv-helper/backups/`，或在 `README_en.md` 中生成完整英文说明。

