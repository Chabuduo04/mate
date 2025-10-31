package services

import (
	"github.com/Chabuduo04/mate/back_end/dto/request"
	"github.com/Chabuduo04/mate/back_end/dto/respond"
	"github.com/Chabuduo04/mate/back_end/zlog"
)

type chatService struct {
}

var ChatService = new(chatService)

// Chat 聊天接口
func (m *chatService) Chat(req request.ChatRequest, svc *Services) (string, *respond.ChatTextRespond, int) {
	// 检查角色是否存在
	role, ok := svc.RoleService.GetRole(req.RoleID)
	if !ok {
		return "角色不存在", nil, -2
	}

	//获取聊天上下文
	var messages []ChatMessage
	messages = append(messages, ChatMessage{Role: "system", Content: role.Prompt})

	_, recordList, _ := MessageService.GetMessageList(req.UserID, req.RoleID)
	for _, record := range recordList {
		messages = append(messages, ChatMessage{Role: "user", Content: record.Message})
		messages = append(messages, ChatMessage{Role: "assistant", Content: record.Response})
	}

	// call llm service
	reply, err := svc.LLM.Chat(messages)
	if err != nil {
		zlog.Error(err.Error())
		return "聊天失败，请稍后再试", nil, -1
	}

	chatRsp := respond.ChatTextRespond{}
	chatRsp.ReplyText = reply

	// 调用TTS服务
	voiceType := req.Voice
	if voiceType == "" {
		voiceType = role.VoiceType
	}
	audioBase64, err := svc.TTS.Synthesize(reply, voiceType)
	if err != nil {
		zlog.Error(err.Error())
		// TTS失败不影响文字回复，继续返回文字
		audioBase64 = ""
	} else if audioBase64 != "" {
		audioBase64 = "data:audio/mp3;base64," + audioBase64
	}
	chatRsp.AudioBase64 = audioBase64

	// 保存聊天记录
	if err := MessageService.AppendMessageRecord(req.UserID, req.RoleID, req.Message, reply, recordList); err != nil {
		zlog.Error(err.Error())
	}

	return "聊天成功", &chatRsp, 0
}
