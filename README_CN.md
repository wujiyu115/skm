# SKM — AI 智能体技能管理器

[English](README.md)

统一的 AI 编程智能体技能管理工具。单个 Go 二进制文件同时提供 CLI 和 Web UI。

## 功能特性

- **多智能体支持** — Claude Code、Cursor、Codex（可扩展）
- **技能组** — 批量安装/管理技能集合
- **双作用域** — 全局（`~/.agent/skills/`）和项目级（`.agent/skills/`）
- **Web UI** — 内嵌于二进制文件的 React 管理面板
- **符号链接同步** — 中央库 + 符号链接部署（支持复制回退）
- **多种来源** — GitHub URL、简写（`owner/repo`）、本地目录

## 安装

```bash
go install github.com/wujiyu115/skm/cmd/skm@latest
```

或从源码构建：

```bash
make build
```

## 快速开始

```bash
# 从 GitHub 安装技能
skm install owner/repo -a claude -g

# 列出已安装的技能
skm list

# 同步所有技能到检测到的智能体
skm sync

# 创建和使用技能组
skm group create frontend
skm group add frontend my-skill another-skill
skm group install frontend -a claude

# 启动 Web UI
skm serve --open
```

## CLI 命令

| 命令 | 说明 |
|------|------|
| `skm install <source>` | 从 GitHub 或本地路径安装技能 |
| `skm list` | 列出已安装的技能 |
| `skm show <skill>` | 显示技能详情和内容 |
| `skm remove <skill>` | 移除技能 |
| `skm sync` | 同步技能到智能体目录 |
| `skm update [skill\|--all]` | 更新 Git 来源的技能 |
| `skm group <subcommand>` | 管理技能组 |
| `skm agent list\|add\|remove` | 管理智能体适配器 |
| `skm serve` | 启动 Web UI（默认端口 :3721） |
| `skm version` | 显示版本信息 |

## 来源格式

```
https://github.com/owner/repo/tree/branch/path  # GitHub 目录树 URL
https://github.com/owner/repo                    # GitHub 仓库
owner/repo/subpath                               # 简写 + 子路径
owner/repo                                       # 简写
./local/path                                     # 本地目录
```

## 技能格式

技能是一个包含 `SKILL.md` 文件的目录，使用 YAML frontmatter：

```markdown
---
name: my-skill
description: 这个技能做什么
metadata:
  type: coding
  tags: [react, frontend]
---

# AI 智能体的指令内容
```

## 架构

```
~/.skm/
├── skills/          # 中央技能库
├── skm.db           # SQLite 索引
├── .metadata/       # JSON 镜像（可 Git 备份）
└── cache/           # 克隆缓存
```

- **Go** — Cobra CLI、Fiber HTTP、纯 Go SQLite（modernc.org/sqlite）
- **前端** — React 19、Vite、Tailwind CSS v4，通过 `go:embed` 嵌入
- **同步** — 优先符号链接，复制回退，SHA-256 新鲜度检查

## 开发

```bash
# 终端 1：Vite 开发服务器（支持 HMR）
cd web && npm run dev

# 终端 2：Go 服务器代理到 Vite
SKM_DEV=1 go run ./cmd/skm serve

# 运行测试
make test
```

## 许可证

MIT
