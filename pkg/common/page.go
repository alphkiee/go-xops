package common

// PageInfo ...分页封装
type PageInfo struct {
	Current  uint  `json:"current" form:"current"`   // 当前页码
	PageSize uint  `json:"pageSize" form:"pageSize"` // 每页显示条数
	Total    int64 `json:"total"`                    // 数据总条数
	All      bool  `json:"all" form:"all"`           // 不使用分页
}

// PageData ....带分页数据封装
type PageData struct {
	PageInfo
	DataList interface{} `json:"data"` // 数据列表
}

// GetLimit ...计算limit/offset, 如果需要用到返回的PageSize, PageNum, 务必保证Total值有效
func (s *PageInfo) GetLimit() (int, int) {
	var pageSize int64
	var current int64
	total := s.Total
	if s.PageSize < 1 {
		pageSize = 10
	} else {
		pageSize = int64(s.PageSize)
	}
	if s.Current < 1 {
		current = 1
	} else {
		current = int64(s.Current)
	}

	if total > 0 && current > total {
		current = total
	}

	maxPageNum := total/pageSize + 1
	if total%pageSize == 0 {
		maxPageNum = total / pageSize
	}
	if maxPageNum < 1 {
		maxPageNum = 1
	}

	if current > maxPageNum {
		current = maxPageNum
	}

	limit := pageSize
	offset := limit * (current - 1)

	return int(limit), int(offset)
}
