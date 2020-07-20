package dao

import (
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/model"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type Article struct {
	ID             uint32 `json:"id"`
	TagID          uint32 `json:"tag_id"`
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Content        string `json:"content"`
	ConverImageUrl string `json:"conver_image_url"`
	CreatedBy      string `json:"created_by"`
	ModifiedBy     string `json:"modified_by"`
	State          uint8  `json:"state"`
}

func (d *Dao) GetArticle(id uint32, state uint8) (*model.Article, error) {
	article := model.Article{
		Model: &model.Model{ID: id},
		State: state,
	}

	return article.Get(d.engine)
}

func (d *Dao) GetArticleList(title string, state uint8, page, pageSize int) ([]*model.Article, error) {
	article := model.Article{Title: title, State: state}

	return article.List(d.engine, pageSize, app.GetPageOffset(page, pageSize))
}

func (d *Dao) CountArticle(title string, state uint8) (int, error) {
	article := model.Article{Title: title, State: state}

	return article.Count(d.engine)
}

func (d *Dao) CreateArticle(param *Article) error {
	article := model.Article{
		Title:          param.Title,
		Desc:           param.Desc,
		Content:        param.Content,
		ConverImageUrl: param.ConverImageUrl,
		State:          param.State,
		Model:          &model.Model{CreatedBy: param.CreatedBy},
	}

	return article.Create(d.engine)
}

func (d *Dao) UpdateArticle(param *Article) error {
	article := model.Article{
		Model: &model.Model{ID: param.ID},
	}

	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}

	if param.Title != "" {
		values["title"] = param.Title
	}

	if param.ConverImageUrl != "" {
		values["conver_image_url"] = param.ConverImageUrl
	}

	if param.Content != "" {
		values["content"] = param.Content
	}

	if param.Desc != "" {
		values["desc"] = param.Desc
	}

	return article.Update(d.engine, values)
}

func (d *Dao) DeleteArticle(id uint32) error {
	article := model.Article{Model: &model.Model{ID: id}}

	return article.Delete(d.engine)
}

func (d *Dao) CountArticleListByTagID(tagID uint32, state uint8) (int, error) {
	article := model.Article{State: state}

	return article.CountByTagID(d.engine, tagID)
}

func (d *Dao) GetArticleListByTagID(tagID uint32, state uint8, pager *app.Pager) ([]*model.ArticleRow, error) {
	article := model.Article{State: state}

	return article.ListByTagID(d.engine, tagID, pager.PageSize, app.GetPageOffset(pager.Page, pager.PageSize))
}
