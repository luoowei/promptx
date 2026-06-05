<p align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/English-EN-blue" alt="English"></a>
  <a href="README_CN.md"><img src="https://img.shields.io/badge/简体中文-CN-red" alt="简体中文"></a>
</p>

<h1 align="center">PromptX</h1>
<p align="center"><strong>AI 终端助手 — 问答、生成命令、调试错误。</strong></p>
<p align="center">
  <a href="README.md">English</a> | <a href="README_CN.md">简体中文</a>
</p>
<br>
<p align="center">
  <a href="https://github.com/luoowei/promptx/stargazers"><img src="https://img.shields.io/github/stars/luoowei/promptx?style=social" alt="Stars"></a>
  <a href="https://github.com/luoowei/promptx/blob/main/LICENSE"><img src="https://img.shields.io/github/license/luoowei/promptx" alt="License"></a>
  <a href="https://pkg.go.dev/github.com/luoowei/promptx"><img src="https://pkg.go.dev/badge/github.com/luoowei/promptx.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/luoowei/promptx"><img src="https://goreportcard.com/badge/github.com/luoowei/promptx" alt="Go Report Card"></a>
</p>

---

**PromptX** 将 AI 超能力搬进你的终端。无需离开命令行，就能提问、生成 shell 命令、调试错误、与 AI 模型对话。

- **单二进制文件** — 不需要 Python、Node.js 或 Docker。下载即用。
- **多模型支持** — OpenAI、Anthropic、Ollama、OpenRouter、Groq，以及任何 OpenAI 兼容的 API。
- **精美的 TUI 界面** — 交互式聊天模式，支持 markdown 渲染和语法高亮。
- **隐私优先** — API 密钥不会离开你的电脑。使用本地模型实现 100% 离线使用。
- **Shell 集成** — 通过管道直接调试错误：`command 2>&1 | px ask "修复这个问题"`

## 快速开始

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/luoowei/promptx/main/install.sh | bash

# Homebrew
brew install luoowei/tap/promptx

# Go
go install github.com/luoowei/promptx/cmd/px@latest

# Windows
iwr https://raw.githubusercontent.com/luoowei/promptx/main/install.ps1 | iex
```

设置 API 密钥并开始使用：

```bash
export OPENAI_API_KEY="sk-..."     # 或用 ANTHROPIC_API_KEY
px ask "如何找出目录中最大的文件？"
```

## 使用方式

### 交互式聊天模式
```bash
px
# 在终端中打开精美的 TUI 界面，与 AI 聊天
```

### 一次性问答
```bash
px ask "解释这个正则：^(?=.*[A-Z])(?=.*[a-z])(?=.*\d).{8,}$"
px ask "撤销最后一次 commit 但保留更改的 git 命令是什么？"
```

### Shell 命令生成
```bash
px sh "查找最近24小时内修改的所有 .log 文件并压缩它们"
px sh "在 ubuntu 上创建一个有 sudo 权限的新用户"
```

### 错误调试
```bash
# 通过管道传递错误
make build 2>&1 | px ask "出了什么问题？怎么修复？"

# 或直接粘贴错误信息
px ask "什么原因导致 'bind: address already in use'？如何修复？"
```

### 配置管理
```bash
px config show                          # 查看当前配置
px config set-provider anthropic        # 切换到 Claude
px config set-key openai sk-...         # 保存 API 密钥
px --provider ollama --model llama3.2   # 使用本地 Ollama 模型
```

## 交互模式功能

```
> 如何查看哪个进程占用了 3000 端口？

PromptX
要查找占用 3000 端口的进程：

Linux/macOS:
  lsof -i :3000
  或: ss -tlnp | grep :3000

Windows:
  netstat -ano | findstr :3000

要终止该进程：
  kill -9 <PID>            # Linux/macOS
  taskkill /PID <PID> /F   # Windows
```

命令：`/clear` `/model` `/help` · `Ctrl+D` 退出

## 支持的提供商

| 提供商 | 模型示例 | 设置方式 |
|----------|---------------|-------|
| **OpenAI** | gpt-4o, gpt-4.1 | `export OPENAI_API_KEY=sk-...` |
| **Anthropic** | claude-opus-4-8, claude-sonnet-4-6 | `export ANTHROPIC_API_KEY=sk-ant-...` |
| **Ollama** | llama3.2, mistral, qwen | `ollama pull llama3.2`（无需密钥） |
| **OpenRouter** | 250+ 模型 | `export OPENROUTER_API_KEY=...` |
| **Groq** | llama-3.3-70b | `export GROQ_API_KEY=...` |
| **任何 OpenAI 兼容** | 自定义端点 | `px config set-key custom <key>` |

## 工作原理

```
  终端             PromptX CLI        AI API
  (你)     -->     (Go 二进制)   -->  (OpenAI/
                    |                 Anthropic/
                    v                 Ollama)
                  响应   <--    后端
```

PromptX 是一个静态编译的 Go 二进制文件：
1. 从 CLI 参数、管道或 TUI 交互读取输入
2. 将请求发送到配置的 AI 提供商
3. 以 markdown 渲染格式返回响应

**运行时零依赖。** 不需要 Python venv。不需要 npm install。不需要 Docker pull。

## 为什么选择 PromptX？

| | PromptX | LLM CLI 封装 | Web 聊天 UI |
|---|:---:|:---:|:---:|
| **单一二进制** | ✅ | ❌ (需要 Python/Node) | ❌ (仅限浏览器) |
| **管道支持** | ✅ | ⚠️ 部分 | ❌ |
| **Shell 集成** | ✅ | ⚠️ | ❌ |
| **离线运行（本地模型）** | ✅ | ✅ | ❌ |
| **多提供商** | ✅ | ⚠️ | ✅ |
| **TUI 模式** | ✅ | ✅ | ❌ |
| **启动速度** | <50ms | 200-2000ms | 2-5s |
| **二进制大小** | ~15MB | 50-500MB+ | N/A |

## 从源码构建

```bash
# 需要 Go 1.22+
git clone https://github.com/luoowei/promptx.git
cd promptx
go build -o px ./cmd/px
./px ask "你好，世界"
```

## 项目结构

```
promptx/
├── cmd/px/main.go          # 入口文件
├── internal/
│   ├── ai/                 # AI 提供商客户端
│   │   ├── client.go       # 客户端接口 + 工厂
│   │   ├── openai.go       # OpenAI 兼容客户端
│   │   └── anthropic.go    # Anthropic 客户端
│   ├── cli/root.go         # Cobra CLI 命令
│   ├── config/config.go    # 配置管理
│   ├── shell/handler.go    # Shell 命令生成
│   └── tui/
│       ├── tui.go          # Bubble Tea UI 模型
│       └── renderer.go     # Markdown 渲染器
├── docs/                   # 展示页面（GitHub Pages）
└── install.sh              # 一键安装脚本
```

## 路线图

- [ ] TUI 模式中的流式响应
- [ ] Shell 自动补全脚本（bash, zsh, fish, powershell）
- [ ] 会话历史搜索
- [ ] 自定义系统提示
- [ ] MCP 服务端集成
- [ ] 额外提供商的插件系统
- [ ] Windows 终端主题

## 参与贡献

欢迎贡献！详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 开源协议

MIT © PromptX 贡献者 — 详见 [LICENSE](LICENSE)

---

<p align="center">
  <b>⭐ 如果觉得有用，给这个仓库点个 Star！</b><br>
  <sub>使用 Go、Bubble Tea 和热爱为终端而建</sub>
</p>
