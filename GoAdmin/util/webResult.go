package util

type WebResultOption struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type WebResultCommonList struct {
	Count     int64       `json:"count"`     // 总计记录条数
	List      interface{} `json:"list"`      // 数据列表
	PageIndex int         `json:"pageIndex"` // 原封不动还回去，没鸟用
	PageSize  int         `json:"pageSize"`  // 原封不动还回去，没鸟用
}
