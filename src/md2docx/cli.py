"""Click CLI entry point for md2docx."""

from __future__ import annotations

from pathlib import Path

import click

from md2docx import __version__
from md2docx.convert import convert, convert_merged
from md2docx.pandoc import check_pandoc
from md2docx.style import list_styles, load_style


@click.group()
@click.version_option(version=__version__, prog_name="md2docx")
def main() -> None:
    """md2docx - Markdown to Word converter with template system."""


@main.command(name="convert")
@click.argument("source", type=click.Path(exists=True, path_type=Path))
@click.option("-o", "--output", type=click.Path(path_type=Path), help="输出 .docx 路径")
@click.option("-s", "--style", "style_name", default=None, help="样式名称")
@click.option("--no-post", is_flag=True, help="跳过后处理，只用 pandoc 原始输出")
def convert_cmd(
    source: Path,
    output: Path | None,
    style_name: str | None,
    no_post: bool,
) -> None:
    """转换单个 Markdown 文件为 Word 文档。"""
    try:
        result = convert(
            source=source,
            output=output,
            style_name=style_name,
            no_post=no_post,
        )
        click.echo(f"已导出: {result}")
    except Exception as e:
        click.echo(f"错误: {e}", err=True)
        raise SystemExit(1)


@main.command(name="merge")
@click.argument("sources", nargs=-1, type=str)
@click.option("-o", "--output", type=click.Path(path_type=Path), help="输出 .docx 路径")
@click.option("-s", "--style", "style_name", default=None, help="样式名称")
@click.option("--no-post", is_flag=True, help="跳过后处理")
def merge_cmd(
    sources: tuple[str, ...],
    output: Path | None,
    style_name: str | None,
    no_post: bool,
) -> None:
    """合并多个 Markdown 文件并转换为 Word 文档。

    SOURCES 可以是 contents.yaml 路径，或多个 glob 模式。
    """
    if not sources:
        click.echo("错误: 需要提供 contents.yaml 或 glob 模式", err=True)
        raise SystemExit(1)

    try:
        # Check if first source is a contents.yaml
        first = Path(sources[0])
        if first.suffix in (".yaml", ".yml") and first.exists():
            result = convert_merged(
                contents=first,
                output=output,
                style_name=style_name,
                no_post=no_post,
            )
        else:
            result = convert_merged(
                patterns=list(sources),
                output=output,
                style_name=style_name,
                no_post=no_post,
            )
        click.echo(f"已导出: {result}")
    except Exception as e:
        click.echo(f"错误: {e}", err=True)
        raise SystemExit(1)


@main.group(name="styles")
def styles_group() -> None:
    """样式管理。"""


@styles_group.command(name="list")
def styles_list_cmd() -> None:
    """列出所有可用样式。"""
    available = list_styles()
    if not available:
        click.echo("没有找到任何样式。")
        return

    for name in available:
        try:
            style = load_style(name)
            ref = "有模板" if style.reference_doc else "无模板"
            filters = len(style.lua_filters)
            click.echo(f"  {name:20s} {ref}, {filters} 个 Lua filter")
        except Exception:
            click.echo(f"  {name:20s} (加载失败)")


@styles_group.command(name="show")
@click.argument("name")
def styles_show_cmd(name: str) -> None:
    """显示指定样式的详细配置。"""
    try:
        style = load_style(name)
    except FileNotFoundError as e:
        click.echo(f"错误: {e}", err=True)
        raise SystemExit(1)

    click.echo(f"样式: {style.name}")
    click.echo(f"模板: {style.reference_doc or '(无)'}")
    click.echo(f"Lua filters: {[str(f) for f in style.lua_filters] or '(无)'}")
    click.echo("后处理:")
    click.echo(f"  目录: {'是' if style.post.toc else '否'} (深度 {style.post.toc_depth})")
    click.echo(f"  标题编号: {'是' if style.post.heading_numbering else '否'}")
    click.echo(f"  图片宽度: {style.post.image_width_pct}%")
    click.echo(f"  页眉左: {style.post.header_left or '(无)'}")
    click.echo(f"  页眉右: {style.post.header_right or '(无)'}")
    click.echo(f"  页脚中: {style.post.footer_center or '(无)'}")
    click.echo(f"  封面: {'是' if style.post.cover.enabled else '否'}")


@styles_group.command(name="dir")
def styles_dir_cmd() -> None:
    """打印样式目录路径。"""
    user_dir = Path.home() / ".md2docx" / "styles"
    pkg_dir = Path(__file__).resolve().parent.parent.parent / "styles"
    click.echo(f"用户样式目录: {user_dir}")
    click.echo(f"内置样式目录: {pkg_dir}")


@styles_group.command(name="init")
@click.argument("name")
def styles_init_cmd(name: str) -> None:
    """创建新样式骨架。"""
    user_dir = Path.home() / ".md2docx" / "styles" / name
    if user_dir.exists():
        click.echo(f"错误: 样式 '{name}' 已存在于 {user_dir}", err=True)
        raise SystemExit(1)

    user_dir.mkdir(parents=True)
    filters_dir = user_dir / "filters"
    filters_dir.mkdir()

    style_yaml = user_dir / "style.yaml"
    style_yaml.write_text(
        f"""name: {name}
reference_doc: reference.docx

post_processing:
  toc: false
  toc_depth: 3
  heading_numbering: false
  image_width_pct: 80
  header_left: ""
  header_right: ""
  footer_center: ""
  cover:
    enabled: false
    fields:
      title: "{{title}}"
      author: "{{author}}"
      date: "{{date}}"
""",
        encoding="utf-8",
    )

    click.echo(f"已创建样式骨架: {user_dir}")
    click.echo("请将 reference.docx 模板文件放入该目录。")


@main.command(name="check")
def check_cmd() -> None:
    """检查 pandoc 环境。"""
    version = check_pandoc()
    if version:
        click.echo(f"pandoc 可用: {version}")
    else:
        click.echo("pandoc 未找到。请安装: https://pandoc.org/installing.html", err=True)
        raise SystemExit(1)
