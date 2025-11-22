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

// Upsert inserts a key/value pair into the tree if the key
// doesn't exit, or, updates the value if the key does exist already
// Since the act of inserting a key may return an error, Upsert also
// returns an error, indicating if inseration failed.
func (tree *Critbit[T]) Upsert(key string, value T) error {
	has, refNum := tree.findRef(key)
	if has {
		tree.externalRefs[refNum].value = value
		return nil
	} else {
		inserted, err := tree.Insert(key, value)
		if !inserted {
			panic("Insert should have succeeded")
		}
		return err
	}
}
