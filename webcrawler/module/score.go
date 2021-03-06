package module

import "github.com/summerKK/go-code-snippet-library/webcrawler/module/base"

// CalculateScoreSimple 代表简易的组件评分计算函数。
func CalculateScoreSimple(counts base.Counts) uint64 {
	return counts.CalledCount + counts.AcceptedCount<<1 + counts.CompletedCount<<2 + counts.HandlingNumber<<4
}

// SetScore 用于设置给定组件的评分。
// 结果值代表是否更新了评分。
func SetScore(module base.IModule) bool {
	calculator := module.ScoreCalculator()
	if calculator == nil {
		calculator = CalculateScoreSimple
	}
	newScore := calculator(module.Counts())
	if newScore == module.Score() {
		return false
	}
	module.SetScore(newScore)
	return true
}
