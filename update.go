package critbit

// Update changes the value for the given key. If the key is
// not stored in the tree, the returned bool value is false.
func (tree *Critbit[T]) Update(key string, value T) bool {
	has, refNum := tree.findRef(key)
	if !has {
		return false
	}
	tree.externalRefs[refNum].value = value
	return true
}
