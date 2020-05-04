package internal

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"net/http"
	"net/url"
	"path"
	"strings"
)

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
		doc.Find("a").Each(func(index int, selection *goquery.Selection) {
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

		doc.Find("img").Each(func(index int, selection *goquery.Selection) {
			// 前期过滤。
			imgSrc, exists := selection.Attr("src")
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
