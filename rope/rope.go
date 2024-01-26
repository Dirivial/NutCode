package rope

import "fmt"

type Rope struct {
	Head *Node
}

type Node struct {
	Left    *Node
	Right   *Node
	Content string
	Weight  int
}

func New(s string) *Rope {
	r := &Rope{Head: createRope(s)}
	return r
}

// createRope recursively creates a rope from a string.
func createRope(s string) *Node {
	if len(s) <= 5 { // You can adjust this threshold based on your needs
		return &Node{Content: s, Weight: len(s)}
	}

	mid := (len(s) - 1) / 2
	leftSubString := s[:mid]
	rightSubString := s[mid:]

	return &Node{
		Left:   createRope(leftSubString),
		Right:  createRope(rightSubString),
		Weight: mid,
	}
}

// Insert a string into the rope structure
func (*Rope) Insert(index int, str string) {
}

// Delete part of the rope structure
func (*Rope) Delete(start, length int) {
}

// Get character at a position.
// NOTE: 1 indexed
func (r *Rope) Index(index int) string {
	return Index(r.Head, index)
}

// index recursively searches for the character at the specified index.
func Index(node *Node, index int) string {
	if node == nil {
		return ""
	}

	if index > node.Weight && node.Right != nil {
		return Index(node.Right, index-node.Weight)
	}
	if node.Left != nil {
		return Index(node.Left, index)
	}
	return string(node.Content[index-1])
}

// Collect all leaves of the rope structure
func (*Rope) CollectLeaves() []Node {
	return []Node{}
}

// Build a string from the entire rope structure
func (*Rope) Report(start, length int) string {
	return ""
}

// Rebalance the rope structure
func (*Rope) Rebalance() Rope {
	return Rope{}
}

// Concatenate a rope with another
func (r *Rope) Concat(rope *Rope) *Rope {
	return &Rope{Head: Concatenate(r.Head, rope.Head)}
}

// concatenate combines two rope nodes into a new rope node.
func Concatenate(node1, node2 *Node) *Node {
	if node1 == nil {
		return node2
	}
	if node2 == nil {
		return node1
	}

	return &Node{
		Left:   node1,
		Right:  node2,
		Weight: node1.ComputeTotalWeight(),
	}
}

// Compute the total weight of a node, useful in concat
func (n *Node) ComputeTotalWeight() int {
	if n.Left == nil && n.Right == nil {
		return n.Weight
	}
	weight := 0
	if n.Left != nil {
		weight += n.Left.ComputeTotalWeight()
	}
	if n.Right != nil {
		weight += n.Right.ComputeTotalWeight()
	}
	return weight
}

// Split a rope into two
func (r *Rope) Split(index int) *Rope {
	removedNodes := Split(r.Head, index)
	if len(removedNodes) == 0 {
		return &Rope{}
	}

	rope := &Rope{Head: &Node{Left: removedNodes[0], Weight: removedNodes[0].ComputeTotalWeight()}}
	for i := 1; i < len(removedNodes); i++ {
		if removedNodes[i] != nil {
			toConcat := &Rope{removedNodes[i]}
			rope = rope.Concat(toConcat)
		} else {
			fmt.Printf("ERROR: Nil node saved...")
		}
	}
	// Recompute weights
	if r.Head.Left != nil {
		r.Head.Weight = r.Head.Left.ComputeTotalWeight()
	}
	if rope.Head.Left != nil {
		rope.Head.Weight = rope.Head.Left.ComputeTotalWeight()
	}
	return rope
}

func Split(node *Node, index int) []*Node {
	if node == nil {
		return []*Node{}
	}

	if index > node.Weight && node.Right != nil {
		return Split(node.Right, index-node.Weight)
	}
	if node.Left != nil {
		n := node.Right
		node.Right = nil
		return append(Split(node.Left, index), n)
	}
	// Check if the split should occurr somewhere within the content
	if index > 1 && index < node.Weight {
		// Create a new node and fill it with content
		movedContent := node.Content[index:]
		newNode := &Node{Content: movedContent, Weight: len(movedContent)}

		// Remove the moved content from this node
		node.Content = node.Content[:index]
		node.Weight = len(node.Content)

		// Return the newly created node
		return []*Node{newNode}
	} else if index == 1 {
		// Return this node
		return []*Node{node}
	}
	// Return empty node, as the split occurrs after this node
	return []*Node{}
}

func (r *Rope) printRope() {
	n := r.Head
	n.printNode(0)
	n.printChildren(1)
}

func (n *Node) printChildren(depth int) {
	if n.Left != nil {
		n.Left.printNode(depth)
		n.Left.printChildren(depth + 1)
	}

	if n.Right != nil {
		n.Right.printNode(depth)
		n.Right.printChildren(depth + 1)
	}
}

func (n *Node) printNode(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("Weight: %d. Left: %t, Right: %t, Content: '%s'\n", n.Weight, n.Left != nil, n.Right != nil, n.Content)
}
