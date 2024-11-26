package models

import "io"

type File struct {
	Name string
	Data io.Reader
}
