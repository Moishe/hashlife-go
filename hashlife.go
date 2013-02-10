package hashlife

import (
	"errors"
	"log"
	"math"
)

type NodeChildren struct {
	ul, ur, ll, lr *Node
}

type Node struct {
	// Nodes are nil if a leaf node.
	nc NodeChildren

	// Unused if not a leaf node.
	value byte

	// Level = 1 for leaf node, 2 for node consisting of four leaves, etc.
	level uint

	// The next node from this node. nil if not yet computed.
	next *Node
}

type NodeMap map[NodeChildren]*Node

var (
	OnLeaf  = &Node{value: 1, level: 1}
	OffLeaf = &Node{value: 0, level: 1}
)

var CacheHits int32 = 0
var CacheMisses int32 = 0
var SkippedCache int32 = 0

// TODO (moishel): for testing, let this function be passed to NextGeneration
func NodeFromNeighborCount(c byte, current *Node) *Node {
	switch c {
	case 3:
		return OnLeaf
	case 2:
		return current
	}
	return OffLeaf
}

var nodeMap NodeMap = make(NodeMap)

func FindNode(ul, ur, ll, lr *Node) *Node {
	nc := NodeChildren{ul, ur, ll, lr}
	node, ok := nodeMap[nc]
	if ok {
		return node
	}
	node = &Node{nc, 0, ul.level + 1, nil}
	nodeMap[nc] = node
	return node
}

func CountsFromLeaves(n *Node) (ulc, urc, llc, lrc byte) {
	ulc = (n.nc.ul.nc.ul.value +
		n.nc.ul.nc.ur.value +
		n.nc.ul.nc.ll.value +
		n.nc.ur.nc.ul.value +
		n.nc.ur.nc.ll.value +
		n.nc.ll.nc.ul.value +
		n.nc.ll.nc.ur.value +
		n.nc.lr.nc.ul.value)

	urc = (n.nc.ur.nc.ur.value +
		n.nc.ur.nc.ul.value +
		n.nc.ur.nc.lr.value +
		n.nc.ul.nc.ur.value +
		n.nc.ul.nc.lr.value +
		n.nc.lr.nc.ur.value +
		n.nc.lr.nc.ul.value +
		n.nc.ll.nc.ur.value)

	llc = (n.nc.ll.nc.ll.value +
		n.nc.ll.nc.ul.value +
		n.nc.ll.nc.lr.value +
		n.nc.ul.nc.ll.value +
		n.nc.ul.nc.lr.value +
		n.nc.lr.nc.ll.value +
		n.nc.lr.nc.ul.value +
		n.nc.ur.nc.ll.value)

	lrc = (n.nc.lr.nc.ll.value +
		n.nc.lr.nc.lr.value +
		n.nc.lr.nc.ur.value +
		n.nc.ll.nc.ur.value +
		n.nc.ll.nc.lr.value +
		n.nc.ur.nc.ll.value +
		n.nc.ur.nc.lr.value +
		n.nc.ul.nc.lr.value)

	return
}

type NodeLoc struct {
	i int
	n *Node
}

type NextFunc func(c chan NodeLoc, i int, ul *Node, ur *Node, ll *Node, lr *Node)
type SimpleNextFunc func(c chan NodeLoc, i int, node *Node)

func NextGenLevel3(n *Node) *Node {
	ulc, urc, llc, lrc := CountsFromLeaves(n)
	ul := NodeFromNeighborCount(ulc, n.nc.ul.nc.lr)
	ur := NodeFromNeighborCount(urc, n.nc.ur.nc.ll)
	ll := NodeFromNeighborCount(llc, n.nc.ll.nc.ur)
	lr := NodeFromNeighborCount(lrc, n.nc.lr.nc.ul)

	newNode := FindNode(ul, ur, ll, lr)

	n.next = newNode
	return newNode
}

