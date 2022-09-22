package tree

import (
	"fmt"
	"testing"
)


//go test -v -run RedBlueTre .
func TestRedBlueTree(t *testing.T) {
	Num := 20
	values := make([]int, Num)
	for i := 0;i < Num;i++ {
		values[i] = i
	}

	//rand.Shuffle(Num, func(i, j int) {
	//	values[i], values[j] = values[j], values[i]
	//})

	if Num <= 20 {
		fmt.Println("======values:", values)
	}

	tree := &RedBlackTree{}
	for _, v := range values {
		tree.Insert(v)
		fmt.Println("----------------------------v:", v)
		tree.Print()
	}


	for _, v := range values {
		//fmt.Println("--------------remove:", v, "----------------")
		tree.Remove(v)
		//tree.Print()
	}
}
