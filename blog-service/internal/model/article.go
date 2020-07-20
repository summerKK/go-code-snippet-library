package model

import (
	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

type Article struct {
	*Model
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Content        string `json:"content"`
	ConverImageUrl string `json:"conver_image_url"`
	State          uint8  `json:"state"`
}

type ArticleRow struct {
	ArticleID      uint32 `json:"article_id"`
	TagID          uint32 `json:"tag_id"`
	TagName        string `json:"tag_name"`
	ArticleTitle   string `json:"article_title"`
	ArticleDesc    string `json:"article_desc"`
	Content        string `json:"content"`
	ConverImageUrl string `json:"conver_image_url"`
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Get(db *gorm.DB) (*Article, error) {
	var article Article
	err := db.Where("id = ? and state = ? and id_del = ?", a.ID, a.State, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &article, err
	}

	return &article, nil
}

func (a Article) Count(db *gorm.DB) (int, error) {
	var count int
	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("state = ?", a.State)
	if err := db.Model(&a).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (a Article) List(db *gorm.DB, pageSize, pageOffset int) ([]*Article, error) {
	var articles []*Article
	if pageSize >= 0 && pageOffset >= 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("state = ?", a.State)

	if err := db.Where("is_del = ?", 0).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageSize, pageOffset int) ([]*ArticleRow, error) {
	fields := []string{"a.title as article_title", "a.desc as article_desc", "a.id as article_id", "a.conver_image_url", "a.content"}
	fields = append(fields, []string{"b.name as tag_name", "b.id as tag_id"}...)
	if pageSize >= 0 && pageOffset >= 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := db.Select(fields).Table(ArticleTag{}.TableName()+" as c").
		Joins("left join "+Article{}.TableName()+" as a on a.id = c.article_id").
		Joins("left join "+Tag{}.TableName()+" as b on b.id = c.tag_id").
		Where("a.state = ? and a.is_del = ? and c.tag_id = ?", a.State, 0, tagID).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		err := rows.Scan(r, r.ArticleTitle, r.ArticleDesc, r.ArticleID, r.ConverImageUrl, r.Content, r.TagName, r.TagID)
		if err != nil {
			return nil, err
		}
		articles = append(articles, r)
	}

	return articles, nil
}

func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" as c ").
		Joins("left join "+Article{}.TableName()+" as a on c.article_id = a.id").
		Joins("left join"+Tag{}.TableName()+" as b on c.tag_id = b.id").
		Where("a.state = ? and a.is_del = ? and c.tag_id = ?", a.State, 0, tagID).
		Count(&count).Error

	return count, err
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	err := db.Create(&a).Error

	return &a, err
}

func (a Article) Update(db *gorm.DB, vales interface{}) error {
	return db.Model(&a).Where("id = ? and is_del = ?", a.ID, 0).Updates(vales).Error
}

func (a Article) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", a.ID, 0).Delete(&a).Error
}
