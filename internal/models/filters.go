package models

import "github.com/namikaze-dev/bluelight/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 1000, "page_size", "must be a maximum of 1000")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
