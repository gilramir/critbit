package critbit

import (
	"fmt"
	"math"
)

// A KeyValueTuple is returned during an iteration.
type KeyValueTuple[T any] struct {
	Key   string
	Value T
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

func (tree *Critbit[T]) newWalkerStack() *walkerStack {
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

func (tree *Critbit[T]) createWalkerItemFromNodeNum(nodeNum uint32) *walkerItem {
	return &walkerItem{
		itemType: kChildIntNode,
		itemID:   nodeNum,
	}
}

func (tree *Critbit[T]) createWalkerItemFromRefNum(refNum uint32) *walkerItem {
	return &walkerItem{
		itemType: kChildExtRef,
		itemID:   refNum,
	}
}

// Keys returns a string slice containing all the keys in the tree.
// The keys are in sorted order.
func (tree *Critbit[T]) Keys() []string {
	// Get the keys
	var keys []string
	tupleChan := tree.GetKeyValueTupleChan()
	for keyTuple := range tupleChan {
		keys = append(keys, keyTuple.Key)
	}
	return keys
}

// GetKeyValueTuplesCHan returns a channel that can be read from which contains
// each key-value pair, in sorted order by the keys.
func (tree *Critbit[T]) GetKeyValueTupleChan() chan *KeyValueTuple[T] {
	tupleChan := make(chan *KeyValueTuple[T])

	go tree._iterateKeyTuples(tupleChan)
	return tupleChan
}

// Returns all the KeyValueTuples in key-sorted order.
func (tree *Critbit[T]) GetKeyValueTuples() []*KeyValueTuple[T] {
	kvts := make([]*KeyValueTuple[T], tree.Length())

	tupleChan := make(chan *KeyValueTuple[T])
	go tree._iterateKeyTuples(tupleChan)

	i := 0
	for kvt := range tupleChan {
		kvts[i] = kvt
		i++
	}
	if i != tree.Length() {
		panic(fmt.Sprintf("After GetKeyValueTuples, i=%d not %d",
			i, tree.Length()))
	}

	return kvts
}

func (tree *Critbit[T]) _iterateKeyTuples(tupleChan chan *KeyValueTuple[T]) {
	defer close(tupleChan)
	switch tree.rootItemType() {
	case kChildNil:
		// Empty tree?
		return

	case kChildExtRef:
		// One ref?
		tree.sendKeyTuple(tree.rootItem, tupleChan)
		return
	}

	// Push the first item in the stack
	stack := tree.newWalkerStack()
	stack.push(tree.createWalkerItemFromNodeNum(tree.rootItem))

	// Walk the tree
	for stack.Len() > 0 {
		// Pop
		walker := stack.pop()

		// leaf?
		if walker.itemType == kChildExtRef {
			tree.sendKeyTuple(walker.itemID, tupleChan)

		} else {
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
}

func (tree *Critbit[T]) sendKeyTuple(refNum uint32, tupleChan chan *KeyValueTuple[T]) {
	ref := &tree.externalRefs[refNum]
	tupleChan <- &KeyValueTuple[T]{
		Key:   ref.key,
		Value: ref.value,
	}
}
