"""Pandoc environment check, command building, and execution."""

from __future__ import annotations

import shutil
import subprocess
from pathlib import Path

from md2docx.style import StyleConfig


def check_pandoc() -> str | None:
    """Return pandoc version string if available, None otherwise."""
    pandoc = shutil.which("pandoc")
    if pandoc is None:
        return None
    try:
        result = subprocess.run(
            [pandoc, "--version"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        first_line = result.stdout.strip().split("\n")[0]
        return first_line
    except (subprocess.SubprocessError, IndexError):
        return None


def require_pandoc() -> str:
    """Ensure pandoc is available. Raises RuntimeError if not."""
    version = check_pandoc()
    if version is None:
        raise RuntimeError(
            "pandoc 未找到。请先安装:\n"
            "  macOS:  brew install pandoc\n"
            "  Linux:  sudo apt install pandoc\n"
            "  其他:   https://pandoc.org/installing.html"
        )
    return version


def build_pandoc_args(
    input_path: Path,
    output_path: Path,
    style: StyleConfig,
) -> list[str]:
    """Build the pandoc command arguments list."""
    args = ["pandoc", "-f", "markdown", str(input_path), "-o", str(output_path)]

    if style.reference_doc is not None:
        args.extend(["--reference-doc", str(style.reference_doc)])

    for lua in style.lua_filters:
        args.extend(["--lua-filter", str(lua)])

    return args


def run_pandoc(
    input_path: Path,
    output_path: Path,
    style: StyleConfig,
) -> Path:
    """Run pandoc to convert markdown to docx. Returns output path."""
    args = build_pandoc_args(input_path, output_path, style)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    result = subprocess.run(args, capture_output=True, text=True, timeout=120)
    if result.returncode != 0:
        raise RuntimeError(f"pandoc 转换失败:\n{result.stderr}")

    return output_path
