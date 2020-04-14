package downloader

import (
	"bufio"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
	"net/http"
	"testing"
)

func TestDownloader_Download(t *testing.T) {
	genMid, err := module.GenMid(module.TYPE_DOWNLOADER, module.DefaultSNGen.Next(), nil)
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
	req := data.NewRequest(request, 0)
	response, err := downloader.Download(req)
	if err != nil {
		t.Fatal(err)
	}
	body := response.Resp().Body
	if body == nil {
		t.Fatal(err)
	}
	reader := bufio.NewReader(body)
	line, _, err := reader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(line))

}
