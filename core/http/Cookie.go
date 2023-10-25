package http

import "time"

type Cookie struct {
	Name     string
	Value    string
	Expires  time.Time
	Secure   bool
	HTTPOnly bool
}