func NextGeneration(n *Node) *Node {
	if n.next != nil {
		return n.next
	}

	if n.level == 3 {
		return NextGenLevel3(n)
	}

	var nodes [10]*Node
	next := func(c chan NodeLoc, i int, ul *Node, ur *Node, ll *Node, lr *Node) {
		nodes[i] = NextGeneration(FindNode(ul, ur, ll, lr))
	}

	simple_next := func(c chan NodeLoc, i int, node *Node) {
		nodes[i] = NextGeneration(node)
	}

	Get9Grid(nil, n, next, simple_next)

	piece_map := [16]int{
		1, 2, 4, 5,
		2, 3, 5, 6,
		4, 5, 7, 8,
		5, 6, 8, 9}

	var quadrant_nodes [4]*Node
	for i := 0; i < 4; i++ {
		quadrant_nodes[i] = NextGeneration(FindNode(nodes[piece_map[i*4]],
			nodes[piece_map[i*4+1]],
			nodes[piece_map[i*4+2]],
			nodes[piece_map[i*4+3]]))
	}

	newNode := FindNode(
		quadrant_nodes[0],
		quadrant_nodes[1],
		quadrant_nodes[2],
		quadrant_nodes[3])

	n.next = newNode
	return newNode
}

func Get9Grid(c chan NodeLoc, n *Node, next NextFunc, simple_next SimpleNextFunc) {
	simple_next(c, 1, n.nc.ul)
	next(c, 2, n.nc.ul.nc.ur,
		n.nc.ur.nc.ul,
		n.nc.ul.nc.lr,
		n.nc.ur.nc.ll)
	simple_next(c, 3, n.nc.ur)
	next(c, 4, n.nc.ul.nc.ll,
		n.nc.ul.nc.lr,
		n.nc.ll.nc.ul,
		n.nc.ll.nc.ur)
	next(c, 5, n.nc.ul.nc.lr,
		n.nc.ur.nc.ll,
		n.nc.ll.nc.ur,
		n.nc.lr.nc.ul)
	next(c, 6, n.nc.ur.nc.ll,
		n.nc.ur.nc.lr,
		n.nc.lr.nc.ul,
		n.nc.lr.nc.ur)
	simple_next(c, 7, n.nc.ll)
	next(c, 8, n.nc.ll.nc.ur,
		n.nc.lr.nc.ul,
		n.nc.ll.nc.lr,
		n.nc.lr.nc.ll)
	simple_next(c, 9, n.nc.lr)
}

func TreeFromBitmapBase(bits []byte) (n *Node, err error) {
	size := math.Sqrt(float64(len(bits)))
	f, e := math.Frexp(size)
	if f != 0.5 {
		log.Println("Invalid board size: ", len(bits), size, f, e)
		return nil, errors.New("Invalid board size.")
	}
	return TreeFromBitmap(bits, uint(e), uint(e), 0, 0), nil
}

func TreeFromBitmap(bits []byte, level uint, sublevel uint, x uint, y uint) *Node {
	var bitmap_width uint = 1 << (level - 1)
	if sublevel == 1 {
		if bits[x+y*bitmap_width] == 1 {
			return OnLeaf
		}
		return OffLeaf
	}
	sublevel -= 1
	var child_width uint = 1 << (sublevel - 1)
	return FindNode(
		TreeFromBitmap(bits, level, sublevel, x, y),
		TreeFromBitmap(bits, level, sublevel, x+child_width, y),
		TreeFromBitmap(bits, level, sublevel, x, y+child_width),
		TreeFromBitmap(bits, level, sublevel, x+child_width, y+child_width))
}

type Location struct {
	x, y  int
	value byte
}

func DumpNode(node *Node) []byte {
	size := 1 << (node.level - 1)
	result := make([]byte, size*size)
	c := make(chan Location)
	go func() {
		GetValues(c, node, 0, 0)
		close(c)
	}()

	for i := range c {
		result[i.x+i.y*size] = i.value
	}
	return result
}

func GetValues(c chan Location, node *Node, x int, y int) {
	if node.level == 1 {
		c <- Location{x, y, node.value}
		return
	}
	done := make(chan bool)
	get := func(c chan Location, node *Node, x int, y int) {
		GetValues(c, node, x, y)
		done <- true
	}
	halfsize := 1 << (node.level - 2)
	go get(c, node.nc.ul, x, y)
	go get(c, node.nc.ur, x+halfsize, y)
	go get(c, node.nc.ll, x, y+halfsize)
	go get(c, node.nc.lr, x+halfsize, y+halfsize)

	for i := 0; i < 4; i++ {
		<-done
	}
}
