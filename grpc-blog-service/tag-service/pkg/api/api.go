package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("%s/%s", a.Url, path))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
