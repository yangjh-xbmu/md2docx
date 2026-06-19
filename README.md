---
purpose: 将 Markdown 文件转换为 Word (.docx) 格式，提供单一二进制文件、零依赖的转换工具。
status: active
next_steps: []
capabilities:
  - markdown-to-docx
  - single-binary
  - zero-dependency
---
# md2docx

Markdown to Word (.docx) converter. Single binary, zero dependencies.

Powered by [goldmark](https://github.com/yuin/goldmark) (Markdown parsing) + custom OOXML writer (docx generation). No pandoc, no LibreOffice, no external tools needed.

## Features

- CJK dual-font support (auto ascii/eastAsia font selection)
- Table of contents with field codes
- Multi-level heading numbering (Chinese format support)
- Cover page with declarative layout
- Headers and footers with variable substitution
- GFM tables with borders, header styling, inline formatting
- Image embedding with aspect-ratio-preserving scaling
- Multi-file merge (direct files, glob patterns, contents.yaml)
- Multiple built-in styles, custom styles via YAML
- YAML frontmatter for per-document configuration

## Install

Download the binary for your platform from [Releases](https://github.com/yangjh-xbmu/md2docx/releases).

Or build from source:

```bash
go install github.com/yangjh-xbmu/md2docx/cmd/md2docx@latest
```

## Quick Start

```bash
# Basic conversion (output to ~/Desktop/)
md2docx input.md

# Specify output path
md2docx input.md -o output.docx

# Use academic style with cover, TOC, heading numbering
md2docx input.md --style academic-cn

# Skip cover page and TOC
md2docx input.md --no-cover --no-toc

# Merge multiple files into one document
md2docx merge ch1.md ch2.md ch3.md -o book.docx

# List available styles
md2docx styles list

# Show style configuration details
md2docx styles show academic-cn
```

## Example

Given a Markdown file `paper.md`:

```markdown
---
title: Research on News Production
author: John Doe
date: 2026-04-03
style: academic-cn
toc: true
cover: true
heading_numbering: true
---

# Introduction

## Background

Digital technology has fundamentally altered how news is produced,
distributed, and consumed. The concept of **convergent media** refers
to deep integration across organization, workflow, technology, and
distribution channels.

## Research Objectives

1. Identify key transition points from traditional to converged newsrooms
2. Analyze the effectiveness of the "central kitchen" model
3. Assess the impact of AI on news editing workflows

# Literature Review

| Research Area | Key Finding |
|---------------|-------------|
| Organization  | Flat management enables convergence |
| Technology    | AI reshapes editorial judgment |
| User Participation | UGC/PGC boundaries are blurring |

> Media convergence is not simple addition, but a deep, chemistry-like
> integration that requires rethinking every step of news production.

# Methods

The study uses a `grounded theory` approach:

    Open coding → Axial coding → Selective coding

Code example for analysis pipeline:

​```python
def analyze_news_flow(articles):
    stages = ["collect", "edit", "review", "publish", "feedback"]
    for article in articles:
        for stage in stages:
            duration = article.get_stage_duration(stage)
            yield stage, duration
​```

# Conclusion

1. Convergence is a gradual process
2. Technology drives change, but **human judgment** remains decisive
3. Small media should pursue differentiated convergence strategies
```

Convert it:

```bash
md2docx paper.md
```

This generates a Word document with:

- A cover page showing title, author, and date
- An auto-generated table of contents
- Numbered headings (1, 1.1, 1.2, 2, 2.1 ...)
- Formatted tables with borders and header styling
- Block quotes with left border and indent
- Code blocks with monospace font and background
- Page numbers in the footer

## Styles

Three built-in styles:

| Style | Description | Cover | TOC | Numbering |
|-------|-------------|:-----:|:---:|:---------:|
| `default` | General-purpose, Song body / Hei headings, A4, 1.5x spacing | - | - | - |
| `academic-cn` | Chinese academic paper, first-line indent, formal layout | yes | yes | yes |
| `simple` | Minimal, no decoration, compact spacing | - | - | - |

### Custom Styles

Create `~/.md2docx/styles/<name>.yaml` to define your own style. View built-in styles as reference:

```bash
md2docx styles show default
```

## Frontmatter

Use YAML frontmatter to set metadata and override style settings:

```yaml
---
title: My Document
author: Author Name
date: 2026-01-01
style: academic-cn      # style template
toc: true               # table of contents
cover: true             # cover page
heading_numbering: true # heading numbers
---
```

## Multi-file Merge

Merge multiple Markdown files into a single document:

```bash
# Direct file list
md2docx merge ch1.md ch2.md ch3.md -o book.docx

# Glob pattern
md2docx merge "chapters/*.md" -o book.docx
```

Or use a `contents.yaml` for fine-grained control:

```yaml
files:
  - path: preface.md
  - path: ch1.md
    heading_offset: 1    # demote headings by 1 level
  - path: ch2.md
    heading_offset: 1
```

```bash
md2docx merge --config contents.yaml -o book.docx
```

## CLI Reference

```
md2docx [file.md]              Convert Markdown to Word
md2docx merge <files...>       Merge multiple files into one document
md2docx styles list            List available styles
md2docx styles show <name>     Show style configuration

Flags:
  -o, --output string    Output .docx path (default: ~/Desktop/<name>.docx)
  -s, --style string     Style name (default: "default")
      --no-cover         Skip cover page
      --no-toc           Skip table of contents
      --no-numbering     Skip heading numbering
  -v, --version          Show version
```

## Limitations

- Math formulas and footnotes are not yet supported
- Code block syntax highlighting is not supported
- Fonts must be installed on the system (no font embedding)
- TOC requires updating fields after opening in Word (Ctrl+A, F9)

## License

MIT
