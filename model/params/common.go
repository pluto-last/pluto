package params

// Paging common input parameter structure
type PageInfo struct {
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

// Find by id structure
type GetById struct {
	ID float64 `json:"id" form:"id"` // 主键ID
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

type PageResult struct {
	List   interface{} `json:"list"`
	Total  int64       `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}
