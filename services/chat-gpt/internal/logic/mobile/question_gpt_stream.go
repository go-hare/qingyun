package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/common/openai"
	"qingyun/services/chat-gpt/models"
	chat_gpt "qingyun/services/chat-gpt/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileChatGptService) QuestionGptStream(ctx context.Context, request *chat_gpt.QuestionGptStreamRequest, stream chat_gpt.MobileChatGptService_QuestionGptStreamStream) (err error) {
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "QuestionGptStream",
	})
	response := &chat_gpt.QuestionGptStreamResponse{}
	if request.ModelId < 1 || request.Question == "" {
		logger.WithFields(log.Fields{
			"ErrorType": "Param",
			"Function":  "QuestionGptStream",
		}).Error("Param Error")
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_param_error)
		stream.Send(response)
		return nil
	}
	model, err := models.GetModel(func(session *xorm.Session) *xorm.Session {
		return session.Where("id = ?", request.ModelId).Where("if_del = ?", 0)
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetModel",
		}).Error(err)
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_internal_error)
		stream.Send(response)
		return nil
	}

	if model == nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "GetModel",
		}).Error("Model Not Found")
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_param_error)
		stream.Send(response)
		return nil
	}
	resp, err := m.SvcContext.OpenAiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: request.Question,
				},
			},
		},
	)

	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "OpenAi",
			"Function":  "CreateChatCompletion",
		}).Error(err)
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_internal_error)
		stream.Send(response)
		return nil
	}
	response.Answer = resp.Choices[0].Message.Content
	stream.Send(response)
	return nil
}
