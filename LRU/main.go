package main

import "fmt"

type Node struct {
	key, value int
	next, prev *Node
}

type DoubleList struct {
	Head     *Node
	Tail     *Node
	Cache    map[int]*Node
	Capacity int
}

func NewDoubleList(capacity int) *DoubleList {
	head := &Node{}
	tail := &Node{}

	head.next = tail
	tail.prev = head
	return &DoubleList{
		Head:     head,
		Tail:     tail,
		Capacity: capacity,
		Cache:    make(map[int]*Node),
	}
}

func (dll *DoubleList) remove(node *Node) {
	if node == nil || node == dll.Head || node == dll.Tail {
		return
	}
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (dll *DoubleList) AddToHead(node *Node) {
	node.next = dll.Head.next
	node.prev = dll.Head
	dll.Head.next = node
	node.next.prev = node
}

func (dll *DoubleList) Get(key int) int {
	if dll.Head.next == dll.Tail {
		return -1
	}

	node, ok := dll.Cache[key]
	if !ok {
		return -1
	}

	dll.remove(node)
	dll.AddToHead(node)
	return node.value
}

func (dll *DoubleList) Put(key, value int) {
	if dll.Capacity <= 0 {
		return
	}

	node, exists := dll.Cache[key]
	if exists {
		node.value = value
		dll.remove(node)
		dll.AddToHead(node)
		return
	}

	if len(dll.Cache) >= dll.Capacity {
		tailPrev := dll.Tail.prev
		if tailPrev != dll.Head {
			dll.remove(tailPrev)
			delete(dll.Cache, tailPrev.key)
		}
	}

	newNode := &Node{
		key:   key,
		value: value,
	}
	dll.AddToHead(newNode)
	dll.Cache[key] = newNode
}

func (dll *DoubleList) Display() {
	fmt.Println("Displaying ")
	cur := dll.Head.next
	for {
		if cur == dll.Tail {
			break
		}
		fmt.Println(cur.key, cur.value)
		cur = cur.next
	}
	fmt.Println("Done")
}

func main() {
	fmt.Println("Creating LRU cache with capacity 3")
	cache := NewDoubleList(3)

	fmt.Println("\nAdding key-value pairs:")
	fmt.Println("Put(1, 1)")
	cache.Put(1, 1)
	fmt.Println("Put(2, 4)")
	cache.Put(2, 4)
	fmt.Println("Put(3, 9)")
	cache.Put(3, 9)

	fmt.Println("\nCurrent cache state:")
	cache.Display() // Should show: 3->2->1

	fmt.Println("\nAccessing key 2 (should move it to front):")
	val := cache.Get(2)
	fmt.Printf("Get(2) = %d\n", val) // Should print 4

	fmt.Println("\nCache after accessing key 2:")
	cache.Display() // Should show: 2->3->1

	fmt.Println("\nAdding new key (should evict least recently used - key 1):")
	cache.Put(4, 16)
	fmt.Println("Put(4, 16)")

	fmt.Println("\nCache after adding key 4:")
	cache.Display() // Should show: 4->2->3

	fmt.Println("\nUpdating existing key 3:")
	cache.Put(3, 27)
	fmt.Println("Put(3, 27)")

	fmt.Println("\nCache after updating key 3:")
	cache.Display() // Should show: 3->4->2

	fmt.Println("\nTesting non-existent key:")
	val = cache.Get(1)
	fmt.Printf("Get(1) = %d (expected -1 as it was evicted)\n", val)
}
