# Feature Specification: md2docx Go Rewrite

**Feature Branch**: `001-go-rewrite`
**Created**: 2026-04-03
**Status**: Approved

## User Scenarios & Testing

### User Story 1 - Basic Conversion (Priority: P1) 🎯 MVP

零编码背景用户下载 md2docx 可执行文件，运行 `md2docx 论文.md`，在桌面得到一个排版合理的 Word 文档。

**Why this priority**: 这是产品的核心价值。没有这个，其他一切没有意义。

**Independent Test**: 用一个包含标题、段落、列表、粗体斜体的 Markdown 文件，运行命令，用 Word 打开输出文件验证格式正确。

**Acceptance Scenarios**:

1. **Given** 一个包含 h1/h2/h3、正文段落、有序/无序列表的 .md 文件, **When** 运行 `md2docx input.md`, **Then** 在 ~/Desktop/ 生成 input.docx，内容完整，标题/正文/列表样式正确区分
2. **Given** 一个包含 CJK 中文文本和 Latin 英文的 .md 文件, **When** 转换, **Then** 中文使用宋体、英文使用 Times New Roman，双字体正确分离
3. **Given** 一个包含粗体、斜体、行内代码、超链接的 .md 文件, **When** 转换, **Then** 行内格式正确保留
4. **Given** 运行 `md2docx input.md -o /tmp/output.docx`, **Then** 输出到指定路径

---

### User Story 2 - Style Selection (Priority: P2)

用户通过 `--style academic-cn` 参数选择内置样式，得到符合中文学术论文排版规范的 Word 文档。

**Why this priority**: 样式系统是产品差异化的核心。内置样式让用户零配置得到专业排版。

**Independent Test**: 用同一个 Markdown 文件，分别用 default 和 academic-cn 样式转换，对比两个 docx 的字体、字号、行距差异。

**Acceptance Scenarios**:

1. **Given** `md2docx input.md --style academic-cn`, **When** 转换, **Then** 标题使用黑体、正文宋体 12pt、首行缩进 2em、1.5 倍行距、A4 纸张
2. **Given** `md2docx input.md --style simple`, **When** 转换, **Then** 使用简洁风格，无封面无目录
3. **Given** `md2docx styles list`, **Then** 输出所有内置样式名称和简介
4. **Given** `md2docx input.md --style nonexistent`, **Then** 输出清晰的错误信息，列出可用样式

---

### User Story 3 - Page Setup & Document Structure (Priority: P3)

用户通过样式文件控制纸张大小、页边距、目录、标题编号、页眉页脚。

**Why this priority**: 完整的文档结构是学术/公文场景的刚需。

**Independent Test**: 使用带 TOC、标题编号、页眉页脚配置的样式文件转换，验证 Word 文档中各元素存在且格式正确。

**Acceptance Scenarios**:

1. **Given** 样式配置 `page.size: A4, page.margin.top: 25.4mm`, **When** 转换, **Then** docx 页面设置为 A4、上边距 25.4mm
2. **Given** 样式配置 `heading_numbering.enabled: true`, **When** 转换包含 h1/h2/h3 的 md, **Then** 标题前有编号 "1", "1.1", "1.1.1"（已有编号的标题先去重）
3. **Given** 样式配置 `toc.enabled: true`, **When** 转换, **Then** 文档开头有 TOC 域代码
4. **Given** 样式配置 `header.left: "{title}"`, frontmatter 中 `title: 我的论文`, **When** 转换, **Then** 页眉左侧显示 "我的论文"
5. **Given** 样式配置 `footer.center: "第 {page} 页"`, **When** 转换, **Then** 页脚有 PAGE 域代码

---

### User Story 4 - Cover Page (Priority: P4)

用户通过 frontmatter 提供 title/author/date/institution，样式定义封面布局，自动生成封面页。

**Why this priority**: 封面是学术文档的标配，但实现复杂度低于核心管线。

**Independent Test**: 用带封面配置的样式和包含 frontmatter 的 md 文件转换，验证第一页是封面。

**Acceptance Scenarios**:

1. **Given** frontmatter 含 title/author/date, 样式 cover.enabled: true, **When** 转换, **Then** 第一页居中显示标题（大字加粗）、作者、日期，封面后有分页符
2. **Given** 样式 cover.enabled: false, **When** 转换, **Then** 无封面页
3. **Given** `md2docx input.md --no-cover`, **When** 转换, **Then** CLI 参数覆盖样式设置，无封面

---

### User Story 5 - Tables & Images (Priority: P5)

Markdown 中的 GFM 表格和图片正确渲染到 Word 中，支持样式控制。

