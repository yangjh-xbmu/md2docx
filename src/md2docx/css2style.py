"""CSS to Word style converter: parse CSS and generate reference.docx with mapped styles."""

from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path

import tinycss2
from docx import Document
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.shared import Cm, Pt, RGBColor
from docx.oxml import OxmlElement
from docx.oxml.ns import qn


# Mapping from CSS selectors to Word built-in style names
_SELECTOR_TO_WORD_STYLE: dict[str, str] = {
    "h1": "Heading 1",
    "h2": "Heading 2",
    "h3": "Heading 3",
    "h4": "Heading 4",
    "h5": "Heading 5",
    "h6": "Heading 6",
    "p": "Normal",
    "blockquote": "Quote",
    "code": "No Spacing",
    "pre": "No Spacing",
    "li": "List Paragraph",
    "a": "Default Paragraph Font",
    "table": "Table Grid",
}


@dataclass
class ParsedStyle:
    """Parsed CSS properties for a single selector."""

    selector: str
    word_style_name: str
    font_family: str | None = None
    font_family_eastasia: str | None = None
    font_size: Pt | None = None
    bold: bool | None = None
    italic: bool | None = None
    underline: bool | None = None
    color: RGBColor | None = None
    text_align: WD_ALIGN_PARAGRAPH | None = None
    line_height: float | None = None
    space_before: Pt | None = None
    space_after: Pt | None = None
    first_line_indent: Cm | None = None
    margin_left: Cm | None = None
    margin_right: Cm | None = None
    keep_with_next: bool | None = None
    page_break_before: bool | None = None


def parse_css(css_text: str) -> list[ParsedStyle]:
    """Parse CSS text and return a list of ParsedStyle objects."""
    rules = tinycss2.parse_stylesheet(css_text, skip_whitespace=True, skip_comments=True)
    styles: list[ParsedStyle] = []

    for rule in rules:
        if rule.type != "qualified-rule":
            continue

        selector = tinycss2.serialize(rule.prelude).strip()
        declarations = tinycss2.parse_declaration_list(rule.content, skip_whitespace=True)

        # Determine Word style name
        word_name = _resolve_word_style_name(selector)

        parsed = ParsedStyle(selector=selector, word_style_name=word_name)

        for decl in declarations:
            if decl.type != "declaration":
                continue
            _apply_declaration(parsed, decl.name, decl.value)

        styles.append(parsed)

    return styles


def _resolve_word_style_name(selector: str) -> str:
    """Map a CSS selector to a Word style name."""
    sel = selector.strip()

    # Class selector: .abstract → "abstract"
    if sel.startswith("."):
        return sel[1:].strip()

    # Element selector: h1 → "Heading 1"
    tag = sel.split()[0].split(":")[0].split(".")[0].split("#")[0].lower()
    return _SELECTOR_TO_WORD_STYLE.get(tag, sel)


def _parse_color(value_tokens: list) -> RGBColor | None:
    """Parse color from CSS tokens."""
    text = tinycss2.serialize(value_tokens).strip()

    # #rrggbb or #rgb
    if text.startswith("#"):
        hex_str = text[1:]
        if len(hex_str) == 3:
            hex_str = "".join(c * 2 for c in hex_str)
        if len(hex_str) == 6:
            try:
                r = int(hex_str[0:2], 16)
                g = int(hex_str[2:4], 16)
                b = int(hex_str[4:6], 16)
                return RGBColor(r, g, b)
            except ValueError:
                return None

    # rgb(r, g, b)
    match = re.match(r"rgb\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)", text)
    if match:
        return RGBColor(int(match.group(1)), int(match.group(2)), int(match.group(3)))

    # Named colors (common ones)
    named = {
        "black": RGBColor(0, 0, 0),
        "white": RGBColor(255, 255, 255),
        "red": RGBColor(255, 0, 0),
        "blue": RGBColor(0, 0, 255),
        "green": RGBColor(0, 128, 0),
        "gray": RGBColor(128, 128, 128),
        "grey": RGBColor(128, 128, 128),
    }
    return named.get(text.lower())


