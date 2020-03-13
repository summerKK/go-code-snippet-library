package logic

import (
	"math"
	"summer/blogger/dal/db"
	"summer/blogger/model"
)

func GetArticleRecordList(page int, pageSize int) (list []*model.ArticleRecord, err error) {
	// 获取文章列表
	articleList, err := db.ArticleList(page, pageSize)
	if err != nil {
		return
	}
	// 获取分类id
	categoryIds := getCategoryIds(articleList)
	categoryList, err := db.CategoryList(categoryIds)
	if err != nil {
		return
	}

	for _, info := range articleList {
		temp := &model.ArticleRecord{
			ArticleInfo: *info,
		}
		for _, category := range categoryList {
			if info.CategoryId == category.Id {
				temp.Category = *category
				break
			}
		}
		list = append(list, temp)
	}

	return
}

func ArticleInsert(username, title, content string, categoryId int64) (err error) {
	contentLen := []rune(content)
	min := int(math.Min(float64(len(contentLen)), 128))
	summary := string(contentLen[:min])
	articleDetail := &model.ArticleDetail{
		ArticleInfo: model.ArticleInfo{
			CategoryId: categoryId,
			Title:      title,
			Summary:    summary,
		},
		Content: content,
	}

	_, err = db.ArticleInsert(articleDetail)
	return
}

func getCategoryIds(list []*model.ArticleInfo) (categoryIds []int64) {
loop:
	for _, record := range list {
		for _, id := range categoryIds {
			if record.CategoryId == id {
				continue loop
			}
		}
		categoryIds = append(categoryIds, record.CategoryId)
	}
	return
}
