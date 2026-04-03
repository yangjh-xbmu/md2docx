"""Multi-file merge: combine multiple markdown files into one before conversion."""

from __future__ import annotations

import glob as glob_mod
from dataclasses import dataclass
from pathlib import Path

import yaml


@dataclass(frozen=True)
class MergeEntry:
    path: Path
    heading_offset: int = 0


def parse_contents_yaml(contents_path: Path) -> list[MergeEntry]:
    """Parse a contents.yaml file into a list of MergeEntry."""
    raw = yaml.safe_load(contents_path.read_text(encoding="utf-8"))
    if not isinstance(raw, dict) or "files" not in raw:
        raise ValueError("contents.yaml 格式错误: 需要顶层 'files' 键")

    base_dir = contents_path.parent
    global_offset = raw.get("heading_offset", 0)
    entries: list[MergeEntry] = []

    for item in raw["files"]:
        if isinstance(item, str):
            path = base_dir / item
            entries.append(MergeEntry(path=path, heading_offset=global_offset))
        elif isinstance(item, dict):
            file_path = base_dir / item["path"]
            offset = item.get("heading_offset", global_offset)
            entries.append(MergeEntry(path=file_path, heading_offset=offset))
        else:
            raise ValueError(f"contents.yaml 中的文件条目格式错误: {item}")

    return entries


def resolve_glob_patterns(patterns: list[str], base_dir: Path) -> list[MergeEntry]:
    """Resolve glob patterns into MergeEntry list."""
    entries: list[MergeEntry] = []
    seen: set[Path] = set()

    for pattern in patterns:
        if base_dir and not Path(pattern).is_absolute():
            full_pattern = str(base_dir / pattern)
        else:
            full_pattern = pattern

        matches = sorted(glob_mod.glob(full_pattern))
        for match in matches:
            p = Path(match).resolve()
            if p not in seen and p.suffix == ".md":
                seen.add(p)
                entries.append(MergeEntry(path=p))

    return entries


def _offset_headings(text: str, offset: int) -> str:
    """Add heading level offset by prepending '#' characters."""
    if offset <= 0:
        return text

    lines = text.split("\n")
    result = []
    for line in lines:
        stripped = line.lstrip()
        if stripped.startswith("#"):
            # Count existing heading level
            hashes = 0
            for ch in stripped:
                if ch == "#":
                    hashes += 1
                else:
                    break
            if hashes <= 6 and (len(stripped) == hashes or stripped[hashes] == " "):
                new_level = min(hashes + offset, 6)
                result.append("#" * new_level + stripped[hashes:])
                continue
        result.append(line)

    return "\n".join(result)


def merge_files(entries: list[MergeEntry]) -> str:
    """Merge multiple markdown files into a single markdown string."""
    parts: list[str] = []
    missing: list[str] = []

    for entry in entries:
        if not entry.path.exists():
            missing.append(str(entry.path))
            continue

        content = entry.path.read_text(encoding="utf-8").rstrip()
        if entry.heading_offset > 0:
            content = _offset_headings(content, entry.heading_offset)
        parts.append(content)

    if missing:
        raise FileNotFoundError("以下文件不存在:\n" + "\n".join(f"  - {m}" for m in missing))

    return "\n\n".join(parts) + "\n"
