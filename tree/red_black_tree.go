package tree

/*
	红黑树
*/

/*
 - 1.节点是红色或黑色；
 - 2.根节点是黑色；
 - 3.所有叶子节点是黑色；（叶子节点是NULL节点）
 - 4.每个红色节点的两个子节点都是黑色；（从根节点到每个叶子节点的路径上不能有两个连续的红节点）
 - 5.从任何一个节点到每个叶子节点的所有路径都包含相同数目的黑色节点；
*/

type rbColor byte
const (
	rbColor_black = rbColor(0) + iota
	rbColor_red
)

type redBlackNode struct {
	parent, left, right *redBlackNode
	color rbColor //颜色
	value	int   //值
}

func (this *redBlackNode) find(value int) bool {
	if this == nil {
		return false
	}

	if value == this.value {
		return true
	} else if value < this.value {
		return this.left.find(value)
	} else {
		return this.right.find(value)
	}
}

func (this *redBlackNode) isBlack() bool {
	return this != nil && this.color == rbColor_black
}

func (this *redBlackNode) isRed() bool {
	return this != nil && this.color == rbColor_red
}

func (this *redBlackNode) insert(value int, parent *redBlackNode) *redBlackNode {
	if this.find(value) { //查找
		return this
	}

	if this == nil {
		this = &redBlackNode{
			parent: parent,
			color: rbColor_black,
			value: value,
		}

		if parent == nil { //根节点
			return this
		}

		this.color = rbColor_red
		return this
	}

	if value < this.value {
		this.left = this.left.insert(value, this)
	} else {
		this.right = this.right.insert(value, this)
	}

	if this.isBlack() {
		return this
	}

	if this.left == nil {

	} else {

	}

	return this
}

type RedBlackTree struct {
	root *redBlackNode
}

func (this *RedBlackTree) Print() {

}

func (this *RedBlackTree) Insert(value int) {
	this.root = this.root.insert(value, nil)
}

func (this *RedBlackTree) Remove(value int) bool {

}


