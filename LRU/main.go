package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	Key   int
	Value int
	TTL   int64
	Left  *Node
	Right *Node
}

type DoubleList struct {
	mu       sync.RWMutex
	Head     *Node
	Tail     *Node
	Capacity int
	Cache    map[int]*Node
}

func CreateNewDL(capacity int) *DoubleList {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	dl := &DoubleList{
		Head:     &Node{},
		Tail:     &Node{},
		Capacity: capacity,
		Cache:    make(map[int]*Node),
	}
	dl.Head.Right = dl.Tail
	dl.Tail.Left = dl.Head
	return dl
}

func (dl *DoubleList) Get(key int) int {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	node, ok := dl.Cache[key]
	if !ok {
		return -1
	}

	// Check if the key has expired
	if node.TTL > 0 && node.TTL < time.Now().Unix() {
		dl.RemoveNode(node)
		delete(dl.Cache, key)
		return -1
	}

	dl.RemoveNode(node)
	dl.AddToFront(node)
	return node.Value
}

func (dl *DoubleList) RemoveNode(node *Node) {
	left := node.Left
	right := node.Right
	left.Right = right
	right.Left = left

}

func (dl *DoubleList) AddToFront(node *Node) {
	next := dl.Head.Right
	prev := dl.Head

	node.Right = next
	node.Left = prev
	prev.Right = node
	next.Left = node
}

func (dl *DoubleList) Set(key, value int, seconds int) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	ttl := int64(0)
	if seconds > 0 {
		ttl = time.Now().Unix() + int64(seconds)
	}

	node, ok := dl.Cache[key]
	if ok {
		node.Value = value
		node.TTL = ttl // Update TTL for existing keys
		dl.RemoveNode(node)
		dl.AddToFront(node)
		return
	}

	if dl.Capacity == len(dl.Cache) {
		lastNode := dl.Tail.Left
		dl.RemoveNode(lastNode)
		delete(dl.Cache, lastNode.Key)
	}

	node = &Node{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	dl.AddToFront(node)
	dl.Cache[key] = node
}

// Size returns the current number of items in the cache
func (dl *DoubleList) Size() int {
	dl.mu.RLock()
	defer dl.mu.RUnlock()
	return len(dl.Cache)
}

func (dl *DoubleList) Display() {
	cur := dl.Head.Right
	for {
		if cur != dl.Tail {
			fmt.Println("key : ", cur.Key, " Value: ", cur.Value)
			cur = cur.Right
		} else {
			break
		}
	}
	fmt.Println()
}

func (dl *DoubleList) FreeKeys() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		dl.mu.Lock()

		cur := dl.Head.Right
		for cur != dl.Tail {
			if cur.TTL > 0 && cur.TTL < time.Now().Unix() {
				fmt.Printf("freeing key %d\n", cur.Key)
				next := cur.Right
				dl.RemoveNode(cur)
				delete(dl.Cache, cur.Key)
				cur = next
			} else {
				cur = cur.Right
			}
		}

		dl.mu.Unlock()
	}
}

func main() {

	dl := CreateNewDL(5)
	go dl.FreeKeys()

	dl.Set(1, 1, 5)
	dl.Set(2, 4, 6)
	dl.Set(3, 9, 7)
	dl.Set(4, 16, 8)
	dl.Set(5, 25, 50)

	dl.Display()

	for i := 1; i < 8; i++ {
		value := dl.Get(i)
		fmt.Printf("key %d value %d\n", i, value)
	}

	dl.Set(6, 36, 60)
	dl.Display()
	time.Sleep(10 * time.Second)
	dl.Display()

	dl.Set(1, 1, 5)
	dl.Set(2, 4, 6)
	dl.Set(3, 9, 7)
	dl.Set(4, 16, 8)
	dl.Set(5, 25, 50)

	dl.Display()

}
