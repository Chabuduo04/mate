package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Chabuduo04/mate/back_end/constants"
	"github.com/Chabuduo04/mate/back_end/dao"
	"github.com/Chabuduo04/mate/back_end/db"
	"github.com/Chabuduo04/mate/back_end/dto/respond"
	myredis "github.com/Chabuduo04/mate/back_end/redis"
	"github.com/Chabuduo04/mate/back_end/zlog"
	"github.com/redis/go-redis/v9"

	"github.com/Chabuduo04/mate/back_end/models"
)

type messageService struct {
}

var MessageService = new(messageService)

// GetMessageList 获取聊天记录
func (m *messageService) GetMessageList(userId, roleId string) (string, []respond.ChatRecordRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("message_list_" + userId + "_" + roleId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info(err.Error())
			zlog.Info(fmt.Sprintf("%s %s", userId, roleId))
			var messageList []models.ChatRecord
			if messageList, err = dao.GetChatRecordsByUserAndRole(db.GetConn(), userId, roleId); err != nil {
				zlog.Error(err.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			var rspList []respond.ChatRecordRespond
			for _, message := range messageList {
				rspList = append(rspList, respond.ChatRecordRespond{
					UserID:    message.UserID,
					RoleID:    message.RoleID,
					Message:   message.Message,
					Response:  message.Response,
					CreatedAt: message.CreatedAt.Unix(),
				})
			}
			rspString, err := json.Marshal(rspList)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("message_list_"+userId+"_"+roleId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取聊天记录成功", rspList, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp []respond.ChatRecordRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取聊天记录成功", rsp, 0
}

// AppendMessageRecord 添加聊天记录
func (m *messageService) AppendMessageRecord(userId, roleId, message, response string, recordList []respond.ChatRecordRespond) error {
	record := &models.ChatRecord{
		UserID:    userId,
		RoleID:    roleId,
		Message:   message,
		Response:  response,
		CreatedAt: time.Now(),
	}
	if err := dao.CreateChatRecord(db.GetConn(), record); err != nil {
		return err
	}
	// 更新缓存
	recordList = append(recordList, respond.ChatRecordRespond{
		UserID:    userId,
		RoleID:    roleId,
		Message:   message,
		Response:  response,
		CreatedAt: record.CreatedAt.Unix(),
	})
	rspString, err := json.Marshal(recordList)
	if err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.SetKeyEx("message_list_"+userId+"_"+roleId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
		zlog.Error(err.Error())
	}

	return nil
}