def _parse_length(value_tokens: list) -> tuple[float, str] | None:
    """Parse a CSS length value. Returns (number, unit) or None."""
    text = tinycss2.serialize(value_tokens).strip()
    match = re.match(r"([\d.]+)\s*(pt|px|cm|mm|em|in|rem)", text, re.IGNORECASE)
    if match:
        return float(match.group(1)), match.group(2).lower()
    return None


def _length_to_pt(value: float, unit: str) -> float:
    """Convert a CSS length to points."""
    conversions = {
        "pt": 1.0,
        "px": 0.75,  # 96dpi → 72pt/inch
        "cm": 28.3465,
        "mm": 2.83465,
        "in": 72.0,
        "em": 12.0,  # Assume 12pt base
        "rem": 12.0,
    }
    return value * conversions.get(unit, 1.0)


def _length_to_cm(value: float, unit: str) -> float:
    """Convert a CSS length to centimeters."""
    pt = _length_to_pt(value, unit)
    return pt / 28.3465


def _apply_declaration(parsed: ParsedStyle, prop: str, value_tokens: list) -> None:
    """Apply a single CSS declaration to a ParsedStyle."""
    text = tinycss2.serialize(value_tokens).strip()

    match prop:
        case "font-family":
            families = [f.strip().strip("'\"") for f in text.split(",")]
            # First CJK font goes to eastasia, first Latin font to main
            for f in families:
                if _is_cjk_font(f):
                    if parsed.font_family_eastasia is None:
                        parsed.font_family_eastasia = f
                else:
                    if parsed.font_family is None:
                        parsed.font_family = f

        case "font-size":
            length = _parse_length(value_tokens)
            if length:
                parsed.font_size = Pt(_length_to_pt(length[0], length[1]))

        case "font-weight":
            if text in ("bold", "bolder", "700", "800", "900"):
                parsed.bold = True
            elif text in ("normal", "400"):
                parsed.bold = False

        case "font-style":
            parsed.italic = text == "italic"

        case "text-decoration":
            if "underline" in text:
                parsed.underline = True

        case "color":
            parsed.color = _parse_color(value_tokens)

        case "text-align":
            align_map = {
                "left": WD_ALIGN_PARAGRAPH.LEFT,
                "center": WD_ALIGN_PARAGRAPH.CENTER,
                "right": WD_ALIGN_PARAGRAPH.RIGHT,
                "justify": WD_ALIGN_PARAGRAPH.JUSTIFY,
            }
            parsed.text_align = align_map.get(text.lower())

        case "line-height":
            if text.replace(".", "").isdigit():
                parsed.line_height = float(text)
            else:
                length = _parse_length(value_tokens)
                if length and length[1] == "pt":
                    # Absolute line height → convert to multiple of font size
                    parsed.line_height = length[0] / 12.0

        case "margin-top":
            length = _parse_length(value_tokens)
            if length:
                parsed.space_before = Pt(_length_to_pt(length[0], length[1]))

        case "margin-bottom":
            length = _parse_length(value_tokens)
            if length:
                parsed.space_after = Pt(_length_to_pt(length[0], length[1]))

        case "margin-left":
            length = _parse_length(value_tokens)
            if length:
                parsed.margin_left = Cm(_length_to_cm(length[0], length[1]))

        case "margin-right":
            length = _parse_length(value_tokens)
            if length:
                parsed.margin_right = Cm(_length_to_cm(length[0], length[1]))

        case "text-indent":
            length = _parse_length(value_tokens)
            if length:
                parsed.first_line_indent = Cm(_length_to_cm(length[0], length[1]))

        case "page-break-before":
            parsed.page_break_before = text == "always"

        case "page-break-after" | "break-after":
            if text == "page" or text == "always":
                parsed.keep_with_next = False

        case "orphans" | "widows":
            pass  # Recognized but not mapped

        # Shorthand: margin
        case "margin":
            parts = text.split()
            if len(parts) >= 1:
                length = _parse_length(
                    tinycss2.parse_component_value_list(parts[0], skip_whitespace=True)
                )
                if length:
                    parsed.space_before = Pt(_length_to_pt(length[0], length[1]))
            if len(parts) >= 2:
                # right margin
                pass
            if len(parts) >= 3:
                length = _parse_length(
                    tinycss2.parse_component_value_list(parts[2], skip_whitespace=True)
                )
                if length:
                    parsed.space_after = Pt(_length_to_pt(length[0], length[1]))


