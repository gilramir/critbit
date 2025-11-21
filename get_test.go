package critbit

import (
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestGet01(c *C) {
	trie := New(0)

	trie.Insert("zoo", 1)
	trie.Insert("green", 1)
	trie.Insert("green boy", 1)
	trie.Insert("gremlin", 1)
	trie.Insert("gray", 1)
	trie.Insert("gas", 1)
	trie.Insert("apple", 1)
	trie.Insert("fan", 1)

	// Exact matches
	has, _ := trie.Get("gas")
	c.Assert(has, Equals, true)

	has, _ = trie.Get("gr")
	c.Assert(has, Equals, false)

	// Non-exact matches
	kvt := trie.GetHasPrefix("gas")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "gas")

	kvt = trie.GetHasPrefix("gr")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "gray")

	kvt = trie.GetHasPrefix("gra")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "gray")

	kvt = trie.GetHasPrefix("gre")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "green")

	kvt = trie.GetHasPrefix("g")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "gas")

	kvt = trie.GetHasPrefix("b")
	c.Assert(kvt, IsNil)

	kvt = trie.GetHasPrefix("y")
	c.Assert(kvt, IsNil)

	kvt = trie.GetHasPrefix("z")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "zoo")

	kvt = trie.GetHasPrefix("a")
	c.Assert(kvt, NotNil)
	c.Assert(kvt.Key, Equals, "apple")
}
