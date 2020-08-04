package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
)

const (
	APP_KEY    = "summer"
	APP_SECRET = "summer"
)

type API struct {
	Url string
}

type AccessToken struct {
	Token string `json:"token"`
}

func NewApi(url string) *API {
	return &API{Url: url}
}

func (a *API) GetAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpGet(ctx, fmt.Sprintf("api/auth?app_key=%s&app_secret=%s", APP_KEY, APP_SECRET))
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return "", err
	}

	return accessToken.Token, nil
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	return a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
}

func (a *API) GetArticleList(ctx context.Context, tagId uint32) ([]byte, error) {
	token, err := a.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	return a.httpGet(ctx, fmt.Sprintf("%s?token=%s&tag_id=%d", "api/v1/articles", token, tagId))
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	urlstr := fmt.Sprintf("%s/%s", a.Url, path)
	req, err := http.NewRequest(http.MethodGet, urlstr, nil)
	if err != nil {
		return nil, err
	}

	span, newCtx := opentracing.StartSpanFromContext(ctx, "HTTP GET: "+a.Url,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	)
	span.SetTag("url", urlstr)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	req = req.WithContext(newCtx)
	// 设置抓包代理
	//fixedUrl, _ := url.Parse("http://127.0.0.1:8888")

	client := http.Client{
		Timeout: time.Second * 60,
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(fixedUrl),
		//},
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	defer span.Finish()

	if response.StatusCode != http.StatusOK {
		log.Warnf("get %s 失败,状态码:%d", path, response.StatusCode)
		return nil, errors.New("httpGet请求失败")
	}

	return ioutil.ReadAll(response.Body)
}
