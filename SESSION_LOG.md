# SESSION LOG

## 完成
- 2026-04-03 仓库设为公开，重写 README（含学术论文完整案例），新增 release.yml 自动构建 4 平台二进制，制定开源推广计划（PROMOTION.md，已 gitignore）
- 2026-04-03 打 v0.1.0 tag 并创建 GitHub Release（annotated tag + gh release create，含完整 release notes）
- 2026-04-03 完成 Go 版 Phase 9：质量收尾和发布。BOM 处理、CLI Long/Example help、ldflags 版本注入、GitHub Actions CI、Makefile sha256 校验、README.md、测试覆盖率从 30-48% 提升到 84-91%（cli 89.6%, ooxml 87.5%, render 84.4%），三样式端到端 Word 验收通过
- 2026-04-03 完成 Go 版 Phase 8：merge 子命令，支持直接文件参数、glob 模式、contents.yaml 配置（per-file heading_offset），10 个单元测试通过，端到端验证 Word 打开正常
- 2026-04-03 完成 Go 版 Phase 7：GFM 表格渲染（w:tbl/w:tr/w:tc），支持边框、表头加粗、表头背景色、单元格内 bold/italic/code 格式，7 个单元测试通过
- 2026-04-03 完成 Go 版 Phase 7：图片嵌入（wp:inline drawing），支持相对路径解析、max_width_pct 缩放（保持宽高比）、word/media/ 嵌入 + relationship，5 个单元测试通过
- 2026-04-03 替换 renderer.go 中简化的 tab 分隔表格渲染为真正 OOXML 表格元素
- 2026-04-03 新增 ooxml/types.go 中 27 个 OOXML 类型定义（Table/Drawing/Picture 等）
- 2026-04-03 完成 Go 版 Phase 6：封面页生成（render/cover.go），支持 CoverConfig.Layout 声明式布局（text/spacer），frontmatter 变量绑定，CJK 字体，literal: 前缀直出文本，分页符，6 个单元测试通过
- 2026-04-03 修复 YAML 日期解析问题：yaml.Unmarshal 将日期字符串解析为 time.Time，在 resolveText 中格式化为 "2006-01-02"
- 2026-04-03 完成 Go 版 Phase 5：TOC 域代码生成、标题编号（多级 + 中文格式）、页眉页脚（变量替换 + PAGE 域代码），18 个新测试全部通过
- 2026-04-03 生成 AST（抽象语法树）学习资料，保存到 Obsidian，结合 md2docx 项目实例讲解 beginner 级别
- 2026-04-03 生成三份端到端测试文档（学术论文/简洁备忘/默认混排），用户验证效果远超预期
- 2026-04-03 完成 Go 版 Phase 4：多样式系统，新增 academic-cn（学术论文）和 simple（简洁）两个内嵌样式，实现 styles list/show 子命令
- 2026-04-03 修复 setHeadingDefaults bool 零值问题：将 EmbeddedStyles 从 embed.FS 改为 fs.FS 接口，PageBreakBefore 不再在 applyDefaults 中覆盖
- 2026-04-03 完成 Go 重写 MVP：goldmark + 自建 OOXML writer，单二进制 5.4MB，零依赖，可生成含标题/段落/列表/粗体/斜体/代码/引用/表格的 docx
- 2026-04-03 建立 Spec Kit 工作流：constitution → spec → plan → tasks，产出完整开发文档（specs/001-go-rewrite/）
- 2026-04-03 实现 YAML 样式系统：style/types.go 定义完整 schema，style/loader.go 支持 go:embed 内嵌 + 用户目录加载
- 2026-04-03 实现 OOXML writer 核心模块：types.go/writer.go/styles.go/units.go/numbering.go，直接用 encoding/xml + archive/zip 生成 docx
- 2026-04-03 实现 goldmark AST → OOXML 渲染器：render/renderer.go，支持 Heading/Paragraph/List/Emphasis/CodeSpan/CodeBlock/Link/Blockquote/Table
- 2026-04-03 创建 default.yaml 内嵌样式：宋体正文、黑体标题、A4 纸张、CJK 双字体
- 2026-04-03 调研 GitHub 同类项目：Achuan-2/pandoc_docx_template(624★)、nihole/md2docx(190★)、docxcompose(128★) 等
- 2026-04-03 调研 AI 开发工作流工具：Spec Kit(85k★)、GSD(47k★)、BMAD(43k★)、CCPM(8k★)，选定 Spec Kit
- 2026-04-03 修复 Python 版两个 bug：标题编号重复（正则去重）、默认输出路径改为桌面
- 2026-04-03 调研现有三代 md2docx 实现（works-used-python/merge、.claude/skills/md2docx、hongo/tools/md2docx），确定加强版方案
- 2026-04-03 完成技术设计（PLAN.md）：两阶段管线、多模板系统、frontmatter 驱动、多文件合并
- 2026-04-03 创建项目骨架：pyproject.toml、src/md2docx/ 七个模块、styles/default/、tests/
- 2026-04-03 实现核心转换管线：pandoc 调用 + python-docx 后处理（TOC、标题编号、图片缩放、页眉页脚、封面）
- 2026-04-03 实现多模板系统：style.yaml 配置、用户/内置样式目录、样式查找链
- 2026-04-03 实现多文件合并：contents.yaml 驱动、glob 模式、heading_offset 标题层级偏移
- 2026-04-03 实现 Click CLI：convert / merge / styles(list/show/dir/init/build) / check 子命令
- 2026-04-03 实现 CSS → Word 样式转换（css2style 模块）：解析 CSS、CJK 字体自动识别、生成 reference.docx
- 2026-04-03 CLI 集成 --css 参数和 styles build 子命令
- 2026-04-03 54 个 pytest 全部通过，端到端验证 CLI 所有命令

