package critbit

import (
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestCreate(c *C) {
	// Create it
	tree := New(0)
	tree.Dump()
}
