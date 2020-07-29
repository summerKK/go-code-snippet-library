package errcode

var (
	ErrorGetTagListFail     = NewError(2001001, "获取标签列表失败")
	ErrorGetArticleListFail = NewError(2001002, "获取文章列表失败")
)
