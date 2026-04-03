# Analyze: Phase 9 - Polish & Release

**Input**: [tasks.md](./tasks.md) Phase 9 (T059-T065)
**Date**: 2026-04-03

## Pre-Implementation Checklist

### 错误处理 (T059)

- [ ] 样式 YAML 格式错误时，报错信息包含文件名和具体错误
- [ ] 图片文件不存在时，stderr 警告并继续转换（不中断）
- [ ] 输出目录不存在时，自动创建
- [ ] 空 Markdown 文件不报错，生成仅含页面设置的 docx
- [ ] BOM 开头的文件正确处理
- [ ] Markdown 文件不存在时，清晰报错并退出码非零

### CLI help 文本 (T060)

- [ ] root 命令 help 包含用法示例
- [ ] merge 命令 help 已有详细示例（已完成）
- [ ] styles 命令 help 包含子命令说明
- [ ] --version 输出版本号
- [ ] 不存在的样式名报错时列出可用样式（已完成）

### CI (T061)

- [ ] GitHub Actions workflow：push/PR 触发
- [ ] 步骤：checkout → Go setup → test → lint
- [ ] Go 版本矩阵或固定最新稳定版

### Release (T062)

- [ ] Makefile release 已有 3 平台交叉编译（已完成）
- [ ] 添加 sha256 校验和生成
- [ ] 添加版本号注入（ldflags -X）

### README (T063)

- [ ] 安装方式（下载二进制 / go install）
- [ ] 快速开始（基本用法）
- [ ] 样式系统说明
- [ ] CLI 参数参考
- [ ] 已知限制

### 测试覆盖 (T064)

当前覆盖率：
| 包 | 覆盖率 | 目标 |
|----|--------|------|
| cli | 22.4% | 60%+ |
| merge | 91.2% | ✅ |
| ooxml | 30.2% | 60%+ |
| parser | 68.4% | ✅ |
| render | 48.8% | 60%+ |
| style | 59.2% | ✅ |

重点补测：cli（转换主流程）、ooxml（writer/styles）、render（renderer 主逻辑）

### 手工验收 (T065)

- [ ] 用 testdata/academic-sample.md 转换，Word 打开全面检查
- [ ] 检查项：标题编号、TOC、封面、页眉页脚、表格、图片、CJK 字体

## Execution Order

1. T059 错误处理 + T060 CLI help（并行）
2. T062 版本号注入 + sha256
3. T064 补测试覆盖率
4. T061 CI workflow
5. T063 README
6. T065 手工验收
