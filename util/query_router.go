package util

import (
	"regexp"

	"github.com/quokki/quokki/types"
)

// Router provides handlers for each transaction type.
type QueryRouter interface {
	AddRoute(r string, h types.QueryHandler) (rtr QueryRouter)
	Route(path string) (h types.QueryHandler)
}

// map a transaction type to a handler
type queryRoute struct {
	r string
	h types.QueryHandler
}

type queryRouter struct {
	routes []queryRoute
}

func NewQueryRouter() *queryRouter {
	return &queryRouter{
		routes: make([]queryRoute, 0),
	}
}

var isAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func (rtr *queryRouter) AddRoute(r string, h types.QueryHandler) QueryRouter {
	if !isAlpha(r) {
		panic("route expressions can only contain alphanumeric characters")
	}
	rtr.routes = append(rtr.routes, queryRoute{r, h})

	return rtr
}

// Route - TODO add description
// TODO handle expressive matches.
func (rtr *queryRouter) Route(path string) (h types.QueryHandler) {
	for _, route := range rtr.routes {
		if route.r == path {
			return route.h
		}
	}
	return nil
}
