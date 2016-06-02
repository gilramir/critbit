package critbit

import (
	"github.com/pkg/errors"
)

// Insert inserts a key/value pair into the tree. It returns
// a boolean indicating if the key was inserted. If the key already exists,
// in the tree, the value will not be changed, and false will be returned.
// An error value is also returned. An error will occur if the tree
// is full, or if the string is too long to be inserted.
// If an error is returned, the boolean value returned will be false.
func (tree *Critbit) Insert(key string, value interface{}) (bool, error) {
	// Sanity check
	if len(key) > kMaxStringLength {
		return false, errors.Errorf("Maximum string length is %d", kMaxStringLength)
	}

	// Is the tree empty? Insert the first ref
	if tree.numExternalRefs == 0 {
		err := tree.insertFirstString(key, value)
		if err != nil {
			return false, errors.Wrap(err, "Insert() first key")
		}
		return true, nil
	}

	// Find the best external reference
	bestRefNum := tree.findBestExternalReference(key)

	// find critical bit
	identical, off, bit, ndir := tree.findCriticalBit(bestRefNum, key)

	// Is it already in the tree?
	if identical {
		return false, nil
	}

	// If there is only one external ref, then there are no internal nodes.
	// Insert the first node (and a new ref)
	if tree.numExternalRefs == 1 {
		err := tree.insertSecondString(key, value, off, bit, ndir)
		if err != nil {
			return false, errors.Wrap(err, "Insert() second key")
		}
		return true, nil
	}

	// Find the node from which to branch
	branchNodeNum, parentNodeNum, prevDirection, insertAtRoot,
		finalChildType := tree.findBranchNode(off, bit, key)

	// Add the new ref
	newRefNum, err := tree.addExternalRef(key, value)
	if err != nil {
		return false, errors.Wrap(err, "Insert() adding an external ref")
	}

	// Insert a new node, which be inserted where the branching node currently is.
	newNodeNum, newNode := tree.addInternalNode()
	newNode.setChild(1-ndir, newRefNum, kChildExtRef)
	newNode.offset = off
	newNode.bit = bit

	if insertAtRoot {
		// The new node becomes the new tree root
		newNode.setChild(ndir, tree.rootItem, kChildIntNode)
		tree.rootItem = newNodeNum
	} else {
		// The branch node's parent points to the new node, and the new node
		// subsumes the branch node. The parent-child connection from the new
		// node to the branch node must indicate the correct child type.
		parentNode := &tree.internalNodes[parentNodeNum]
		parentNode.setChild(prevDirection, newNodeNum, kChildIntNode)
		newNode.setChild(ndir, branchNodeNum, finalChildType)
	}

	return true, nil
}

// Adds the first ref but no node
func (tree *Critbit) insertFirstString(key string, value interface{}) error {
	refNum, err := tree.addExternalRef(key, value)
	if err != nil {
		return err
	}
	tree.rootItem = refNum
	return nil
}

// Adds the first node, and sets the existing single ref as a child,
// and adds another ref for the other child.
func (tree *Critbit) insertSecondString(key string, value interface{}, off uint16, bit byte, ndir byte) error {
	refNum, err := tree.addExternalRef(key, value)
	if err != nil {
		return err
	}
	nodeNum, node := tree.addInternalNode()
	node.offset = off
	node.bit = bit
	node.setChild(1-ndir, refNum, kChildExtRef)
	node.setChild(ndir, tree.rootItem, kChildExtRef)
	tree.rootItem = nodeNum
	return nil
}
