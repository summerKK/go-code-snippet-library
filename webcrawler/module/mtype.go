package module

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/analyzer"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/downloader"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/pipeline"
)

func Legalletter(mtype MType) bool {
	if _, ok := legalletterMap[mtype]; ok {
		return true
	}
	return false
}

func CheckType(mType MType, module IModule) bool {
	if mType == "" || module == nil {
		return false
	}
	switch mType {
	case TYPE_PIPELINE:
		if _, ok := module.(pipeline.IPipeline); ok {
			return true
		}
	case TYPE_ANALYZER:
		if _, ok := module.(analyzer.IAnalyzer); ok {
			return true
		}
	case TYPE_DOWNLOADER:
		if _, ok := module.(downloader.IDownloader); ok {
			return true
		}
	}

	return false
}
