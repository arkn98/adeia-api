/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package middleware

import (
	"context"
	"net/http"

	"adeia"
	"adeia/pkg/constants"
	"adeia/pkg/log"
	"adeia/pkg/util/httputil"
)

// QueryParamParser is a middleware that parses the request URL for query-params
// depending on how it is configured. The parsed QueryParams is stored in the
// context and passed along for consumption by the controllers.
func QueryParamParser(
	log log.Logger,
	searchable bool,
	sortable bool,
	filterable bool,
	paginatable bool,
) Func {
	parser := httputil.ParseQueryParams(searchable, sortable, filterable, paginatable)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			qp, err := parser(r)
			if err != nil {
				log.Debugf("error parsing query params: %v", err)
				httputil.LogWriteErr(log, httputil.RespondWithErr(w, adeia.ErrInvalidRequest))
				return
			}

			if qp != nil {
				log.Debugf("saving parsed query params to context: %v", qp)
				ctx := context.WithValue(r.Context(), constants.CtxQueryParamsKey, qp)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			log.Debug("no query params (of interest) to parse and save")
			next.ServeHTTP(w, r)
		})
	}
}
