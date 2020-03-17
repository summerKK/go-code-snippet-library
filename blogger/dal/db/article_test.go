package db

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/blogger/model"
	"testing"
)

func init() {
	dns := "root:root@tcp(127.0.0.1)/blogger?parseTime=true"
	err := Init(dns)
	if err != nil {
		panic(err)
	}
}

func TestArticleInsert(t *testing.T) {

	articleDetail := model.ArticleDetail{
		ArticleInfo: model.ArticleInfo{
			CategoryId:   1,
			Title:        "summer",
			Summary:      "summer",
			ViewCount:    10,
			CommentCount: 10,
			Username:     "summer",
		},
		Content: "hello,summer",
	}
	articleId, err := ArticleInsert(&articleDetail)
	if err != nil {
		t.Errorf("insert into article faild,error:%v", err)
	}

	fmt.Printf("got article id:%d", articleId)
}

func TestArticleList(t *testing.T) {
	list, err := ArticleList(0, 10)
	if err != nil {
		t.Errorf("got list faild,error:%v", err)
	}
	t.Logf("got item,len:%d", len(list))
}

func TestArticleInfo(t *testing.T) {
	articleInfo, err := ArticleInfo(3)
	if err != nil {
		t.Errorf("got article info failed,error:%v", err)
	}
	t.Logf("got item %+v", articleInfo)
}
