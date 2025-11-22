package critbit

import (
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestDelete(c *C) {
	// Create it
	table := make([]string, 7)
	tree := New[int64](0)

	table[0] = "@@@"
	table[1] = "AAA"
	table[2] = "BBB"
	table[3] = "ZZZ"
	table[4] = "DDD"
	table[5] = "CCC"
	table[6] = "zzz"

	var ok bool
	var err error
	var i int64
	for i = 0; i < 7; i++ {
		ok, err = tree.Insert(table[i], i)
		c.Assert(err, IsNil)
		c.Check(ok, Equals, true)
	}
	// The final tree looks like this
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1       12    13     6     7     8     9    4         5     2  3 10 11
		//     ----  -------  -------  -  -------  -  -------  -  -------  -------  -  - -- --
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0})

	ok = tree.Delete(table[6])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1        6     7     8     9    4         5     2  3 10 11
		//     ----  -------  -------  -  -------  -  -------  -------  -  - -- --
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0})

	// Can't delete nox-existing key
	ok = tree.Delete(table[6])
	c.Check(ok, Equals, false)

	ok = tree.Delete(table[5])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1        6     7     8     9    4      5  2  3
		//     ----  -------  -------  -  -------  -  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0})

	ok = tree.Delete(table[4])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1        6     7     4     5  2  3
		//     ----  -------  -------  -  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0})

	ok = tree.Delete(table[3])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1        4     5  2  3
		//     ----  -------  -------  -  -  -
		[]byte{1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0})

	ok = tree.Delete(table[2])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR      1     2  3
		//     ----  -------  -  -
		[]byte{1, 0, 1, 1, 0, 0, 0})

	ok = tree.Delete(table[1])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		//      SR   1
		//     ----  -
		[]byte{1, 0, 0})

	ok = tree.Delete(table[0])
	c.Check(ok, Equals, true)
	c.Check(tree.Louds().ToBytes(), DeepEquals, []byte{0})
}
