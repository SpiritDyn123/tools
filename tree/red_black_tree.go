package tree

import "fmt"

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
 - 工具网站：https://www.cs.usfca.edu/~galles/visualization/RedBlack.html
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

func (this *redBlackNode) valueWithColor() string {
	/*
	 // 前景 背景 颜色
	    // ---------------------------------------
	    // 30  40  黑色
	    // 31  41  红色
	    // 32  42  绿色
	    // 33  43  黄色
	    // 34  44  蓝色
	    // 35  45  紫红色
	    // 36  46  青蓝色
	    // 37  47  白色
	    //
	    // 代码 意义
	    // -------------------------
	    //  0  终端默认设置
	    //  1  高亮显示
	    //  4  使用下划线
	    //  5  闪烁
	    //  7  反白显示
	    //  8  不可见
	*/
	color := 30
	if this.color == rbColor_red {
		color = 31
	}
	return fmt.Sprintf("%c[1;40;%dm%d%c[0m", 0x1B, color, this.value, 0x1B)
}

func (this *redBlackNode) print(pos string, ceng int, prefix string) {
	if this == nil {
		fmt.Println("[NIL TREE]")
		return
	}

	fmt.Printf("%s- %s 层%d %s\n", prefix, this.valueWithColor(), ceng, pos)
	if this.left != nil {
		this.left.print("左", ceng+1, prefix + print_prefix)
	}

	if this.right != nil {
		this.right.print("右", ceng+1, prefix + print_prefix)
	}
}

/*
	1、单黑色节点（类似234树2节点）
	2、黑色节点+单红色子节点（左右）（类似234树3节点）
	3、黑色节点+2红色子节点（类似234树4节点）
*/
func (this *redBlackNode) case2() bool {//2节点
	return this != nil && this.color == rbColor_black && this.left == nil && this.right == nil
}

func (this *redBlackNode) case3() bool {//3节点
	return this != nil && this.color == rbColor_black &&
		(this.left != nil && (this.right == nil || this.right.color == rbColor_black) ||
			(this.left == nil || this.left.color == rbColor_black) && this.right != nil)
}

func (this *redBlackNode) case4() bool {////4节点
	return this != nil && this.color == rbColor_black &&
		this.left != nil && this.left.color == rbColor_red &&
		this.right != nil  && this.right.color == rbColor_red
}

func (this *redBlackNode) caseV() int {
	if this.case2() {
		return 2
	} else if this.case3() {
		return 3
	} else if this.case4() {
		return 4
	} else {
		return 0
	}
}

func (this *redBlackNode) ll() *redBlackNode {
	newNode := this.right
	this.right = newNode.left
	newNode.left = this

	if this.right != nil {
		this.right.parent = this
	}

	newNode.parent = this.parent
	if newNode.left != nil {
		newNode.left.parent = newNode
	}
	if newNode.right != nil {
		newNode.right.parent = newNode
	}

	newNode.color = rbColor_black
	if newNode.left != nil {
		newNode.left.color = rbColor_red
	}
	if newNode.right != nil {
		newNode.right.color = rbColor_red
	}
	return newNode
}

func (this *redBlackNode) rr() *redBlackNode {
	newNode := this.left
	this.left = newNode.right
	newNode.right = this

	if this.left != nil {
		this.left.parent = this
	}

	newNode.parent = this.parent
	if newNode.left != nil {
		newNode.left.parent = newNode
	}
	if newNode.right != nil {
		newNode.right.parent = newNode
	}

	newNode.color = rbColor_black
	if newNode.left != nil {
		newNode.left.color = rbColor_red
	}
	if newNode.right != nil {
		newNode.right.color = rbColor_red
	}

	return newNode
}

func (this *redBlackNode) lr() *redBlackNode {
	this.left = this.left.ll()
	return this.rr()
}

func (this *redBlackNode) rl() *redBlackNode {
	this.right = this.right.rr()
	return this.ll()
}

func (this *redBlackNode) insert(value *int, parent *redBlackNode) *redBlackNode {
	if this == nil {
		this = &redBlackNode{
			parent: parent,
			color: rbColor_red,
			value: *value,
		}

		//根节点
		if parent == nil {
			this.color = rbColor_black
			return this
		}

		return this
	}

	childNode := &this.right
	if *value < this.value {
		childNode = &this.left
	}

	case_v := this.caseV()

	newNode := (*childNode).insert(value, this)
	*childNode = newNode
	if case_v == 0 || case_v == 2 || newNode.value == *value {//回退到黑色根节点（爷爷节点），情况3是防止case3直接插入
		return this
	}

	b_left := *value < this.value
	switch case_v {
	case 3:
		childNode = &newNode.right
		if *value < newNode.value {
			childNode = &newNode.left
		}

		if (*childNode).color == rbColor_black {
			return this
		}

		if b_left { //左子树
			if *value < this.left.value { //RR
				this = this.rr()
			} else { //LR
				this = this.lr()
			}
		} else { //右子树
			if *value > this.right.value { //LL
				this = this.ll()
			} else { //RL
				this = this.rl()
			}
		}
	case 4:
		if this.parent == nil { //根节点是黑色
			this.color = rbColor_black
		} else {
			this.color = rbColor_red
		}

		this.left.color = rbColor_black
		this.right.color = rbColor_black
	default:
		panic(fmt.Sprintf("invalid case:%d", case_v))
	}

	*value = this.value

	return this
}

type RedBlackTree struct {
	root *redBlackNode
}

func (this *RedBlackTree) Print() {
	this.root.print("根", 1, "")
}

func(this *RedBlackTree) Find(value int) bool {
	return this.root.find(value)
}

func (this *RedBlackTree) Insert(value int) bool {
	if this.root.find(value) { //查找
		return false
	}

	this.root = this.root.insert(&value, nil)
	return true
}

func (this *RedBlackTree) Remove(value int) bool {
	return false
}