## 发现
- 2026-04-03 Go merge 多文件合并的核心设计：Markdown 级拼接（非 docx 级），第一个文件 frontmatter 保留，后续剥离，然后整体走 convert 管线。比 docx 级合并简单且可靠
- 2026-04-03 OOXML 图片嵌入需要同时处理三个层面：word/media/ 中存放文件、word/_rels/document.xml.rels 中添加 relationship、document.xml 中用 wp:inline + pic:pic 引用。DocPr.id 可以用 relID 字符串（Word 能容忍）
- 2026-04-03 goldmark GFM Table 的 TableHeader 节点是一个完整的行（包含 TableCell 子节点），不是单独的单元格标记。遍历时 TableHeader 和 TableRow 平级处理即可
- 2026-04-03 OOXML 表格宽度设 type="auto" + w="0" 让 Word 自动计算列宽，比手动计算 twips 更可靠
- 2026-04-03 YAML yaml.Unmarshal 会将 "2026-04-03" 等日期字符串自动解析为 time.Time，在封面等需要文本显示的场景中需要特殊处理格式化
- 2026-04-03 OOXML SectionProperties 中 headerReference/footerReference 必须在 pgSz/pgMar 之前，否则 encoding/xml 序列化后 Word 可能忽略引用。Go struct 字段顺序决定 XML 输出顺序
- 2026-04-03 Go style YAML 的 bool 零值陷阱：yaml 解析 `false` 后得到 Go 零值，applyDefaults 中 `if !field` 无法区分「未设置」和「显式 false」，解法是不在 defaults 中覆盖 bool 字段
- 2026-04-03 OOXML CJK 双字体不需要逐 run 检测字符：只需在 styles.xml 的 w:rFonts 同时设置 w:ascii 和 w:eastAsia，Word 会根据字符 Unicode range 自动选择字体
- 2026-04-03 go:embed 不跟随符号链接，嵌入外部目录的文件需要 cp 复制而非 ln -s
- 2026-04-03 goldmark 的 GFM Table 扩展没有 TableBody 类型，表头用 TableHeader、数据行用 TableRow，都是 Table 的直接子节点
- 2026-04-03 AI 编码时代自建 OOXML writer 可行：docx 本质是 zip+XML，encoding/xml 直接生成，AI 擅长这种结构化代码
- 2026-04-03 pandoc --reference-doc 只能控制样式定义，无法控制 TOC/封面/页眉等结构性元素，必须用 python-docx 二次处理才能达到"零手动调格式"
- 2026-04-03 CSS font-family 中 CJK 字体需要映射到 Word 的 w:eastAsia 属性而非 w:ascii，否则中文字体不生效
- 2026-04-03 python-docx 插入 TOC 是通过 Word 域代码（field code），文档打开后需按 Ctrl+A 再 F9 才能更新实际目录内容

## 待办
1. 开源推广：GitHub 加 topics，README 顶部加转换效果截图/GIF，提交 awesome-go
2. 开源推广：写中文技术社区推广文（少数派/掘金/V2EX/知乎），学术论文场景 before/after 对比
3. 开源推广：提交 Homebrew tap，降低安装门槛
4. 将 Achuan-2/pandoc_docx_template 的 Lua filters（preserve_font_color、add-inline-code）移植到 Go 版渲染器
5. 考虑将 course-toolkit 的 docx 生成统一迁移到 md2docx Go 版
