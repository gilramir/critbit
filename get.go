package critbit

// Get finds the key and returns its value. The boolean
// indicates if it was found or not.
func (tree *Critbit) Get(key string) (bool, interface{}) {
	has, refNum := tree.findRef(key)
	if !has {
		return false, nil
	}
	return true, tree.externalRefs[refNum].value
}

// Returns: ok, refNum
func (tree *Critbit) findRef(key string) (bool, uint32) {
	// Is the tree empty? Nothing to find.
	if len(tree.externalRefs) == 0 {
		return false, 0
	}

	// Find the best external reference
	bestRefNum := tree.findBestExternalReference(key)

	// find critical bit, but more importantly, is there
	// an identical match?
	identical, _, _, _ := tree.findCriticalBit(bestRefNum, key)

	// Is it already in the tree?
	if identical {
		return true, bestRefNum
	} else {
		return false, 0
	}
}
