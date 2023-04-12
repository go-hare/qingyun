package mobile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qingyun/services/fishing/models"
	"strconv"
	"testing"
	"time"
	"xorm.io/xorm"
)

func TestTheme(t *testing.T) {
	list, _ := models.ListArticleThemes(func(session *xorm.Session) *xorm.Session {
		return session
	})
	for i := 0; i < len(list); i++ {
		if list[i].UsedCount > 0 || list[i].ViewCount > 0 {
			models.UpdateArticleTag(func(session *xorm.Session) *xorm.Session {
				return session.Where("id = ?", list[i].TagId).Cols("used_count,view_count")
			}, &models.ArticleTag{ViewCount: list[i].ViewCount, UsedCount: list[i].UsedCount})
		}
	}

}

func Post(url string, data interface{}, contentType string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

//{"fkArticleTypeTid":"10001","sort":"create_date","page":1,"limit":10,"token":"","aid":"1297804114846748673","pid":"1285147583764086786","secret":"H2gDKM06H8QvNCNF"}

func TestActive(t *testing.T) {
	url := "https://gw.youyong.co/article/article/articleListType"
	//m := make(map[string]string)
	//m["fkArticleTypeTid"] = "10001"
	//m["sort"] = "create_date"
	//m["page"] = "1"
	//m["limit"] = "10"
	//m["aid"] = "1297804114846748673"
	//m["pid"] = "1285147583764086786"
	//m["secret"] = "H2gDKM06H8QvNCNF"
	//post := Post(url, m, "application/json")
	//dataMap := make(map[string]interface{})
	//json.Unmarshal([]byte(post), &dataMap)
	//dataListStr := dataMap["data"]
	//js, _ := json.Marshal(dataListStr)
	//list := []map[string]interface{}{}
	//json.Unmarshal(js, &list)
	//fmt.Println(list, len(list))
	//for j := 0; j < len(list); j++ {
	//	tid := list[j]["tid"].(string)
	//	fmt.Println(tid)
	//	GetActiveDetail(tid)
	//}

	ifhas := true
	page := 1
	for ifhas {
		m := make(map[string]string)
		m["fkArticleTypeTid"] = "10001"
		m["sort"] = "create_date"
		m["page"] = strconv.Itoa(page)
		m["limit"] = "10"
		m["aid"] = "1297804114846748673"
		m["pid"] = "1285147583764086786"
		m["secret"] = "H2gDKM06H8QvNCNF"
		post := Post(url, m, "application/json")
		dataMap := make(map[string]interface{})
		json.Unmarshal([]byte(post), &dataMap)
		dataListStr := dataMap["data"]
		js, _ := json.Marshal(dataListStr)
		list := []map[string]interface{}{}
		json.Unmarshal(js, &list)
		fmt.Println(list, len(list))
		for j := 0; j < len(list); j++ {
			tid := list[j]["tid"].(string)
			fmt.Println(tid)
			GetActiveDetail(tid)
		}
		page += 1
		if len(list) < 10 {
			ifhas = false
		}
	}

}

type Article struct {
	CollectCount     string `json:"collectCount"`
	CommentCount     string `json:"commentCount"`
	Content          string `json:"content"`
	ImageHeight      int64  `json:"imageHeight"`
	ImageUrl         string `json:"imageUrl"`
	ImageWidth       int64  `json:"imageWidth"`
	LikeCount        string `json:"likeCount"`
	ShareCount       string `json:"shareCount"`
	Title            string `json:"title"`
	ViewCount        string `json:"viewCount"`
	FkArticleTypeTid string `json:"fkArticleTypeTid"`
}

type Tag struct {
	Name      string `json:"name"`
	UsedCount string `json:"usedCount"`
	ViewCount string `json:"viewCount"`
}

type ArticleInfo struct {
	Data struct {
		Article *Article `json:"article"`
		TagList []*Tag   `json:"tagList"`
	} `json:"data"`
}

func GetActiveDetail(tid string) {
	url := "https://gw.youyong.co/article/article/detail"
	m := make(map[string]string)
	m["fkArticleTid"] = tid
	m["aid"] = "1297804114846748673"
	m["pid"] = "1285147583764086786"
	m["secret"] = "H2gDKM06H8QvNCNF"
	m["token"] = ""
	post := Post(url, m, "application/json")
	var articleInfo ArticleInfo
	json.Unmarshal([]byte(post), &articleInfo)
	tagList := articleInfo.Data.TagList
	tagName := []string{}
	for i := 0; i < len(tagList); i++ {
		tagName = append(tagName, tagList[i].Name)
	}
	fmt.Println(tagName)

	ListArticleTags, _ := models.ListArticleTags(func(session *xorm.Session) *xorm.Session {
		return session.In("name", tagName)
	})
	tagNameMap := make(map[string]struct{})
	tagIds := []int64{}
	tagNames := []string{}
	for i := 0; i < len(ListArticleTags); i++ {
		tagNameMap[ListArticleTags[i].Name] = struct{}{}
		tagIds = append(tagIds, ListArticleTags[i].Id)
		tagNames = append(tagNames, ListArticleTags[i].Name)
	}
	for i := 0; i < len(tagList); i++ {
		if _, ok := tagNameMap[tagList[i].Name]; !ok {
			usedCount, _ := strconv.Atoi(tagList[i].UsedCount)
			viewCount, _ := strconv.Atoi(tagList[i].ViewCount)
			articleTag := &models.ArticleTag{
				Name:       tagList[i].Name,
				CreateTime: time.Now().Unix(),
				UsedCount:  int64(usedCount),
				ViewCount:  int64(viewCount),
			}
			models.CreateArticleTag(articleTag)
			tagIds = append(tagIds, articleTag.Id)
			tagNames = append(tagNames, articleTag.Name)
		}
	}
	categoryId, _ := strconv.Atoi(articleInfo.Data.Article.FkArticleTypeTid)
	likeCount, _ := strconv.Atoi(articleInfo.Data.Article.LikeCount)
	commentCount, _ := strconv.Atoi(articleInfo.Data.Article.CommentCount)
	collectCount, _ := strconv.Atoi(articleInfo.Data.Article.CollectCount)
	shareCount, _ := strconv.Atoi(articleInfo.Data.Article.ShareCount)
	viewCount, _ := strconv.Atoi(articleInfo.Data.Article.ViewCount)
	article := &models.Article{
		CategoryId:        int64(categoryId),
		UserId:            1,
		Title:             articleInfo.Data.Article.Title,
		CreateTime:        time.Now().Unix(),
		ArticleStatus:     1,
		AuthorTime:        time.Now().Unix(),
		RejectTime:        0,
		AuthorAdminId:     1,
		LikeCount:         int64(likeCount),
		CommentCount:      int64(commentCount),
		CollectCount:      int64(collectCount),
		ShareCount:        int64(shareCount),
		ViewCount:         int64(viewCount),
		ListTags:          tagNames,
		ImageUrl:          articleInfo.Data.Article.ImageUrl,
		ImageWidth:        articleInfo.Data.Article.ImageWidth,
		ImageHeight:       articleInfo.Data.Article.ImageHeight,
		Content:           articleInfo.Data.Article.Content,
		ListArticleTagIds: tagIds,
	}

	models.CreateArticle(article)
	ListArticleTagDetail := []*models.ArticleTagDetail{}
	for i := 0; i < len(tagIds); i++ {
		ListArticleTagDetail = append(ListArticleTagDetail, &models.ArticleTagDetail{
			ArticleId:    article.Id,
			ArticleTagId: tagIds[i],
			CreateTime:   time.Now().Unix(),
		})
	}
	models.CreateListArticleTagDetails(ListArticleTagDetail)

}

func TestInfomation(t *testing.T) {

}