def _is_cjk_font(name: str) -> bool:
    """Check if a font name is a CJK font."""
    cjk_indicators = [
        "SimSun",
        "SimHei",
        "FangSong",
        "KaiTi",
        "Microsoft YaHei",
        "Source Han",
        "Noto Sans CJK",
        "Noto Serif CJK",
        "宋体",
        "黑体",
        "仿宋",
        "楷体",
        "微软雅黑",
        "华文",
        "方正",
        "思源",
        "PingFang",
        "Hiragino",
        "STSong",
        "STHeiti",
        "STKaiti",
        "STFangsong",
    ]
    return any(indicator in name for indicator in cjk_indicators)


def build_reference_docx(
    styles: list[ParsedStyle],
    base_doc: Path | None = None,
    output: Path | None = None,
) -> Path:
    """Generate a reference.docx with Word styles derived from parsed CSS.

    Args:
        styles: List of ParsedStyle from parse_css().
        base_doc: Optional existing .docx to use as base (preserves its styles).
        output: Output path. Defaults to ./reference.docx.
    """
    if base_doc and base_doc.exists():
        doc = Document(str(base_doc))
    else:
        doc = Document()

    for parsed in styles:
        _apply_style_to_doc(doc, parsed)

    out = output or Path("reference.docx")
    out.parent.mkdir(parents=True, exist_ok=True)
    doc.save(str(out))
    return out


def _apply_style_to_doc(doc: Document, parsed: ParsedStyle) -> None:
    """Apply a ParsedStyle to a Word document's style definitions."""
    style_name = parsed.word_style_name

    # Find or create the style
    style = None
    for s in doc.styles:
        if s.name == style_name:
            style = s
            break

    if style is None:
        # Create new paragraph style
        from docx.enum.style import WD_STYLE_TYPE

        try:
            style = doc.styles.add_style(style_name, WD_STYLE_TYPE.PARAGRAPH)
        except ValueError:
            # Style already exists with different casing
            return

    # Apply font properties
    font = style.font
    if parsed.font_family is not None:
        font.name = parsed.font_family

    if parsed.font_family_eastasia is not None:
        # Set East Asian font via XML
        rPr = style.element.get_or_add_rPr()
        rFonts = rPr.find(qn("w:rFonts"))
        if rFonts is None:
            rFonts = OxmlElement("w:rFonts")
            rPr.insert(0, rFonts)
        rFonts.set(qn("w:eastAsia"), parsed.font_family_eastasia)

    if parsed.font_size is not None:
        font.size = parsed.font_size

    if parsed.bold is not None:
        font.bold = parsed.bold

    if parsed.italic is not None:
        font.italic = parsed.italic

    if parsed.underline is not None:
        font.underline = parsed.underline

    if parsed.color is not None:
        font.color.rgb = parsed.color

    # Apply paragraph format properties
    pf = style.paragraph_format

    if parsed.text_align is not None:
        pf.alignment = parsed.text_align

    if parsed.line_height is not None:
        pf.line_spacing = parsed.line_height

    if parsed.space_before is not None:
        pf.space_before = parsed.space_before

    if parsed.space_after is not None:
        pf.space_after = parsed.space_after

    if parsed.first_line_indent is not None:
        pf.first_line_indent = parsed.first_line_indent

    if parsed.margin_left is not None:
        pf.left_indent = parsed.margin_left

    if parsed.margin_right is not None:
        pf.right_indent = parsed.margin_right

    if parsed.keep_with_next is not None:
        pf.keep_with_next = parsed.keep_with_next

    if parsed.page_break_before is not None:
        pf.page_break_before = parsed.page_break_before


def css_to_reference_docx(
    css_path: Path,
    base_doc: Path | None = None,
    output: Path | None = None,
) -> Path:
    """High-level API: read CSS file, generate reference.docx.

    Args:
        css_path: Path to .css file.
        base_doc: Optional base .docx to extend.
        output: Output .docx path.

    Returns:
        Path to the generated reference.docx.
    """
    css_text = css_path.read_text(encoding="utf-8")
    styles = parse_css(css_text)

    if not styles:
        raise ValueError(f"CSS 文件中没有找到有效的样式规则: {css_path}")

    return build_reference_docx(styles, base_doc, output)
