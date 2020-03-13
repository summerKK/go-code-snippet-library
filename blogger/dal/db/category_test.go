package db

import "testing"

func init() {
	dns := "root:root@tcp(127.0.0.1)/blogger?parseTime=true"
	err := Init(dns)
	if err != nil {
		panic(err)
	}
}

func TestCategoryList(t *testing.T) {
	list, err := CategoryList([]int64{1, 2, 3})
	if err != nil {
		t.Errorf("got list failed,error:%v", err)
	}
	t.Logf("got items,len:%d", len(list))
}
