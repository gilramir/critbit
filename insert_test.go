package critbit

import (
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestInsert(c *C) {
	// Create it
	table := make([]string, 7)
	tree := New(0)
	c.Check(tree.LoudsSlice(), DeepEquals, []byte{0})

	var ok bool
	var err error
	table[0] = "@@@"
	ok, err = tree.Insert(table[0], 0)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR   1
		//     ----  -
		[]byte{1, 0, 0})

	table[1] = "AAA"
	ok, err = tree.Insert(table[1], 1)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1     2  3
		//     ----  -------  -  -
		[]byte{1, 0, 1, 1, 0, 0, 0})

	table[2] = "BBB"
	ok, err = tree.Insert(table[2], 2)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1        4     5  2  3
		//     ----  -------  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0})

	table[3] = "ZZZ"
	ok, err = tree.Insert(table[3], 3)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1        6     7     4     5  2  3
		//     ----  -------  -------  -  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0})

	table[4] = "DDD"
	ok, err = tree.Insert(table[4], 4)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1        6     7     8     9    4      5  2  3
		//     ----  -------  -------  -  -------  -  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0})

	table[5] = "CCC"
	ok, err = tree.Insert(table[5], 5)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1        6     7     8     9    4         5     2  3 10 11
		//     ----  -------  -------  -  -------  -  -------  -------  -  - -- --
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0})

	// Can't insert a duplicate
	ok, err = tree.Insert(table[5], 5)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, false)

	table[6] = "zzz"
	ok, err = tree.Insert(table[6], 6)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR      1       12    13     6     7     8     9    4         5     2  3 10 11
		//     ----  -------  -------  -  -------  -  -------  -  -------  -------  -  - -- --
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0})
}

// There was a bug in insertSecondString that this test found.
// The offset and bit in the new node were being set to 0 and 1 for some reason,
// instead of the true offset nad bit.
func (s *MySuite) TestInsertSecondString(c *C) {
	var ok bool
	var err error
	// Create it
	table := make([]string, 3)
	tree := New(0)
	c.Check(tree.LoudsSlice(), DeepEquals, []byte{0})
	keys := tree.Keys()
	c.Assert(err, IsNil)
	c.Check(keys, IsNil)

	table[0] = "CCC"
	ok, err = tree.Insert(table[0], 0)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR  1
		//     ---  -
		[]byte{1, 0, 0})
	keys = tree.Keys()
	c.Assert(err, IsNil)
	c.Check(keys, DeepEquals, []string{"CCC"})

	table[1] = "@@@"
	ok, err = tree.Insert(table[1], 1)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR    0    1  0
		//     ---  -----  -  -
		[]byte{1, 0, 1, 1, 0, 0, 0})
	keys = tree.Keys()
	c.Assert(err, IsNil)
	c.Check(keys, DeepEquals, []string{"@@@", "CCC"})

	table[2] = "AAA"
	ok, err = tree.Insert(table[2], 2)
	c.Assert(err, IsNil)
	c.Check(ok, Equals, true)
	c.Check(tree.LoudsSlice(), DeepEquals,
		//      SR    0      1    0  1  2
		//     ---  -----  -----  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0})

	keys = tree.Keys()
	c.Assert(err, IsNil)
	c.Check(keys, DeepEquals, []string{"@@@", "AAA", "CCC"})
}
