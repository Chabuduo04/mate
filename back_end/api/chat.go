package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/Chabuduo04/mate/back_end/constants"
	"github.com/Chabuduo04/mate/back_end/dto/request"
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/Chabuduo04/mate/back_end/zlog"
	"github.com/gin-gonic/gin"
)

type ChatApi struct {
	Service services.Services
}

func NewChatApi(service services.Services) ChatApi {
	return ChatApi{Service: service}
}

func (c *ChatApi) ChatWithText(ctx *gin.Context) {
	var req request.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}

	if req.UserID == "" {
		req.UserID = ctx.GetString("user_id")
	}

	message, chatRsp, ret := services.ChatService.Chat(req, &c.Service)

	JsonBack(ctx, message, ret, chatRsp)
}

// todo: 处理webm格式的音频文件,目前只支持wav和mp3，统一返回格式
func (c *ChatApi) ChatWithAudio(ctx *gin.Context) {
	// 获取上传的音频文件
	file, header, err := ctx.Request.FormFile("audio")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
		return
	}
	defer file.Close()

	// 上传音频文件
	audioURL, err := c.Service.Storage.UploadAudio(file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload audio"})
		return
	}
	c.Service.Logger.Sugar().Infof("audio uploaded to: %s", audioURL)
	// 将音频文件转换为文本
	audioText, err := c.Service.ASR.TranscribeFromURL(audioURL)
	if err != nil {
		c.Service.Logger.Sugar().Errorf("ASR error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "audio to text error"})
		return
	}

	// 使用转换后的文本进行聊天
	ctx.Request.Body = io.NopCloser(strings.NewReader(audioText))
	c.ChatWithText(ctx)
	ctx.JSON(200, gin.H{"transcribed_text": audioText})
}

func (c *ChatApi) ChatList(ctx *gin.Context) {
	var req request.GetMessageListRequest
	if err := ctx.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	if req.UserId == "" {
		req.UserId = ctx.GetString("user_id")
	}
	message, rsp, ret := services.MessageService.GetMessageList(req.UserId, req.RoleId)
	JsonBack(ctx, message, ret, rsp)
}
