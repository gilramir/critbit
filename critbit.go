// Package critbit provides an implementation of a Critbit tree
// that stores nodes in arrays rather than allocating them individually
// and relying on pointers.
package critbit

const (
	kNilNode = (1<<32 - 1)
	kNilRef  = kNilNode

	kChildNil     = 0x00 // 00000000
	kChildIntNode = 0x01 // 00000001
	kChildExtRef  = 0x02 // 00000010
	kChildBitmask = 0x03 // 00000011

	kLeftMask  = 0xcf // 11001111
	kRightMask = 0xfc // 11111100

	kDirectionLeft  = 0
	kDirectionRight = 1

	// Since we use a uint32 to keep track of internal nodes and
	// external references, we can store up to 2^32 of each. One value,
	// 0xffffffff is used as a "nil" value. If we have N external refs,
	// we need N-1 interal nodes to differentiate them. So, the max number
	// of external refs (strings) is (2^32)-1, and accordingly, the max
	// number of internal nodes we would use would be (2^32)-2.
	kMaxStrings = (1 << 32) - 1

	// A uint16 value is used to store the offset within a string,
	// so the maximum allowed string length is 65,536.
	MaxStringLength = 65536

	kMaxStringLength = MaxStringLength
)

// A Critbit represents one Critbit tree.
type Critbit struct {
	totalStringSize int

	internalNodes []internalNode
	externalRefs  []externalRef

	numInternalNodes int    // num used, not num allocated
	numExternalRefs  int    // num used, not num allocated
	rootItem         uint32 // root node, or if no nodes, root ref
	firstDeletedNode uint32 // kNilNode if none are deleted
	firstDeletedRef  uint32 // kNilRef if none are deleted
}

type internalNode struct {
	offset uint16
	bit    uint8
	flags  uint8     // dirty, leftChildType=(nil|int|ext), rightChildType=(nil|int|ext)
	child  [2]uint32 // if deleted, child[1] = nextDeleted
}

type externalRef struct {
	key   string
	value interface{} // if deleted, uint32 pointing to nextDeleted
}

// New allocates a new Critbit tree and returns a pointer to it.
// The capacityStrings argument allows smart allocation of the internal arrays,
// if you happen to know that a tree will contain a certain amount of
// strings. This is for efficiency only; capacityStrings does not impose any
// limit on the number of strings.
func New(capacityStrings int) *Critbit {
	// For every external ref (string), we need one branching (internal) node,
	// except for the very first external ref (hence, the minus one).
	var capacityInternalNodes int = 0
	if capacityStrings > 0 {
		capacityInternalNodes = capacityStrings - 1
	}

	return &Critbit{
		internalNodes:    make([]internalNode, 0, capacityInternalNodes),
		externalRefs:     make([]externalRef, 0, capacityStrings),
		firstDeletedNode: kNilNode,
		firstDeletedRef:  kNilRef,
	}
}

// Length returns the number of keys currently stored in the tree. More space for
// keys may have been allocated, if keys were deleted and no other
// keys were inserted.
func (tree *Critbit) Length() int {
	return tree.numExternalRefs
}

// MemorySize returns the approximate number of bytes used used
// by the tree.
func (tree *Critbit) MemorySizeBytes() int {
	return 4*8 + 12 + // static part of Critbit
		(20 * len(tree.internalNodes)) +
		(8 * len(tree.externalRefs)) +
		tree.totalStringSize
}
