package module

type SNGenerator interface {
	// 最小序列号
	Min() uint64
	// 最大序列号
	Max() uint64
	Next() uint64
	// 获取循环计数
	CycleCount() uint64
	// 用于获取序列号并且准备下个序列号
	Get() uint64
}
