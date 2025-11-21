package critbit

import (
	"fmt"
	"math"
	"strings"
)

// A KeyValueTuple is returned during an iteration.
type KeyValueTuple struct {
	Key   string
	Value interface{}
}

// maybe this should just be externalRef?

type walkerItem struct {
	itemType uint8
	itemID   uint32
}

type walkerStack struct {
	array []*walkerItem
	size  int
	top   int // where the next entry will be written to

	largestTop int
}

func (s *walkerStack) Len() int {
	return s.top
}

func (s *walkerStack) push(walker *walkerItem) {
	if s.top == s.size {
		s.array = append(s.array, make([]*walkerItem, s.size/2)...)
		s.size += s.size / 2
	}
	s.array[s.top] = walker
	s.top++
	if s.top > s.largestTop {
		s.largestTop = s.top
	}
}

func (s *walkerStack) pop() *walkerItem {
	if s.top > 0 {
		walker := s.array[s.top-1]
		s.top--
		return walker
	}
	panic("pop() of empty stack")
}

func (tree *Critbit) newWalkerStack() *walkerStack {
	// The maximum size of the stack is the number of nodes we still
	// need to visit, which is the height of the tree. The max height
	// of the tree is bounded by the length of the keys, which goes up
	// to 64k characters. However, we are also bounded by the current number of
	// refs in the tree. We can't understand the topology, but we can make a good
	// guess.  So, we preallocate a good amount which should cover most trees,
	// and allow for dynamic growth for the corner cases.
	// A good estimate seems to be log2(refs) * 1.5
	stackSize := int(math.Log2(float64(tree.numExternalRefs)) * 1.5)
	if stackSize < 3 {
		stackSize = 3
	}

	return &walkerStack{
		array:      make([]*walkerItem, stackSize),
		size:       stackSize,
		top:        0,
		largestTop: 0,
	}
}

func (tree *Critbit) createWalkerItemFromNodeNum(nodeNum uint32) *walkerItem {
	return &walkerItem{
		itemType: kChildIntNode,
		itemID:   nodeNum,
	}
}

func (tree *Critbit) createWalkerItemFromRefNum(refNum uint32) *walkerItem {
	return &walkerItem{
		itemType: kChildExtRef,
		itemID:   refNum,
	}
}

// Keys returns a string slice containing all the keys in the tree.
// The keys are in sorted order.
func (tree *Critbit) Keys() []string {
	// Get the keys
	var keys []string
	tupleChan := tree.GetKeyValueTuples()
	for keyTuple := range tupleChan {
		keys = append(keys, keyTuple.Key)
	}
	return keys
}

// GetKeyValueTuples returns a channel that can be read from which contains
// each key-value pair, in sorted order by the keys.
func (tree *Critbit) GetKeyValueTuples() chan *KeyValueTuple {
	tupleChan := make(chan *KeyValueTuple)

	// Push the first item in the stack
	stack := tree.newWalkerStack()
	stack.push(tree.createWalkerItemFromNodeNum(tree.rootItem))

	go tree._iterateKeyTuples(stack, tupleChan)
	return tupleChan
}

