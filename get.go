package critbit

import "strings"

// Get finds the key and returns its value. The boolean
// indicates if it was found or not.
func (tree *Critbit) Get(key string) (bool, interface{}) {
	has, refNum := tree.findRef(key)
	if !has {
		return false, nil
	}
	return true, tree.externalRefs[refNum].value
}

// Returns the first key that starts with a string, and returns
// the KeyValueTuple, or nil
func (tree *Critbit) GetHasPrefix(key string) *KeyValueTuple {

	has, refNum := tree.findRef(key)

	if !has {
		// Not an exact match, but, did we find something that does start
		// with our string?
		foundRefKey := tree.externalRefs[refNum].key
		if !strings.HasPrefix(foundRefKey, key) {
			// No, the best ref does not start with the user's key
			return nil
		}
		// keep going!
	}
	return &KeyValueTuple{
		Key:   tree.externalRefs[refNum].key,
		Value: tree.externalRefs[refNum].value,
	}
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

	return identical, bestRefNum
}

// Returns identicalMatch?, refNum, parentNodeNum, parentDirection
func (tree *Critbit) findRefWithAncestry(key string) (bool, uint32, uint32, byte) {
	// Is the tree empty? Nothing to find.
	if len(tree.externalRefs) == 0 {
		return false, 0, 0, 0
	}

	// Find the best external reference
	bestRefNum, _, _, parentNodeNum, parentDirection, _ := tree.findBestExternalReferenceWithAncestry(key)

	// find critical bit, but more importantly, is there
	// an identical match?
	identical, _, _, _ := tree.findCriticalBit(bestRefNum, key)

	return identical, bestRefNum, parentNodeNum, parentDirection
}
