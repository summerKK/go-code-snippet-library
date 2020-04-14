package module

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
)

func Legalletter(mtype base.MType) bool {
	if _, ok := base.LegalletterMap[mtype]; ok {
		return true
	}
	return false
}

func CheckType(mType base.MType, module base.IModule) bool {
	if mType == "" || module == nil {
		return false
	}
	switch mType {
	case base.TYPE_PIPELINE:
		if _, ok := module.(base.IPipeline); ok {
			return true
		}
	case base.TYPE_ANALYZER:
		if _, ok := module.(base.IAnalyzer); ok {
			return true
		}
	case base.TYPE_DOWNLOADER:
		if _, ok := module.(base.IDownloader); ok {
			return true
		}
	}

	return false
}
