package main

import (
	"fmt"
	"github.com/TheAlgorithms/Go/structure/hashmap"
)

func main() {
	myMap := hashmap.DefaultNew() // Capacity 1024

	// 1. Insert 800 fake users using a loop
	for i := 1; i <= 1000000; i++ {
		keyName := fmt.Sprintf("User_%d", i)
		myMap.Put(keyName, "Some Data")

	}

	// 2. Insert our specific targets last
	myMap.Put("Alice", "Engineer")
	myMap.Put("Charlie", "Manager")

	// 3. Print the internal structure (Warning: This will print A LOT of text!)
	myMap.Print()

	// 4. See how many steps it takes to find them in a crowded map
	val, steps := myMap.GetWithSteps("User_519")
	fmt.Printf("Found %v in %d steps\n", val, steps)

	val, steps = myMap.GetWithSteps("User_9468")
	fmt.Printf("Found %v in %d steps\n", val, steps)
	val, steps = myMap.GetWithSteps("User_9285")
	fmt.Printf("Found %v in %d steps\n", val, steps)

	myMap.FindWorstChain()
}
