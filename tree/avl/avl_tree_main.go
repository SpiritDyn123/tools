package main

import (
	"fmt"
	"math/rand"
)

/*
工具网站：https://www.cs.usfca.edu/~galles/visualization/AVLtree.html
*/
type TreeSortType int
const (
	TreeSortType_First = TreeSortType(1) + iota
	TreeSortType_Midlle
	TreeSortType_Last
)

type AvlTree struct {
	V int
	Left, Right *AvlTree //parent增加会导致难度
}

func (this *AvlTree) Height() int {
	if this == nil {
		return 0
	}

	h := this.Left.Height()
	rh := this.Right.Height()
	if rh > h {
		h = rh
	}

	return h + 1
}

func (this *AvlTree) Data(st TreeSortType) (data []int) {
	if this == nil {
		return
	}

	ldata := this.Left.Data(st)
	rdata := this.Right.Data(st)
	switch st {
	case TreeSortType_First:
		data = append([]int{ this.V }, ldata...)
		data = append(data, rdata...)
	case TreeSortType_Midlle:
		data = append(ldata, this.V)
		data = append(data, rdata...)
	case TreeSortType_Last:
		data = append(ldata, rdata...)
		data = append(data, this.V)
	}

	return
}

func (this *AvlTree) LL() *AvlTree {
	rNode := this.Right
	this.Right = rNode.Left
	rNode.Left = this
	return rNode
}


func (this *AvlTree) RR() *AvlTree  {
	rNode := this.Left
	this.Left = rNode.Right
	rNode.Right = this

	return rNode
}

func (this *AvlTree) LR() *AvlTree {
	this.Left = this.Left.LL()
	this = this.RR()
	return this
}

func (this *AvlTree) RL() *AvlTree {
	this.Right = this.Right.RR()
	this = this.LL()
	return this
}

func (this *AvlTree) fitHeight() *AvlTree {
	rh, lh := this.Right.Height(), this.Left.Height()
	if rh - lh > 1 {
		if this.Right.Left.Height() >= this.Right.Right.Height() { //insert的时候直接用v <= this.Right.V 更高效
			this = this.RL()
		} else {
			this = this.LL()
		}
	} else if lh - rh > 1 {
		if this.Left.Right.Height() >= this.Left.Left.Height() {  //insert的时候直接用v > this.Left.V 更高效
			this = this.LR()
		} else {
			this = this.RR()
		}
	}

	return this
}

func (this *AvlTree) print(depth int, desc string) {
	if this == nil {
		return
	}

	var prefix = ""
	for i := 0;i < depth;i++ {
		prefix += "\t"
	}

	fmt.Printf("%s 层:%d%s value:%d\n",prefix, depth, desc, this.V)
	if this.Left != nil {
		this.Left.print(depth + 1, "左")
	}

	if this.Right != nil {
		this.Right.print(depth + 1, "右")
	}
}

func (this *AvlTree) PrintTree() {
	if this == nil {
		fmt.Println("[NIL TREE]")
		return
	}

	fmt.Printf("ROOT value:%d\n", this.V)
	if this.Left != nil {
		this.Left.print(1, "左")
	}

	if this.Right != nil {
		this.Right.print(1, "右")
	}
}

func (this *AvlTree) Insert(v int) *AvlTree {
	if this == nil {
		this = &AvlTree{
			V: v,
		}
		return this
	}

	if v <= this.V {
		this.Left = this.Left.Insert(v)
	} else {
		this.Right = this.Right.Insert(v)
	}

	//rh, lh := this.Right.Height(), this.Left.Height()
	//if rh - lh > 1 {
	//	if v <= this.Right.V {
	//		this = this.RL()
	//	} else {
	//		this = this.LL()
	//	}
	//} else if lh - rh > 1 {
	//	if v > this.Left.V {
	//		this = this.LR()
	//	} else {
	//		this = this.RR()
	//	}
	//}

	return this.fitHeight()
}

/*
就是找到叶子结点进行替换，分三种情况
 1、左右子树都不为空
 2、左子树不为空或者右子树不为空
 3、左右子树都为空（已经是叶子节点）
*/
func (this *AvlTree) Remove(v int) *AvlTree {
	if this == nil {
		return nil
	}

	eq := false
	if v < this.V {
		this.Left = this.Left.Remove(v)
	} else if v > this.V {
		this.Right = this.Right.Remove(v)
	} else {
		eq = true
	}

	//旋转操作
	if !eq {
		return this.fitHeight()
	}

	//condition 1 直接删除
	if this.Left == nil && this.Right == nil {
		return nil
	}

	if this.Left != nil && this.Right != nil {
		//取左子树最大值替换
		node := this.Left
		for ;node.Right != nil; node = node.Right {}

		//交换值
		node.V, this.V = this.V, node.V
		//再继续删除
		this.Left = this.Left.Remove(v)
		this = this.fitHeight()
	} else if this.Left != nil {
		this.V = this.Left.V
		this.Left = nil
	} else if this.Right != nil {
		this.V = this.Right.V
		this.Right = nil
	}

	return this
}

func main() {
	var avlTree *AvlTree
	fmt.Println(avlTree.Data(TreeSortType_Last))
	data := []int{}
	for i := 0;i < 10;i++ {
		avlTree = avlTree.Insert(i)
		data = append(data, i)
	}

	fmt.Println(avlTree.Data(TreeSortType_Last))
	avlTree.PrintTree()

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})


	for _, d := range data {
		avlTree = avlTree.Remove(d)
		fmt.Println("remove:", d)
		fmt.Println(avlTree.Data(TreeSortType_Last))
	}

	fmt.Println("end success")
}
