package internal

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var globalPage = 1

func genResponseParsers() []base.ParseResponse {
	parseLink := func(httpResp *http.Response, respDepth uint32) ([]base.IData, []error) {
		dataList := make([]base.IData, 0)
		if httpResp == nil {
			return nil, []error{fmt.Errorf("nil http response")}
		}
		httpReq := httpResp.Request
		if httpReq == nil {
			return nil, []error{fmt.Errorf("nul http request")}
		}
		reqUrl := httpReq.URL
		if httpResp.StatusCode != 200 {
			err := fmt.Errorf("unsupported status code %d (requestUrl:%s)", httpResp.StatusCode, reqUrl)
			return nil, []error{err}
		}
		body := httpResp.Body
		if body == nil {
			err := fmt.Errorf("nil http response body (requestUrl:%s)", reqUrl)
			return nil, []error{err}
		}
		var matchedContentType bool
		if httpResp.Header != nil {
			contentTypes := httpResp.Header["Content-Type"]
			for _, contentType := range contentTypes {
				if strings.HasPrefix(contentType, "text/html") {
					matchedContentType = true
					break
				}
			}
		}
		if !matchedContentType {
			return dataList, nil
		}
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return dataList, []error{err}
		}
		errs := make([]error, 0)
		doc.Find("#main div.thumb.pos-r a").Each(func(index int, selection *goquery.Selection) {
			href, exists := selection.Attr("href")
			if !exists || href == "" || href == "#" || href == "/" {
				return
			}
			href = strings.TrimSpace(href)
			lowerHref := strings.ToLower(href)
			if href == "" || strings.HasPrefix(lowerHref, "javascript") {
				return
			}
			aUrl, err := url.Parse(href)
			if err != nil {
				logger.Logger.Warnf("an error occurs when parsing attribute %q in tag %q :%s (href:%s)", err, "href", "a", href)
				return
			}
			if !aUrl.IsAbs() {
				aUrl = reqUrl.ResolveReference(aUrl)
			}
			httpReq, err := http.NewRequest("GET", aUrl.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				req := base.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		})

		// 获取下一页信息
		nextPage(&dataList, &errs, reqUrl, respDepth)

		// 查找img标签并提取地址。
		doc.Find("#content-innerText img").Each(func(index int, sel *goquery.Selection) {
			// 前期过滤。
			imgSrc, exists := sel.Attr("src")
			if !exists || imgSrc == "" || imgSrc == "#" || imgSrc == "/" {
				return
			}
			imgSrc = strings.TrimSpace(imgSrc)
			imgURL, err := url.Parse(imgSrc)
			if err != nil {
				errs = append(errs, err)
				return
			}
			if !imgURL.IsAbs() {
				imgURL = reqUrl.ResolveReference(imgURL)
			}
			httpReq, err := http.NewRequest("GET", imgURL.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
				httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
				httpReq.Header.Set("Referer", "http://www.jdlingyu.mobi/tuji/")
				httpReq.Header.Set("Accept-Encoding", "gzip, deflate")
				httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
				req := base.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		})

		return dataList, errs
	}

	parseImg := func(httpResp *http.Response, respDepth uint32) ([]base.IData, []error) {
		// 检查响应。
		if httpResp == nil {
			return nil, []error{fmt.Errorf("nil HTTP response")}
		}
		httpReq := httpResp.Request
		if httpReq == nil {
			return nil, []error{fmt.Errorf("nil HTTP request")}
		}
		reqURL := httpReq.URL
		if httpResp.StatusCode != 200 {
			err := fmt.Errorf("unsupported status code %d (requestURL: %s)",
				httpResp.StatusCode, reqURL)
			return nil, []error{err}
		}
		httpRespBody := httpResp.Body
		if httpRespBody == nil {
			err := fmt.Errorf("nil HTTP response body (requestURL: %s)",
				reqURL)
			return nil, []error{err}
		}
		// 检查HTTP响应头中的内容类型。
		dataList := make([]base.IData, 0)
		var pictureFormat string
		if httpResp.Header != nil {
			contentTypes := httpResp.Header["Content-Type"]
			var contentType string
			for _, ct := range contentTypes {
				if strings.HasPrefix(ct, "image") {
					contentType = ct
					break
				}
			}
			index1 := strings.Index(contentType, "/")
			index2 := strings.Index(contentType, ";")
			if index1 > 0 {
				if index2 < 0 {
					pictureFormat = contentType[index1+1:]
				} else if index1 < index2 {
					pictureFormat = contentType[index1+1 : index2]
				}
			}
		}
		if pictureFormat == "" {
			return dataList, nil
		}
		// 生成条目。
		item := make(map[string]interface{})
		item["reader"] = httpRespBody
		item["name"] = path.Base(reqURL.Path)
		item["ext"] = pictureFormat
		dataList = append(dataList, base.Item(item))
		return dataList, nil
	}

	return []base.ParseResponse{parseLink, parseImg}
}

func nextPage(dataList *[]base.IData, errs *[]error, reqUrl *url.URL, respDepth uint32) {

	form := url.Values{}
	form.Add("type", "collection1495")
	form.Add("paged", strconv.Itoa(globalPage))
	globalPage += 1
	req, err := http.NewRequest("POST", "https://www.jdlingyu.mobi/wp-admin/admin-ajax.php?action=zrz_load_more_posts", strings.NewReader(form.Encode()))
	if err != nil {
		*errs = append(*errs, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	req.Header.Set("cookie", "wpjam-compare=0; UM_distinctid=171f36377b77b0-088cd4c986d089-30677c00-1aeaa0-171f36377b8a99; Hm_lvt_4c0b4bd72dc090c1c4a836b68d5c4d4b=1588926315,1588926389; CNZZDATA1274771516=1489723408-1588921214-https%253A%252F%252Fwww.google.com%252F%7C1588945388; Hm_lpvt_4c0b4bd72dc090c1c4a836b68d5c4d4b=1588948988")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		*errs = append(*errs, err)
		return
	}

	type apiResponse struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}
	var resp apiResponse
	b, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(b, &resp)
	if err != nil {
		*errs = append(*errs, err)
		return
	}

	if resp.Status != 200 {
		return
	}

	fromReader, err := goquery.NewDocumentFromReader(strings.NewReader(resp.Msg))
	if err != nil {
		*errs = append(*errs, err)
		return
	}

	fromReader.Find("div.thumb.pos-r a").Each(func(index int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if !exists || href == "" || href == "#" || href == "/" {
			return
		}
		href = strings.TrimSpace(href)
		lowerHref := strings.ToLower(href)
		if href == "" || strings.HasPrefix(lowerHref, "javascript") {
			return
		}
		aUrl, err := url.Parse(href)
		if err != nil {
			logger.Logger.Warnf("an error occurs when parsing attribute %q in tag %q :%s (href:%s)", err, "href", "a", href)
			return
		}
		if !aUrl.IsAbs() {
			aUrl = reqUrl.ResolveReference(aUrl)
		}
		httpReq, err := http.NewRequest("GET", aUrl.String(), nil)
		if err != nil {
			*errs = append(*errs, err)
		} else {
			req := base.NewRequest(httpReq, respDepth)
			*dataList = append(*dataList, req)
		}
	})
}
