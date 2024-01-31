package rope

import (
	"fmt"
	"strings"
)

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
func (r *Rope) Insert(index int, str string) *Rope {
	ropeEnd := r.Split(index)
	ropeEnd.Rebalance()
	ropeMiddle := New(str)
	fullRope := r.Concat(ropeMiddle)
	fullRope.Rebalance()
	fullRope = fullRope.Concat(ropeEnd)
	fullRope.Rebalance()
	return fullRope
}

// Delete part of the rope structure
func (r *Rope) Delete(start, length int) *Rope {
	// Ignore unusable calls
	if length <= 0 || start < 0 {
		return r
	}
	intermediate := r.Split(start)
	intermediate.Rebalance()
	right := intermediate.Split(length)
	right.Rebalance()
	ret := r.Concat(right)
	ret.Rebalance()
	return ret
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
	if index > node.Weight {
		return ""
	}
	return string(node.Content[index-1])
}

// Collect all leaves of the rope structure
func (r *Rope) CollectLeaves() []Node {
	return []Node{}
}

// Build a (sub)string from the entire rope structure
func (r *Rope) Report(start, length int) string {
	content := Report(r.Head, start, start+length)
	return content
}

func Report(n *Node, start, end int) string {
	if n == nil {
		return ""
	}

	content := ""
	if (start > n.Weight || end > n.Weight) && n.Right != nil {
		content += Report(n.Right, max(start-n.Weight, 1), end-n.Weight)
	}
	if start < n.Weight && n.Left != nil {
		content = Report(n.Left, start, end) + content
	}

	if n.Left == nil && n.Right == nil {
		if n.Weight < start {
			return ""
		} else if n.Weight >= end {
			return n.Content[start-1 : end-1]
		}
		return n.Content[start-1:]
	}
	return content
}

// Rebalance the rope structure
func (r *Rope) Rebalance() *Rope {
	// I will begin by just reducing the amount of "linked lists" that appear
	r.Head.Left.Rebalance(r.Head)
	r.Head.Right.Rebalance(r.Head)
	return r
}

func (n *Node) Rebalance(parent *Node) {
	if n == nil {
		return
	}
	if n.Left == nil && n.Right != nil {
		if parent.Left == n {
			parent.Left = n.Right
		} else {
			parent.Right = n.Right
		}
		n.Right.Rebalance(parent)
	} else if n.Left != nil && n.Right == nil {
		if parent.Left == n {
			parent.Left = n.Left
		} else {
			parent.Right = n.Left
		}
		n.Left.Rebalance(parent)
	}
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
	// TODO: Think if this is really necessary.
	// Shouldn't I just add n.Weight with n.Right.ComputeTotalWeight()?
	// (ofc, assuming n.Right != nil)
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
	// This should just move the entire rope
	if index == 0 {
		movedRope := &Rope{Head: r.Head}
		r.Head = &Node{Weight: 0, Content: ""}

		return movedRope
	}
	// Split the rope, resulting in a list of nodes
	removedNodes := Split(r.Head, index)
	if len(removedNodes) == 0 {
		return &Rope{}
	}

	// Create a new rope with the removed nodes
	rope := &Rope{Head: &Node{Left: removedNodes[0], Weight: removedNodes[0].ComputeTotalWeight()}}
	for i := 1; i < len(removedNodes); i++ {
		if removedNodes[i] != nil {
			toConcat := &Rope{removedNodes[i]}
			rope = rope.Concat(toConcat)
		} else {
			// fmt.Println("ERROR: Nil node saved...")
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
		// fmt.Println("Is this also a problem?")
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
	if index >= 1 && index < node.Weight {
		// Create a new node and fill it with content
		movedContent := node.Content[index:]
		newNode := &Node{Content: movedContent, Weight: len(movedContent)}

		// Remove the moved content from this node
		node.Content = node.Content[:index]
		node.Weight = len(node.Content)

		// Return the newly created node
		return []*Node{newNode}
	} else if index == 0 {
		// Return this node
		return []*Node{node}
	}
	// fmt.Println("Bruhmode")
	// Return empty node, as the split occurrs after this node
	return []*Node{{Content: "", Weight: 0}}
}

func (r *Rope) printRope() {
	n := r.Head
	n.printChildren(1)
}

func (n *Node) printChildren(depth int) {
	if n.Left != nil {
		n.Left.printChildren(depth + 1)
	}

	n.printNode(depth)
	if n.Right != nil {
		n.Right.printChildren(depth + 1)
	}
}

func (n *Node) printNode(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("\t")
	}
	if len(n.Content) != 0 {
		fmt.Printf("Weight: %d. Content: %s\n", n.Weight, n.Content)
	} else {
		fmt.Printf("Weight: %d\n", n.Weight)
	}
}

func (r *Rope) GetContent() string {
	return r.Head.GetContent()
}

func (n *Node) GetContent() string {
	content := ""
	if n.Left != nil {
		content = content + n.Left.GetContent()
	}

	if n.Right != nil {
		content = content + n.Right.GetContent()
	}

	if n.Right == nil && n.Left == nil {
		return n.Content
	}
	return content
}

func (r *Rope) SearchChar(char rune, startFrom int) int {
	if startFrom < 0 {
		return -1
	}

	return r.Head.SearchChar(char, startFrom, 1)
}

func (n *Node) SearchChar(char rune, index, globalIndex int) int {
	if n == nil {
		// fmt.Println("Cringe")
		return -1
	}

	if index > n.Weight && n.Right != nil {
		ret := n.Right.SearchChar(char, index-n.Weight, 0)
		if ret == -1 {
			return -1
		}
		return ret + n.Weight
	}
	if n.Left != nil {
		ret := n.Left.SearchChar(char, index, globalIndex)
		// If we could not find it to the left, check right path
		if ret == -1 {
			ret = n.Right.SearchChar(char, 1, n.Weight)
			if ret == -1 {
				return -1
			}
			return ret + n.Weight
		}
		return ret
	}
	if index >= n.Weight {
		return -1
	}
	c := strings.IndexRune(n.Content[index-1:], char)
	if c == -1 {
		return -1
	}
	return c + index
}
