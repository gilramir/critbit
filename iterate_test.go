package critbit

import (
	. "gopkg.in/check.v1"
)

func testIterateKeyTuples(c *C, table []string) []string {
	// Create the tree
	tree := New(len(table))
	var ok bool
	var err error
	i := 0
	for value, key := range table {
		ok, err = tree.Insert(key, value)
		c.Assert(err, IsNil)
		c.Check(ok, Equals, true)
		i++
	}
	// Get the keys
	var keys []string
	tupleChan := tree.GetKeyValueTuples()
	for keyTuple := range tupleChan {
		keys = append(keys, keyTuple.Key)
	}
	return keys
}

func (s *MySuite) TestIterate0(c *C) {
	table := make([]string, 4)
	table[0] = "@@@"
	table[1] = "AAA"
	table[2] = "BBB"
	table[3] = "CCC"

	keys := testIterateKeyTuples(c, table)
	c.Check(keys, DeepEquals, []string{"@@@", "AAA", "BBB", "CCC"})
}

func (s *MySuite) TestIterate1(c *C) {
	table := make([]string, 7)
	table[0] = "@@@"
	table[1] = "AAA"
	table[2] = "BBB"
	table[3] = "CCC"
	table[4] = "DDD"
	table[5] = "ZZZ"
	table[6] = "zzz"

	keys := testIterateKeyTuples(c, table)
	c.Check(keys, DeepEquals, []string{"@@@", "AAA", "BBB", "CCC", "DDD", "ZZZ", "zzz"})
}

func (s *MySuite) TestIterate2(c *C) {
	table := make([]string, 8)
	table[0] = "a"
	table[1] = "b"
	table[2] = "c"
	table[3] = "d"
	table[4] = "k"
	table[5] = "l"
	table[6] = "m"
	table[7] = "naa"

	keys := testIterateKeyTuples(c, table)
	c.Check(keys, DeepEquals, []string{"a", "b", "c", "d", "k", "l", "m", "naa"})
}

func (s *MySuite) TestIterate3(c *C) {
	table := make([]string, 14)
	table[0] = "a"
	table[1] = "b"
	table[2] = "c"
	table[3] = "d"
	table[4] = "k"
	table[5] = "l"
	table[6] = "m"
	table[7] = "naa"
	table[8] = "nab"
	table[9] = "nac"
	table[10] = "nad"
	table[11] = "nba"
	table[12] = "o"
	table[13] = "p"

	keys := testIterateKeyTuples(c, table)
	c.Check(keys, DeepEquals, []string{"a", "b", "c", "d", "k", "l", "m", "naa",
		"nab", "nac", "nad", "nba", "o", "p"})
}

func (s *MySuite) TestGetKeyValueTuplesFrom01(c *C) {
	// Create it
	trie := New(0)
	trie.Insert("red", 1)
	trie.Insert("red apple", 1)
	trie.Insert("red box", 1)
	trie.Insert("red crayon", 1)
	trie.Insert("blue", 1)
	trie.Insert("blue arrow", 1)
	trie.Insert("blue boy", 1)
	trie.Insert("blue car", 1)
	trie.Insert("green", 1)
	trie.Insert("gremlin", 1)
	trie.Insert("green action", 1)
	trie.Insert("green babble", 1)
	trie.Insert("green crown", 1)

	//trie.Dump()
	has, _ := trie.Get("green")
	c.Check(has, Equals, true)

	// Starting from leaf, only 1
	tuples := trie.GetKeyValueTuplesFrom("blue car", true, 3)
	c.Assert(len(tuples), Equals, 1)
	c.Check(tuples[0].Key, Equals, "blue car")

	tuples = trie.GetKeyValueTuplesFrom("green", true, 3)
	c.Assert(len(tuples), Equals, 3)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")

	tuples = trie.GetKeyValueTuplesFrom("green", true, 4)
	c.Assert(len(tuples), Equals, 4)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")
	c.Check(tuples[3].Key, Equals, "green crown")

	tuples = trie.GetKeyValueTuplesFrom("green", true, 5)
	c.Assert(len(tuples), Equals, 4)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")
	c.Check(tuples[3].Key, Equals, "green crown")

	tuples = trie.GetKeyValueTuplesFrom("green", true, 0)
	c.Assert(len(tuples), Equals, 4)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")
	c.Check(tuples[3].Key, Equals, "green crown")

	// "gr" doesn't match exactly
	tuples = trie.GetKeyValueTuplesFrom("gr", true, 0)
	c.Assert(len(tuples), Equals, 0)
	c.Check(len(tuples), Equals, 0)
}

func (s *MySuite) TestGetKeyValueTuplesFrom02(c *C) {
	// Create it
	trie := New(0)
	trie.Insert("red", 1)
	trie.Insert("red apple", 1)
	trie.Insert("red box", 1)
	trie.Insert("red crayon", 1)
	trie.Insert("blue", 1)
	trie.Insert("blue arrow", 1)
	trie.Insert("blue boy", 1)
	trie.Insert("blue car", 1)
	trie.Insert("green", 1)
	trie.Insert("gremlin", 1)
	trie.Insert("green action", 1)
	trie.Insert("green babble", 1)
	trie.Insert("green crown", 1)

	//	trie.Dump()
	has, _ := trie.Get("green")
	c.Check(has, Equals, true)

	// Check !showExact against actual exact matches
	tuples := trie.GetKeyValueTuplesFrom("blue car", false, 3)
	c.Assert(len(tuples), Equals, 1)
	c.Check(tuples[0].Key, Equals, "blue car")

	tuples = trie.GetKeyValueTuplesFrom("green", false, 3)
	c.Assert(len(tuples), Equals, 3)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")

	// Check non-exact matches
	tuples = trie.GetKeyValueTuplesFrom("gr", false, 0)
	c.Assert(len(tuples), Equals, 4)
	c.Check(tuples[0].Key, Equals, "green")
	c.Check(tuples[1].Key, Equals, "green action")
	c.Check(tuples[2].Key, Equals, "green babble")
	c.Check(tuples[3].Key, Equals, "green crown")

	trie.Dump()
}
