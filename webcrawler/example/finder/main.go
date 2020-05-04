package main

import (
	"flag"
	"fmt"
	lib "github.com/summerKK/go-code-snippet-library/webcrawler/example/finder/internal"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/scheduler"
	"net/http"
	"os"
	"strings"
)

var (
	firstUrl string
	domains  string
	depth    uint
	dirPath  string
)

func init() {
	flag.StringVar(&firstUrl, "first", "http://zhihu.sogou.com/zhihu?query=golang+logo", "the first url which you want to access")
	flag.StringVar(&domains, "domains", "zhihu.com", "the primary domains which you accepted. please using comma-separated multiple domains")
	flag.UintVar(&depth, "depth", 3, "the depth for crawling")
	flag.StringVar(&dirPath, "dir", "./pictures", "the path which you want to save the image files.")
}

func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\tfinder [flags] \n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	// 创建调度器
	sched := scheduler.NewScheduler()
	// 准备调度器的初始化参数
	domainParts := strings.Split(domains, ",")
	var acceptedDomains []string
	for _, domain := range domainParts {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			acceptedDomains = append(acceptedDomains, domain)
		}
	}
	requestArgs := scheduler.RequestArgs{
		AcceptedDomains: acceptedDomains,
		MaxDepth:        uint32(depth),
	}

	dataArgs := scheduler.DataArgs{
		ReqBufCap:     50,
		ReqMaxBufNum:  1000,
		RespBufCap:    50,
		RespMaxBufNum: 10,
		ItemBufCap:    50,
		ItemMaxBufNum: 100,
		ErrBufCap:     50,
		ErrMaxBufNum:  1,
	}
	downloaders, err := lib.GetDownloaders(1)
	if err != nil {
		logger.Logger.Fatalf("an error occurs when creating donwloaders:%s", err)
	}
	analyzers, err := lib.GetAnalyzers(1)
	if err != nil {
		logger.Logger.Fatalf("an error occurs when creating analyzers:%s", err)
	}
	pipelines, err := lib.GetPipelines(1, dirPath)
	if err != nil {
		logger.Logger.Fatalf("an error occurs when creating pipelines:%s", err)
	}
	modulesArgs := scheduler.ModuleArgs{
		Downloaders: downloaders,
		Analyzers:   analyzers,
		Pipelines:   pipelines,
	}
	// 初始化调度器
	err = sched.Init(requestArgs, dataArgs, modulesArgs)
	if err != nil {
		logger.Logger.Fatalf("an error occurs when initializing scheduler:%s", err)
	}
	firstHttpReq, err := http.NewRequest("GET", firstUrl, nil)
	if err != nil {
		logger.Logger.Fatalln(err)
	}
	err = sched.Start(firstHttpReq)
	if err != nil {
		logger.Logger.Fatalf("an error occurs when starting scheduler:%s", err)
	}

}
