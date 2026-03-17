# Demystifying the Hash Map (Go Implementation)

Welcome! If you are new to data structures, reading a Hash Map implementation in Go can be intimidating. This guide breaks down exactly how this code works using a simple real-world analogy, visualizes the data, and provides a sandbox for you to test its performance limits.

---

## 1. The "Post Office" Analogy

To understand how this code works, imagine a Hash Map is a **Post Office**. 
* You have a wall of exactly 1,024 empty **PO Boxes**.
* You want to store **Letters** (Values) inside envelopes that have names written on them (Keys).
* You use a **Magic Calculator** (Hash Function) to turn a name like "Alice" into a specific PO Box number.

### Mapping the Analogy to the Code:
* **The Wall of Boxes (`table []*node`)**: An array of pointers. It starts completely empty (`nil`).
* **The Envelope (`node struct`)**: This holds your `key` ("Alice") and your `value` (her data).
* **The Piece of String (`next *node`)**: What happens if "Alice" and "Bob" both get assigned to Box #42 by the magic calculator? We can't throw one away! Instead, we put Alice in the box, and use a "piece of string" (`next`) to tie Bob's envelope to Alice's. This is called **Separate Chaining** (handling collisions using a linked list).

---

## 2. Why `1 << 10`? (The Magic of Bitwise Shifts)

In the code, you will see the starting capacity defined like this:
`var defaultCapacity uint64 = 1 << 10`

Why use a bitwise shift instead of just writing `1024`? And why can't we use a clean number like `1000`?

Hash Maps rely heavily on their capacity being a perfect **power of 2** so the math runs lightning fast. `1 << 10` takes the binary number `1` and pushes it left by 10 zeroes, which equals exactly 1024. Writing it this way is a massive neon sign to other programmers: *"Do not change this to 1000! It must be a perfect power of 2!"*

### Why 1000 Breaks the Algorithm
Look at the `hash()` function. It uses this formula to find the right box:
`return (capacity - 1) & hashValue`

Instead of using slow division (`%`), it uses a **Bitwise AND (`&`)**. This trick *only* works with powers of 2.
* **If capacity is 1024:** `1024 - 1 = 1023`. In binary, 1023 is solid ones (`1111111111`). This perfectly distributes the items across all 1024 boxes.
* **If capacity is 1000:** `1000 - 1 = 999`. In binary, 999 has zeroes in it (`1111100111`). The `&` operator forces those zero spots to *always* be zero. As a result, certain boxes will **never** receive an item, sitting permanently empty, while other boxes get massively overcrowded. The map is broken!

---

## 3. Performance (Big O Complexity)

The entire purpose of a Hash Map is to solve the **Search Problem**. If you have 1 million users, finding "Diana" shouldn't take 1 million steps.

| Operation | Average Case | Worst Case (Overloaded/Bad Hash) |
| :--- | :--- | :--- |
| **Insert (`Put`)** | **O(1)** Instant | **O(n)** Slow |
| **Search (`Get`)** | **O(1)** Instant | **O(n)** Slow |

### The Impact of Size & Resizing
* **The O(1) Best Case:** When there are plenty of empty boxes, the magic calculator puts "Diana" in her own private box. Searching for her takes exactly **1 step**.
* **The O(n) Worst Case:** If you cram 100,000 items into 1,024 boxes, the "chains" of envelopes tied together become massive. If Diana is at the end of a chain of 100 envelopes, the computer must manually check all 100. The search time grows linearly with the data size (**O(n)**).
* **The Savior (`resize()`):** To prevent O(n), the Hash Map tracks its "Load Factor". If it gets more than 75% full, `hm.resize()` kicks in, doubles the amount of boxes to 2048, and redistributes everyone. This keeps chains short and speed permanently at O(1)!

---

## 4. Test it Yourself (The Local Sandbox)

Want to physically see the difference between O(1) and O(n)? You can run this test suite on your local machine.

### Step 1: Add the "Detective Tools" to `hashmap.go`
Paste these three helper methods at the very bottom of `structure/hashmap/hashmap.go`. They allow us to peek inside the memory and count the search steps.

```go
import "fmt"

// 1. Print shows the internal chains of the HashMap
func (hm *HashMap) Print() {
	fmt.Printf("\n--- Current HashMap State (Size: %d | Capacity: %d) ---\n", hm.size, hm.capacity)
	for index, headNode := range hm.table {
		if headNode != nil {
			fmt.Printf("Bucket %d: ", index)
			current := headNode
			for current != nil {
				fmt.Printf("[%v] -> ", current.key)
				current = current.next
			}
			fmt.Println("nil")
		}
	}
	fmt.Println("-----------------------------------------------------")
}

// 2. GetWithSteps returns the value AND the exact number of nodes it had to check
func (hm *HashMap) GetWithSteps(key any) (value any, steps int) {
	index := hm.hash(key)
	current := hm.table[index]
	steps = 0 
	for current != nil {
		steps++ 
		if current.key == key {
			return current.value, steps
		}
		current = current.next
	}
	return nil, steps
}

// 3. FindWorstChain hunts for the bucket with the longest collision chain
func (hm *HashMap) FindWorstChain() {
	maxChainLength := 0
	worstBucketIndex := 0
	var worstLastKey any

	for index, headNode := range hm.table {
		chainLength := 0
		current := headNode
		var currentChainLastKey any 
		
		for current != nil {
			chainLength++
			currentChainLastKey = current.key 
			current = current.next
		}

		if chainLength > maxChainLength {
			maxChainLength = chainLength
			worstBucketIndex = index
			worstLastKey = currentChainLastKey 
		}
	}

	if maxChainLength == 0 {
		return
	}

	fmt.Printf("\n🚨 WORST CASE SCENARIO FOUND 🚨\n")
	fmt.Printf("Bucket %d is the most crowded with %d items!\n", worstBucketIndex, maxChainLength)
	fmt.Printf("Testing Search Speed for the last item in that chain...\n")
	_, steps := hm.GetWithSteps(worstLastKey)
	fmt.Printf("Result: Found '%v' in exactly %d steps.\n\n", worstLastKey, steps)
}
