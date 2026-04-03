# md2docx

Markdown to Word (.docx) converter. Single binary, zero dependencies.

Supports CJK dual fonts, table of contents, heading numbering, cover page, headers/footers, GFM tables, image embedding, and multiple built-in styles.

## Install

Download the binary for your platform from [Releases](https://github.com/yangjh-xbmu/md2docx/releases), or build from source:

```bash
go install github.com/yangjh-xbmu/md2docx/cmd/md2docx@latest
```

## Quick Start

```bash
# Basic conversion (output to ~/Desktop/)
md2docx input.md

# Specify output path
md2docx input.md -o output.docx

# Use academic style
md2docx input.md --style academic-cn

# Skip cover page and TOC
md2docx input.md --no-cover --no-toc

# Merge multiple files
md2docx merge ch1.md ch2.md ch3.md -o book.docx

# List available styles
md2docx styles list
```

## Styles

Built-in styles:

| Style | Description |
|-------|-------------|
| `default` | Default style with Song/Hei fonts, A4, 1.5x line spacing |
| `academic-cn` | Chinese academic paper format (Song 12pt, Hei headings, first-line indent) |
| `simple` | Minimal style, no cover, no TOC |

Custom styles can be placed in `~/.md2docx/styles/<name>.yaml`.

View a style's full configuration:

```bash
md2docx styles show academic-cn
```

## Frontmatter

Use YAML frontmatter to set metadata and override style settings:

```yaml
---
title: My Paper
author: Author Name
date: 2026-01-01
style: academic-cn
toc: false
cover: false
heading_numbering: false
---
```

## CLI Reference

```
md2docx [file.md]          Convert Markdown to Word
md2docx merge <files...>   Merge multiple files into one document
md2docx styles list         List available styles
md2docx styles show <name>  Show style configuration

Flags:
  -o, --output string    Output .docx path (default: ~/Desktop/<name>.docx)
  -s, --style string     Style name (default: "default")
      --no-cover         Skip cover page
      --no-toc           Skip table of contents
      --no-numbering     Skip heading numbering
  -v, --version          Show version
```

## Limitations

- Math formulas and footnotes are not supported in v1
- Code block syntax highlighting is not supported
- Fonts must be installed on the system (no font embedding)
- TOC requires updating fields after opening in Word (Ctrl+A, F9)

## License

MIT
