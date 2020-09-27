/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package ioutil

import "io"

// CheckCloseErr is a util that checks for errors returned when closing an io.Closer.
func CheckCloseErr(c io.Closer, err *error) {
	cErr := c.Close()
	if *err == nil {
		*err = cErr
	}
}
