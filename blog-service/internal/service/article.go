package service

import (
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/dao"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/model"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	Title string `form:"title" binding:"min=2,max=100"`
	TagID uint32 `form:"tag_id" binding:"gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"required,min=2,max=100"`
	Desc          string `form:"desc" binding:"required,min=2,max=255"`
	Content       string `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
	CreatedBy     string `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32 `form:"id" binding:"required,gte=1"`
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"min=2,max=100"`
	Desc          string `form:"desc" binding:"min=2,max=255"`
	Content       string `form:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"url"`
	ModifiedBy    string `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*model.Article, error) {
	return s.dao.GetArticleList(param.Title, param.State, pager.Page, pager.PageSize)
}

func (s *Service) CreateArticle(param *CreateArticleRequest) error {
	article := &dao.Article{
		TagID:          param.TagID,
		Title:          param.Title,
		Desc:           param.Desc,
		Content:        param.Content,
		ConverImageUrl: param.CoverImageUrl,
		CreatedBy:      param.CreatedBy,
		State:          param.State,
	}

	return s.dao.CreateArticle(article)
}

func (s *Service) UpdateArticle(param *UpdateArticleRequest) error {
	article := &dao.Article{
		ID:             param.ID,
		TagID:          param.TagID,
		Title:          param.Title,
		Desc:           param.Desc,
		Content:        param.Content,
		ConverImageUrl: param.CoverImageUrl,
		ModifiedBy:     param.ModifiedBy,
		State:          param.State,
	}

	return s.dao.UpdateArticle(article)
}

func (s *Service) DeleteArticle(param *DeleteArticleRequest) error {
	return s.dao.DeleteArticle(param.ID)
}
