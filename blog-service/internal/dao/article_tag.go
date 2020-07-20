package dao

import "github.com/summerKK/go-code-snippet-library/blog-service/internal/model"

func (d *Dao) GetArticleTagByAID(articleID uint32) (*model.ArticleTag, error) {
	articleTag := model.ArticleTag{ArticleId: articleID}

	return articleTag.GetByAID(d.engine)
}

func (d *Dao) GetArticleListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{TagId: tagID}

	return articleTag.ListByTID(d.engine)
}

func (d *Dao) GetArticleListByAIDs(aIDs []uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{}

	return articleTag.ListByAIDs(d.engine, aIDs)
}

func (d *Dao) CreateArticleTag(articleID uint32, tagID uint32, createdBy string) error {
	articleTag := model.ArticleTag{
		TagId:     tagID,
		ArticleId: articleID,
		Model: &model.Model{
			CreatedBy: createdBy,
		},
	}

	return articleTag.Create(d.engine)
}

func (d *Dao) UpdateArticleTag(articleID, tagID uint32, modifiedBy string) error {
	articleTag := model.ArticleTag{
		ArticleId: articleID,
	}

	values := map[string]interface{}{
		"tag_id":      tagID,
		"modified_by": modifiedBy,
	}

	return articleTag.UpdateOne(d.engine, values)
}

func (d *Dao) DeleteArticleTag(articleID uint32) error {
	articleTag := model.ArticleTag{
		ArticleId: articleID,
	}

	return articleTag.DeleteOne(d.engine)
}
