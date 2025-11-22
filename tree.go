package critbit

import (
	"fmt"

	"github.com/pkg/errors"
)

// Returns the node type of the root node
func (tree *Critbit[T]) rootItemType() byte {
	switch tree.numExternalRefs {
	case 0:
		return kChildNil
	case 1:
		return kChildExtRef
	default:
		return kChildIntNode
	}
}

func (tree *Critbit[T]) addExternalRef(key string, value T) (uint32, error) {
	var refNum uint32
	if tree.firstDeletedRef == kNilRef {
		if refNum == kMaxStrings {
			return 0, errors.Errorf("Critbit is full")
		}
		refNum = uint32(len(tree.externalRefs))
		tree.externalRefs = append(tree.externalRefs, externalRef[T]{key, value, 0})
	} else {
		refNum = tree.firstDeletedRef
		tree.firstDeletedRef = tree.externalRefs[refNum].nextDeletedRef
		tree.externalRefs[int(refNum)].key = key
		tree.externalRefs[int(refNum)].value = value
		tree.externalRefs[int(refNum)].nextDeletedRef = 0
	}
	tree.numExternalRefs++
	return refNum, nil
}

func (tree *Critbit[T]) deleteExternalRef(refNum uint32) {
	var nilVal T
	tree.numExternalRefs--
	tree.externalRefs[refNum].key = ""
	tree.externalRefs[refNum].value = nilVal
	tree.externalRefs[refNum].nextDeletedRef = tree.firstDeletedRef
	tree.firstDeletedRef = refNum
}

func (tree *Critbit[T]) addInternalNode() (uint32, *internalNode) {
	var nodeNum uint32
	if tree.firstDeletedNode == kNilNode {
		nodeNum = uint32(len(tree.internalNodes))
		tree.internalNodes = append(tree.internalNodes, internalNode{})
	} else {
		nodeNum = tree.firstDeletedNode
		tree.firstDeletedNode = tree.internalNodes[nodeNum].child[1]
		tree.internalNodes[int(nodeNum)] = internalNode{}
		tree.internalNodes[int(nodeNum)].child[0] = kNilNode
		tree.internalNodes[int(nodeNum)].child[1] = kNilNode
	}
	tree.numInternalNodes++
	return nodeNum, &tree.internalNodes[nodeNum]
}

func (tree *Critbit[T]) deleteInternalNode(nodeNum uint32) {
	tree.numInternalNodes--
	tree.internalNodes[nodeNum].child[1] = tree.firstDeletedNode
	tree.firstDeletedNode = nodeNum
}

// The caller must ensure that rootItem is valid (either a ref or a node)
func (tree *Critbit[T]) findBestExternalReference(key string) uint32 {
	// If there is only one ref, then it must be the best choice
	if tree.numExternalRefs == 1 {
		return tree.rootItem
	}

	nodeNum := tree.rootItem
	for {
		node := &tree.internalNodes[nodeNum]
		direction := node.direction(key)
		childType := node.getChildType(direction)
		switch childType {
		case kChildIntNode:
			nodeNum = node.child[direction]
		case kChildExtRef:
			return node.child[direction]
		default:
			panic(fmt.Sprintf("Child %d of nodeNum %d has unexpected type 0x%02x",
				direction, nodeNum, childType))
		}
	}
}

// The caller must ensure that rootItem is valid (either a ref or a node)
// Returns extRefNum, grandparentNodeNum, grandparentDirection, parentNodeNum, parentDirection, parentIsRoot
func (tree *Critbit[T]) findBestExternalReferenceWithAncestry(key string) (uint32, uint32, byte, uint32, byte, bool) {
	// If there is only one ref, then it must be the best choice
	if tree.numExternalRefs == 1 {
		return tree.rootItem, 0, 0, 0, 0, false
	}

	var parentIsRoot bool = true
	var grandparentNodeNum uint32
	var grandparentDirection byte
	var parentNodeNum uint32
	var parentDirection byte

	nodeNum := tree.rootItem
	for {
		node := &tree.internalNodes[nodeNum]

		grandparentDirection = parentDirection
		parentDirection = node.direction(key)
		grandparentNodeNum = parentNodeNum
		parentNodeNum = nodeNum

		childType := node.getChildType(parentDirection)
		switch childType {
		case kChildIntNode:
			nodeNum = node.child[parentDirection]
			parentIsRoot = false
		case kChildExtRef:
			return node.child[parentDirection], grandparentNodeNum, grandparentDirection, parentNodeNum, parentDirection, parentIsRoot
		default:
			panic(fmt.Sprintf("Child %d of nodeNum %d has unexpected type 0x%02x",
				parentDirection, nodeNum, childType))
		}
	}
}

// Returns identical, off, bit, ndir, err
func (tree *Critbit[T]) findCriticalBit(refNum uint32, newKey string) (bool, uint16, byte, byte) {
	return findCriticalBit(tree.externalRefs[refNum].key, newKey)
}

// Returns identical, off, bit, ndir
func findCriticalBit(storedKey string, newKey string) (bool, uint16, byte, byte) {
	// find critical bit
	var off uint16
	var ch, bit byte
	// find differing byte
	for off = 0; off < uint16(len(newKey)); off++ {
		if ch = 0; off < uint16(len(storedKey)) {
			ch = storedKey[off]
		}
		if keych := newKey[off]; ch != keych {
			bit = ch ^ keych
			goto ByteFound
		}
	}
	if off < uint16(len(storedKey)) {
		ch = storedKey[off]
		bit = ch
		goto ByteFound
	}
	return true, 0, 0, 0

ByteFound:
	// find differing bit
	bit |= bit >> 1
	bit |= bit >> 2
	bit |= bit >> 4
	bit = bit &^ (bit >> 1)
	var ndir byte
	if ch&bit != 0 {
		ndir++
	}
	return false, off, bit, ndir
}

// The caller must ensure that there is at least one internal node
// Returns nodeNum, parentNode, prevDirection, insertAtRoot, finalChildType
func (tree *Critbit[T]) findBranchNode(off uint16, bit byte,
	key string) (uint32, uint32, byte, bool, byte) {
	var parentNodeNum uint32 = 0
	var prevDirection byte
	var insertAtRoot bool = true

	nodeNum := tree.rootItem

	for {
		node := &tree.internalNodes[nodeNum]
		if node.offset > off || node.offset == off && node.bit < bit {
			return nodeNum, parentNodeNum, prevDirection, insertAtRoot, kChildIntNode
		}
		// try the next node
		direction := node.direction(key)
		childType := node.getChildType(direction)
		switch childType {
		case kChildIntNode:
			parentNodeNum = nodeNum
			prevDirection = direction
			insertAtRoot = false
			nodeNum = node.child[direction]
		case kChildExtRef:
			return node.child[direction], nodeNum, direction, false, kChildExtRef
		default:
			panic(fmt.Sprintf("Child %d of nodeNum %d has unexpected type 0x%02x",
				direction, nodeNum, childType))
		}
	}
	panic("should not reach here")
	return 0, 0, 0, false, 0
}
