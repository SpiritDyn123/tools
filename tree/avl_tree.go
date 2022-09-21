package tree

import "fmt"

/*
	平衡二叉排序树
	(任何一个节点的左子树或者右子树都是「平衡二叉树」（左右高度差小于等于 1）)
	参考：
	- https://blog.csdn.net/jarvan5/article/details/112428036
	- https://blog.csdn.net/weixin_57023347/article/details/118565474
*/

type avlNode struct {
	value int
	left, right *avlNode
}

func (this *avlNode) height() int {
	if this == nil {
		return 0
	}

	lh := this.left.height()
	rh := this.right.height()
	if lh < rh {
		lh =  rh
	}
	return lh + 1
}

func (this *avlNode) find(value int) bool {
	if this == nil {
		return false
	}

	if this.value == value {
		return true
	}

	if value < this.value {
		return this.left.find(value)
	}

	return this.right.find(value)
}

func (this *avlNode) print(pos string, ceng int, prefix string) {
	if this == nil {
		fmt.Println(prefix + "[NIL TREE]")
		return
	}

	fmt.Printf("%s- %d 层%d %s\n", prefix, this.value,ceng, pos)
	if this.left != nil {
		this.left.print("左", ceng+1, prefix + print_prefix)
	}

	if this.right != nil {
		this.right.print("右", ceng+1, prefix + print_prefix)
	}
}

//左旋
func (this *avlNode) ll() *avlNode {
	newNode := this.right
	this.right = newNode.left
	newNode.left = this
	return newNode
}

//右旋
func (this *avlNode) rr() *avlNode {
	newNode := this.left
	this.left = newNode.right
	newNode.right = this
	return newNode
}

//左右
func (this *avlNode) lr() *avlNode {
	this.left = this.left.ll()
	return this.rr()
}

//右左
func (this *avlNode) rl() *avlNode {
	this.right = this.right.rr()
	return this.ll()
}

func (this *avlNode) insert(value int) *avlNode {
	if this == nil {
		this = &avlNode{
			value: value,
		}
		return this
	}

	if value <= this.value {
		this.left = this.left.insert(value)
	} else {
		this.right = this.right.insert(value)
	}

	//判断高度差
	lh := this.left.height()
	rh := this.right.height()
	if lh - rh > 1 { //插入左边了
		if value < this.left.value { //RR
			this = this.rr()
		} else { //LR
			this = this.lr()
		}
	} else if lh - rh < -1 { //插入右边了
		if value > this.right.value { //LL
			this = this.ll()
		} else { //RL
			this = this.rl()
		}
	}

	return this
}

func (this *avlNode) GetMaxNode() *avlNode {
	if this == nil {
		return nil
	}

	if this.right == nil {
		return this
	}

	return this.right.GetMaxNode()
}

func (this *avlNode) GetMinNode() *avlNode {
	if this == nil {
		return nil
	}

	if this.left == nil {
		return this
	}

	return this.left.GetMinNode()
}

/*
	1、叶子
	2、left != nil && right == nil
	3、left == nil && right != nil
	4、left != nil && right != nil
*/

func (this *avlNode) remove(value int) *avlNode {
	if this.value == value {
		if this.left != nil && this.right != nil {
			max_node := this.left.GetMaxNode() //左节点最大值替换
			this.value = max_node.value
			this.left = this.left.remove(max_node.value)
			if this.right.height() - this.left.height() > 1 {
				rightNode := this.right
				if rightNode.left.height() > rightNode.right.height() { //RL
					this = this.rl()
				} else { //LL
					this = this.ll()
				}
			}

		} else if this.left != nil { //左节点替代
			this.value = this.left.value
			this.right = this.left.right
			this.left = this.left.left
		} else if this.right != nil { //右节点替代
			this.value = this.right.value
			this.left = this.right.left
			this.right = this.right.right
		} else {
			return nil
		}
	} else if value < this.value {
		this.left = this.left.remove(value)
		if this.right.height() - this.left.height() > 1 {
			rightNode := this.right
			if rightNode.left.height() > rightNode.right.height() { //RL
				this = this.rl()
			} else { //LL
				this = this.ll()
			}
		}
	} else {
		this.right = this.right.remove(value)
		if this.left.height() - this.right.height() > 1 {
			leftNode := this.right
			if leftNode.left.height() > leftNode.right.height() { //RR
				this = this.rr()
			} else { //LR
				this = this.lr()
			}
		}
	}

	return this
}

type AvlTree struct {
	root *avlNode
}

func (this *AvlTree) Height() int {
	return this.root.height()
}

func (this *AvlTree) Print() {
	this.root.print("根", 1, "")
}

func (this *AvlTree) Find(value int) bool {
	return this.root.find(value)
}

func (this *AvlTree) Insert(value int) bool {
	if this.Find(value) {
		return false
	}
	this.root = this.root.insert(value)
	return true
}

func (this *AvlTree) Remove(value int) bool {
	if !this.Find(value) {
		return false
	}

	this.root  = this.root.remove(value)
	return true
}