// GetKeyValueTuples return a slice of KeyValueTuples starting from the node
// corresponding to the key passed in, and then walking the Critbit tree after
// that. Up at maxCount items, including the first key/value itself, are
// returned. If maxCount is 0, all items are returned.
func (tree *Critbit) GetKeyValueTuplesFrom(key string, startExact bool, maxCount int) []*KeyValueTuple {
	results := make([]*KeyValueTuple, 0)

	has, refNum, parentNodeNum, parentDirection := tree.findRefWithAncestry(key)
	if !has {
		if startExact {
			// User asked for an exact match. It's not. Return now
			return results
		} else {
			// Not an exact match, but, did we find something that does start
			// with our string?
			foundRefKey := tree.externalRefs[refNum].key
			if !strings.HasPrefix(foundRefKey, key) {
				// No, the best ref does not start with the user's key
				return results
			}
			// keep going!
		}
	}
	tupleChan := make(chan *KeyValueTuple)

	// Push the first item in the stack
	stack := tree.newWalkerStack()

	// If the ref is on the left side, add the right side sibling too, in
	// case need to walk down that side too
	if parentDirection == kDirectionLeft {
		parentNode := tree.internalNodes[parentNodeNum]
		siblingNodeNum := parentNode.child[kDirectionRight]
		siblingNodeType := parentNode.getChildType(kDirectionRight)
		if siblingNodeType == kChildIntNode {
			stack.push(tree.createWalkerItemFromNodeNum(siblingNodeNum))
			//			fmt.Printf("added nodenum %d => %+v\n", siblingNodeNum, stack.array[0])
		} else {
			stack.push(tree.createWalkerItemFromRefNum(siblingNodeNum))
			//			fmt.Printf("added refnum %d => %+v\n", siblingNodeNum, stack.array[0])
		}
	}

	wi := tree.createWalkerItemFromRefNum(refNum)
	stack.push(wi)
	//	fmt.Printf("created walker from refnum %d => %+v\n", refNum, wi)

	go tree._iterateKeyTuplesWithStack(stack, tupleChan, maxCount)

	for kvt := range tupleChan {
		results = append(results, kvt)
	}

	return results
}

func (tree *Critbit) _iterateKeyTuples(stack *walkerStack, tupleChan chan *KeyValueTuple) {
	switch tree.rootItemType() {
	case kChildNil:
		// Empty tree?
		close(tupleChan)
		return

	case kChildExtRef:
		// One ref?
		tree.sendKeyTuple(tree.rootItem, tupleChan)
		close(tupleChan)
		return
	}
	tree._iterateKeyTuplesWithStack(stack, tupleChan, 0)
}

func (tree *Critbit) _iterateKeyTuplesWithStack(stack *walkerStack, tupleChan chan *KeyValueTuple, maxCount int) {
	defer close(tupleChan)

	numSent := 0
	// Walk the tree
	for stack.Len() > 0 {
		/*
			fmt.Printf("Stack len: %d\n", stack.Len())
			for i, si := range stack.array {
				fmt.Printf("\t%d: %+v\n", i, si)
			}
		*/

		// Pop
		walker := stack.pop()
		//		fmt.Printf("Popped walker: %+v  isExtRef? %v newLen=%d\n",
		//			walker, walker.itemType == kChildExtRef, stack.Len())

		// leaf?
		if walker.itemType == kChildExtRef {
			//			fmt.Printf("is leaf\n")
			tree.sendKeyTuple(walker.itemID, tupleChan)

			// Did we reach our max?
			numSent++
			if maxCount > 0 && numSent == maxCount {
				return
			}
		} else {
			//			fmt.Printf("has children\n")
			// Push each child
			node := &tree.internalNodes[walker.itemID]
			// Right side pushed first
			switch node.getChildType(kDirectionRight) {
			case kChildIntNode:
				stack.push(tree.createWalkerItemFromNodeNum(node.child[kDirectionRight]))
			case kChildExtRef:
				stack.push(tree.createWalkerItemFromRefNum(node.child[kDirectionRight]))
			default:
				panic(fmt.Sprintf("Node %d has child[1] type = %d", walker.itemID,
					node.getChildType(1)))
			}
			// Then left side
			switch node.getChildType(kDirectionLeft) {
			case kChildIntNode:
				stack.push(tree.createWalkerItemFromNodeNum(node.child[kDirectionLeft]))
			case kChildExtRef:
				stack.push(tree.createWalkerItemFromRefNum(node.child[kDirectionLeft]))
			default:
				panic(fmt.Sprintf("Node %d has child[0] type = %d", walker.itemID,
					node.getChildType(0)))
			}
		}
	}
	// fmt.Printf("finished walking tree\n")
}

func (tree *Critbit) sendKeyTuple(refNum uint32, tupleChan chan *KeyValueTuple) {
	ref := &tree.externalRefs[refNum]
	//	fmt.Printf("sending key=%s value=%v\n", ref.key, ref.value)
	tupleChan <- &KeyValueTuple{ref.key, ref.value}
}
