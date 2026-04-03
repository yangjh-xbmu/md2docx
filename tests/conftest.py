"""Shared test fixtures."""

from __future__ import annotations

from pathlib import Path

import pytest

FIXTURES_DIR = Path(__file__).parent / "fixtures"


@pytest.fixture
def fixtures_dir() -> Path:
    return FIXTURES_DIR


@pytest.fixture
def simple_md() -> Path:
    return FIXTURES_DIR / "simple.md"


@pytest.fixture
def frontmatter_md() -> Path:
    return FIXTURES_DIR / "with_frontmatter.md"


@pytest.fixture
def merge_dir() -> Path:
    return FIXTURES_DIR / "merge_project"


@pytest.fixture
def tmp_output(tmp_path: Path) -> Path:
    return tmp_path / "output.docx"
