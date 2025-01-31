package testfiles

// Recursive node types
type Node[T any] struct {
	Value     T
	LeftNode  *Node[T]
	RightNode *Node[T]
}

type Tree[T any] struct {
	Value     string
	LeftTree  *Node[T]
	RightTree *Node[T]
}

type Role = int

const (
	Admin = iota
	User
)
