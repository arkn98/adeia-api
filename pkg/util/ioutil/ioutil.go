package ioutil

import "io"

// CheckCloseErr is a util that checks for errors returned when closing an io.Closer.
func CheckCloseErr(c io.Closer, err *error) {
	cErr := c.Close()
	if *err == nil {
		*err = cErr
	}
}
