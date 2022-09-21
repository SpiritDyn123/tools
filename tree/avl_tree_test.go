package tree

import (
	"fmt"
	"math/rand"
	"testing"
)


//go test -v -run AvlTre .
func TestAvlTree(t *testing.T) {
	Num := 12
	values := make([]int, Num)
	for i := 0;i < Num;i++ {
		values[i] = i
	}

	rand.Shuffle(Num, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	if Num <= 20 {
		fmt.Println("======values:", values)
	}

	tree := &AvlTree{}
	for _, v := range values {
		tree.Insert(v)
	}

	fmt.Println("Height:", tree.Height())
	tree.Print()

	for _, v := range values {
		fmt.Println("--------------remove:", v, "----------------")
		tree.Remove(v)
		tree.Print()
	}
}
