# SESSION LOG

## 完成
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
- 2026-04-03 编写 19 个人工验收测试用例（TESTING.md）

## 发现
- 2026-04-03 Go style YAML 的 bool 零值陷阱：yaml 解析 `false` 后得到 Go 零值，applyDefaults 中 `if !field` 无法区分「未设置」和「显式 false」，解法是不在 defaults 中覆盖 bool 字段
- 2026-04-03 OOXML CJK 双字体不需要逐 run 检测字符：只需在 styles.xml 的 w:rFonts 同时设置 w:ascii 和 w:eastAsia，Word 会根据字符 Unicode range 自动选择字体
- 2026-04-03 go:embed 不跟随符号链接，嵌入外部目录的文件需要 cp 复制而非 ln -s
- 2026-04-03 goldmark 的 GFM Table 扩展没有 TableBody 类型，表头用 TableHeader、数据行用 TableRow，都是 Table 的直接子节点
- 2026-04-03 AI 编码时代自建 OOXML writer 可行：docx 本质是 zip+XML，encoding/xml 直接生成，AI 擅长这种结构化代码
- 2026-04-03 pandoc --reference-doc 只能控制样式定义，无法控制 TOC/封面/页眉等结构性元素，必须用 python-docx 二次处理才能达到"零手动调格式"
- 2026-04-03 CSS font-family 中 CJK 字体需要映射到 Word 的 w:eastAsia 属性而非 w:ascii，否则中文字体不生效
- 2026-04-03 python-docx 插入 TOC 是通过 Word 域代码（field code），文档打开后需按 Ctrl+A 再 F9 才能更新实际目录内容

## 待办
1. 继续 Go 版 Phase 5-9：TOC、标题编号、封面、页眉页脚、图片嵌入、merge、发布
2. 将 Achuan-2/pandoc_docx_template 的 Lua filters（preserve_font_color、add-inline-code）移植到 Go 版渲染器
3. 考虑将 course-toolkit 的 docx 生成统一迁移到 md2docx Go 版
