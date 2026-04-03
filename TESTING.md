# md2docx 人工验收测试手册

## 准备工作

```bash
cd ~/Desktop/repos/md2docx
uv sync
```

确认环境：

```bash
uv run md2docx --version     # 应输出: md2docx, version 0.1.0
uv run md2docx check          # 应输出: pandoc 可用: pandoc x.x.x
```

---

## 一、基础转换（convert）

### 测试 1.1：最简转换

```bash
uv run md2docx convert tests/fixtures/simple.md -o /tmp/test_1_1.docx
```

**验收标准**：
- 终端输出 `已导出: /tmp/test_1_1.docx`
- 用 Word/WPS 打开 `/tmp/test_1_1.docx`
- 检查内容包含 "Hello World"、"Section One"、"Section Two"
- 列表项（Item one/two/three）正确显示

### 测试 1.2：默认输出路径

```bash
cp tests/fixtures/simple.md /tmp/test_1_2.md
uv run md2docx convert /tmp/test_1_2.md
```

**验收标准**：
- 输出文件自动生成在 `/tmp/test_1_2.docx`（与源文件同目录、同名）

### 测试 1.3：跳过后处理

```bash
uv run md2docx convert tests/fixtures/simple.md -o /tmp/test_1_3.docx --no-post
```

**验收标准**：
- 转换成功
- 打开文件，样式是 pandoc 原始输出（无额外处理）

---

## 二、Frontmatter 驱动

### 测试 2.1：带 frontmatter 的转换

```bash
uv run md2docx convert tests/fixtures/with_frontmatter.md -o /tmp/test_2_1.docx
```

**验收标准**：
- 打开文件，检查：
  - 文档开头有 **"目录"** 两个字（TOC 域代码已插入）
  - 标题带有层次编号（如 "1 第一章 概述"、"1.1 1.1 背景"）
  - frontmatter 中的 `---` 块内容不出现在正文中

### 测试 2.2：自定义 frontmatter 文档

创建 `/tmp/test_2_2.md`：

```markdown
---
title: 我的测试报告
author: 测试员
date: 2026-04-03
toc: true
heading_numbering: true
header_left: "{title}"
header_right: "{author}"
footer_center: "第 {page} 页"
---

# 引言

这是引言内容。

## 背景

这是背景内容。

## 目的

这是目的描述。

# 方法

## 数据来源

描述数据来源。

## 分析方法

描述分析方法。

# 结论

最终结论。
```

```bash
uv run md2docx convert /tmp/test_2_2.md -o /tmp/test_2_2.docx
```

**验收标准**：
- 打开文件，检查：
  - 有目录（显示"目录"标题 + TOC 域代码）
  - 标题自动编号：1 引言、1.1 背景、1.2 目的、2 方法 …
  - 页眉左侧："我的测试报告"
  - 页眉右侧："测试员"
  - 页脚居中：显示页码（在 Word 中按 Ctrl+A 后 F9 更新域可看到实际页码）

---

## 三、多文件合并（merge）

### 测试 3.1：通过 contents.yaml 合并

```bash
uv run md2docx merge tests/fixtures/merge_project/contents.yaml -o /tmp/test_3_1.docx
```

**验收标准**：
- 文件包含三章内容："第一章 引言"、"第二章 方法"、"第三章 结论"
- 第三章的标题层级被降了一级（原来是 `#`，应变为 `##`，因为 heading_offset=1）

### 测试 3.2：通过 glob 模式合并

```bash
uv run md2docx merge "tests/fixtures/merge_project/chapter*.md" -o /tmp/test_3_2.docx
```

**验收标准**：
- 文件包含所有三章内容
- 文件顺序按字母排序（chapter1 → chapter2 → chapter3）

### 测试 3.3：自定义 contents.yaml

创建 `/tmp/merge_test/` 目录和文件：

```bash
mkdir -p /tmp/merge_test
```

创建 `/tmp/merge_test/preface.md`：
```markdown
# 前言

这是前言。
```

创建 `/tmp/merge_test/ch1.md`：
```markdown
# 第一章

正文第一章。

## 第一节

第一节内容。
```

创建 `/tmp/merge_test/ch2.md`：
```markdown
# 第二章

正文第二章。
```

