package db

import (
	simplejson "github.com/bitly/go-simplejson"
	"net/url"
	"strconv"
	"strings"
)

type Filter struct {
	Where   *simplejson.Json
	Page    int
	PerPage int
}

func NewFilterFromQueryString(q string) (*Filter, error) {
	filter := &Filter{}
	filter.Where = simplejson.New()
	filter.loadInitialValues()

	urlValues, err := url.ParseQuery(q)
	if err != nil {
		return nil, err
	}

	for key, value := range urlValues {
		filter.putUrlValue(key, value[0])
	}

	return filter, nil
}

func (filter *Filter) loadInitialValues() {
	filter.Page = 1
	filter.PerPage = 10
}

func (filter *Filter) Skip() int {
	return (filter.Page - 1) * filter.PerPage
}

func (filter *Filter) putUrlValue(key, value string) {
	path := []string{}

	for _, part := range strings.Split(key, "[") {
		if last := part[len(part)-1]; last == ']' {
			part = part[:len(part)-1]
		}

		path = append(path, part)

	}

	if path[0] == "filter" && len(path) > 1 {
		switch path[1] {
		case "perPage":
			filter.setPerPageFromString(value)
			return
		case "page":
			filter.setPageFromString(value)
			return
		case "where":
			if len(path) > 2 {
				filter.putWhere(path[2:], value)
			}
		}
	}
}

func (filter *Filter) setPerPageFromString(perPage string) {
	if s, err := strconv.Atoi(perPage); err == nil {
		if s > 1000 {
			s = 1000
		}
		filter.PerPage = s
	}
}

func (filter *Filter) setPageFromString(page string) {
	if s, err := strconv.Atoi(page); err == nil {
		if s < 1 {
			s = 1
		}
		filter.Page = s
	}
}

func (filter *Filter) putWhere(path []string, value string) {
	filter.Where.SetPath(path, value)
}
