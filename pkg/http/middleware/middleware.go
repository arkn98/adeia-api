/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package middleware

import "net/http"

// Func represents a middleware function.
type Func func(http.Handler) http.Handler

// FuncChain is a slice of Funcs, representing the middleware chain.
type FuncChain struct {
	funcs []Func
}

// NewChain creates a new middleware chain.
func NewChain(funcs ...Func) FuncChain {
	return FuncChain{append(([]Func)(nil), funcs...)}
}

// Compose applies/composes all the middleware funcs in-order on the provided handler.
func (c *FuncChain) Compose(f http.Handler) http.Handler {
	for _, m := range c.funcs {
		f = m(f)
	}
	return f
}

// Append appends the passed-in Funcs to the FuncChain.
func (c *FuncChain) Append(funcs ...Func) FuncChain {
	newChain := make([]Func, 0, len(c.funcs)+len(funcs))
	newChain = append(newChain, c.funcs...)
	newChain = append(newChain, funcs...)

	return FuncChain{newChain}
}

// Chain returns the slice of funcs stored in the FuncChain.
func (c *FuncChain) Chain() []Func {
	return c.funcs
}
