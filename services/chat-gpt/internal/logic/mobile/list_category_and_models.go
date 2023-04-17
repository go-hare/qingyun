package mobile

import (
	"context"
	log "github.com/sirupsen/logrus"
	"qingyun/services/chat-gpt/models"
	chat_gpt "qingyun/services/chat-gpt/proto/mobile"
	"xorm.io/xorm"
)

func (m *MobileChatGptService) ListCategoryAndModels(ctx context.Context, request *chat_gpt.ListCategoryAndModelsRequest, response *chat_gpt.ListCategoryAndModelsResponse) (err error) {
	response.Status = &chat_gpt.Status{}
	logger := log.WithFields(log.Fields{
		"Module": "Service",
		"Method": "ListCategoryAndModels",
	})
	ListCategorys, err := models.ListCategorys(func(session *xorm.Session) *xorm.Session {
		return session.Where("if_del = ?", 0).OrderBy("wgiht desc")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListCategorys",
		}).Error(err)
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_internal_error)
		return nil
	}
	if len(ListCategorys) < 1 {
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_param_error)
		return nil
	}

	categoryId := []int64{}
	for i := 0; i < len(ListCategorys); i++ {
		categoryId = append(categoryId, ListCategorys[i].Id)
	}
	ListModels, err := models.ListModels(func(session *xorm.Session) *xorm.Session {
		return session.Where("if_del = ?", 0).In("category_id", categoryId).OrderBy("wgiht desc")
	})
	if err != nil {
		logger.WithFields(log.Fields{
			"ErrorType": "Database",
			"Function":  "ListModels",
		}).Error(err)
		chat_gpt.ServiceStatus(response.Status, chat_gpt.StatusCode_status_internal_error)
		return nil
	}
	ListModelsMap := map[int64][]*chat_gpt.Model{}
	for i := 0; i < len(ListModels); i++ {
		listTutorial := []*chat_gpt.Tutorial{}
		for j := 0; j < len(ListModels[i].Tutorial); j++ {
			listTutorial = append(listTutorial, &chat_gpt.Tutorial{
				Name: ListModels[i].Tutorial[j].Name,
				Desc: ListModels[i].Tutorial[j].Desc,
			})
		}
		model := &chat_gpt.Model{
			ModelId:      ListModels[i].Id,
			Name:         ListModels[i].Name,
			CategoryId:   ListModels[i].CategoryId,
			Icon:         ListModels[i].Icon,
			Title:        ListModels[i].Title,
			Desc:         ListModels[i].Desc,
			CreateTime:   ListModels[i].CreateTime,
			ListTutorial: listTutorial,
		}
		ListModelsMap[ListModels[i].CategoryId] = append(ListModelsMap[ListModels[i].CategoryId], model)
	}
	for i := 0; i < len(ListCategorys); i++ {
		if v, ok := ListModelsMap[ListCategorys[i].Id]; ok {
			response.List = append(response.List, &chat_gpt.Category{
				CategoryId: ListCategorys[i].Id,
				Name:       ListCategorys[i].Name,
				ListModel:  v,
			})
		}
	}
	return nil
}
