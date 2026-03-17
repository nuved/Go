package hashmap

import (
	"fmt"
	"hash/fnv"
)

var defaultCapacity uint64 = 1 << 10

type node struct {
	key   any
	value any
	next  *node
}

// HashMap is a Golang implementation of a hashmap
type HashMap struct {
	capacity uint64
	size     uint64
	table    []*node
}

// DefaultNew returns a new HashMap instance with default values
func DefaultNew() *HashMap {
	return &HashMap{
		capacity: defaultCapacity,
		table:    make([]*node, defaultCapacity),
	}
}

// New creates a new HashMap instance with the specified size and capacity
func New(size, capacity uint64) *HashMap {
	return &HashMap{
		size:     size,
		capacity: capacity,
		table:    make([]*node, capacity),
	}
}

// Get returns the value associated with the given key
func (hm *HashMap) Get(key any) any {
	node := hm.getNodeByKey(key)
	if node != nil {
		return node.value
	}
	return nil
}

// Put inserts a new key-value pair into the hashmap
func (hm *HashMap) Put(key, value any) {
	index := hm.hash(key)
	if hm.table[index] == nil {
		hm.table[index] = &node{key: key, value: value}
	} else {
		current := hm.table[index]
		for {
			if current.key == key {
				current.value = value
				return
			}
			if current.next == nil {
				break
			}
			current = current.next
		}
		current.next = &node{key: key, value: value}
	}
	hm.size++
	if float64(hm.size)/float64(hm.capacity) > 0.75 {
		hm.resize()
	}
}

// Contains checks if the given key is stored in the hashmap
func (hm *HashMap) Contains(key any) bool {
	return hm.getNodeByKey(key) != nil
}

// getNodeByKey finds the node associated with the given key
func (hm *HashMap) getNodeByKey(key any) *node {
	index := hm.hash(key)
	current := hm.table[index]
	for current != nil {
		if current.key == key {
			return current
		}
		current = current.next
	}
	return nil
}

// resize doubles the capacity of the hashmap and rehashes all existing entries
func (hm *HashMap) resize() {
	oldTable := hm.table
	hm.capacity <<= 1
	hm.table = make([]*node, hm.capacity)
	hm.size = 0

	for _, head := range oldTable {
		for current := head; current != nil; current = current.next {
			hm.Put(current.key, current.value)
		}
	}
}

// hash generates a hash value for the given key
func (hm *HashMap) hash(key any) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(fmt.Sprintf("%v", key)))
	hashValue := h.Sum64()
	return (hm.capacity - 1) & (hashValue ^ (hashValue >> 16))
}

// GetWithSteps returns the value AND the exact number of nodes it had to check.
func (hm *HashMap) GetWithSteps(key any) (value any, steps int) {
	// 1. Find the bucket
	index := hm.hash(key)
	current := hm.table[index]

	steps = 0 // Start our counter at 0

	// 2. Walk through the chain (the "string" of envelopes)
	for current != nil {
		steps++ // We are looking at a new envelope, so add 1 step!

		if current.key == key {
			// We found it! Return the value and how many steps it took.
			return current.value, steps
		}

		// Not a match, move to the next envelope in the chain
		current = current.next
	}

	// 3. If we get here, the key doesn't exist.
	return nil, steps
}


// Print shows the internal structure of the HashMap, including the chains.
func (hm *HashMap) Print() {
	fmt.Println("--- Current HashMap State ---")
	fmt.Printf("Size: %d | Capacity: %d\n", hm.size, hm.capacity)

	// 1. Loop through every single "PO Box"
	for index, headNode := range hm.table {

		// 2. Only print if the box is NOT empty
		if headNode != nil {
			fmt.Printf("Bucket %d: ", index)

			// 3. Start at the first envelope
			current := headNode

			// 4. Follow the "string" until we hit the end (nil)
			for current != nil {
				fmt.Printf("[%v: %v] -> ", current.key, current.value)
				current = current.next // Move to the next envelope
			}

			// 5. Mark the end of the chain
			fmt.Println("nil")
		}
	}
	fmt.Println("-----------------------------")
}

// FindWorstChain hunts for the bucket with the longest collision chain and prints it.
func (hm *HashMap) FindWorstChain() {
	maxChainLength := 0
	worstBucketIndex := 0
	var worstLastKey any

	for index, headNode := range hm.table {
		chainLength := 0
		current := headNode
		var currentChainLastKey any // Temp variable for this specific bucket
		
		for current != nil {
			chainLength++
			currentChainLastKey = current.key // Track the last key of THIS bucket
			current = current.next
		}

		// THE FIX: Only save the key if this bucket breaks the record!
		if chainLength > maxChainLength {
			maxChainLength = chainLength
			worstBucketIndex = index
			worstLastKey = currentChainLastKey 
		}
	}

	if maxChainLength == 0 {
		fmt.Println("The map is completely empty!")
		return
	}

	fmt.Printf("\n🚨 WORST CASE SCENARIO FOUND 🚨\n")
	fmt.Printf("Bucket %d is the most crowded with %d items!\n", worstBucketIndex, maxChainLength)
	
	fmt.Printf("Bucket %d: ", worstBucketIndex)
	current := hm.table[worstBucketIndex]
	for current != nil {
		fmt.Printf("[%v] -> ", current.key)
		current = current.next
	}
	fmt.Println("nil\n")

	fmt.Printf("Testing O(n) Search Speed for the last item in the chain...\n")
	_, steps := hm.GetWithSteps(worstLastKey)
	fmt.Printf("Result: Found '%v' in exactly %d steps.\n", worstLastKey, steps)
}
