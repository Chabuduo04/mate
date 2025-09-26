# Mate - AI角色扮演语音聊天

一个支持语音交互的AI角色扮演聊天平台，包含后端Go服务和前端React应用。

## 项目结构

```
mate/
├── back_end/          # Go后端服务
│   ├── handlers/         # HTTP处理器
│   ├── models/           # 数据模型
│   ├── services/         # 业务逻辑服务
│   ├── config/           # 配置文件
│   ├── cmd/              # 主程序入口
│   └── roles.json        # 角色配置
├── frontend/          # React前端应用
│   ├── src/              # 源代码
│   ├── public/           # 静态资源
│   └── package.json      # 依赖配置
└── README.md          # 项目说明
```

## 功能特性

- 🎭 **多角色支持**: 哈利·波特、苏格拉底、虚拟面试官等
- 💬 **文本聊天**: 实时文本对话
- 🎤 **语音输入**: 语音录制和自动识别
- 🔊 **语音输出**: AI回复语音合成
- 🎙️ **语音聊天**: 完整的语音交互流程（上传→ASR→LLM→TTS→返回）
- 📱 **响应式设计**: 支持桌面和移动设备
- 🎨 **现代化UI**: 美观的用户界面

### 语音聊天流程

1. **语音录制**: 用户点击麦克风按钮录制语音
2. **文件上传**: 音频文件上传到云存储，生成唯一URL
3. **语音识别**: 调用ASR服务将语音转换为文字
4. **AI处理**: 调用LLM服务处理文字并生成回复
5. **语音合成**: 调用TTS服务将回复转换为语音
6. **结果返回**: 同时返回文字回复和语音文件给前端

## 快速开始

### 后端启动

```bash
cd back_end
go mod tidy
go run cmd/main.go
```

后端服务将在 http://localhost:8080 启动。

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端应用将在 http://localhost:3000 启动。

或者使用启动脚本：

```bash
cd frontend
./start.sh
```

## API接口

- `GET /api/roles` - 获取角色列表
- `POST /api/llm` - 发送聊天消息
- `POST /api/asr` - 语音识别
- `POST /api/tts` - 语音合成
- `POST /api/upload` - 文件上传
- `POST /api/voice-chat` - 语音聊天（完整流程：上传→ASR→LLM→TTS→返回）

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- 自定义服务层

### 前端
- React 18
- TypeScript
- Vite
- Tailwind CSS
- Axios

## 测试

运行测试脚本验证语音聊天功能：

```bash
./test_voice_chat.sh
```

测试脚本会检查：
- 后端和前端服务状态
- 角色列表API
- 文本聊天API
- 语音聊天API

## 开发说明

1. 确保已安装Go 1.21+和Node.js 18+
2. 后端和前端需要同时运行
3. 前端通过代理访问后端API
4. 语音功能需要HTTPS环境或localhost
5. 需要配置云存储服务（七牛云）用于文件上传
6. 需要配置ASR和TTS服务的API密钥

## 环境变量配置

后端需要配置以下环境变量：

```bash
# 服务配置
PORT=8080
REDIS_ADDR=localhost:6379

# API配置
API_KEY=your_api_key
API_URL=your_api_url
LLM_MODEL=your_llm_model

# ASR和TTS服务
ASR_ENDPOINT=your_asr_endpoint
TTS_ENDPOINT=your_tts_endpoint

# 七牛云存储
KODO_HOST=your_kodo_host
ACCESS_KEY=your_access_key
SECRET_KEY=your_secret_key
BUCKET=your_bucket_name
```

## 许可证

MIT License
