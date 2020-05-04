package internal

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/analyzer"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/downloader"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/pipeline"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/sn"
)

var snGen = sn.NewGenerator(1, 0)

func GetDownloaders(number uint8) ([]base.IDownloader, error) {
	var downloaders []base.IDownloader
	if number == 0 {
		return downloaders, nil
	}
	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMid(base.TYPE_DOWNLOADER, snGen.Get(), nil)
		if err != nil {
			return downloaders, err
		}
		d, err := downloader.New(mid, module.CalculateScoreSimple, genHTTPClient())
		if err != nil {
			return downloaders, err
		}
		downloaders = append(downloaders, d)
	}

	return downloaders, nil
}

func GetAnalyzers(number uint8) ([]base.IAnalyzer, error) {
	var analyzers []base.IAnalyzer
	if number == 0 {
		return analyzers, nil
	}
	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMid(base.TYPE_ANALYZER, snGen.Get(), nil)
		if err != nil {
			return analyzers, err
		}
		a, err := analyzer.New(mid, module.CalculateScoreSimple, genResponseParsers())
		if err != nil {
			return analyzers, err
		}
		analyzers = append(analyzers, a)
	}

	return analyzers, nil
}

func GetPipelines(number uint8, dirPath string) ([]base.IPipeline, error) {
	var pipelines []base.IPipeline
	if number == 0 {
		return pipelines, nil
	}
	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMid(base.TYPE_PIPELINE, snGen.Get(), nil)
		if err != nil {
			return pipelines, err
		}
		p, err := pipeline.New(mid, module.CalculateScoreSimple, genItemProcessors(dirPath))
		if err != nil {
			return pipelines, err
		}
		pipelines = append(pipelines, p)
	}

	return pipelines, nil
}
