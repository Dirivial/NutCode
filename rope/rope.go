package rope

type Rope struct {
	Head Node
}

type Node struct {
	Left    *Node
	Right   *Node
	Content string
	Weight  int
}

func New() *Rope {
	r := &Rope{}
	return r
}

// Insert a string into the rope structure
func (*Rope) Insert(index int, str string) {
}

// Delete part of the rope structure
func (*Rope) Delete(start, length int) {
}

// Get character at a position
func (*Rope) Index(index int) int {
	return 0
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
func (*Rope) Concat(rope *Rope) Rope {
	return Rope{}
}

// Split a rope into two ropes
func (*Rope) Split(index int) []Rope {
	return []Rope{}
}
