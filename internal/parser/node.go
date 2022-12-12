package parser

import "reflect"

type NodeType uint
type ValType uint

const (
	NodeCommand NodeType = iota
	NodeVar              = iota
)
const (
	ValSInt    ValType = iota
	ValUInt            = iota
	ValSInt8           = iota
	ValUInt8           = iota
	ValSInt16          = iota
	ValUInt16          = iota
	ValSInt32          = iota
	ValUInt32          = iota
	ValSInt64          = iota
	ValUInt64          = iota
	ValFloat32         = iota
	ValFloat64         = iota
	ValByte            = iota
	ValString          = iota
)

type Node struct {
	Kind        NodeType
	Class       ValType
	Val         reflect.Value
	Children    int
	FirstChild  *Node
	NextSibling *Node
	PrevSibling *Node
}

func NewNode(t NodeType) *Node {
	n := Node{
		Kind: t,
	}

	return &n
}

func (n *Node) String() string {
	var out string
	return n._string(&out, 0)
}

func (n *Node) _string(out *string, level int) string {
	for i := 0; i < level; i++ {
		*out += " "
	}

	if n.Kind == NodeCommand {
		*out += "COMMAND\n"
	} else {
		*out += n.Val.String() + "\n"
	}

	if n.NextSibling != nil {
		n.NextSibling._string(out, level)
	}

	if n.FirstChild != nil {
		n.FirstChild._string(out, level+1)
	}

	return *out
}

func AddChild(parent *Node, child *Node) {
	if parent == nil || child == nil {
		return
	}

	if parent.FirstChild == nil {
		parent.FirstChild = child
		parent.Children += 1
		return
	}

	sibling := parent.FirstChild
	for ; sibling.NextSibling != nil; sibling = sibling.NextSibling {
	}
	sibling.NextSibling = child
	child.PrevSibling = sibling
	parent.Children += 1
}

func SetNodeValueStr(node *Node, val string) {
	node.Class = ValString

	node.Val = reflect.ValueOf(val)
}
