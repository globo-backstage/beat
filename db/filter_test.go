package db

import (
	"gopkg.in/check.v1"
	"testing"
)

var _ = check.Suite(&S{})

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestNewFilterFromQueryStringEmpty(c *check.C) {
	filter, err := NewFilterFromQueryString("")
	c.Assert(err, check.IsNil)
	c.Assert(filter, check.NotNil)
	c.Assert(filter.PerPage, check.Equals, 10)

	whereMap, err := filter.Where.Map()
	c.Assert(err, check.IsNil)
	c.Assert(len(whereMap), check.Equals, 0)
}

func (s *S) TestNewFilterFromQueryPerPage(c *check.C) {
	filter, err := NewFilterFromQueryString("filter[perPage]=15")
	c.Assert(err, check.IsNil)
	c.Assert(filter.PerPage, check.Equals, 15)
}

func (s *S) TestNewFilterFromQueryPerPageOverFlow(c *check.C) {
	filter, err := NewFilterFromQueryString("filter[perPage]=10000")
	c.Assert(err, check.IsNil)
	c.Assert(filter.PerPage, check.Equals, 1000)
}

func (s *S) TestNewFilterFromQueryWhere(c *check.C) {
	filter, err := NewFilterFromQueryString("filter[where][name]=wilson")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Where.Get("name").MustString(), check.Equals, "wilson")

	filter, err = NewFilterFromQueryString("filter[where][name]=wilson&filter[where][title][like]=juju")
	c.Assert(err, check.IsNil)

	c.Assert(filter.Where.Get("name").MustString(), check.Equals, "wilson")
	c.Assert(filter.Where.GetPath("title", "like").MustString(), check.Equals, "juju")
}

func (s *S) TestNewFilterFromQueryPage(c *check.C) {
	filter, err := NewFilterFromQueryString("filter[page]=2")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Page, check.Equals, 2)

	filter, err = NewFilterFromQueryString("filter[page]=0")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Page, check.Equals, 1)
}

func (s *S) TestNewFilterSkip(c *check.C) {
	filter, err := NewFilterFromQueryString("")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Skip(), check.Equals, 0)

	filter, err = NewFilterFromQueryString("filter[page]=2&filter[perPage]=100")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Skip(), check.Equals, 100)

	filter, err = NewFilterFromQueryString("filter[page]=1&filter[perPage]=100")
	c.Assert(err, check.IsNil)
	c.Assert(filter.Skip(), check.Equals, 0)
}
