package critbit

type LOUDS []byte

func (s LOUDS) ToBytes() []byte {
	return []byte(s)
}

// LoudsSlice returns a slice of bytes, all of which are either 1 or 0,
// which represent the tree structure in the LOUDS (level-order unary
// degree separation) representation, a succinct representation of the tree.
// Given N nodes (internal nodes + external refs), 2N+1 bytes will
// be returned in the slice. See
// https://memoria-framework.dev/docs/data-zoo/louds-tree/
// for an introduction to LOUDS.
func (tree *Critbit) Louds() LOUDS {
	n := tree.numInternalNodes + tree.numExternalRefs
	if n == 0 {
		return LOUDS([]byte{0})
	}
	bits := 2*n + 1
	answer := make([]byte, 2, bits)
	// Fake super-root
	answer[0] = 1
	answer[1] = 0

	byteChan := make(chan byte)

	go tree._loudSlice(byteChan)

	for b := range byteChan {
		answer = append(answer, b)
	}
	return LOUDS(answer)
}

type itemTuple struct {
	itemID   uint32
	itemType byte
}

// Walk breadth-first
func (tree *Critbit) _loudSlice(byteChan chan byte) {
	defer close(byteChan)
	rootLayer := make([]itemTuple, 1)
	rootLayer[0].itemID = tree.rootItem
	rootLayer[0].itemType = tree.rootItemType()
	tree._loudSliceLayer(byteChan, rootLayer)
}

func (tree *Critbit) _loudSliceLayer(byteChan chan byte, layer []itemTuple) {
	nextLayer := make([]itemTuple, 0)

	for _, item := range layer {
		children := tree._loudSliceNode(byteChan, item)
		nextLayer = append(nextLayer, children...)
	}
	if len(nextLayer) > 0 {
		tree._loudSliceLayer(byteChan, nextLayer)
	}
}

func (tree *Critbit) _loudSliceNode(byteChan chan byte, item itemTuple) []itemTuple {
	switch item.itemType {
	case kChildNil:
		panic("not reached")
	case kChildExtRef:
		byteChan <- 0
		return nil
	case kChildIntNode:
		// By definition, and internal node has 2 children.
		byteChan <- 1
		byteChan <- 1
		byteChan <- 0
		node := tree.internalNodes[item.itemID]
		children := make([]itemTuple, 2)
		children[0].itemID = node.child[0]
		children[0].itemType = node.getChildType(0)
		children[1].itemID = node.child[1]
		children[1].itemType = node.getChildType(1)
		return children
	}
	return nil
}
