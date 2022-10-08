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
 - https://zhuanlan.zhihu.com/p/269069974 (用234树来对比操作，比较初级，方便理解)
 - https://blog.csdn.net/weixin_29163797/article/details/113332639
 - https://blog.csdn.net/cy973071263/article/details/122543826(注: 史上最全，但是没看，不太容易理解)
 - 工具网站：https://www.cs.usfca.edu/~galles/visualization/RedBlack.html

 个人理解：
	1、 红黑树一组及节点操作 等同于 234树
		父（黑）+ 左儿子（红） + 右儿子(红)  == 4节点
		父（黑）+ 左儿子（红）|| 父（黑）+ 右儿子（红）  == 3节点
		父（黑）  == 2节点
	2、红黑树的insert remove操作原理也是同234树基本一致
	3、红黑树的旋转、着色过程 == 234树的分离（transfer)和合并(merge)
*/

type rbColor byte

const (
	rbColor_black = rbColor(0) + iota
	rbColor_red
)

type redBlackNode struct {
	parent, left, right *redBlackNode
	color               rbColor //颜色
	value               int     //值
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

func (this *redBlackNode) String() string {
	str := ""
	if this.left != nil {
		str += fmt.Sprintf("%s<-", this.left.valueWithColor())
	}

	str += fmt.Sprintf("%s", this.valueWithColor())

	if this.right != nil {
		str += fmt.Sprintf("->%s", this.right.valueWithColor())
	}

	return str
}
func (this *redBlackNode) print(pos string, ceng int, prefix string) {
	if this == nil {
		fmt.Println("[NIL TREE]")
		return
	}

	fmt.Printf("%s- %s 层%d %s\n", prefix, this.valueWithColor(), ceng, pos)
	if this.left != nil {
		this.left.print("左", ceng+1, prefix+print_prefix)
	}

	if this.right != nil {
		this.right.print("右", ceng+1, prefix+print_prefix)
	}
}

/*
	1、单黑色节点（类似234树2节点）
	2、黑色节点+单红色子节点（左右）（类似234树3节点）
	3、黑色节点+2红色子节点（类似234树4节点）
*/
func (this *redBlackNode) case2() bool { //2节点
	return this != nil && this.color == rbColor_black &&
		( this.left == nil && this.right == nil ||
		  this.left != nil && this.left.color == rbColor_black && this.right != nil && this.right.color == rbColor_black )
}

func (this *redBlackNode) case3() bool { //3节点
	return this != nil &&
		this.color == rbColor_black &&
		(this.left != nil && this.left.color == rbColor_red && (this.right == nil || this.right.color == rbColor_black) ||
			this.right != nil && this.right.color == rbColor_red && (this.left == nil || this.left.color == rbColor_black))
}

func (this *redBlackNode) case4() bool { ////4节点
	return this != nil &&
		this.color == rbColor_black &&
		this.left != nil && this.left.color == rbColor_red &&
		this.right != nil && this.right.color == rbColor_red
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

//设置位置和着色
func (this *redBlackNode) setPosColor(color rbColor, cdiff bool) {
	if this == nil {
		return
	}

	//设置父亲
	if this.left != nil {
		this.left.parent = this
	}
	if this.right != nil {
		this.right.parent = this
	}

	//着色
	this.color = color
	ccolor := color
	if cdiff {
		if color == rbColor_black {
			ccolor = rbColor_red
		} else {
			ccolor = rbColor_black
		}
	}
	if this.left != nil {
		this.left.color = ccolor
	}
	if this.right != nil {
		this.right.color = ccolor
	}
}

func (this *redBlackNode) ll() *redBlackNode {
	newNode := this.right
	this.right = newNode.left
	newNode.left = this

	//父亲
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

	return newNode
}

func (this *redBlackNode) rr() *redBlackNode {
	newNode := this.left
	this.left = newNode.right
	newNode.right = this

	//父亲
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

/*
	1、插入节点初始必须是红色
	2、父亲是黑色节点直接插入
	3、判断爷爷节点的case3 case4 旋转变色逻辑即可
	注意：
		1、一定是3层【爷父子】绑定做一次caseV判断

	====更精简易懂的实现=======
	！！！爷父子三层一次操作（换种理解：遇到黑色（除了最底下的父亲层）就操作）
 */
func (this *redBlackNode) insert2(value int, parent *redBlackNode) *redBlackNode {
	if this == nil {
		node := &redBlackNode{
			value: value,
			color: rbColor_red,
		}

		if parent == nil {
			node.color = rbColor_black
		}

		return node
	}

	b_left := false
	var newNode *redBlackNode
	var brother *redBlackNode
	if value < this.value {
		this.left = this.left.insert2(value, this)
		b_left = true
		newNode = this.left
		brother = this.right
	} else {
		this.right = this.right.insert2(value, this)
		newNode = this.right
		brother = this.left
	}


	//----------------------最核心-------------------------
	//插入节点为黑的时候跳过
	if newNode.isBlack() {
		return this
	}

	//遇到黑色,表示到爷爷层
	if !this.isBlack() {
		return this
	}

	//最下面一层 父亲层刚好是黑色 不作操作
	if newNode.value == value {
		return this
	}

	bb_left := false
	gsNode := newNode.right
	if value <= newNode.value {
		bb_left = true
		gsNode = newNode.left
	}

	if gsNode.color == rbColor_black {
		return this
	}
	//--------------------------------------------------


	if brother == nil || brother.isBlack() { //3节点
		if b_left {
			if bb_left {
				this = this.rr()
			} else {
				this = this.lr()
			}
		} else {
			if bb_left {
				this = this.rl()
			} else {
				this = this.ll()
			}
		}

		this.color = rbColor_black
		if this.left != nil {
			this.left.color = rbColor_red
		}

		if this.right != nil {
			this.right.color = rbColor_red
		}
	} else if brother.isRed() { //4节点 分裂
		//直接上色
		if this.parent == nil { //根节点
			this.color = rbColor_black
		}else {
			this.color = rbColor_red //变红就相当于234树向上merge的过程
		}

		newNode.color = rbColor_black
		brother.color = rbColor_black
	} else { // 异常情况
		panic(fmt.Sprintf("node:%v insert case err", brother))
	}

	return this
}

//删除
/*
		思路一：
		1、如果是红色，则直接删除，不用后续调整；
		2、黑色
			2.1、2节点情况
				- 父黑
					- 兄红 必存在双黑子 旋转
					- 兄黑
						- 不存在子节点，兄变红向上递归(最难情况)
						- 存在红子节点 旋转
				- 父红
					- 兄黑
						- 不存在子节点，父变黑兄变红
						- 存在红子节点 旋转
			2.2、3节点情况
				- 直接替换子节点删除
			2.3、4节点情况
				- 不存在


		思路二：
		1 如果是红色，则直接删除，不用后续调整；
		2 如果是黑色，则需要考虑其兄弟节点颜色，以及兄弟节点的儿子情况：
		  - 2.1 如果兄弟节点是红色，则要满足红黑树第5点，兄弟节点必有两个黑色的儿子，则修改兄弟节点的左儿子为红色，
			兄弟节点为黑色，对父节点左旋，调整完毕；
		  - 2.2 如果兄弟节点是黑色（如果有儿子，则一定是红色，黑色则不满足红黑树第5点）：
			  - 1 兄弟节点有一个右儿子：将父节点颜色给兄弟节点，修改父节点和兄弟右儿子节点为红色，对父节点左旋，调整完毕；
			  - 2 兄弟节点有一个左儿子，互换兄弟与其左儿子节点颜色，对兄弟节点右旋，此时和 2.2.1 一样，执行即可；
			  - 3 兄弟节点有两个儿子，无视兄弟左儿子节点，则该情况其实和 2.2.1 一样，执行 2.2.1 流程；
			  - 4 兄弟节点没有儿子，因为删除的节点为黑色，为了动态平衡，直接修改兄弟节点的颜色为红色，
				但此时可能不满足红黑树第4点（父节点可能是红色），因此，将待调整节点指为父节点，继续执行第 2 点；

*/

//思路一
func (this *redBlackNode) removeTrans1(node *redBlackNode, case_v int) *redBlackNode {
	//红色节点直接删除
	if node.isRed() {
		return this
	}

	parent := node.parent
	//跟节点
	if parent == nil {
		return nil
	}

	pparent := parent.parent
	b_left := node.value <= parent.value

	var newNode = parent
	switch case_v {
	case 2:
		brother := parent.left //兄弟节点必然为黑
		if b_left {
			brother = parent.right
		}
		bcase_v := brother.caseV()

		if parent.isBlack() { // 父亲节点是2节点
			if bcase_v == 2 { // 父黑兄黑，需要向上借（递归了）
				brother.color = rbColor_red

				//向上调整
				return this.removeTrans1(parent, 2)
			} else if bcase_v == 3 {
				//兄弟节点借一个
				if b_left {
					if brother.left == nil {
						newNode = parent.ll()
					} else {
						newNode = parent.rl()
					}
				} else {
					if brother.right == nil {
						newNode = parent.rr()
					} else {
						newNode = parent.lr()
					}
				}

				//着色
				newNode.color = rbColor_black
				if newNode.left != nil {
					newNode.left.color = rbColor_black
				}

				if newNode.right != nil {
					newNode.right.color = rbColor_black
				}

			} else if bcase_v == 4 { //
				if b_left { // 左旋转
					newNode = parent.ll()
				} else {
					newNode = parent.rr()
				}

				//着色
				newNode.color = rbColor_black
				if newNode.left != nil {
					newNode.left.color = rbColor_black
				}

				if newNode.right != nil {
					newNode.right.color = rbColor_black
				}

			} else { //父黑兄红（必有2黑子节点）
				if b_left { // 左旋转
					newNode = parent.ll()
				} else {
					newNode = parent.rr()
				}

				//着色
				newNode.color = rbColor_black
				if newNode.left != nil {
					//TODO 判断需要旋转
					newNode.left.color = rbColor_black
					if newNode.left.right != nil {
						newNode.left.right.color = rbColor_red
					}
				}

				if newNode.right != nil {
					newNode.right.color = rbColor_black
					if newNode.right.left != nil {
						newNode.right.left.color = rbColor_red
					}
				}
			}
		} else { //父节点为3,4节点
			if bcase_v == 2 {
				//合并成4节点
				parent.color = rbColor_black
				brother.color = rbColor_red
			} else if bcase_v == 3 {
				//兄弟节点借一个
				if b_left {
					if brother.left == nil {
						newNode = parent.ll()
					} else {
						newNode = parent.rl()
					}
				} else {
					if brother.right == nil {
						newNode = parent.rr()
					} else {
						newNode = parent.lr()
					}
				}

				//着色
				newNode.color = rbColor_red
				if newNode.left != nil {
					newNode.left.color = rbColor_black
				}

				if newNode.right != nil {
					newNode.right.color = rbColor_black
				}

			} else if bcase_v == 4 {
				if b_left { // 左旋转
					newNode = parent.ll()
				} else {
					newNode = parent.rr()
				}

				//着色
				newNode.color = rbColor_red
				if newNode.left != nil {
					newNode.left.color = rbColor_black
				}

				if newNode.right != nil {
					newNode.right.color = rbColor_black
				}
			} else {
				panic(fmt.Sprintf("brother:%d invalid case:%d", brother.value, bcase_v))
			}
		}
	case 3: //直接找到子节点替换删除
	default:
		panic(fmt.Sprintf("invalid case value:%d", case_v))
	}

	if pparent != nil {
		if newNode.value < pparent.value {
			pparent.left = newNode
		} else {
			pparent.right = newNode
		}
	}

	if newNode != nil && newNode.parent == nil { //已经成为根节点
		return newNode
	}

	return this
}

func (this *redBlackNode) remove(value int) *redBlackNode {
	curNode := this
	var dstNode *redBlackNode
	for curNode != nil {
		if value == curNode.value {
			dstNode = curNode
			curNode = curNode.left //取左子树最右节点
		} else if value < curNode.value {
			curNode = curNode.left
		} else {
			if dstNode != nil && curNode.left == nil {
				break
			}

			curNode = curNode.right
		}
	}

	if dstNode == nil {
		panic(fmt.Sprintf("value:%d not exist", value))
	}

	if curNode != nil { //交换到叶子节点上
		repValue := curNode.value
		curNode.value = dstNode.value
		dstNode.value = repValue
		curNode = dstNode
	} else { //直接就是叶子节点
		curNode = dstNode
	}

	//删除节点
	_delNode := func(node *redBlackNode) {
		parent := node.parent
		if parent == nil { //根节点
			return
		}

		//最多只可能有一个节点
		sunNode := node.left
		if sunNode == nil {
			sunNode = node.right
		}

		//用子节点替换当前点，修改颜色
		if sunNode != nil {
			sunNode.color = node.color
		}

		b_left := node == parent.left
		if b_left {
			parent.left = sunNode
		} else {
			parent.right = sunNode
		}
	}

	_delNode(curNode)
	case_v := curNode.caseV()
	return this.removeTrans1(curNode, case_v)
}

type RedBlackTree struct {
	root *redBlackNode
}

func (this *RedBlackTree) Print() {
	this.root.print("根", 1, "")
}

func (this *RedBlackTree) Find(value int) bool {
	return this.root.find(value)
}

func (this *RedBlackTree) Insert(value int) bool {
	if this.root.find(value) { //查找
		return false
	}

	this.root = this.root.insert2(value, nil)
	return true
}

func (this *RedBlackTree) Remove(value int) bool {
	if !this.root.find(value) {
		return false
	}

	this.root = this.root.remove(value)
	return true
}
