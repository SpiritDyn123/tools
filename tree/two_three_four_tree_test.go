package tree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_TwoThreeFourTree(t *testing.T) {
	//初始化乱序数组
	Num := 54
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

	tree := &TwoThreeFourTree{}
	for _, v := range values {
		tree.Insert(v)
	}

	fmt.Println("======height:", tree.Height())
	tree.Print()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 2;i++ {
		rv := values[rand.Intn(len(values))]
		fmt.Println("value:", rv, "find count:", tree.Find(rv))
		fmt.Println("value:", rv, "path:", tree.Path(rv))
	}

	//tv := 35
	//fmt.Println("value:", tv, "path:", tree.Path(tv))
	//tree.Remove(tv)
	//tree.Print()
	for _, rv := range values {
		fmt.Println("value:", rv, "remove:", tree.Remove(rv))
		//tree.Print()
	}

	tree.Print()

	for _, v := range values {
		tree.Insert(v)
	}

	tree.Print()
}

