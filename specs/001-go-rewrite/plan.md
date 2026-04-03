# Implementation Plan: md2docx Go Rewrite

**Branch**: `001-go-rewrite` | **Date**: 2026-04-03 | **Spec**: [spec.md](./spec.md)

## Summary

用 Go 重写 md2docx，实现单二进制 Markdown→Word 转换工具。goldmark 解析 Markdown AST，自建 OOXML writer 生成 docx。样式系统用 YAML 声明式定义，内嵌样式通过 go:embed。

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: goldmark (Markdown AST), cobra (CLI), gopkg.in/yaml.v3 (样式解析)
**Storage**: N/A (CLI 工具，无持久化)
**Testing**: go test + testify
**Target Platform**: macOS (arm64/amd64), Windows (amd64), Linux (amd64)
**Project Type**: CLI tool
**Performance Goals**: 转换 < 1 秒 (typical markdown file)
**Constraints**: 单二进制 < 20MB, 零运行时依赖
**Scale/Scope**: 单用户 CLI

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Single Binary | ✅ | Go 静态编译 + go:embed |
| II. Style-Driven | ✅ | YAML 样式系统 |
| III. AI-Native | ✅ | CLI 设计对 AI 友好 |
| IV. Zero Config Default | ✅ | 内嵌 default 样式 |
| V. OOXML Self-Owned | ✅ | encoding/xml + archive/zip |
| VI. Incremental | ✅ | 按 User Story 优先级逐步交付 |

## Project Structure

### Documentation

```text
specs/001-go-rewrite/
├── spec.md              # 需求规格
├── plan.md              # 本文件
└── tasks.md             # 任务列表
```

### Source Code

```text
cmd/
└── md2docx/
    └── main.go                  # 入口

internal/
├── cli/
│   ├── root.go                  # cobra root，默认行为 = convert
│   ├── convert.go               # convert 命令
│   ├── merge.go                 # merge 命令
│   └── styles.go                # styles list/show 命令
├── parser/
│   ├── markdown.go              # goldmark 配置 + AST 解析
│   └── frontmatter.go           # YAML frontmatter 提取
├── style/
│   ├── types.go                 # Style YAML 结构体定义
│   ├── loader.go                # 加载样式：embed → 用户目录
│   └── defaults.go              # 默认值填充
├── ooxml/
│   ├── writer.go                # 顶层：组装各 part → zip 输出
│   ├── document.go              # word/document.xml 生成
│   ├── styles.go                # word/styles.xml 生成
│   ├── numbering.go             # word/numbering.xml（标题编号+列表）
│   ├── header.go                # word/header1.xml
│   ├── footer.go                # word/footer1.xml
│   ├── image.go                 # 图片嵌入（word/media/ + relationships）
│   ├── table.go                 # 表格 XML 生成
│   ├── rels.go                  # .rels 关系文件
│   ├── content_types.go         # [Content_Types].xml
│   └── types.go                 # OOXML 元素类型定义
├── render/
│   ├── renderer.go              # goldmark AST → ooxml 调用（核心渲染器）
│   ├── inline.go                # 行内元素渲染（bold/italic/code/link）
│   ├── block.go                 # 块级元素渲染（paragraph/heading/list/quote）
│   ├── table.go                 # 表格渲染
│   ├── image.go                 # 图片渲染
│   ├── cover.go                 # 封面生成
│   ├── toc.go                   # 目录生成
│   └── heading_number.go        # 标题编号逻辑
└── merge/
    └── merge.go                 # 多文件合并（Markdown 级拼接）

styles/                          # go:embed 内嵌
├── default.yaml
├── academic-cn.yaml
└── simple.yaml

testdata/                        # 测试用 fixture
├── basic.md
├── cjk.md
├── with_frontmatter.md
├── with_table.md
├── with_image.md
├── golden/                      # golden file 测试用的期望 XML
│   ├── basic_document.xml
│   └── ...
└── images/
    └── test.png

tests/                           # 集成测试
└── e2e_test.go
```

## Key Design Decisions

### 1. OOXML Writer 架构

不使用任何第三方 docx 库。直接用 `encoding/xml` 生成 XML，`archive/zip` 打包。

一个 docx 文件的最小结构：

```
[Content_Types].xml
_rels/.rels
word/document.xml          # 文档主体
word/styles.xml            # 样式定义
word/numbering.xml         # 编号定义（标题编号+列表）
word/header1.xml           # 页眉（可选）
word/footer1.xml           # 页脚（可选）
word/_rels/document.xml.rels  # 文档关系
word/media/                # 嵌入图片（可选）
```

每个 part 由对应的 Go 文件生成。Writer 负责组装和 zip 打包。

### 2. CJK 双字体处理

OOXML 的 `<w:rFonts>` 支持 4 个属性：`w:ascii`（Latin）、`w:eastAsia`（CJK）、`w:hAnsi`（高位 Latin）、`w:cs`（复杂脚本）。

策略：在 styles.xml 中定义样式级别的字体。不需要逐 run 检测 CJK 字符，因为 Word 会根据字符的 Unicode range 自动选择 eastAsia 或 ascii 字体。只需在样式中同时设置 ascii 和 eastAsia 即可。

### 3. goldmark 渲染器

实现 goldmark 的 `renderer.NodeRenderer` 接口。遍历 AST 时，调用 ooxml 包的 API 构建文档结构：

```go
// 伪代码
func (r *DocxRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) {
    if entering {
        r.doc.StartParagraph(style)
    } else {
        r.doc.EndParagraph()
    }
}
```

### 4. 样式加载优先级

1. CLI `--style` 参数
2. Frontmatter `style:` 字段
3. 默认 "default"

样式查找顺序：
1. `~/.md2docx/styles/<name>.yaml`（用户自定义）
2. 内嵌样式（go:embed）

### 5. Frontmatter 覆盖

Frontmatter 中的特定字段可覆盖样式设置：

```yaml
---
title: 我的论文
toc: false           # 覆盖样式的 toc.enabled
heading_numbering: false
cover: false
page_size: Letter    # 覆盖样式的 page.size
---
```

## Complexity Tracking

无 Constitution 违规，不需要额外复杂度论证。
