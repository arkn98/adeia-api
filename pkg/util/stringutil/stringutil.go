/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package stringutil

import "strings"

// Escaper replaces all occurrences of a slice of strings, with the corresponding
// escaped versions.
//
// For example, the inputs
//     string: "hello, this is a string"
//     prefix: "\"
//     chars: "l", "n"
// will produce the output
//     "he\l\lo, this is a stri\ng".
type Escaper func(string) string

// NewEscaper creates a new Escaper with the specified prefix and chars.
func NewEscaper(prefix string, chars ...string) Escaper {
	var oldnew []string
	for _, char := range chars {
		oldnew = append(oldnew, char, prefix+char)
	}
	r := strings.NewReplacer(oldnew...)

	return func(s string) string {
		return r.Replace(s)
	}
}
