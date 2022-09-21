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

参考：
 - https://blog.csdn.net/ystyaoshengting/article/details/121423014
 - https://zhuanlan.zhihu.com/p/269069974
 - https://blog.csdn.net/weixin_29163797/article/details/113332639
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
	}

	if value < this.value {
		return this.left.find(value)
	}

	return this.right.find(value)
}

func (this *redBlackNode) isBlack() bool {
	return this != nil && this.color == rbColor_black
}

func (this *redBlackNode) isRed() bool {
	return this != nil && this.color == rbColor_red
}
/*
	1、单黑色节点（类似234树2节点）
	2、黑色节点+单红色子节点（左右）（类似234树3节点）
	3、黑色节点+2红色子节点（类似234树4节点）
*/
func (this *redBlackNode) insert(value int, parent *redBlackNode) *redBlackNode {
	if this == nil {
		this = &redBlackNode{
			parent: parent,
			color: rbColor_red,
			value: value,
		}

		//根节点
		if parent == nil {
			this.color = rbColor_black
			return this
		}

		brotherNode := this.parent.right
		if value > this.parent.value {
			brotherNode = this.parent.left
		}

		if this.parent.isBlack() {
			if brotherNode == nil { //2节点

			} else { //3节点

			}
		} else { //4节点
			pparent := parent.parent
			if pparent == nil || !pparent.isBlack() { //不能连续红色节点
				panic("pparent must not be nil and be black")
			}

			

		}


		return this
	}

	if value < this.value {
		this.left = this.left.insert(value, this)
	} else {
		this.right = this.right.insert(value, this)
	}

	if this.isBlack() { //黑节点
		if this.left == nil || this.right == nil { //黑色单节点

		}

	} else { //黑节点

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

func(this *RedBlackTree) Find(value int) bool {
	return this.root.find(value)
}

func (this *RedBlackTree) Insert(value int) bool {
	if this.root.find(value) { //查找
		return false
	}

	this.root = this.root.insert(value, nil)
	return true
}

func (this *RedBlackTree) Remove(value int) bool {
	return false
}


