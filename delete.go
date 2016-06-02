package critbit

import (
	"fmt"
)

// Delete removes the key from the tree. The boolean return value
// indicates if the key was in the tree.
func (tree *Critbit) Delete(key string) bool {
	// Is the tree empty? Do nothing
	if tree.numExternalRefs == 0 {
		return false
	}

	// Find the best external reference
	bestRefNum, grandparentNodeNum, grandparentDirection, parentNodeNum, parentDirection,
		parentIsRoot := tree.findBestExternalReferenceWithAncestry(key)

	// find critical bit
	identical, _, _, _ := tree.findCriticalBit(bestRefNum, key)

	// Is it NOT in the tree?
	if !identical {
		return false
	}

	// delete from tree
	tree.deleteExternalRef(bestRefNum)

	// Was that the last ref? Then there are no internal nodes to worry about.
	if tree.numExternalRefs == 0 {
		tree.rootItem = 0
		return true
	}

	if parentIsRoot {
		if parentNodeNum != tree.rootItem {
			panic(fmt.Sprintf("parent is root, parent=%d, root=%d", parentNodeNum, tree.rootItem))
		}
		parentNode := &tree.internalNodes[parentNodeNum]
		siblingItemNum := parentNode.child[1-parentDirection]
		// Our sibling becomes the root item
		tree.deleteInternalNode(parentNodeNum)
		tree.rootItem = siblingItemNum
		return true
	}

	parentNode := &tree.internalNodes[parentNodeNum]
	siblingItemNum := parentNode.child[1-parentDirection]

	// the parent gets elided, and the child replaces the parent
	grandparentNode := &tree.internalNodes[grandparentNodeNum]
	grandparentNode.setChild(grandparentDirection, siblingItemNum, parentNode.getChildType(1-parentDirection))
	tree.deleteInternalNode(parentNodeNum)

	return true
}
