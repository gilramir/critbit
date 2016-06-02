package critbit

import (
	"fmt"
	"sync"
)

// Split splits a tree into two trees, each having one half of the key-value pairs.
// If there is an odd number of keys, the right tree (the second returned tree)
// will have the extra key-value pair.
func (tree *Critbit) Split() (*Critbit, *Critbit) {
	// Empty, or other trivial cases?
	switch tree.numExternalRefs {
	case 0:
		return tree, New(0)
	case 1:
		return tree, New(0)
	case 2:
		return tree.splitTwoExternalRefs()
	}

	leftNumKeys := tree.numExternalRefs / 2
	return tree.SplitAt(leftNumKeys)
}

func (tree *Critbit) splitTwoExternalRefs() (*Critbit, *Critbit) {
	rootNode := &tree.internalNodes[tree.rootItem]

	left := New(1)
	leftRef := &tree.externalRefs[rootNode.child[0]]
	leftRefNum, err := tree.addExternalRef(leftRef.key, leftRef.value)
	// An error should not happen because of the size of the tree
	if err != nil {
		panic(err.Error())
	}
	left.rootItem = leftRefNum

	right := New(1)
	rightRef := &tree.externalRefs[rootNode.child[1]]
	rightRefNum, err := tree.addExternalRef(rightRef.key, rightRef.value)
	// An error should not happen because of the size of the tree
	if err != nil {
		panic(err.Error())
	}
	right.rootItem = rightRefNum
	return left, right
}

// Split splits a tree into two arbitrarily sized trees. The leftNumKeys
// arguments indicates how many treees the left tree (the first returned tree)
// should have. The right tree (the second returned tree) will have the rest.
func (tree *Critbit) SplitAt(leftNumKeys int) (*Critbit, *Critbit) {
	leftItemChan := make(chan *splitItem)
	rightItemChan := make(chan *splitItem)

	rightNumKeys := tree.numExternalRefs - leftNumKeys
	if rightNumKeys < 0 {
		rightNumKeys = 0
	}
	leftTree := New(leftNumKeys)
	rightTree := New(rightNumKeys)

	go tree.splitWalkTree(leftNumKeys, leftItemChan, rightItemChan)

	var wg sync.WaitGroup
	wg.Add(2)
	go createLeftSplit(&wg, leftTree, leftItemChan)
	go createRightSplit(&wg, rightTree, rightItemChan)
	wg.Wait()

	return leftTree, rightTree
}

func (tree *Critbit) splitWalkTree(leftNumKeys int,
	leftItemChan chan *splitItem, rightItemChan chan *splitItem) {

	defer close(leftItemChan)
	defer close(rightItemChan)

	state := &splitWalkerState{
		// It's impossible to approximate the longest path in the tree,
		// but we can use the # of external refs as a pseuco max
		path:                 make([]*splitItem, 0, tree.numExternalRefs),
		numLeftKeysRemaining: leftNumKeys,
		channels:             make([]chan *splitItem, 2),
		feedingRight:         leftNumKeys == 0,
	}

	state.channels[0] = leftItemChan
	state.channels[1] = rightItemChan

	tree.splitWalkTreeRecurse(state)
}

type splitWalkerState struct {
	visitedRoot          bool
	path                 []*splitItem
	numLeftKeysRemaining int
	channels             []chan *splitItem
	feedingRight         bool
	channel              chan *splitItem
}

func (tree *Critbit) splitWalkTreeRecurse(state *splitWalkerState) {
	sendPopup := true
	if state.feedingRight {
		state.channel = state.channels[1]
	} else {
		state.channel = state.channels[0]
	}

	// Just started?
	if !state.visitedRoot {
		state.path = append(state.path, &splitItem{
			metaType: kSplitItemTreeData,
			itemType: kChildIntNode,
			itemID:   tree.rootItem,
			offset:   tree.internalNodes[tree.rootItem].offset,
			bit:      tree.internalNodes[tree.rootItem].bit,
		})
		state.visitedRoot = true
	}

	item := state.path[len(state.path)-1]
	state.channel <- item

	switch item.itemType {
	case kChildIntNode:
		state.path = append(state.path, tree.createSplitItemFromNodeChild(item.itemID, 0))
		tree.splitWalkTreeRecurse(state)
		state.path = state.path[:len(state.path)-1]

		state.path = append(state.path, tree.createSplitItemFromNodeChild(item.itemID, 1))
		tree.splitWalkTreeRecurse(state)
		state.path = state.path[:len(state.path)-1]

	case kChildExtRef:
		if !state.feedingRight {
			state.numLeftKeysRemaining--
			if state.numLeftKeysRemaining == 0 {
				state.feedingRight = true
				state.channel = state.channels[1]
				// Need to feed 'path' up to, but not including, the ext ref to the right tree,
				// so popups make sense.
				// But we need to make copies of each splitItem, as the leftTree will need
				// its own copy

				for _, pathItem := range state.path[0 : len(state.path)-1] {
					clonedItem := pathItem.Clone()
					state.channel <- clonedItem
				}
				sendPopup = false
			}
		}
	}
	// The left reader gets popups only until the # keys hasn't been reached;
	if sendPopup && (state.feedingRight || state.numLeftKeysRemaining > 0) {
		state.channel <- &popupItem
	}
}

