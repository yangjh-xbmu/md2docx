# Tasks: md2docx Go Rewrite

**Input**: [plan.md](./plan.md), [spec.md](./spec.md)
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Setup

**Purpose**: Go 项目初始化和基础结构

- [ ] T001 初始化 Go module：`go mod init github.com/yangjh-xbmu/md2docx`，添加依赖 goldmark、cobra、yaml.v3、testify
- [ ] T002 创建目录结构：cmd/md2docx/、internal/{cli,parser,style,ooxml,render,merge}/、styles/、testdata/
- [ ] T003 [P] 创建 cobra CLI 骨架：root.go（默认行为=convert）、convert.go、styles.go，支持 `--version`
- [ ] T004 [P] 创建 Makefile：build、test、lint、release（交叉编译 3 平台）

---

## Phase 2: Foundational (OOXML Writer + Markdown Parser)

**Purpose**: 核心基础设施，所有 User Story 依赖此阶段

⚠️ CRITICAL: 无此阶段，任何 User Story 无法开始

- [ ] T005 实现 ooxml/types.go：定义 OOXML 核心类型（Body、Paragraph、Run、RunProperties、Text、ParagraphProperties 等）
- [ ] T006 实现 ooxml/content_types.go：生成 [Content_Types].xml
- [ ] T007 [P] 实现 ooxml/rels.go：生成 _rels/.rels 和 word/_rels/document.xml.rels
- [ ] T008 实现 ooxml/document.go：生成 word/document.xml（body + paragraphs + runs）
- [ ] T009 实现 ooxml/styles.go：生成 word/styles.xml（Normal、Heading1-6、Quote、Code 等样式定义，含 CJK 双字体 rFonts）
- [ ] T010 实现 ooxml/writer.go：组装所有 part → archive/zip → .docx 输出
- [ ] T011 测试：生成最小 docx（"Hello World" 段落），用 `unzip -l` 验证结构，用 Word 打开验证
- [ ] T012 实现 parser/markdown.go：goldmark 配置（CommonMark + GFM extensions）
- [ ] T013 [P] 实现 parser/frontmatter.go：提取 YAML frontmatter，返回 metadata map + body string
- [ ] T014 实现 style/types.go：YAML 样式结构体定义（Page、Fonts、Styles、HeadingNumbering、TOC、Cover、Header、Footer）
- [ ] T015 [P] 实现 style/loader.go：从 embed FS 或用户目录加载样式，合并默认值
- [ ] T016 [P] 创建 styles/default.yaml：默认样式（宋体正文、黑体标题、A4、1.5 倍行距、首行缩进 2em）

**Checkpoint**: 能生成包含正确样式定义的 docx，goldmark 能解析 Markdown AST

---

## Phase 3: User Story 1 - Basic Conversion (P1) 🎯 MVP

**Goal**: `md2docx input.md` 生成可用的 Word 文档

**Independent Test**: 转换含标题/段落/列表/粗体斜体的 md，Word 打开验证

### Tests

- [ ] T017 [P] [US1] 单元测试 ooxml/document_test.go：验证段落、标题、列表生成的 XML 结构
- [ ] T018 [P] [US1] 单元测试 render/renderer_test.go：验证 AST 节点→OOXML 映射
- [ ] T019 [US1] 集成测试 e2e_test.go：转换 testdata/basic.md，验证输出 docx 合法（可 unzip、document.xml 结构正确）

### Implementation

- [ ] T020 实现 render/renderer.go：goldmark NodeRenderer 接口，注册所有 block/inline 渲染函数
- [ ] T021 实现 render/block.go：Heading、Paragraph、List（ordered+unordered）、ListItem、BlockQuote、CodeBlock、ThematicBreak
- [ ] T022 实现 render/inline.go：Text、Emphasis（bold/italic）、CodeSpan、Link、HardLineBreak、SoftLineBreak
- [ ] T023 实现 ooxml/document.go 扩展：支持 list paragraph（numId 引用 numbering.xml）
- [ ] T024 实现 ooxml/numbering.go：基础列表编号定义（bullet + ordered）
- [ ] T025 实现 cli/convert.go 完整逻辑：读文件 → 解析 frontmatter → 加载样式 → goldmark 解析 → 渲染 → 写 docx
- [ ] T026 实现默认输出路径逻辑：无 -o 时输出到 ~/Desktop/<同名>.docx
- [ ] T027 创建 testdata/basic.md（标题+段落+列表+格式）和 testdata/cjk.md（中英文混排）
- [ ] T028 端到端验证：`go run ./cmd/md2docx testdata/basic.md`，Word 打开检查

**Checkpoint**: MVP 可用。`md2docx input.md` 生成格式正确的 Word 文档

---

## Phase 4: User Story 2 - Style Selection (P2)

**Goal**: `--style academic-cn` 切换内置样式

**Independent Test**: 同一 md 用不同样式转换，对比输出差异

### Implementation

- [ ] T029 创建 styles/academic-cn.yaml：黑体标题、宋体正文 12pt、首行缩进 2em、1.5 行距、A4
- [ ] T030 [P] 创建 styles/simple.yaml：简洁风格，无封面无目录
- [ ] T031 实现 cli/styles.go：`styles list`（列出内嵌+用户样式）、`styles show <name>`（打印样式详情）
- [ ] T032 实现 style/defaults.go：为未设置的样式字段填充合理默认值
- [ ] T033 实现 ooxml/styles.go 扩展：根据 Style YAML 动态生成 styles.xml（字体、字号、颜色、对齐、行距、缩进）
- [ ] T034 实现 ooxml/document.go 扩展：页面设置（w:pgSz 纸张大小、w:pgMar 页边距）
- [ ] T035 [US2] 测试：academic-cn 样式输出的 styles.xml 中 Heading1 font-family 包含"黑体"

