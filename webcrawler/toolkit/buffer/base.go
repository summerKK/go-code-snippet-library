package buffer

type IPool interface {
	// 单个缓冲器的容量
	BufCap() uint32
	// 获取缓冲池中所有缓冲器的最大数量
	MaxBufNum() uint32
	// 获取缓冲池缓冲器的数量
	BufNum() uint32
	Total() uint64
	Put(datum interface{}) error
	Get() (datum interface{}, err error)
	Close() bool
	// 是否关闭
	Closed() bool
}

type IBuf interface {
	Cap() uint32
	Len() uint32
	Put(datum interface{}) (bool, error)
	// get方法需要判断返回的interface{}是否是nil
	Get() (interface{}, error)
	Close() bool
	Closed() bool
}
