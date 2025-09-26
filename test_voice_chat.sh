#!/bin/bash

# 测试语音聊天功能的脚本

echo "=== Mate 语音聊天功能测试 ==="

# 检查后端是否运行
echo "1. 检查后端服务状态..."
if curl -s http://localhost:8080/api/roles > /dev/null; then
    echo "✅ 后端服务运行正常"
else
    echo "❌ 后端服务未运行，请先启动后端服务"
    echo "   运行命令: cd back_end && go run cmd/main.go"
    exit 1
fi

# 检查前端是否运行
echo "2. 检查前端服务状态..."
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ 前端服务运行正常"
else
    echo "❌ 前端服务未运行，请先启动前端服务"
    echo "   运行命令: cd frontend && npm run dev"
    exit 1
fi

# 测试角色列表API
echo "3. 测试角色列表API..."
ROLES_RESPONSE=$(curl -s http://localhost:8080/api/roles)
if echo "$ROLES_RESPONSE" | grep -q "harry"; then
    echo "✅ 角色列表API正常"
else
    echo "❌ 角色列表API异常"
    echo "响应: $ROLES_RESPONSE"
fi

# 测试文本聊天API
echo "4. 测试文本聊天API..."
CHAT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/llm \
  -H "Content-Type: application/json" \
  -d '{"role_id": "harry", "message": "你好", "user_id": "test_user"}')

if echo "$CHAT_RESPONSE" | grep -q "reply_text"; then
    echo "✅ 文本聊天API正常"
else
    echo "❌ 文本聊天API异常"
    echo "响应: $CHAT_RESPONSE"
fi

# 测试语音聊天API（使用模拟音频文件）
echo "5. 测试语音聊天API..."
# 创建一个小的测试音频文件（实际项目中应该使用真实的音频文件）
echo "测试音频数据" > /tmp/test_audio.wav

VOICE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/voice-chat \
  -F "audio=@/tmp/test_audio.wav" \
  -F "role_id=harry" \
  -F "user_id=test_user")

if echo "$VOICE_RESPONSE" | grep -q "reply_text"; then
    echo "✅ 语音聊天API正常"
else
    echo "❌ 语音聊天API异常"
    echo "响应: $VOICE_RESPONSE"
fi

# 清理测试文件
rm -f /tmp/test_audio.wav

echo ""
echo "=== 测试完成 ==="
echo "如果所有测试都通过，说明语音聊天功能已正确实现！"
echo ""
echo "功能流程："
echo "1. 前端录制语音 → 上传到后端"
echo "2. 后端上传到云存储 → 生成唯一URL"
echo "3. 后端调用ASR → 将语音转为文字"
echo "4. 后端调用LLM → 处理文字并生成回复"
echo "5. 后端调用TTS → 将回复转为语音"
echo "6. 后端返回文字和语音给前端"
echo "7. 前端显示文字并播放语音"