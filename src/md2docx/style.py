"""Style system: load, resolve, and list available styles."""

from __future__ import annotations

from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

import yaml

# Style directories, searched in order (first match wins)
_PACKAGE_STYLES_DIR = Path(__file__).resolve().parent.parent.parent / "styles"
_USER_STYLES_DIR = Path.home() / ".md2docx" / "styles"


@dataclass(frozen=True)
class CoverConfig:
    enabled: bool = False
    fields: dict[str, str] = field(default_factory=dict)


@dataclass(frozen=True)
class PostProcessConfig:
    toc: bool = False
    toc_depth: int = 3
    heading_numbering: bool = False
    image_width_pct: int = 80
    header_left: str = ""
    header_right: str = ""
    footer_center: str = ""
    cover: CoverConfig = field(default_factory=CoverConfig)


@dataclass(frozen=True)
class StyleConfig:
    name: str
    reference_doc: Path | None
    lua_filters: list[Path]
    post: PostProcessConfig


def _style_dirs() -> list[Path]:
    """Return style directories in priority order."""
    dirs = []
    if _USER_STYLES_DIR.is_dir():
        dirs.append(_USER_STYLES_DIR)
    if _PACKAGE_STYLES_DIR.is_dir():
        dirs.append(_PACKAGE_STYLES_DIR)
    return dirs


def find_style_dir(name: str) -> Path | None:
    """Find the directory for a named style."""
    for base in _style_dirs():
        candidate = base / name
        if candidate.is_dir():
            return candidate
    return None


def load_style(name: str) -> StyleConfig:
    """Load a style by name. Raises FileNotFoundError if not found."""
    style_dir = find_style_dir(name)
    if style_dir is None:
        available = list_styles()
        raise FileNotFoundError(
            f"样式 '{name}' 不存在。可用样式: {', '.join(available) or '(无)'}"
        )

    config_path = style_dir / "style.yaml"
    raw: dict[str, Any] = {}
    if config_path.exists():
        raw = yaml.safe_load(config_path.read_text(encoding="utf-8")) or {}

    # Reference doc
    ref_doc_name = raw.get("reference_doc", "reference.docx")
    ref_doc = style_dir / ref_doc_name
    reference_doc = ref_doc if ref_doc.exists() else None

    # Lua filters
    filters_dir = style_dir / "filters"
    lua_filters: list[Path] = []
    if filters_dir.is_dir():
        lua_filters = sorted(filters_dir.glob("*.lua"))
    for extra in raw.get("lua_filters", []):
        p = style_dir / extra
        if p.exists() and p not in lua_filters:
            lua_filters.append(p)

    # Post-processing config
    post_raw = raw.get("post_processing", {})
    cover_raw = post_raw.pop("cover", None) or {}
    cover = CoverConfig(
        enabled=cover_raw.get("enabled", False),
        fields=cover_raw.get("fields", {}),
    )
    post = PostProcessConfig(
        toc=post_raw.get("toc", False),
        toc_depth=post_raw.get("toc_depth", 3),
        heading_numbering=post_raw.get("heading_numbering", False),
        image_width_pct=post_raw.get("image_width_pct", 80),
        header_left=post_raw.get("header_left", ""),
        header_right=post_raw.get("header_right", ""),
        footer_center=post_raw.get("footer_center", ""),
        cover=cover,
    )

    return StyleConfig(
        name=name,
        reference_doc=reference_doc,
        lua_filters=lua_filters,
        post=post,
    )


def list_styles() -> list[str]:
    """List all available style names."""
    seen: set[str] = set()
    styles: list[str] = []
    for base in _style_dirs():
        if not base.is_dir():
            continue
        for child in sorted(base.iterdir()):
            if child.is_dir() and child.name not in seen:
                seen.add(child.name)
                styles.append(child.name)
    return styles


def override_post_config(
    base: PostProcessConfig,
    overrides: dict[str, Any],
) -> PostProcessConfig:
    """Create a new PostProcessConfig with overrides from frontmatter/CLI."""
    fields = {
        "toc": overrides.get("toc", base.toc),
        "toc_depth": overrides.get("toc_depth", base.toc_depth),
        "heading_numbering": overrides.get("heading_numbering", base.heading_numbering),
        "image_width_pct": overrides.get("image_width_pct", base.image_width_pct),
        "header_left": overrides.get("header_left", base.header_left),
        "header_right": overrides.get("header_right", base.header_right),
        "footer_center": overrides.get("footer_center", base.footer_center),
        "cover": base.cover,
    }
    return PostProcessConfig(**fields)
