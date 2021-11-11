package check

import (
	"sort"
)

type Result struct {
	Name        string
	Type        Type
	Data        Data
	IsMalicious bool
	Error       *ResultError
}

type Results []Result

func (r Results) SortByName() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].Name < r[j].Name
	})
}
