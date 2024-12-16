package data

import (
	"strings"

	"github.com/wangyaodream/greenlight/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
    SortSafelist []string
}

type Metadata struct {
    CurrentPage int `json:"current_page,omitempty"`
    PageSize int `json:"page_size,omitempty"`
    FirstPage int `json:"first_page,omitempty"`
    LastPage int `json:"last_page,omitempty"`
    TotalRecords int `json:"total_records,omitempty"`
}

func ValidateFilters(v *validator.Validator, f Filters) {
    // 检查page和PageSize字段是否有效
    v.Check(f.Page > 0, "page", "must be greater than zero")
    v.Check(f.Page <= 10_000, "page", "must be a maximum of 10,000")
    v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
    v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

    v.Check(validator.PermiteedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
    for _, safeValue := range f.SortSafelist {
        if f.Sort == safeValue {
            // 去掉排序字段前面的负号
            return strings.TrimPrefix(f.Sort, "-")
        }
    }

    panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
    if strings.HasPrefix(f.Sort, "-") {
        return "DESC"
    }

    return "ASC"
}


func (f Filters) limit() int {
    return f.PageSize
}

func (f Filters) Offset() int {
    return (f.Page - 1) * f.PageSize
}
