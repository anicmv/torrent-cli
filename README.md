# 🌊 torrent-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/anicmv/torrent-cli/pulls)

**一个功能强大的 BitTorrent 种子文件命令行工具**

[功能特性](#-功能特性) •
[快速开始](#-快速开始) •
[使用指南](#-使用指南) •
[开发文档](#-开发文档) •
[贡献指南](#-贡献指南)

---

## 📖 目录

- [功能特性](#-功能特性)
- [快速开始](#-快速开始)
    - [环境要求](#环境要求)
    - [安装](#安装)
    - [快速使用](#快速使用)
- [使用指南](#-使用指南)
    - [创建种子](#创建种子)
    - [查看信息](#查看信息)
    - [生成磁力链接](#生成磁力链接)
    - [编辑种子](#编辑种子)
- [命令参考](#-命令参考)
- [项目结构](#-项目结构)
- [开发文档](#-开发文档)
- [常见问题](#-常见问题)
- [贡献指南](#-贡献指南)
- [许可证](#-许可证)

---

## ✨ 功能特性

- 🚀 **快速创建** - 支持单文件和多文件夹种子创建
- 🔍 **信息查看** - 详细展示种子文件元信息
- 🧲 **磁力链接** - 一键生成 Magnet URI
- ✏️ **种子编辑** - 修改种子文件的各种属性
- 📊 **智能分片** - 自动计算最优 piece length
- 🎨 **友好界面** - 清晰的命令行输出
- 🔒 **隐私保护** - 支持 Private Flag 设置
- 🌐 **跨平台** - 支持 Windows、Linux、macOS

---

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- Git

### 安装

#### 方式一：从源码编译

```bash
# 克隆仓库
git clone https://github.com/anicmv/torrent-cli.git
cd torrent-cli

# 安装依赖
go mod download

# 编译
go build -o torrent-cli

# 验证安装
./torrent-cli --version


#### 方式二：使用 go install

```bash
go install github.com/anicmv/torrent-cli@latest
```

#### 方式三：下载预编译二进制

前往 [Releases](https://github.com/anicmv/torrent-cli/releases) 页面下载对应平台的二进制文件。

### 快速使用

```bash
# 创建种子文件
./torrent-cli create myfile.txt -o myfile.torrent

# 查看种子信息
./torrent-cli info myfile.torrent

# 生成磁力链接
./torrent-cli magnet myfile.torrent
```

---

## 📚 使用指南

### 创建种子

#### 基础用法

```bash
# 为单个文件创建种子
./torrent-cli create file.txt

# 为文件夹创建种子
./torrent-cli create /path/to/directory

# 指定输出路径
./torrent-cli create file.txt -o output.torrent
```

#### 高级选项

```bash
./torrent-cli create file.txt \
  -o output.torrent \
  -a "http://tracker.example.com:8080/announce" \
  -c "这是一个测试种子" \
  -n "自定义名称" \
  -s "MySource" \
  --created-by "torrent-cli v1.0" \
  --publisher "anicmv" \
  -p
```

**参数说明：**

| 参数                   | 简写   | 说明          | 默认值                   |
|----------------------|------|-------------|-----------------------|
| `--output`           | `-o` | 输出文件路径      | `<name>.torrent`      |
| `--announce`         | `-a` | Tracker URL | `https://example.com` |
| `--comment`          | `-c` | 种子注释        | `torrent-cli`        |
| `--name`             | `-n` | 种子名称        | 输入文件名                 |
| `--source`           | `-s` | 来源标签        | `anicmv :)`           |
| `--private`          | `-p` | 私有种子标志      | `false`               |
| `--created-by`       | -    | 创建者信息       | `torrent-cli`        |
| `--publisher`        | -    | 发布者信息       | `anicmv :)`           |
| `--no-created-by`    | -    | 不包含创建者      | `false`               |
| `--no-creation-date` | -    | 不包含创建日期     | `false`               |
| `--no-publisher`     | -    | 不包含发布者      | `false`               |
| `--no-source`        | -    | 不包含来源标签     | `false`               |

#### 示例场景

**场景 1：创建私有种子**

```bash
./torrent-cli create movie.mp4 \
  -a "https://private-tracker.com/announce" \
  -p \
  -s "PrivateTracker" \
  -c "高清电影"
```

**场景 2：批量创建种子**

```bash
#!/bin/bash
for file in *.mp4; do
  ./torrent-cli create "$file" -o "torrents/${file%.mp4}.torrent"
done
```

**场景 3：创建多 Tracker 种子**

```bash
./torrent-cli create file.txt \
  -a "http://tracker1.com/announce,http://tracker2.com/announce"
```

---

### 查看信息

```bash
# 查看种子基本信息
./torrent-cli info myfile.torrent

# 输出示例：
# Metafile:         myfile.torrent
# InfoHash:         v1: a1b2c3d4e5f6...
# Announce:         http://tracker.example.com:8080/announce
# Name:             myfile.txt
# Piece size:       32768
# Length:           1024
```

---

### 生成磁力链接

```bash
# 生成磁力链接
./torrent-cli magnet myfile.torrent

# 输出示例：
# magnet:?xt=urn:btih:a1b2c3d4e5f6...&dn=myfile.txt&tr=http%3A%2F%2Ftracker.example.com%3A8080%2Fannounce
```

**使用磁力链接：**

```bash
# 复制到剪贴板（macOS）
./torrent-cli magnet myfile.torrent | pbcopy

# 复制到剪贴板（Linux）
./torrent-cli magnet myfile.torrent | xclip -selection clipboard

# 保存到文件
./torrent-cli magnet myfile.torrent > magnet.txt
```

---

### 编辑种子

```bash
# 修改 Announce URL
./torrent-cli edit myfile.torrent -a "http://new-tracker.com/announce"

# 修改种子名称
./torrent-cli edit myfile.torrent -n "新名称"

# 添加注释
./torrent-cli edit myfile.torrent -c "更新的注释"

# 设置为私有种子
./torrent-cli edit myfile.torrent -p yes

# 移除来源标签
./torrent-cli edit myfile.torrent --no-source
```

---

## 📋 命令参考

### 全局选项

```bash
./torrent-cli [command] [flags]
```

| 选项              | 说明     |
|-----------------|--------|
| `-h, --help`    | 显示帮助信息 |
| `-v, --version` | 显示版本信息 |

### 命令列表

| 命令       | 说明     |
|----------|--------|
| `create` | 创建种子文件 |
| `info`   | 查看种子信息 |
| `magnet` | 生成磁力链接 |
| `edit`   | 编辑种子文件 |

### 获取帮助

```bash
# 查看全局帮助
./torrent-cli --help

# 查看特定命令帮助
./torrent-cli create --help
./torrent-cli info --help
```

---

## 📁 项目结构

```
torrent-cli/
├── cmd/                    # 命令行命令
│   ├── root.go            # 根命令
│   ├── create.go          # 创建命令
│   ├── info.go            # 信息命令
│   ├── magnet.go          # 磁力链接命令
│   └── edit.go            # 编辑命令
├── pkg/                    # 核心包
│   ├── bencode/           # BEncode 编解码
│   │   ├── bencode.go     # 主接口
│   │   ├── decoder.go     # 解码器
│   │   └── encoder.go     # 编码器
│   ├── torrent/           # 种子处理
│   │   ├── torrent.go     # 种子结构
│   │   ├── creator.go     # 创建器
│   │   ├── reader.go      # 读取器
│   │   └── metainfo.go    # 元信息
│   └── utils/             # 工具函数
│       └── tree_printer.go # 目录树打印
├── main.go                 # 程序入口
├── go.mod                  # Go 模块文件
├── go.sum                  # 依赖校验文件
├── README.md               # 项目文档
├── LICENSE                 # 许可证
└── .gitignore             # Git 忽略文件
```

---

## 🛠 开发文档

### 环境搭建

```bash
# 克隆项目
git clone https://github.com/anicmv/torrent-cli.git
cd torrent-cli

# 安装依赖
go mod download

# 运行测试
go test ./...

# 本地运行
go run main.go create test.txt
```

### 使用 GoLand 开发

1. **打开项目**
    - File → Open → 选择 `torrent-cli` 目录

2. **配置运行**
    - Run → Edit Configurations
    - 添加 Go Build 配置
    - Program arguments: `create test.txt -o test.torrent`

3. **调试**
    - 在代码行号左侧点击设置断点
    - 右键选择 Debug 'go build main.go'

### 添加新命令

1. 在 `cmd/` 目录创建新文件，例如 `verify.go`

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [torrent file]",
	Short: "Verify torrent file integrity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 实现验证逻辑
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
```

2. 重新编译运行

```bash
go build -o torrent-cli
./torrent-cli verify test.torrent
```

### 代码规范

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

```bash
# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test -v ./...
```

---

## ❓ 常见问题

### Q1: 如何创建多 Tracker 种子？

**A:** 使用逗号分隔多个 Tracker URL：

```bash
./torrent-cli create file.txt -a "http://tracker1.com/announce,http://tracker2.com/announce"
```

### Q2: 种子文件的 piece length 如何确定？

**A:** 程序会根据文件大小自动计算最优的 piece length（16KB - 16MB 之间）。你也可以手动指定：

```bash
./torrent-cli create file.txt --piece-size 262144  # 256KB
```

### Q3: 如何创建私有种子？

**A:** 使用 `-p` 或 `--private` 参数：

```bash
./torrent-cli create file.txt -p
```

### Q4: 支持哪些操作系统？

**A:** 支持所有 Go 语言支持的平台：

- Windows (amd64, 386, arm64)
- Linux (amd64, 386, arm, arm64)
- macOS (amd64, arm64)

### Q5: 如何批量处理文件？

**A:** 使用 Shell 脚本：

```bash
#!/bin/bash
for file in *.mp4; do
    ./torrent-cli create "$file" -o "torrents/${file%.mp4}.torrent"
done
```

### Q6: 生成的种子文件在哪里？

**A:**

- 如果指定了 `-o` 参数，在指定的路径
- 否则在当前目录，文件名为 `<输入文件名>.torrent`

### Q7: 如何验证种子文件是否正确？

**A:** 使用 `info` 命令查看种子信息：

```bash
./torrent-cli info myfile.torrent
```

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork 项目**

2. **创建特性分支**
   ```bash
   git checkout -b feature/AmazingFeature
   ```

3. **提交更改**
   ```bash
   git commit -m 'Add some AmazingFeature'
   ```

4. **推送到分支**
   ```bash
   git push origin feature/AmazingFeature
   ```

5. **提交 Pull Request**

### 贡献类型

- 🐛 报告 Bug
- 💡 提出新功能
- 📝 改进文档
- 🎨 优化代码
- ✅ 添加测试

### 代码审查标准

- 代码符合 Go 语言规范
- 包含必要的测试
- 更新相关文档
- 通过所有 CI 检查
---

## 🗺 Roadmap

- [x] 基础种子创建功能
- [x] 种子信息查看
- [x] 磁力链接生成
- [ ] 种子编辑功能
- [ ] 种子验证功能
- [ ] 支持 BT v2 协议
- [ ] Web UI 界面
- [ ] 批量处理模式
- [ ] 配置文件支持
- [ ] 多语言支持

---

## 📜 更新日志

### v1.0.0 (2025-10-18)

**新增功能：**

- ✨ 支持单文件和多文件种子创建
- ✨ 种子信息查看功能
- ✨ 磁力链接生成
- ✨ 自动计算最优 piece length

**改进：**

- 🎨 优化命令行输出
- 📝 完善文档

**修复：**

- 🐛 修复输出路径创建错误的问题

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) - 强大的 CLI 框架
- [BitTorrent Protocol](https://www.bittorrent.org/beps/bep_0003.html) - 协议规范
- 所有贡献者和支持者

---

## 📞 联系方式

- 作者：anicmv
- GitHub: [anicmv](https://github.com/anicmv)
- 问题反馈：[Issues](https://github.com/anicmv/torrent-cli/issues)

---


**如果这个项目对你有帮助，请给一个 ⭐️ Star！**

Made with ❤️ by anicmv