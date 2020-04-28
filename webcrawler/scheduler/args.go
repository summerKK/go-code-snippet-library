package scheduler

import "github.com/summerKK/go-code-snippet-library/webcrawler/module/base"

// 请求参数
type RequestArgs struct {
	// 代表可以接受的URL的主域名的列表,不在该列表的URL都会被忽略(限定爬虫爬取的范围)
	AcceptedDomains []string `json:"accepted_primary_domains"`
	// 爬取的最大深度
	MaxDepth uint32 `json:"max_depth"`
}

func (r *RequestArgs) Check() error {
	if r.AcceptedDomains == nil {
		return genError("nil accepted primary domain list")
	}

	return nil
}

// 定义数据缓冲池容量
type DataArgs struct {
	ReqBufCap     uint32 `json:"req_buf_cap"`
	ReqMaxBufNum  uint32 `json:"req_max_buf_num"`
	RespBufCap    uint32 `json:"resp_buf_cap"`
	RespMaxBufNum uint32 `json:"resp_max_buf_num"`
	ItemBufCap    uint32 `json:"item_buf_cap"`
	ItemMaxBufNum uint32 `json:"item_max_buf_num"`
	ErrBufCap     uint32 `json:"err_buf_cap"`
	ErrMaxBufNum  uint32 `json:"err_max_buf_num"`
}

func (d *DataArgs) Check() error {
	if d.ReqBufCap == 0 {
		return genError("zero request buffer capacity")
	}
	if d.ReqMaxBufNum == 0 {
		return genError("zero max request buffer number")
	}
	if d.RespBufCap == 0 {
		return genError("zero response buffer capacity")
	}
	if d.RespMaxBufNum == 0 {
		return genError("zero max response buffer number")
	}
	if d.ItemBufCap == 0 {
		return genError("zero item buffer capacity")
	}
	if d.ItemMaxBufNum == 0 {
		return genError("zero max item buffer number")
	}
	if d.ErrBufCap == 0 {
		return genError("zero error buffer capacity")
	}
	if d.ErrMaxBufNum == 0 {
		return genError("zero max error buffer number")
	}

	return nil
}

// ModuleArgsSummary 代表组件相关的参数容器的摘要类型。
type ModuleArgsSummary struct {
	DownloaderListSize int `json:"downloader_list_size"`
	AnalyzerListSize   int `json:"analyzer_list_size"`
	PipelineListSize   int `json:"pipeline_list_size"`
}

// module 模块
type ModuleArgs struct {
	Downloaders []base.IDownloader
	Analyzers   []base.IAnalyzer
	Pipelines   []base.IPipeline
}

// Check 用于当前参数容器的有效性。
func (m *ModuleArgs) Check() error {
	if len(m.Downloaders) == 0 {
		return genError("empty downloader list")
	}
	if len(m.Analyzers) == 0 {
		return genError("empty analyzer list")
	}
	if len(m.Pipelines) == 0 {
		return genError("empty pipeline list")
	}

	return nil
}

func (m *ModuleArgs) Summary() ModuleArgsSummary {
	return ModuleArgsSummary{
		DownloaderListSize: len(m.Downloaders),
		AnalyzerListSize:   len(m.Analyzers),
		PipelineListSize:   len(m.Pipelines),
	}
}
