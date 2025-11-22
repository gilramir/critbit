package critbit

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Dump prints the structure of the entire tree to stdout.
func (tree *Critbit[T]) Dump() {
	fmt.Printf("Tree length=%d\n", tree.numExternalRefs)
	if tree.numExternalRefs == 0 {
		return
	} else if tree.numExternalRefs == 1 {
		tree.dumpExternalRef("Root:", tree.rootItem, "")
		return
	}
	tree.dumpInternalNode("Root:", tree.rootItem, "")
}

func (tree *Critbit[T]) dumpExternalRef(title string, refNum uint32, indent string) {
	// One ref, and it's the root and leaf (no internal nodes)
	fmt.Printf("%s%s refNum=%d (EXT) key=%s\n", indent,
		title, refNum, tree.externalRefs[refNum].key)
}

func (tree *Critbit[T]) dumpInternalNode(title string, nodeNum uint32, indent string) {
	node := &tree.internalNodes[nodeNum]
	fmt.Printf("%s%s nodeNum=%d (INT) off=%d bit=0x%01x\n", indent,
		title, nodeNum, node.offset, node.bit)

	indent += "  "
	switch node.getChildType(0) {
	case kChildNil:
		fmt.Printf("%sLeft  type is nil, value=%d\n", indent, node.child[0])
	case kChildIntNode:
		tree.dumpInternalNode("Left ", node.child[0], indent)
	case kChildExtRef:
		tree.dumpExternalRef("Left ", node.child[0], indent)
	default:
		fmt.Printf("%sUnexpected left childType=%d value=%d\n",
			indent, node.getChildType(0), node.child[0])
	}
	switch node.getChildType(1) {
	case kChildNil:
		fmt.Printf("%sRight type is nil, value=%d\n", indent, node.child[1])
	case kChildIntNode:
		tree.dumpInternalNode("Right", node.child[1], indent)
	case kChildExtRef:
		tree.dumpExternalRef("Right", node.child[1], indent)
	default:
		fmt.Printf("%sUnexpected right childType=%d value=%d\n",
			indent, node.getChildType(1), node.child[1])
	}
}

// SaveDot save the tree structure to a graphviz/dot file with the
// given name. You can run 'dot' on the file to see the graphical
// representation of the tree.
func (tree *Critbit[T]) SaveDot(filename string) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "Opening %s for writing", filename)
	}

	// Ignore errors
	defer outputFile.Close()

	_, err = outputFile.WriteString("digraph xbtrie_critbit {\n")
	if err != nil {
		return err
	}

	if tree.numExternalRefs == 0 {
		// no-op
	} else if tree.numExternalRefs == 1 {
		tree.saveDotExternalRef(outputFile, tree.rootItem)
	} else {
		tree.saveDotInternalNode(outputFile, tree.rootItem)
	}

	_, err = outputFile.WriteString("}\n")
	if err != nil {
		return err
	}
	return nil
}

func (tree *Critbit[T]) saveDotExternalRef(outputFile *os.File, refNum uint32) error {
	var err error
	name := fmt.Sprintf("ref_%d", refNum)

	_, err = outputFile.WriteString(
		fmt.Sprintf("\t%s [label=\"%s\\nrefNum=%d\" shape=\"box\"]\n",
			name, tree.externalRefs[refNum].key, refNum))
	if err != nil {
		return err
	}

	return nil
}

func (tree *Critbit[T]) saveDotInternalNode(outputFile *os.File, nodeNum uint32) error {
	var err error
	name := fmt.Sprintf("node_%d", nodeNum)
	node := &tree.internalNodes[nodeNum]

	// Internal node
	_, err = outputFile.WriteString(
		fmt.Sprintf("\t%s [label=\"off=0x%02x\\nbit=%d\\nnodeNum=%d\"]\n",
			name, node.offset, node.bit, nodeNum))
	if err != nil {
		return err
	}

	var leftName string
	var rightName string

	switch node.getChildType(0) {
	case kChildNil:
		// no-op
	case kChildIntNode:
		leftName = fmt.Sprintf("node_%d", node.child[0])
		err = tree.saveDotInternalNode(outputFile, node.child[0])
	case kChildExtRef:
		leftName = fmt.Sprintf("ref_%d", node.child[0])
		err = tree.saveDotExternalRef(outputFile, node.child[0])
	default:
		return errors.Errorf("Unexpected left childType=%d value=%d",
			node.getChildType(0), node.child[0])
	}
	if err != nil {
		return err
	}
	switch node.getChildType(1) {
	case kChildNil:
		// no-op
	case kChildIntNode:
		rightName = fmt.Sprintf("node_%d", node.child[1])
		err = tree.saveDotInternalNode(outputFile, node.child[1])
	case kChildExtRef:
		rightName = fmt.Sprintf("ref_%d", node.child[1])
		err = tree.saveDotExternalRef(outputFile, node.child[1])
	default:
		return errors.Errorf("Unexpected right childType=%d value=%d",
			node.getChildType(1), node.child[1])
	}
	if err != nil {
		return err
	}

	if leftName == "" {
		_, err = outputFile.WriteString(fmt.Sprintf("\t%s -> %s_lnil\n", name, name))
		if err != nil {
			return err
		}
	} else {
		_, err = outputFile.WriteString(fmt.Sprintf("\t%s -> %s\n", name, leftName))
		if err != nil {
			return err
		}
	}
	if rightName == "" {
		_, err = outputFile.WriteString(fmt.Sprintf("\t%s -> %s_rnil\n", name, name))
		if err != nil {
			return err
		}
	} else {
		_, err = outputFile.WriteString(fmt.Sprintf("\t%s -> %s\n", name, rightName))
		if err != nil {
			return err
		}
	}
	return nil
}
