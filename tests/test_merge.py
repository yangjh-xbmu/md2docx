"""Tests for multi-file merge."""

from __future__ import annotations

from pathlib import Path

import pytest

from md2docx.merge import (
    MergeEntry,
    merge_files,
    parse_contents_yaml,
    _offset_headings,
)


class TestParseContentsYaml:
    def test_parse_basic(self, merge_dir: Path) -> None:
        entries = parse_contents_yaml(merge_dir / "contents.yaml")
        assert len(entries) == 3
        assert entries[0].path == merge_dir / "chapter1.md"
        assert entries[0].heading_offset == 0
        assert entries[2].heading_offset == 1

    def test_missing_files_key(self, tmp_path: Path) -> None:
        bad = tmp_path / "bad.yaml"
        bad.write_text("foo: bar", encoding="utf-8")
        with pytest.raises(ValueError, match="files"):
            parse_contents_yaml(bad)


class TestOffsetHeadings:
    def test_offset_one(self) -> None:
        text = "# Title\n\nSome text\n\n## Sub"
        result = _offset_headings(text, 1)
        assert result.startswith("## Title")
        assert "### Sub" in result

    def test_offset_zero(self) -> None:
        text = "# Title"
        assert _offset_headings(text, 0) == text

    def test_max_level_six(self) -> None:
        text = "###### Deep"
        result = _offset_headings(text, 2)
        assert result == "###### Deep"  # Capped at 6


class TestMergeFiles:
    def test_merge_basic(self, merge_dir: Path) -> None:
        entries = [
            MergeEntry(path=merge_dir / "chapter1.md"),
            MergeEntry(path=merge_dir / "chapter2.md"),
        ]
        result = merge_files(entries)
        assert "第一章 引言" in result
        assert "第二章 方法" in result

    def test_merge_with_offset(self, merge_dir: Path) -> None:
        entries = [
            MergeEntry(path=merge_dir / "chapter3.md", heading_offset=1),
        ]
        result = merge_files(entries)
        assert result.startswith("## 第三章 结论")

    def test_missing_file(self, merge_dir: Path) -> None:
        entries = [
            MergeEntry(path=merge_dir / "nonexistent.md"),
        ]
        with pytest.raises(FileNotFoundError, match="不存在"):
            merge_files(entries)