创建 `/tmp/merge_test/contents.yaml`：
```yaml
heading_offset: 0

files:
  - preface.md
  - ch1.md
  - path: ch2.md
    heading_offset: 1
```

```bash
uv run md2docx merge /tmp/merge_test/contents.yaml -o /tmp/test_3_3.docx
```

**验收标准**：
- 前言是一级标题
- 第一章是一级标题，第一节是二级标题
- 第二章因为 heading_offset=1，变成了二级标题

---

## 四、CSS 样式

### 测试 4.1：从 CSS 生成 reference.docx

```bash
uv run md2docx styles build tests/fixtures/academic.css -o /tmp/test_4_1_ref.docx
```

**验收标准**：
- 终端输出 `已生成模板: /tmp/test_4_1_ref.docx`
- 用 Word 打开，查看样式面板（Word 的"样式"窗格）：
  - "Heading 1" 样式：黑体/Arial、22pt、加粗、居中
  - "Normal" 样式：宋体/Times New Roman、12pt、首行缩进
  - 应存在自定义样式 "abstract"、"keywords"

### 测试 4.2：用 CSS 直接转换

```bash
uv run md2docx convert /tmp/test_2_2.md --css tests/fixtures/academic.css -o /tmp/test_4_2.docx
```

**验收标准**：
- 转换成功
- 打开文件，标题应使用黑体（或 Arial），正文应使用宋体（或 Times New Roman）
- 与测试 2.2 的输出对比，字体和格式明显不同

### 测试 4.3：自定义 CSS 文件

创建 `/tmp/test_4_3.css`：

```css
h1 {
    font-family: "Microsoft YaHei", "Helvetica";
    font-size: 26pt;
    font-weight: bold;
    text-align: center;
    color: #1a5276;
    margin-bottom: 18pt;
}

h2 {
    font-family: "Microsoft YaHei", "Helvetica";
    font-size: 18pt;
    color: #2c3e50;
    margin-top: 16pt;
}

p {
    font-family: "FangSong", "Georgia";
    font-size: 14pt;
    line-height: 2.0;
    text-indent: 2em;
    text-align: justify;
}

.note {
    font-size: 11pt;
    color: #7f8c8d;
    font-style: italic;
    margin-left: 1.5cm;
}
```

```bash
uv run md2docx styles build /tmp/test_4_3.css -o /tmp/test_4_3_ref.docx
uv run md2docx convert /tmp/test_2_2.md --css /tmp/test_4_3.css -o /tmp/test_4_3.docx
```

