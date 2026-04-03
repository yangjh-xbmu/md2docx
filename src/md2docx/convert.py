"""Core conversion orchestration: frontmatter → style → pandoc → post-process."""

from __future__ import annotations

import tempfile
from pathlib import Path
from typing import Any

import click

from md2docx.frontmatter import parse_frontmatter, resolve_metadata
from md2docx.merge import merge_files, parse_contents_yaml, resolve_glob_patterns
from md2docx.pandoc import require_pandoc, run_pandoc
from md2docx.postprocess import postprocess
from md2docx.style import load_style, override_post_config


def convert(
    source: Path,
    output: Path | None = None,
    style_name: str | None = None,
    no_post: bool = False,
    cli_overrides: dict[str, Any] | None = None,
) -> Path:
    """Convert a single markdown file to docx.

    Returns the output path.
    """
    require_pandoc()

    md_text = source.read_text(encoding="utf-8")
    metadata, body = parse_frontmatter(md_text)
    metadata = resolve_metadata(metadata, cli_overrides)

    # Determine style: CLI > frontmatter > default
    effective_style = style_name or metadata.get("style", "default")
    style = load_style(effective_style)

    # Override post-processing config with frontmatter values
    post_config = override_post_config(style.post, metadata)

    # Determine output path
    if output is None:
        output = source.with_suffix(".docx")

    # Write body (without frontmatter) to temp file for pandoc
    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".md", delete=False, encoding="utf-8"
    ) as tmp:
        tmp.write(body)
        tmp_path = Path(tmp.name)

    try:
        if no_post:
            run_pandoc(tmp_path, output, style)
        else:
            # Stage 1: pandoc → rough docx
            rough_output = output.with_name(output.stem + "_rough.docx")
            run_pandoc(tmp_path, rough_output, style)

            # Stage 2: python-docx post-processing
            try:
                postprocess(rough_output, post_config, metadata, output)
                rough_output.unlink(missing_ok=True)
            except Exception as e:
                click.echo(f"警告: 后处理失败 ({e})，使用 pandoc 原始输出", err=True)
                rough_output.rename(output)
    finally:
        tmp_path.unlink(missing_ok=True)

    return output


def convert_merged(
    contents: Path | None = None,
    patterns: list[str] | None = None,
    output: Path | None = None,
    style_name: str | None = None,
    no_post: bool = False,
    cli_overrides: dict[str, Any] | None = None,
) -> Path:
    """Merge multiple markdown files and convert to docx.

    Either contents (path to contents.yaml) or patterns (glob patterns) must be provided.
    """
    require_pandoc()

    if contents is not None:
        entries = parse_contents_yaml(contents)
        base_dir = contents.parent
    elif patterns:
        base_dir = Path.cwd()
        entries = resolve_glob_patterns(patterns, base_dir)
    else:
        raise ValueError("需要提供 contents.yaml 路径或 glob 模式")

    if not entries:
        raise ValueError("没有找到要合并的文件")

    # Merge all files
    merged_md = merge_files(entries)

    # Parse frontmatter from merged content (first file's frontmatter wins)
    metadata, body = parse_frontmatter(merged_md)
    metadata = resolve_metadata(metadata, cli_overrides)

    effective_style = style_name or metadata.get("style", "default")

    if output is None:
        if contents is not None:
            output = contents.with_suffix(".docx")
        else:
            output = Path.cwd() / "merged.docx"

    # Write merged content to temp file
    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".md", delete=False, encoding="utf-8"
    ) as tmp:
        tmp.write(body)
        tmp_path = Path(tmp.name)

    try:
        style = load_style(effective_style)
        post_config = override_post_config(style.post, metadata)

        if no_post:
            run_pandoc(tmp_path, output, style)
        else:
            rough_output = output.with_name(output.stem + "_rough.docx")
            run_pandoc(tmp_path, rough_output, style)

            try:
                postprocess(rough_output, post_config, metadata, output)
                rough_output.unlink(missing_ok=True)
            except Exception as e:
                click.echo(f"警告: 后处理失败 ({e})，使用 pandoc 原始输出", err=True)
                rough_output.rename(output)
    finally:
        tmp_path.unlink(missing_ok=True)

    return output
