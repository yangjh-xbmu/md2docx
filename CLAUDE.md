# CLAUDE.md

## 项目定位

md2docx 是一个 Markdown → Word (.docx) 转换工具，核心特点是两阶段管线：pandoc 做语法解析，python-docx 做格式精调。支持多模板、frontmatter 驱动、多文件合并。

## 常用命令

```bash
# 安装依赖
uv sync

# 运行
uv run md2docx convert input.md
uv run md2docx convert input.md --style academic

# 测试
uv run pytest
uv run pytest -m "not integration"   # 跳过需要 pandoc 的测试
uv run pytest --cov=md2docx

# 格式化
uv run ruff format src/ tests/
uv run ruff check src/ tests/ --fix
```

## 架构

两阶段管线：

```
Stage 1: md → pandoc (--reference-doc + --lua-filter) → rough.docx
Stage 2: rough.docx → python-docx 后处理 → final.docx
```

### 模块职责

| 模块 | 职责 |
|------|------|
| cli.py | Click CLI 入口，子命令：convert / merge / styles / check |
| convert.py | 转换编排：读 frontmatter → 解析 style → 调 pandoc → 后处理 |
| pandoc.py | pandoc 环境检测、命令构建、执行 |
| postprocess.py | python-docx 后处理器：TOC、标题编号、图片、页眉页脚、封面 |
| merge.py | 多文件合并（contents.yaml 或 glob） |
| style.py | 样式系统：加载 style.yaml、查找样式目录、列出可用样式 |
| frontmatter.py | YAML frontmatter 解析，与 style 配置合并 |

### 配置优先级

CLI 参数 > frontmatter > style.yaml > 内置默认值

### 样式目录查找顺序

1. `~/.md2docx/styles/`（用户自定义）
2. 包内 `styles/`（内置默认）

## 约束

- 外部依赖：pandoc 必须在 PATH 中
- 后处理失败时降级输出 pandoc 原始结果
- 样式配置用 YAML，不用 TOML
- 后处理器是独立函数，可单独启用/禁用

## 业务域清单

| 文件 | 功能 |
|------|------|
| src/md2docx/cli.py | CLI 入口 |
| src/md2docx/convert.py | 转换编排 |
| src/md2docx/pandoc.py | pandoc 封装 |
| src/md2docx/postprocess.py | 后处理引擎 |
| src/md2docx/merge.py | 多文件合并 |
| src/md2docx/style.py | 样式系统 |
| src/md2docx/frontmatter.py | frontmatter 解析 |
| styles/default/ | 内置默认样式 |
