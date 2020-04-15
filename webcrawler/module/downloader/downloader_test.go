package downloader

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDownloader_Download(t *testing.T) {
	genMid, err := module.GenMid(base.TYPE_DOWNLOADER, module.DefaultSNGen.Next(), nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	downloader, err := New(genMid, nil, client)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest("GET", "http://www.baidu.com/robots.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := module.NewRequest(request, 0)
	response, err := downloader.Download(req)
	if err != nil {
		t.Fatal(err)
	}
	body := response.Resp().Body
	if body == nil {
		t.Fatal(err)
	}
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}
