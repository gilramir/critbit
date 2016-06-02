package critbit

func (node *internalNode) setChild(direction byte, id uint32, childType uint8) {
	node.child[direction] = id
	node.setChildType(direction, childType)
}

func (node *internalNode) setChildType(direction byte, childType uint8) {
	if direction == 0 {
		childType <<= 4
		childType &^= kLeftMask
		node.flags &= kLeftMask
	} else {
		childType &^= kRightMask
		node.flags &= kRightMask
	}
	node.flags |= childType
}

func (node *internalNode) getChildType(direction byte) uint8 {
	var flags uint8
	if direction == 0 {
		flags = node.flags >> 4
	} else {
		flags = node.flags
	}
	return flags &^ kRightMask
}

// Returns the direction for the given key
func (node *internalNode) direction(key string) byte {
	if node.offset < uint16(len(key)) && key[node.offset]&node.bit != 0 {
		return 1
	}
	return 0
}
