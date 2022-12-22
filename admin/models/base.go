package models

// Login User login structure
type Login struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type Pagination struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

func (m *Pagination) GetPageIndex() int {
	if m.PageNum <= 0 {
		m.PageNum = 1
	}
	return m.PageNum
}

func (m *Pagination) GetPageInfo() (limit, offset int) {
	pageIndex := m.GetPageIndex()
	limit = m.GetPageSize()
	offset = (pageIndex - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

type GetReq struct {
	Id uint64 `json:"id" uri:"id"`
}

type DeleteReq struct {
	Ids []uint64 `json:"ids"`
}

func (s *DeleteReq) GetId() interface{} {
	return s.Ids
}
