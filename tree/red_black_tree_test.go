package tree

import (
	"fmt"
	"math/rand"
	"testing"
)


//go test -v -run RedBlueTre .
func TestRedBlueTree(t *testing.T) {
	Num := 10000
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

	tree := &RedBlackTree{}
	for _, v := range values {
		tree.Insert(v)
		if Num < 20 {
			fmt.Println("------------insert:", v)
			tree.Print()
		}
	}

	tree.Print()
	fmt.Println("============================================")

	//values = []int{3, 7, 10, 9}
	r_cnt := 0
	for _, v := range values {
		if tree.Remove(v) {
			r_cnt++
		}
		if Num < 20 {
			fmt.Println("--------------remove:", v, "----------------")
			tree.Print()
		}
	}

	tree.Print()
	fmt.Println("remove count:", r_cnt)
}
