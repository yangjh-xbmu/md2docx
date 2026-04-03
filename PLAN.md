# md2docx 加强版技术设计

## 1. 目标定义

**核心问题**：用 Markdown 写内容，输出格式精确、可直接提交的 Word 文档，零手动调格式。

**核心用户场景**：

| 场景 | 交付格式 | 模板要求 |
|------|---------|---------|
| 教学大纲/教案/教学日志 | .docx | 学校指定模板 |
| 国社科申报书 | .docx | 基金委模板 |
| 论文/书稿 | .docx 或 PDF | 期刊/出版社模板 |
| 自用资料 | PDF（pandoc → LaTeX） | 自定义 LaTeX 模板 |
| 长文档（多章节） | .docx | 合并多个 md 文件 |

**设计原则**：

- Markdown 是唯一源文件格式
- 一个 md 文件可出 docx 和 PDF 两种格式
- 模板驱动，不同场景切换模板即可
- pandoc 做语法解析重活，python-docx 做格式精调

## 2. 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| 语言 | Python 3.12+ | python-docx 生态成熟，与 course-toolkit 共享经验 |
| CLI 框架 | click | 轻量，子命令支持好 |
| Markdown 解析 | pandoc（外部依赖） | 业界标准，Markdown → docx AST 转换最强 |
| pandoc 调用 | subprocess | 直接调 CLI，不依赖 pypandoc |
| docx 后处理 | python-docx | Word XML 操作的标准库 |
| frontmatter 解析 | python-frontmatter | 解析 YAML frontmatter |
| 配置格式 | YAML | 样式配置、合并清单 |
| 包管理 | uv + pyproject.toml | 现代 Python 包管理 |
| 测试 | pytest | 标准选择 |

## 3. 目录结构

```
md2docx/
├── src/
│   └── md2docx/
│       ├── __init__.py          # 版本号、公开 API
│       ├── cli.py               # Click CLI 入口
│       ├── convert.py           # 核心转换：pandoc 调用 + 后处理编排
│       ├── pandoc.py            # pandoc 环境检测、命令构建、执行
│       ├── postprocess.py       # python-docx 后处理引擎
│       ├── merge.py             # 多文件合并
│       ├── style.py             # 样式系统（加载/解析/查找）
│       └── frontmatter.py       # YAML frontmatter 解析与合并
├── styles/
│   └── default/
│       ├── style.yaml           # 样式配置
│       ├── reference.docx       # pandoc reference-doc 模板
│       └── filters/
│           └── styles.lua       # Lua filter
├── tests/
│   ├── conftest.py
│   ├── test_convert.py
│   ├── test_postprocess.py
│   ├── test_merge.py
│   ├── test_style.py
│   └── fixtures/
│       ├── simple.md
│       ├── with_frontmatter.md
│       └── merge_project/
├── pyproject.toml
├── PLAN.md
├── CLAUDE.md
└── README.md
```

## 4. 数据模型

### ConvertConfig（单次转换的完整配置）

```python
@dataclass(frozen=True)
class ConvertConfig:
    source: Path                    # 输入 md 文件路径
    output: Path                    # 输出 docx 路径
    style_name: str                 # 样式名称
    style: StyleConfig              # 解析后的样式配置
    metadata: dict[str, str]        # frontmatter 元数据
```

### StyleConfig（样式配置）

```python
@dataclass(frozen=True)
class StyleConfig:
    name: str
    reference_doc: Path             # pandoc --reference-doc
    lua_filters: list[Path]         # pandoc --lua-filter
    post: PostProcessConfig         # 后处理配置
```

### PostProcessConfig（后处理配置）

```python
@dataclass(frozen=True)
class PostProcessConfig:
    toc: bool = False               # 插入目录
    toc_depth: int = 3              # 目录深度
    heading_numbering: bool = False # 标题自动编号
    image_width_pct: int = 80       # 图片宽度百分比
    header_left: str = ""           # 页眉左（支持 {title} 等变量）
    header_right: str = ""          # 页眉右
    footer_center: str = ""         # 页脚中
    cover: CoverConfig | None = None
```

### MergeConfig（合并配置）

```python
@dataclass(frozen=True)
class MergeConfig:
    files: list[MergeEntry]         # 待合并文件列表
    heading_offset: int = 0         # 标题层级偏移

@dataclass(frozen=True)
class MergeEntry:
    path: Path
    heading_offset: int = 0         # 单文件标题偏移
```

## 5. CLI 接口

```bash
# 单文件转换（最常用）
md2docx convert input.md
md2docx convert input.md -o output.docx
md2docx convert input.md --style academic
md2docx convert input.md --style academic --no-post  # 跳过后处理

# 多文件合并转换
md2docx merge contents.yaml -o merged.docx
md2docx merge chapters/*.md -o book.docx --style academic

# 样式管理
md2docx styles list                          # 列出可用样式
md2docx styles show default                  # 显示样式配置
md2docx styles init mystyle                  # 创建新样式骨架
md2docx styles dir                           # 打印样式目录路径

# 工具
md2docx check                               # 检查 pandoc 环境
```

### Frontmatter 覆盖

```yaml
---
title: 融合新闻产品策划与制作教学大纲
author: 杨志宏
date: 2026-04-03
style: academic
toc: true
heading_numbering: true
---
```

frontmatter 中的配置优先级高于 style.yaml，低于 CLI 参数。

优先级：CLI 参数 > frontmatter > style.yaml > 内置默认值

## 6. 错误处理策略

- pandoc 不在 PATH：明确提示安装方式（brew/apt/官网）
- 模板文件不存在：列出可用模板
- frontmatter 解析失败：警告并继续（降级为无 frontmatter）
- python-docx 后处理失败：警告并输出 pandoc 原始结果（降级策略）
- 合并时文件不存在：报错并列出缺失文件，不生成部分结果
- 所有错误通过 click.echo + sys.exit(1)，不用 logging 模块

## 7. 测试策略

- 框架：pytest
- 测试分层：
  - 单元测试：frontmatter 解析、style 加载、merge 文件列表解析、postprocess 各处理器
  - 集成测试：需要 pandoc 的完整转换流程（用 `@pytest.mark.integration` 标记）
- fixtures 放在 `tests/fixtures/`
- 覆盖率目标：80%+

```bash
# 运行全部测试
pytest

# 跳过需要 pandoc 的集成测试
pytest -m "not integration"

# 带覆盖率
pytest --cov=md2docx
```

## 8. 两阶段管线详解

```
Stage 1: pandoc
  md + reference.docx + lua_filters → rough.docx

Stage 2: python-docx 后处理
  rough.docx + PostProcessConfig + metadata → final.docx

后处理器（按顺序执行）：
  1. cover_page    — 插入封面页（如果配置了）
  2. toc           — 插入/更新目录域代码
  3. heading_num   — 给标题加编号前缀
  4. image_resize  — 归一化图片宽度
  5. header_footer — 设置页眉页脚（支持变量替换）
```

每个后处理器是独立函数，接收 `Document` + `PostProcessConfig` + `metadata`，返回 `Document`。可单独启用/禁用。
