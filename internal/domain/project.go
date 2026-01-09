package domain

import "errors"

type Project struct {
	id      string
	schemas any
}

func NewProject() (*Project, error) {
	return nil, errors.New("not implemented")
}
