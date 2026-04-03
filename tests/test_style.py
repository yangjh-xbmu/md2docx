"""Tests for style system."""

from __future__ import annotations

import pytest

from md2docx.style import (
    PostProcessConfig,
    list_styles,
    load_style,
    override_post_config,
)


class TestLoadStyle:
    def test_load_default(self) -> None:
        style = load_style("default")
        assert style.name == "default"
        assert style.reference_doc is not None
        assert style.reference_doc.exists()

    def test_load_nonexistent(self) -> None:
        with pytest.raises(FileNotFoundError, match="不存在"):
            load_style("nonexistent_style_xyz")

    def test_default_has_lua_filter(self) -> None:
        style = load_style("default")
        lua_names = [f.name for f in style.lua_filters]
        assert "styles.lua" in lua_names


class TestListStyles:
    def test_includes_default(self) -> None:
        styles = list_styles()
        assert "default" in styles


class TestOverridePostConfig:
    def test_override_toc(self) -> None:
        base = PostProcessConfig()
        assert base.toc is False
        result = override_post_config(base, {"toc": True})
        assert result.toc is True

    def test_override_preserves_unset(self) -> None:
        base = PostProcessConfig(toc=True, heading_numbering=True)
        result = override_post_config(base, {"toc": False})
        assert result.toc is False
        assert result.heading_numbering is True

    def test_empty_overrides(self) -> None:
        base = PostProcessConfig(toc=True)
        result = override_post_config(base, {})
        assert result.toc is True
