"""Tests for pandoc wrapper."""

from __future__ import annotations

from pathlib import Path


from md2docx.pandoc import build_pandoc_args, check_pandoc
from md2docx.style import StyleConfig, PostProcessConfig


class TestCheckPandoc:
    def test_pandoc_available(self) -> None:
        version = check_pandoc()
        # pandoc should be installed on dev machine
        assert version is not None
        assert "pandoc" in version.lower()


class TestBuildPandocArgs:
    def test_basic_args(self) -> None:
        style = StyleConfig(
            name="test",
            reference_doc=None,
            lua_filters=[],
            post=PostProcessConfig(),
        )
        args = build_pandoc_args(Path("/tmp/in.md"), Path("/tmp/out.docx"), style)
        assert args[0] == "pandoc"
        assert "-f" in args
        assert "markdown" in args
        assert "/tmp/in.md" in args
        assert "-o" in args
        assert "/tmp/out.docx" in args

    def test_with_reference_doc(self) -> None:
        style = StyleConfig(
            name="test",
            reference_doc=Path("/tmp/ref.docx"),
            lua_filters=[],
            post=PostProcessConfig(),
        )
        args = build_pandoc_args(Path("/tmp/in.md"), Path("/tmp/out.docx"), style)
        assert "--reference-doc" in args
        assert "/tmp/ref.docx" in args

    def test_with_lua_filters(self) -> None:
        style = StyleConfig(
            name="test",
            reference_doc=None,
            lua_filters=[Path("/tmp/a.lua"), Path("/tmp/b.lua")],
            post=PostProcessConfig(),
        )
        args = build_pandoc_args(Path("/tmp/in.md"), Path("/tmp/out.docx"), style)
        lua_indices = [i for i, a in enumerate(args) if a == "--lua-filter"]
        assert len(lua_indices) == 2
