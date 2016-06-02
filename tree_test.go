package critbit

import (
	"log"
	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestFindCritBit(c *C) {
	// Returns identical, off, bit, ndir
	//func findCriticalBit(storedKey string, newKey string) (bool, uint16, byte, byte) {

	var identical bool
	var off uint16
	var bit byte
	var ndir byte

	// @ = 0x40 A = 0x41
	identical, off, bit, ndir = findCriticalBit("@", "A")
	log.Printf("identical=%v off=%d bit=0x%0x ndir=%d", identical, off, uint8(bit), ndir)
	c.Check(identical, Equals, false)
	c.Check(off, Equals, uint16(0))
	c.Check(bit, Equals, uint8(1))
	c.Check(ndir, Equals, uint8(0))

	// @ = 0x40 A = 0x41
	identical, off, bit, ndir = findCriticalBit("A", "@")
	log.Printf("identical=%v off=%d bit=0x%0x ndir=%d", identical, off, uint8(bit), ndir)
	c.Check(identical, Equals, false)
	c.Check(off, Equals, uint16(0))
	c.Check(bit, Equals, uint8(1))
	c.Check(ndir, Equals, uint8(1))

	identical, off, bit, ndir = findCriticalBit("A", "A")
	log.Printf("identical=%v off=%d bit=0x%0x ndir=%d", identical, off, uint8(bit), ndir)
	c.Check(identical, Equals, true)
}
