"""Tests for frontmatter parsing."""

from __future__ import annotations

from pathlib import Path

from md2docx.frontmatter import parse_frontmatter, resolve_metadata


class TestParseFrontmatter:
    def test_with_frontmatter(self, frontmatter_md: Path) -> None:
        meta, body = parse_frontmatter(frontmatter_md)
        assert meta["title"] == "测试文档"
        assert meta["author"] == "杨志宏"
        assert meta["style"] == "default"
        assert meta["toc"] is True
        assert "# 第一章 概述" in body

    def test_without_frontmatter(self, simple_md: Path) -> None:
        meta, body = parse_frontmatter(simple_md)
        assert meta == {}
        assert "# Hello World" in body

    def test_from_string(self) -> None:
        text = "---\ntitle: Test\n---\n\nBody content"
        meta, body = parse_frontmatter(text)
        assert meta["title"] == "Test"
        assert "Body content" in body


class TestResolveMetadata:
    def test_cli_overrides_frontmatter(self) -> None:
        fm = {"title": "Original", "author": "Author"}
        cli = {"title": "Override", "extra": "new"}
        result = resolve_metadata(fm, cli)
        assert result["title"] == "Override"
        assert result["author"] == "Author"
        assert result["extra"] == "new"

    def test_none_values_not_overridden(self) -> None:
        fm = {"title": "Keep"}
        cli = {"title": None}
        result = resolve_metadata(fm, cli)
        assert result["title"] == "Keep"

    def test_no_overrides(self) -> None:
        fm = {"title": "Keep"}
        result = resolve_metadata(fm)
        assert result["title"] == "Keep"
