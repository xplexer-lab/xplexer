// Package api
// Has to be used as common application layer
package api

// Todo:
// 1. Builder
// - is responsible for building http handler
// - is responsible for generation openapi documentation
// - uses meta information from DTO's for validation and generation information
// - Builder is parametrized with generic. So colled container. Container in general is set of all infrastructure level dependencies. It also has to be used for execution transactions. I see the code next way.

// type InDto struct {
// 	Token string `header:"X-Token"`
// 	Id 		string `param:"id"`
// 	Page 	int  	 `query:"page_nr" default:"0" validation:"required"`
// 	Name 	string `json:"name"`
// }
//
// type OutDto struct {
// 	Name string `json:"name"`
// }
//
// type Deps struct {
// 	Repo1 Repo1
// 	Repo2 Repo2
// }
//
// api.Query(func(in *InDto, deps *Deps) (*Out, error) {
// 	return nil, errors.New("something went wrong")
// }, ... )
