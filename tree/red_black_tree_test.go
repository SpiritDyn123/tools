package tree

import (
	"fmt"
	"testing"
)


//go test -v -run RedBlueTre .
func TestRedBlueTree(t *testing.T) {
	Num := 18
	values := make([]int, Num)
	for i := 0;i < Num;i++ {
		values[i] = i
	}

	if Num <= 20 {
		fmt.Println("======values:", values)
	}

	tree := &RedBlackTree{}
	for _, v := range values {
		tree.Insert(v)
		//fmt.Println("remove:", v)
		//tree.Print()
	}
	tree.Print()

	//rand.Shuffle(Num, func(i, j int) {
	//	values[i], values[j] = values[j], values[i]
	//})

	values = []int{12, 8, 10, 9}
	for _, v := range values {
		fmt.Println("--------------remove:", v, "----------------")
		tree.Remove(v)
		tree.Print()
	}
}