**Why this priority**: 学术文档必有表格和图片，但 OOXML 表格实现复杂度高。

**Independent Test**: 用包含表格和图片的 md 文件转换，验证表格有边框、图片按比例缩放。

**Acceptance Scenarios**:

1. **Given** md 含 GFM 表格, **When** 转换, **Then** Word 表格有边框，表头加粗
2. **Given** md 含 `![alt](path/to/img.png)`, **When** 转换, **Then** 图片嵌入 docx，宽度不超过页面 80%
3. **Given** 图片路径是相对于 md 文件的相对路径, **When** 转换, **Then** 正确解析并嵌入

---

### User Story 6 - Multi-file Merge (Priority: P6)

用户将多个 md 文件合并为一个 Word 文档。

**Why this priority**: 长文档（论文、书籍）的常见需求，但 MVP 之后。

**Independent Test**: 用 3 个 md 文件运行 merge，验证输出包含所有内容且顺序正确。

**Acceptance Scenarios**:

1. **Given** `md2docx merge ch1.md ch2.md ch3.md -o book.docx`, **Then** 输出包含三个文件内容，顺序正确
2. **Given** `md2docx merge ch1.md ch2.md --style academic-cn`, **Then** 合并后的文档应用指定样式

---

### Edge Cases

- Markdown 文件为空 → 生成空 docx（仅有页面设置），不报错
- Markdown 文件有 BOM → 正确处理
- 图片文件不存在 → 跳过图片，stderr 输出警告，继续转换
- 样式 YAML 格式错误 → 明确报错，指出行号
- 超大文件（>10MB markdown）→ 应能处理，不 OOM
- Windows 路径（反斜杠）→ 正确处理
- 文件名含中文/空格 → 正确处理

## Requirements

### Functional Requirements

- **FR-001**: System MUST 解析 CommonMark + GFM Markdown（标题、段落、列表、表格、代码块、粗体、斜体、链接、图片）
- **FR-002**: System MUST 生成有效的 OOXML .docx 文件，可被 Microsoft Word 和 WPS 正常打开
- **FR-003**: System MUST 支持 CJK/Latin 双字体分离（中文字符用 eastAsia 字体，Latin 字符用 ascii 字体）
- **FR-004**: System MUST 支持 YAML frontmatter 解析，提取 title/author/date 等元数据
- **FR-005**: System MUST 内嵌至少 3 个样式（default、academic-cn、simple），通过 go:embed 编译进二进制
- **FR-006**: System MUST 支持页面设置（纸张大小、页边距、页面方向）
- **FR-007**: System MUST 支持目录生成（TOC 域代码）
- **FR-008**: System MUST 支持标题编号（层级编号，去重已有编号）
- **FR-009**: System MUST 支持封面页生成（标题、作者、日期居中排列）
- **FR-010**: System MUST 支持页眉页脚（含 PAGE 页码域代码、frontmatter 变量替换）
- **FR-011**: System MUST 支持图片嵌入（相对路径解析、按比例缩放）
- **FR-012**: System MUST 支持 GFM 表格渲染（边框、表头样式）
- **FR-013**: System MUST 默认输出到 ~/Desktop/，可通过 -o 参数覆盖
- **FR-014**: System MUST 支持多文件合并（merge 子命令）
- **FR-015**: System MUST 作为单个静态链接二进制分发，无运行时依赖

### Key Entities

- **Style**: YAML 样式定义，包含页面设置、字体、段落样式、标题编号、TOC、封面、页眉页脚等全部排版规则
- **Document**: 一次转换的上下文，包含源 Markdown、解析后的 AST、元数据、应用的样式、输出路径
- **OOXML Package**: 最终输出的 docx zip 包，包含 document.xml、styles.xml、numbering.xml 等部件

## Success Criteria

### Measurable Outcomes

- **SC-001**: `md2docx input.md` 从任意 Markdown 文件生成可被 Word 正常打开的 docx，转换时间 < 1 秒
- **SC-002**: 内置 academic-cn 样式生成的文档，中文字体/字号/行距/缩进符合一般学术论文规范
- **SC-003**: 单二进制文件大小 < 20MB
- **SC-004**: 支持 macOS (arm64)、Windows (amd64)、Linux (amd64) 三平台
- **SC-005**: 测试覆盖率 > 80%

## Assumptions

- 用户系统已安装目标字体（宋体、黑体等），md2docx 不负责字体嵌入
- Markdown 输入为 UTF-8 编码
- 图片格式限 PNG、JPEG、GIF（Word 原生支持的格式）
- 数学公式和脚注支持不在 v1 范围内，作为后续增量交付
- 代码块语法高亮不在 v1 范围内