**验收标准**：
- reference.docx 中样式面板可看到：
  - Heading 1：微软雅黑/Helvetica、26pt、蓝色(#1a5276)
  - Normal：仿宋/Georgia、14pt、行距 2.0
  - 存在自定义样式 "note"（灰色、斜体）
- 转换后的文档正文使用仿宋字体、行距明显加大

### 测试 4.4：基于已有模板叠加 CSS

```bash
uv run md2docx styles build /tmp/test_4_3.css --base styles/default/reference.docx -o /tmp/test_4_4_ref.docx
```

**验收标准**：
- 输出的 reference.docx 保留了原始 default 模板的基础样式
- 同时叠加了 CSS 中定义的样式（颜色、字体、字号）

---

## 五、样式管理（styles）

### 测试 5.1：列出样式

```bash
uv run md2docx styles list
```

**验收标准**：
- 至少显示 `default` 样式，标注"有模板, 1 个 Lua filter"

### 测试 5.2：查看样式详情

```bash
uv run md2docx styles show default
```

**验收标准**：
- 显示样式名称、模板路径、Lua filter 路径
- 显示后处理配置（目录:否、标题编号:否、图片宽度:80% 等）

### 测试 5.3：打印样式目录

```bash
uv run md2docx styles dir
```

**验收标准**：
- 显示两行：用户样式目录（`~/.md2docx/styles`）和内置样式目录

### 测试 5.4：创建新样式骨架

```bash
uv run md2docx styles init teststyle
```

**验收标准**：
- 创建了 `~/.md2docx/styles/teststyle/` 目录
- 目录中包含 `style.yaml` 和 `filters/` 子目录
- `style.yaml` 内容是合理的默认配置

测试后清理：
```bash
rm -rf ~/.md2docx/styles/teststyle
```

---

## 六、封面页

### 测试 6.1：带封面的转换

创建 `/tmp/test_6_1.md`：

```markdown
---
title: 2026年度科研项目申报书
author: 杨志宏
date: 2026-04-03
institution: 西北民族大学
---

# 项目概述

这是项目概述内容。

## 研究背景

这是研究背景。

# 研究方案

## 技术路线

这是技术路线。
```

创建 `/tmp/test_6_1_style/` 样式：

```bash
mkdir -p /tmp/test_6_1_style/filters
```

创建 `/tmp/test_6_1_style/style.yaml`：

```yaml
name: cover-test
reference_doc: reference.docx

post_processing:
  toc: true
  heading_numbering: true
  header_right: "{author}"
  footer_center: "第 {page} 页"
  cover:
    enabled: true
    fields:
      title: "{title}"
      author: "{author}"
      date: "{date}"
```

```bash
cp styles/default/reference.docx /tmp/test_6_1_style/reference.docx
cp styles/default/filters/styles.lua /tmp/test_6_1_style/filters/
```

由于自定义样式需要放在样式目录，先复制到用户样式目录：

```bash
mkdir -p ~/.md2docx/styles/cover-test
cp /tmp/test_6_1_style/* ~/.md2docx/styles/cover-test/ 2>/dev/null
cp -r /tmp/test_6_1_style/filters ~/.md2docx/styles/cover-test/
```

```bash
uv run md2docx convert /tmp/test_6_1.md -o /tmp/test_6_1.docx --style cover-test
```

**验收标准**：
- 打开文件，第一页是封面：
  - 标题居中、大字："2026年度科研项目申报书"
  - 作者居中："杨志宏"
  - 日期居中："2026-04-03"
- 第二页开始是目录
- 正文标题有编号
- 页眉右侧："杨志宏"
- 页脚有页码

测试后清理：
```bash
rm -rf ~/.md2docx/styles/cover-test
```

---

## 七、错误处理

### 测试 7.1：不存在的文件

```bash
uv run md2docx convert /tmp/nonexistent_file.md
```

**验收标准**：错误提示文件不存在（Click 自动处理）

### 测试 7.2：不存在的样式

```bash
uv run md2docx convert tests/fixtures/simple.md --style no_such_style
```

**验收标准**：提示 `错误: 样式 'no_such_style' 不存在。可用样式: default`

### 测试 7.3：空的 CSS 文件

```bash
echo "/* nothing */" > /tmp/empty.css
uv run md2docx styles build /tmp/empty.css -o /tmp/empty_ref.docx
```

**验收标准**：提示 `错误: CSS 文件中没有找到有效的样式规则`

### 测试 7.4：合并时缺少文件

创建 `/tmp/bad_merge.yaml`：

```yaml
files:
  - nonexistent1.md
  - nonexistent2.md
```

```bash
uv run md2docx merge /tmp/bad_merge.yaml -o /tmp/bad_merge.docx
```

**验收标准**：提示文件不存在，列出缺失的文件名

---

## 八、综合端到端测试

### 测试 8.1：完整学术论文工作流

这个测试模拟真实使用场景：用 CSS 定义样式，用 frontmatter 控制后处理，合并多章节。

创建工作目录：

```bash
mkdir -p /tmp/paper_test
```

创建 `/tmp/paper_test/paper.css`：

```css
h1 {
    font-family: "SimHei", "Arial";
    font-size: 22pt;
    font-weight: bold;
    text-align: center;
    margin-top: 24pt;
    margin-bottom: 12pt;
}

h2 {
    font-family: "SimHei", "Arial";
    font-size: 16pt;
    font-weight: bold;
    margin-top: 12pt;
    margin-bottom: 6pt;
}

p {
    font-family: "SimSun", "Times New Roman";
    font-size: 12pt;
    text-indent: 2em;
    line-height: 1.75;
    text-align: justify;
}

blockquote {
    font-family: "KaiTi", "Georgia";
    font-size: 10.5pt;
    margin-left: 2cm;
    margin-right: 2cm;
}
```

创建 `/tmp/paper_test/abstract.md`：

```markdown
---
title: 融合媒体环境下新闻生产流程重构研究
author: 杨志宏
date: 2026年4月
toc: true
heading_numbering: true
header_left: "{title}"
footer_center: "第 {page} 页"
---

# 摘要

本文探讨了融合媒体环境下新闻生产流程的重构路径。通过对比传统媒体与新媒体的生产模式，提出了基于技术驱动的流程再造方案。

# 关键词

融合媒体；新闻生产；流程重构；技术驱动
```

创建 `/tmp/paper_test/ch1.md`：

```markdown
# 引言

## 研究背景

随着数字技术的飞速发展，传统媒体面临着前所未有的转型压力。融合媒体成为行业共识，但如何在实践中有效推进，仍是亟待解决的课题。

## 研究目的

本研究旨在系统梳理融合媒体环境下新闻生产流程的变化规律，提出切实可行的重构方案。

## 文献综述

> 融合新闻是指综合运用多媒体手段进行新闻报道的实践形态。（引用格式测试）
```

创建 `/tmp/paper_test/ch2.md`：

```markdown
# 研究方法

## 案例分析

选取国内外 10 家典型融合媒体机构，对其新闻生产流程进行深度剖析。

## 深度访谈

对 20 位资深媒体从业者进行半结构化访谈，了解实际操作中的痛点和创新实践。
```

创建 `/tmp/paper_test/ch3.md`：

```markdown
# 结论与展望

## 主要结论

1. 技术驱动是融合媒体流程重构的核心动力
2. 组织架构调整需要与技术升级同步推进
3. 人才培养模式需要根本性变革

## 未来展望

后续研究将进一步探索人工智能技术在新闻生产中的应用前景。
```

创建 `/tmp/paper_test/contents.yaml`：

```yaml
files:
  - abstract.md
  - ch1.md
  - ch2.md
  - ch3.md
```

执行转换：

```bash
uv run md2docx merge /tmp/paper_test/contents.yaml \
  --css /tmp/paper_test/paper.css \
  -o /tmp/paper_test/论文终稿.docx
```

注意：merge 命令目前不支持 `--css` 参数，改用两步走：

```bash
# 步骤一：先从 CSS 生成 reference.docx
uv run md2docx styles build /tmp/paper_test/paper.css \
  -o /tmp/paper_test/reference.docx

# 步骤二：将 reference.docx 放入自定义样式目录
mkdir -p ~/.md2docx/styles/paper/filters
cp /tmp/paper_test/reference.docx ~/.md2docx/styles/paper/
cat > ~/.md2docx/styles/paper/style.yaml << 'EOF'
name: paper
reference_doc: reference.docx
post_processing:
  toc: false
  heading_numbering: false
  image_width_pct: 80
EOF

# 步骤三：用该样式合并转换
uv run md2docx merge /tmp/paper_test/contents.yaml \
  --style paper \
  -o /tmp/paper_test/论文终稿.docx
```

**验收标准**：

打开 `/tmp/paper_test/论文终稿.docx`，逐项检查：

| 检查项 | 预期结果 |
|--------|---------|
| 目录 | 文档开头有"目录"和 TOC 域代码 |
| 标题编号 | 1 摘要、2 关键词、3 引言、3.1 研究背景 … |
| 一级标题字体 | 黑体（SimHei）或 Arial，22pt，居中 |
| 二级标题字体 | 黑体或 Arial，16pt |
| 正文字体 | 宋体（SimSun）或 Times New Roman，12pt |
| 正文首行缩进 | 约 2 个字符 |
| 正文行距 | 1.75 倍 |
| 引用段落 | 楷体或 Georgia，10.5pt，左右缩进 |
| 页眉左侧 | "融合媒体环境下新闻生产流程重构研究" |
| 页脚 | "第 X 页" |
| 内容完整性 | 摘要 → 关键词 → 引言 → 方法 → 结论，四章齐全 |

测试后清理：

```bash
rm -rf ~/.md2docx/styles/paper
rm -rf /tmp/paper_test
rm -rf /tmp/merge_test
rm -rf /tmp/test_*
rm -f /tmp/empty.css /tmp/bad_merge.yaml
```

---

## 自动测试（可选）

所有自动化测试：

```bash
uv run pytest -v                    # 全部 54 个测试
uv run pytest --cov=md2docx         # 带覆盖率
uv run pytest -m "not integration"  # 仅单元测试（不需要 pandoc）
```
