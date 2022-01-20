package models

// Pagination for page type data
type Pagination struct {
	Page       int         `json:"page,omitempty"`
	PageSize   int         `json:"pageSize,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	TotalRows  int64       `json:"totalRows,omitempty"`
	TotalPages int         `json:"totalPages,omitempty"`
	Rows       interface{} `json:"rows,omitempty"`
}

// GetPage get page
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

// GetPageSize get page size
func (p *Pagination) GetPageSize() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// GetOffset get offset(page * page size)
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// GetSort get sort column name
func (p *Pagination) GetSort() string {
	return p.Sort
}
