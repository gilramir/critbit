package critbit

import (
	"fmt"
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) testSplit(c *C, tree *Critbit[int64], table []string, name string) {
	var ok bool
	var err error
	var i int64
	// Assume the keys are 0-n, and assume that the keys are in alphabetical order!
	numKeys := len(table)
	for splitAt := 0; splitAt < numKeys; splitAt++ {
		leftSplit, _ := tree.SplitAt(splitAt)

		// Make the natural versions of the trees
		leftNatural := New[int64](numKeys)
		for i = 0; i < int64(splitAt); i++ {
			ok, err = leftNatural.Insert(table[i], i)
			c.Assert(err, IsNil)
			c.Check(ok, Equals, true)
		}
		rightNatural := New[int64](numKeys)
		for i = int64(splitAt); i < int64(numKeys); i++ {
			ok, err = rightNatural.Insert(table[i], i)
			c.Assert(err, IsNil)
			c.Check(ok, Equals, true)
		}

		leftSame := compareLouds(leftSplit, leftNatural, name, "left", splitAt)
		c.Check(leftSame, Equals, true)
		rightSame := compareLouds(leftSplit, leftNatural, name, "left", splitAt)
		c.Check(rightSame, Equals, true)
	}
}

func cmpByteSlice(s1 []byte, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func compareLouds[T any](tree1 *Critbit[T], tree2 *Critbit[T], basename string, name string, n int) bool {
	louds1 := tree1.Louds().ToBytes()
	louds2 := tree2.Louds().ToBytes()
	if !cmpByteSlice(louds1, louds2) {
		fmt.Printf("Tree 1: %v\n", louds1)
		fmt.Printf("Tree 2: %v\n", louds2)
		tree1.SaveDot(fmt.Sprintf("%s-%s-%d-split.dot", basename, name, n))
		tree2.SaveDot(fmt.Sprintf("%s-%s-%d-natural.dot", basename, name, n))
		return false
	}
	return true
}

func (s *MySuite) TestSplit1(c *C) {
	// Create it
	table := make([]string, 7)
	tree := New[int64](7)

	table[0] = "@@@"
	table[1] = "AAA"
	table[2] = "BBB"
	table[3] = "CCC"
	table[4] = "DDD"
	table[5] = "ZZZ"
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
		[]byte{1, 0,
			1, 1, 0,
			1, 1, 0, 0,
			1, 1, 0, 0,
			1, 1, 0, 0,
			1, 1, 0, 1, 1, 0,
			0, 0, 0, 0,
		})

	s.testSplit(c, tree, table, "split1")
}

func (s *MySuite) TestSplit2(c *C) {
	// Create it
	table := make([]string, 8)
	tree := New[int64](8)

	table[0] = "a"
	table[1] = "b"
	table[2] = "c"
	table[3] = "d"
	table[4] = "k"
	table[5] = "l"
	table[6] = "m"
	table[7] = "naa"

	var ok bool
	var err error
	var i int64
	for i = 0; i < 8; i++ {
		ok, err = tree.Insert(table[i], i)
		c.Assert(err, IsNil)
		c.Check(ok, Equals, true)
	}
	// The final tree looks like this
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		[]byte{1, 0,
			1, 1, 0,
			1, 1, 0, 1, 1, 0,
			1, 1, 0, 0, 0, 1, 1, 0,
			0, 1, 1, 0, 1, 1, 0, 0,
			0, 0, 0, 0,
		})

	s.testSplit(c, tree, table, "split2")
}

func (s *MySuite) TestSplit3(c *C) {
	// Create it
	table := make([]string, 14)
	tree := New[int64](14)

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

	var ok bool
	var err error
	var i int64
	for i = 0; i < 14; i++ {
		ok, err = tree.Insert(table[i], i)
		c.Assert(err, IsNil)
		c.Check(ok, Equals, true)
	}
	// The final tree looks like this
	c.Check(tree.Louds().ToBytes(), DeepEquals,
		[]byte{1, 0,
			1, 1, 0,
			1, 1, 0, 0,
			1, 1, 0, 1, 1, 0,
			1, 1, 0, 0, 0, 1, 1, 0,
			0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
			0, 0, 0, 0, 1, 1, 0, 0,
			1, 1, 0, 0,
			1, 1, 0, 0,
			0, 1, 1, 0,
			0, 0,
		})

	s.testSplit(c, tree, table, "split3")
}
