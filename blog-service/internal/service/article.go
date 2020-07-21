package service

import (
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/dao"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/model"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type Article struct {
	ID             uint32     `json:"id"`
	Title          string     `json:"title"`
	Desc           string     `json:"desc"`
	Content        string     `json:"content"`
	ConverImageUrl string     `json:"conver_image_url"`
	State          uint8      `json:"state"`
	Tag            *model.Tag `json:"tag"`
}

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	Title string `form:"title"`
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

func (s *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := s.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}
	articleTag, err := s.dao.GetArticleTagByAID(article.ID)
	if err != nil {
		return nil, err
	}
	tag, err := s.dao.GetTag(articleTag.ID, model.STATE_OPEN)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:             article.ID,
		Title:          article.Title,
		Desc:           article.Desc,
		Content:        article.Content,
		ConverImageUrl: article.ConverImageUrl,
		State:          article.State,
		Tag:            tag,
	}, nil
}

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	count, err := s.dao.CountArticleListByTagID(param.TagID, param.State)
	if err != nil {
		return nil, 0, err
	}

	articleList, err := s.dao.GetArticleListByTagID(param.TagID, param.State, pager.PageSize, pager.Page)
	if err != nil {
		return nil, 0, err
	}

	var svcArticle = make([]*Article, 0, len(articleList))
	for _, row := range articleList {
		r := &Article{
			ID:             row.ArticleID,
			Title:          row.ArticleTitle,
			Desc:           row.ArticleDesc,
			Content:        row.Content,
			ConverImageUrl: row.ConverImageUrl,
			Tag: &model.Tag{
				Model: &model.Model{
					ID: row.TagID,
				},
				Name: row.TagName,
			},
		}
		svcArticle = append(svcArticle, r)
	}

	return svcArticle, count, nil
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

	createdArticle, err := s.dao.CreateArticle(article)
	if err != nil {
		return err
	}

	err = s.dao.CreateArticleTag(createdArticle.ID, param.TagID, param.CreatedBy)

	return nil
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

	err := s.dao.UpdateArticle(article)
	if err != nil {
		return err
	}

	err = s.dao.UpdateArticleTag(article.ID, param.TagID, param.ModifiedBy)

	return err
}

func (s *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := s.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}

	err = s.dao.DeleteArticleTag(param.ID)

	return err
}
