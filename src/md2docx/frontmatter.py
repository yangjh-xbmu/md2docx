"""YAML frontmatter parsing and metadata extraction."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import frontmatter


def parse_frontmatter(source: str | Path) -> tuple[dict[str, Any], str]:
    """Parse YAML frontmatter from markdown text or file.

    Returns (metadata_dict, body_text_without_frontmatter).
    """
    if isinstance(source, Path):
        source = source.read_text(encoding="utf-8")

    post = frontmatter.loads(source)
    return dict(post.metadata), post.content


def resolve_metadata(
    fm_meta: dict[str, Any],
    cli_overrides: dict[str, Any] | None = None,
) -> dict[str, Any]:
    """Merge metadata: CLI overrides > frontmatter."""
    merged = dict(fm_meta)
    if cli_overrides:
        merged.update({k: v for k, v in cli_overrides.items() if v is not None})
    return merged