**Checkpoint**: 样式系统完整，`--style` 参数可用

---

## Phase 5: User Story 3 - Document Structure (P3)

**Goal**: TOC、标题编号、页眉页脚

### Implementation

- [ ] T036 实现 render/heading_number.go：遍历 AST 标题节点，生成编号前缀（去重已有编号）
- [ ] T037 实现 render/toc.go：在 document.xml 开头插入 TOC 域代码（fldChar + instrText）
- [ ] T038 实现 ooxml/header.go：生成 word/header1.xml（支持 frontmatter 变量替换）
- [ ] T039 [P] 实现 ooxml/footer.go：生成 word/footer1.xml（支持 PAGE 域代码）
- [ ] T040 实现 ooxml/rels.go 扩展：添加 header/footer 关系
- [ ] T041 实现 ooxml/numbering.go 扩展：标题编号定义（多级列表 abstractNum）
- [ ] T042 实现 style/merge.go：frontmatter 字段覆盖样式配置（toc、heading_numbering、cover 等开关）
- [ ] T043 实现 CLI 参数 --no-toc、--no-numbering
- [ ] T044 [US3] 测试：验证 TOC 域代码结构、标题编号正确性、页眉页脚内容

**Checkpoint**: 完整文档结构（TOC + 编号 + 页眉页脚）可用

---

## Phase 6: User Story 4 - Cover Page (P4)

**Goal**: 自动生成封面页

### Implementation

- [ ] T045 实现 render/cover.go：根据样式 cover.layout 和 frontmatter 元数据生成封面段落
- [ ] T046 实现封面后分页符（w:br w:type="page"）
- [ ] T047 实现 CLI 参数 --no-cover
- [ ] T048 [US4] 测试：验证封面段落存在、字体字号正确、分页符存在

**Checkpoint**: 封面页功能完整

---

## Phase 7: User Story 5 - Tables & Images (P5)

**Goal**: GFM 表格和图片渲染

### Implementation

- [ ] T049 实现 render/table.go：GFM 表格 AST → OOXML 表格（w:tbl、w:tr、w:tc）
- [ ] T050 实现 ooxml/table.go：表格 XML 生成（边框、表头背景色、单元格对齐、表格居中）
- [ ] T051 实现 render/image.go：图片 AST → 读取文件 → 嵌入 docx
- [ ] T052 实现 ooxml/image.go：图片嵌入（word/media/、relationship、inline drawing XML）
- [ ] T053 实现图片缩放逻辑：按样式 image.max_width_pct 限制最大宽度
- [ ] T054 创建 testdata/with_table.md 和 testdata/with_image.md
- [ ] T055 [US5] 测试：验证表格/图片 XML 结构正确

**Checkpoint**: 表格和图片渲染完整

---

## Phase 8: User Story 6 - Multi-file Merge (P6)

**Goal**: merge 子命令

### Implementation

- [x] T056 实现 merge/merge.go：读取多个 md 文件，拼接内容（第一个文件的 frontmatter 生效）
- [x] T057 实现 cli/merge.go：merge 子命令，接受多个文件参数 + -o + --style
- [x] T058 [US6] 测试：3 个 md 文件合并，验证内容完整性和顺序

**Checkpoint**: merge 功能完整

---

## Phase 9: Polish & Release

**Purpose**: 质量收尾和分发

- [x] T059 [P] 完善错误处理：BOM 剥离、空文件处理、输出目录自动创建、图片不存在警告
- [x] T060 [P] 完善 CLI help（Long + Example）、版本号 ldflags 注入
- [x] T061 创建 GitHub Actions CI：test + lint + build
- [x] T062 [P] Makefile release：ldflags 版本注入 + sha256 校验和
- [x] T063 编写 README.md：安装、快速开始、样式系统、CLI 参考、限制
- [x] T064 测试覆盖率：cli 89.6%, ooxml 87.5%, render 84.4%, merge 91.2%
- [x] T065 手工验收：academic-cn/default/simple 三样式端到端 Word 打开验证通过

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup) → Phase 2 (Foundation) → Phase 3 (MVP/US1)
                                        ↓
                            Phase 4 (US2: Styles)
                                        ↓
                            Phase 5 (US3: Structure)
                                        ↓
                            Phase 6 (US4: Cover)
                                        ↓
                            Phase 7 (US5: Tables/Images)
                                        ↓
                            Phase 8 (US6: Merge)
                                        ↓
                            Phase 9 (Polish)
```

### Parallel Opportunities

- Phase 2: T007/T013/T015/T016 可并行
- Phase 3: T017/T018 可并行
- Phase 4: T029/T030 可并行
- Phase 5: T038/T039 可并行
- Phase 7: T049-T050 和 T051-T052 可并行（表格和图片独立）
- Phase 9: T059/T060/T062 可并行

### MVP Delivery Point

Phase 3 完成即为 MVP：`md2docx input.md` 能生成包含标题、段落、列表、格式的 Word 文档。