func (tree *Critbit) createSplitItemFromNodeChild(nodeNum uint32, childDirection byte) *splitItem {
	node := &tree.internalNodes[nodeNum]
	itemType := node.getChildType(childDirection)
	switch itemType {
	case kChildIntNode:
		return &splitItem{
			metaType:  kSplitItemTreeData,
			itemType:  kChildIntNode,
			itemID:    node.child[childDirection],
			direction: childDirection,
			offset:    node.offset,
			bit:       node.bit,
		}
	case kChildExtRef:
		itemID := node.child[childDirection]
		key := tree.externalRefs[itemID].key
		value := tree.externalRefs[itemID].value
		return &splitItem{
			metaType:  kSplitItemTreeData,
			itemType:  kChildExtRef,
			itemID:    itemID,
			direction: childDirection,
			key:       key,
			value:     value,
		}
	default:
		panic(fmt.Sprintf("Node %d has unexpected child type %d in direction %d",
			nodeNum, node.getChildType(childDirection), childDirection))
	}
}

const (
	kSplitItemTreeData = 1
	kSplitItemPopUp    = 2
)

type splitItem struct {
	metaType  int
	newTreeID uint32

	itemType  uint8
	itemID    uint32
	direction byte

	// If internalNode
	offset uint16
	bit    uint8
	// If externalRef
	key   string
	value interface{}
}

// Clones, but clears newTreeID
func (item *splitItem) Clone() *splitItem {
	return &splitItem{
		metaType:  item.metaType,
		itemType:  item.itemType,
		itemID:    item.itemID,
		direction: item.direction,
		offset:    item.offset,
		bit:       item.bit,
		key:       item.key,
		value:     item.value,
	}
}

var popupItem splitItem = splitItem{metaType: kSplitItemPopUp}

func createLeftSplit(wg *sync.WaitGroup, tree *Critbit, itemChan chan *splitItem) {
	defer wg.Done()

	// Populate the tree
	_ = tree.populateFromSplitChannel("left", itemChan)

	// Elide the root and sides
	tree.postSplitElideRootIfNeeded(1)
	tree.postSplitZipSide(1)
}

func createRightSplit(wg *sync.WaitGroup, tree *Critbit, itemChan chan *splitItem) {
	defer wg.Done()

	// Populate the tree
	_ = tree.populateFromSplitChannel("right", itemChan)

	// Elide the root and sides
	tree.postSplitElideRootIfNeeded(0)
	tree.postSplitZipSide(0)
}

func (tree *Critbit) populateFromSplitChannel(side string, itemChan chan *splitItem) []*splitItem {
	var path []*splitItem
	for item := range itemChan {
		if item.metaType == kSplitItemPopUp {
			path = path[:len(path)-1]
			continue
		}

		switch item.itemType {
		case kChildIntNode:
			nodeNum, node := tree.addInternalNode()
			item.newTreeID = nodeNum
			node.offset = item.offset
			node.bit = item.bit
		case kChildExtRef:
			refNum, err := tree.addExternalRef(item.key, item.value)
			// An error should not happen because of the size of the tree
			if err != nil {
				panic(err.Error())
			}
			item.newTreeID = refNum
		}

		if len(path) == 0 {
			if item.itemType != kChildIntNode {
				panic(fmt.Sprintf("(%s) First node has type %d", side, item.itemType))
			}
			tree.rootItem = item.newTreeID
		} else {
			parentNode := &tree.internalNodes[path[len(path)-1].newTreeID]
			parentNode.setChild(item.direction, item.newTreeID, item.itemType)
		}
		path = append(path, item)
	}
	return path
}

func (tree *Critbit) postSplitElideRootIfNeeded(direction byte) {
	// Walk down from the root, eliding the root as necessary.
	if tree.rootItemType() == kChildIntNode {
		rootNode := &tree.internalNodes[tree.rootItem]
		for rootNode.getChildType(direction) == kChildNil {
			prevRootItem := tree.rootItem
			tree.rootItem = tree.internalNodes[prevRootItem].child[1-direction]
			tree.deleteInternalNode(prevRootItem)
			rootNode = &tree.internalNodes[tree.rootItem]
		}
	}
}

func (tree *Critbit) postSplitZipSide(direction byte) {
	// Now that we know we don't need to elide the root, walk down
	// from the root, towards the right, eliding as necessary
	if tree.rootItemType() != kChildIntNode {
		return
	}
	var prevNodeNum uint32 = kNilNode
	var prevNode *internalNode
	var nodeNum uint32 = tree.rootItem
	var node *internalNode
	pathNodeNums := make([]uint32, 0)

	for keepGoing := true; keepGoing; {
		node = &tree.internalNodes[nodeNum]
		// Do something based on the right child type
		switch node.getChildType(direction) {
		case kChildExtRef:
			keepGoing = false
			break
		case kChildIntNode:
			pathNodeNums = append(pathNodeNums, nodeNum)
			prevNodeNum = nodeNum
			prevNode = node
			nodeNum = node.child[direction]
			continue
		case kChildNil:
			lchildType := node.getChildType(1 - direction)
			if lchildType == kChildNil {
				// The node is elided, but really the prevNode needs to be elided as this node
				// is _completely unnecessary
				tree.deleteInternalNode(nodeNum)
				prevNode.setChild(direction, nodeNum, kChildNil)
				// Go up
				nodeNum = prevNodeNum
				node = prevNode
				if len(pathNodeNums) > 1 {
					prevNodeNum = pathNodeNums[len(pathNodeNums)-2]
					prevNode = &tree.internalNodes[prevNodeNum]
					continue
				} else {
					tree.rootItem = prevNode.child[1-direction]
					tree.deleteInternalNode(tree.rootItem)
					keepGoing = false
					break
				}
			}
			lchildID := node.child[1-direction]
			prevNode.setChild(direction, lchildID, lchildType)
			tree.deleteInternalNode(nodeNum)
			// Do something based on the left child type
			switch lchildType {
			case kChildExtRef:
				keepGoing = false
				break
			case kChildIntNode:
				// prevNode stays, but node changes because we just elided ourselves!
				nodeNum = lchildID
				continue
			}
		}
	}
}
