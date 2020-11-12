/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// QueryParams represents all of the query-params parsed from the request URL
// that can be consumed by the controllers. The parsing happens (usually) in a
// middleware, before the request reaches the controllers, and the QueryParams
// is stored in the context and passed along.
type QueryParams struct {
	// searching
	Query       string
	IsSearching bool

	// sorting
	Sort      []*SortItem
	IsSorting bool

	// filtering
	Filters     map[string][]string
	IsFiltering bool

	// pagination
	Page         int
	PerPage      int
	IsPaginating bool
}

// NewQueryParams creates a new *QueryParams.
func NewQueryParams() *QueryParams {
	return &QueryParams{
		PerPage: 20,
		Page:    1,
	}
}

// SortItem represents a single sorting field specified in the query-param.
type SortItem struct {
	Field      string
	Descending bool
}

func parseSearchParams(values map[string][]string, qp *QueryParams) {
	if q, ok := values["q"]; ok {
		qp.IsSearching = true
		qp.Query = q[0] // we ignore other values
		delete(values, "q")
	}
}

func parseSortParams(values map[string][]string, qp *QueryParams) {
	if s, ok := values["sort"]; ok {
		qp.IsSorting = true
		fields := strings.Split(s[0], ",")
		for _, field := range fields {
			sortItem := &SortItem{}
			if strings.HasPrefix(field, "-") {
				sortItem.Descending = true
				sortItem.Field = strings.TrimPrefix(field, "-")
			} else {
				sortItem.Field = field
			}
			qp.Sort = append(qp.Sort, sortItem)
		}
		delete(values, "sort")
	}
}

func parsePaginateParams(values map[string][]string, qp *QueryParams) error {
	page, pageOk := values["page"]
	perPage, perPageOk := values["per_page"]
	if pageOk {
		p, err := strconv.Atoi(page[0])
		if err != nil {
			return fmt.Errorf("invalid query parameter: %v", err)
		}
		qp.Page = p
		qp.IsPaginating = true
		delete(values, "page")
	}
	if perPageOk {
		p, err := strconv.Atoi(perPage[0])
		if err != nil {
			return fmt.Errorf("invalid query parameter: %v", err)
		}
		qp.PerPage = p
		qp.IsPaginating = true
		delete(values, "per_page")
	}
	return nil
}

func parseFilterParams(values map[string][]string, qp *QueryParams) {
	if len(values) > 0 {
		qp.IsFiltering = true
		for k, v := range values {
			temp := make([]string, len(v))
			copy(temp, v)
			qp.Filters[k] = temp
		}
	}
}

// ParseQueryParams is a closure-func that can parse the request URL for various
// query-params depending on what all among {search, sort, filter, paginate} are
// enabled.
func ParseQueryParams(search, sort, filter, paginate bool) func(*http.Request) (*QueryParams, error) {
	return func(r *http.Request) (*QueryParams, error) {
		if !(search || sort || filter || paginate) {
			return nil, nil
		}

		qp := NewQueryParams()

		// make a copy of the query-params map
		values := make(map[string][]string)
		for k, v := range r.URL.Query() {
			temp := make([]string, len(v))
			copy(temp, v)
			values[k] = temp
		}

		if search {
			parseSearchParams(values, qp)
		}

		if sort {
			parseSortParams(values, qp)
		}

		if paginate {
			err := parsePaginateParams(values, qp)
			if err != nil {
				return nil, err
			}
		}

		if filter {
			parseFilterParams(values, qp)
		}

		return qp, nil
	}
}
