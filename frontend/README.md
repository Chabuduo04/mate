# Mate Frontend

AI角色扮演语音聊天平台的前端应用。

## 功能特性

- 🎭 **角色选择**: 支持多个AI角色，每个角色都有独特的个性和技能
- 💬 **文本聊天**: 与AI角色进行实时文本对话
- 🎤 **语音输入**: 支持语音录制和自动语音识别(ASR)
- 🔊 **语音输出**: AI回复支持语音合成(TTS)播放
- 📱 **响应式设计**: 适配桌面和移动设备
- 🎨 **现代化UI**: 使用Tailwind CSS构建的美观界面

## 技术栈

- **React 18** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Tailwind CSS** - 样式框架
- **Axios** - HTTP客户端
- **Lucide React** - 图标库

## 开发

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
src/
├── components/          # React组件
│   ├── RoleSelector.tsx    # 角色选择组件
│   ├── ChatInterface.tsx   # 聊天界面组件
│   ├── MessageBubble.tsx   # 消息气泡组件
│   ├── ChatInput.tsx       # 聊天输入组件
│   └── LoadingSpinner.tsx  # 加载动画组件
├── hooks/              # 自定义Hooks
│   ├── useChat.ts         # 聊天逻辑Hook
│   └── useAudio.ts        # 音频处理Hook
├── services/           # API服务
│   └── api.ts             # API调用封装
├── types/              # TypeScript类型定义
│   └── index.ts           # 类型定义
├── utils/              # 工具函数
│   └── cn.ts              # 样式工具函数
├── App.tsx             # 主应用组件
├── main.tsx            # 应用入口
└── index.css           # 全局样式
```

## API集成

前端与后端API的集成点：

- `GET /api/roles` - 获取角色列表
- `POST /api/llm` - 发送聊天消息
- `POST /api/asr` - 语音识别
- `POST /api/tts` - 语音合成
- `POST /api/upload` - 文件上传

## 浏览器支持

- Chrome 88+
- Firefox 85+
- Safari 14+
- Edge 88+

## 注意事项

- 语音功能需要HTTPS环境或localhost
- 需要用户授权麦克风权限
- 建议使用现代浏览器以获得最佳体验
