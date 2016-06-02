package critbit

import (
	// Bring the symbols in check.v1 into this namespace
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
		keys = append(keys, keyTuple.key)
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
