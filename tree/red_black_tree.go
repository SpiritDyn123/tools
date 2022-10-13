package tree

import (
	"fmt"
	"strconv"
)

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

func (this *redBlackNode) isLeaf() bool {
	return this.left == nil && this.right == nil
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

	parent := newNode.parent
	if parent != nil {
		if newNode.value <= parent.value {
			parent.left = newNode
		} else {
			parent.right = newNode
		}
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

	parent := newNode.parent
	if parent != nil {
		if newNode.value <= parent.value {
			parent.left = newNode
		} else {
			parent.right = newNode
		}
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
	图解：./img/red_blue_insert.png
 */
func (this *redBlackNode) insert2(value int, parent *redBlackNode) *redBlackNode {
	if this == nil {
		node := &redBlackNode{
			value: value,
			color: rbColor_red,
			parent: parent,
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
	//1、插入节点为黑的时候跳过
	if newNode.isBlack() {
		return this
	}

	//2、遇到黑色,表示到爷爷层
	if !this.isBlack() {
		return this
	}

	//3、最下面一层 父亲层刚好是黑色 不作操作
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
	1 如果是红色，则直接删除，不用后续调整；
	2 如果是黑色，则需要考虑其兄弟节点颜色，以及兄弟节点的儿子情况：
	  - 2.1 如果兄弟节点是红色，则要满足红黑树第5点，兄弟节点必有两个黑色的儿子，则修改兄弟节点的左儿子为红色，
		兄弟节点为黑色，对父节点左旋，修改父亲节点为红色 然后继续2.2的逻辑判断（旋转完毕，可能兄弟节点变黑了)
	  - 2.2 如果兄弟节点是黑色（如果有儿子，则一定是红色，黑色则不满足红黑树第5点）：
		  - 1 兄弟节点有一个右儿子：将父节点颜色给兄弟节点，修改父节点和兄弟右儿子节点为红色，对父节点左旋，调整完毕；
		  - 2 兄弟节点有一个左儿子，互换兄弟与其左儿子节点颜色，对兄弟节点右旋，此时和 2.2.1 一样，执行即可；
		  - 3 兄弟节点有两个儿子，无视兄弟左儿子节点，则该情况其实和 2.2.1 一样，执行 2.2.1 流程；
		  - 4 兄弟节点没有儿子，因为删除的节点为黑色，为了动态平衡，直接修改兄弟节点的颜色为红色，
			- 4.1 父亲为红色，直接改父亲为黑色
			- 4.2 父亲为黑色，以父亲节点替代当前节点 递归重复2的处理
	图解：./img/red_blue_remove.xmind 需要安装xmind
*/
//检查一下
func(this *redBlackNode) checkLeaf(pNode *redBlackNode, case_v int) bool {
	if this == nil {
		return true
	}

	if this.isLeaf() {
		if this.parent == nil {
			return true
		}

		b_left := this == this.parent.left
		var brother = this.parent.left
		if b_left {
			brother = this.parent.right
		}

		hit := brother == nil && this.isBlack() ||
			brother != nil && brother.isLeaf() && this.color != brother.color ||
			brother != nil && !brother.isLeaf() && !this.isBlack()

		if hit {
			fmt.Println("checkLeaf：", this, this.parent)
			pNode.print("checkLeaf"+ strconv.Itoa(case_v), 0, "\t")
		}

		return !hit
	}

	if !this.left.checkLeaf(pNode, case_v) {
		return false
	}

	if !this.right.checkLeaf(pNode, case_v) {
		return false
	}

	return true
}

//234树删除法
func (this *redBlackNode) removeTransAs234(curNode *redBlackNode, del bool) (newRoot, delNode *redBlackNode, oprs []string) {
	delNode = curNode
	newRoot = this
	parent := curNode.parent //父亲节点
	if parent == nil { //到达根节点
		newRoot = curNode
		return
	}

	//测试
	//if del && curNode.value == 617 {
	//	fmt.Println(curNode)
	//	curNode.parent.parent.print("BEFORE_DEL:", 0, "\t")
	//}

	if curNode.isRed() {//红色
		oprs = []string{"红色直接删除"}
	} else if curNode.left != nil && del {//存在子节点 交换一下位置即可
		tmp_v := curNode.value
		curNode.value = curNode.left.value
		curNode.left.value = tmp_v
		curNode = curNode.left
		oprs = []string{"替换左节点"}
	} else if curNode.right != nil && del {//存在子节点 交换一下位置即可
		tmp_v := curNode.value
		curNode.value = curNode.right.value
		curNode.right.value = tmp_v
		curNode = curNode.right
		oprs = []string{"替换右边节点"}
	} else {
		brother := parent.right //兄弟节点
		b_left := true
		if curNode == parent.right {
			b_left = false
			brother = parent.left
		}

		//从父亲节点看234节点
		/*
			1、兄弟为红色表示父兄组成3节点， 要把整个兄弟子树（一颗234节点）旋转过来，需要继续递归转换
			2、兄弟为黑色，兄弟本身就是234树一个节点，旋转后只是相当于结了一个值过来，很好处理，
				但是考虑到兄弟如果没有子树（2节点），就只能网上找父节点借（父为红），或者合并（父为黑），
				合并就会导致树降低，需要继续递归转换
		*/
		if brother.isRed() { //必然存在双黑子 父亲和兄节点组成3节点(必为黑)，双黑子为2节点，参照234树，合并一个新的4节点
			brother.color = rbColor_black
			parent.color = rbColor_red
			if b_left {
				newRoot = parent.ll()
				oprs = []string{"兄红左旋"}
			} else {
				newRoot = parent.rr()
				oprs = []string{"兄红右旋"}
			}

			//判断转移后的情况，可能需要继续旋转变色
			var c_oprs []string
			_, _, c_oprs = this.removeTransAs234(curNode, false)
			oprs = append(oprs, c_oprs...)
		} else {
			if brother.left != nil && brother.left.isRed() &&
				brother.right != nil && brother.right.isRed() { //双红子 兄弟为4节点,参照234树，借一个
				brother.color = parent.color
				parent.color = rbColor_black
				if b_left {
					brother.right.color = rbColor_black
					newRoot = parent.ll()
					oprs = []string{"兄黑子满左旋"}
				} else {
					brother.left.color = rbColor_black
					newRoot = parent.rr()
					oprs = []string{"兄黑子满右旋"}
				}
			} else if brother.left != nil && brother.left.isRed() { //单左 兄弟为3节点,参照234树，借一个
				pcolor := parent.color
				parent.color = rbColor_black
				if b_left {
					brother.left.color = pcolor
					newRoot = parent.rl()
					oprs = []string{"兄黑左子右左旋"}
				} else {
					brother.color = pcolor
					brother.left.color = rbColor_black
					newRoot = parent.rr()
					oprs = []string{"兄黑左子右旋"}
				}
			} else if brother.right != nil && brother.right.isRed() {//单右 兄弟为3节点,参照234树，借一个
				pcolor := parent.color
				parent.color = rbColor_black
				if b_left {
					brother.color = pcolor
					brother.right.color = rbColor_black
					newRoot = parent.ll()
					oprs = []string{"兄黑右子左旋"}
				} else {
					brother.right.color = pcolor
					newRoot = parent.lr()
					oprs = []string{"兄黑右子左右旋"}
				}
			} else { //不存在子节点
				brother.color = rbColor_red
				if parent.isRed() { //从父节点借来
					parent.color = rbColor_black
					oprs = []string{"兄黑无子借父红"}
				} else { //黑色相当于234树减层了，要向上递归重复上面的处理
					oprs = []string{"兄黑无子父黑递归"}
					var c_oprs []string
					newRoot, _, c_oprs =  this.removeTransAs234(parent, false) //不接收返回值
					oprs = append(oprs, c_oprs...)
				}
			}
		}
	}

	if del {
		curNode.color = rbColor_red
	}

	//测试
	//if del {
	//	////判断一下
	//	if !newRoot.checkLeaf(newRoot, case_v) {
	//		fmt.Println("DEL_NODE:", curNode, del)
	//		panic("checkLeaf error")
	//	}
	//}


	return newRoot, curNode, oprs
}

//正常逻辑删除
func (this *redBlackNode) removeTrans(curNode *redBlackNode, del bool) (newRoot, delNode *redBlackNode, oprs []string) {
	delNode = curNode
	parent := curNode.parent
	if curNode == nil || parent == nil {
		return
	}

	//一定要让curNode.color变成红色，才能删除
	if curNode.isRed() {
		oprs = []string{"红色直接删除"}
	} else if del && (curNode.left != nil || curNode.right != nil) {
		repNode := curNode.left
		b_left := true
		if repNode == nil  {
			b_left = false
			repNode = curNode.right
		}

		repNode.color = curNode.color
		if b_left {
			newRoot = curNode.rr()
			oprs = []string{"替换左节点"}
		} else {
			newRoot = curNode.ll()
			oprs = []string{"替换右节点"}
		}

	} else {
		b_left := curNode == parent.left
		brother := parent.left
		if b_left {
			brother = parent.right
		}

		if brother.isRed() {//兄弟节点是红色
			brother.color, parent.color = parent.color, brother.color //交换颜色 父必为黑色
			if b_left {
				newRoot = parent.ll()
				oprs = []string{"兄红左旋"}

			} else {
				newRoot = parent.rr()
				oprs = []string{"兄红右旋"}
			}

			//转换成兄弟节点为黑色情况，继续递归
			var c_oprs []string
			_, _, c_oprs = this.removeTrans(curNode, false)
			oprs = append(oprs, c_oprs...)
		} else  { //兄弟节点是黑色
			if brother.isLeaf() || brother.left != nil && brother.left.isBlack() && brother.right != nil && brother.right.isBlack() { //没有子节点
				brother.color = rbColor_red
				if parent.isRed() {
					parent.color = rbColor_black
					oprs = []string{"兄黑无子借父红"}
				} else {
					oprs = []string{"兄黑无子父黑递归"}
					var c_oprs []string
					newRoot, _, c_oprs = this.removeTrans(parent, false)
					oprs = append(oprs, c_oprs...)
				}
			} else {
				opr_desc := "兄黑"
				if b_left {
					if brother.right == nil || brother.right.isBlack() {
						brother = brother.rr()
						opr_desc += "左子右左旋"
					} else {
						opr_desc += "右子左旋"
					}

					brother.color = rbColor_black
					brother.right.color = rbColor_black
					newRoot = parent.ll()
				} else {
					if brother.left == nil || brother.left.isBlack(){
						brother = brother.ll()
						opr_desc += "右子左右旋"
					}else {
						opr_desc += "左子右旋"
					}

					brother.color = rbColor_black
					brother.left.color = rbColor_black

					newRoot = parent.rr()
				}

				oprs = []string{opr_desc}
				//变色 newRoot == brother
				newRoot.color, parent.color = parent.color, newRoot.color
			}
		}
	}

	//只有删除层才能变色
	if del {
		curNode.color = rbColor_red
	}

	delNode = curNode

	return
}

func (this *redBlackNode) remove(value int) *redBlackNode {
	//替换到叶子节点
	curNode := this
	var dstNode *redBlackNode
	for curNode != nil {
		if value == curNode.value {
			dstNode = curNode
			curNode = curNode.left //取左子树最右节点
		} else if value < curNode.value {
			curNode = curNode.left
		} else {
			if dstNode != nil && curNode.right == nil {
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
	} else { //直接就是叶子节点
		curNode = dstNode
	}

	newRoot, delNode, oprs := this.removeTrans(curNode, true)

	fmt.Printf("删除%d操作：%v\n", value, oprs)

	//删除节点
	if delNode.parent == nil { //直接就是根节点
		return nil
	}

	if delNode.parent.left == delNode {
		delNode.parent.left = nil
	} else {
		delNode.parent.right = nil
	}

	//因为涉及旋转可能会让this的root位置变化
	if newRoot != nil && newRoot.parent == nil {
		return newRoot
	}

	return this
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
