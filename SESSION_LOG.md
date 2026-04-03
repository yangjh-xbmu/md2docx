# SESSION LOG

## 完成
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
- 2026-04-03 pandoc --reference-doc 只能控制样式定义，无法控制 TOC/封面/页眉等结构性元素，必须用 python-docx 二次处理才能达到"零手动调格式"
- 2026-04-03 CSS font-family 中 CJK 字体需要映射到 Word 的 w:eastAsia 属性而非 w:ascii，否则中文字体不生效
- 2026-04-03 python-docx 插入 TOC 是通过 Word 域代码（field code），文档打开后需按 Ctrl+A 再 F9 才能更新实际目录内容

## 待办
1. 给 merge 命令也加上 --css 参数支持
2. 增加更多内置样式模板（academic、official、lesson-plan）
3. 在实际教学材料生成场景中试用，收集需要调整的后处理细节
4. 考虑将 course-toolkit 的 docx 生成统一迁移到 md2docx
