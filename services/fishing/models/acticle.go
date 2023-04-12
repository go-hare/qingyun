package models

import "qingyun/common/store/mysql"

type Article struct {
	mysql.Model       `xorm:"extends"`
	CategoryId        int64    `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	UserId            int64    `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	Title             string   `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"title"`
	CreateTime        int64    `xorm:"not null default 0 INT(11) comment('创建时间 ')" json:"create_time"`
	ArticleStatus     int64    `xorm:"not null default 0 INT(11) comment('状态 0，待审核 1，通过 -1 驳回')" json:"article_status"`
	AuthorTime        int64    `xorm:"not null default 0 INT(11) comment('审核时间')" json:"author_time"`
	RejectTime        int64    `xorm:"not null default 0 INT(11) comment('驳回时间 ')" json:"reject_time"`
	AuthorAdminId     int64    `xorm:"not null default 0 INT(11) comment('审核人 ')" json:"author_admin_id"`
	LikeCount         int64    `xorm:"not null default 0 INT(11) comment('点赞数量 ')" json:"like_count"`
	CommentCount      int64    `xorm:"not null default 0 INT(11) comment('评论')" json:"comment_count"`
	CollectCount      int64    `xorm:"not null default 0 INT(11) comment('收藏数量')" json:"collect_count"`
	ShareCount        int64    `xorm:"not null default 0 INT(11) comment('分享')" json:"share_count"`
	ViewCount         int64    `xorm:"not null default 0 INT(11) comment('展示数量')" json:"view_count"`
	ListTags          []string `xorm:"not null default '' VARCHAR(256) comment('标签')" json:"list_tag_names"`
	ImageUrl          string   `xorm:"not null default '' VARCHAR(256) comment('标题')" json:"image_url"`
	ImageWidth        int64    `xorm:"not null default 0 INT(11) comment('展示数量')" json:"image_width"`
	ImageHeight       int64    `xorm:"not null default 0 INT(11) comment('展示数量')" json:"image_height"`
	Content           string   `xorm:"not null  TEXT   comment('内容')" json:"content"`
	ListArticleTagIds []int64  `xorm:"not null default '' VARCHAR(256)  comment('标签id')" json:"list_article_tag_ids"`
	IsDel             int64    `xorm:"not null default 0 INT(11) comment('是否删除 0否 1是')" json:"is_del"`
	DelTime           int64    `xorm:"not null default 0 INT(11) comment('是否删除')" json:"del_time"`
}

func ListArticles(filter mysql.OrmFilter) (list []*Article, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetArticle(filter mysql.OrmFilter) (*Article, error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Article)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func UpdateArticle(filter mysql.OrmFilter, message *Article) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}

func CreateArticle(message *Article) (err error) {
	_, err = mysql.GetDB().Insert(message)
	return
}
