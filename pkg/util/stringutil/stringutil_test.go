/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEscaper(t *testing.T) {
	testcases := []struct {
		prefix string
		chars  []string
		in     string
		want   string
		msg    string
	}{
		{
			prefix: `\`,
			chars:  []string{"n", "t"},
			in:     "nnnntnt",
			want:   `\n\n\n\n\t\n\t`,
			msg:    "escape characters",
		},
		{
			prefix: `\`,
			chars:  []string{`\`},
			in:     `\\\\\`,
			want:   `\\\\\\\\\\`,
			msg:    "escape when chars and prefix are same",
		},
		{
			prefix: `\`,
			chars:  []string{`\`},
			in:     ``,
			want:   ``,
			msg:    "should not touch strings that don't contain chars",
		},
		{
			prefix: `\`,
			chars:  []string{`\`},
			in:     `foobar`,
			want:   `foobar`,
			msg:    "should not touch strings that don't contain chars",
		},
		{
			prefix: `\`,
			chars:  []string{`\`},
			in:     `foobar`,
			want:   `foobar`,
			msg:    "should not touch strings that don't contain chars",
		},
		{
			prefix: `\`,
			chars:  []string{`\`, "ba"},
			in:     `foobar`,
			want:   `foo\bar`,
			msg:    "multi-char escapes",
		},
		{
			prefix: `\`,
			chars:  []string{"baba", "ba"},
			in:     `bababa`,
			want:   `\baba\ba`,
			msg:    "multi-char escapes with overlaps",
		},
		{
			prefix: `\`,
			chars:  []string{"baba", "baba"},
			in:     `bababa`,
			want:   `\bababa`,
			msg:    "multi-char escapes with overlaps",
		},
		{
			prefix: `hello`,
			chars:  []string{"baba", "baba"},
			in:     `bababa`,
			want:   `hellobababa`,
			msg:    "multi-char prefix",
		},
		{
			prefix: `hello`,
			chars:  []string{"baba", "hello"},
			in:     `bababa`,
			want:   `hellobababa`,
			msg:    "multi-char prefix with overlaps",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.msg, func(t *testing.T) {
			t.Parallel()
			e := NewEscaper(tc.prefix, tc.chars...)
			got := e(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
