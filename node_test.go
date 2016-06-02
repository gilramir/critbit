package critbit

import (
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestChildType(c *C) {
	// left internal, right nil
	n := &internalNode{}
	n.setChildType(0, kChildIntNode)
	c.Check(n.getChildType(0), Equals, uint8(kChildIntNode))
	c.Check(n.getChildType(1), Equals, uint8(kChildNil))

	// left nil, right internal
	n = &internalNode{}
	n.setChildType(1, kChildIntNode)
	c.Check(n.getChildType(1), Equals, uint8(kChildIntNode))
	c.Check(n.getChildType(0), Equals, uint8(kChildNil))

	// left internal, right internal
	n = &internalNode{}
	n.setChildType(0, kChildIntNode)
	n.setChildType(1, kChildIntNode)
	c.Check(n.getChildType(0), Equals, uint8(kChildIntNode))
	c.Check(n.getChildType(1), Equals, uint8(kChildIntNode))

	// left external, right external
	n = &internalNode{}
	n.setChildType(0, kChildExtRef)
	n.setChildType(1, kChildExtRef)
	c.Check(n.getChildType(0), Equals, uint8(kChildExtRef))
	c.Check(n.getChildType(1), Equals, uint8(kChildExtRef))

	// left internal, right external
	n = &internalNode{}
	n.setChildType(0, kChildIntNode)
	n.setChildType(1, kChildExtRef)
	c.Check(n.getChildType(0), Equals, uint8(kChildIntNode))
	c.Check(n.getChildType(1), Equals, uint8(kChildExtRef))

	// left external, right internal
	n = &internalNode{}
	n.setChildType(1, kChildIntNode)
	n.setChildType(0, kChildExtRef)
	c.Check(n.getChildType(1), Equals, uint8(kChildIntNode))
	c.Check(n.getChildType(0), Equals, uint8(kChildExtRef))
}

func (s *MySuite) TestChildTypeChanges(c *C) {
	// right starts external, becomes internal
	n := &internalNode{}
	n.setChildType(1, kChildExtRef)
	c.Check(n.getChildType(1), Equals, uint8(kChildExtRef))
	n.setChildType(1, kChildIntNode)
	c.Check(n.getChildType(1), Equals, uint8(kChildIntNode))
}
