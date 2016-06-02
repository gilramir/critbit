package critbit

import (
	"log"
	"testing"

	// Bring the symbols in check.v1 into this namespace
	. "gopkg.in/check.v1"
)

// Hook gocheck into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}
