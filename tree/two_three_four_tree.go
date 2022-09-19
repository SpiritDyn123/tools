package tree

import (
	"fmt"
)

/*
	2-3-4树
*/

const (
	two_three_four_max = 3
	print_prefix = "   " //可选\t
)
type TwoThreeFourNode struct {
	Values []int
	Parent *TwoThreeFourNode
	SubNodes []*TwoThreeFourNode
}

func (this *TwoThreeFourNode) IsNil() bool {
	return this == nil || len(this.Values) == 0
}

func (this *TwoThreeFourNode) IsLeaf() bool {
	return this == nil || len(this.SubNodes) == 0
}

func (this *TwoThreeFourNode) Count() int {
	if this == nil {
		return 0
	}
	return len(this.Values)
}

func (this *TwoThreeFourNode) isTwo() bool {
	return this != nil && len(this.Values) == 1
}

func (this *TwoThreeFourNode) isThree() bool {
	return this != nil && len(this.Values) == 2
}

func (this *TwoThreeFourNode) isFour() bool {
	return this != nil && len(this.Values) == 3
}

func (this *TwoThreeFourNode) IsFull() bool {
	if this == nil {
		return false
	}
	return len(this.Values) == two_three_four_max
}

func (this *TwoThreeFourNode) ChildCount() int {
	if this == nil {
		return 0
	}
	return len(this.SubNodes)
}

func (this *TwoThreeFourNode) IndexSubNode(subNode *TwoThreeFourNode) int {
	if this == nil {
		return -1
	}

	for i, snode := range this.SubNodes {
		if snode == subNode {
			return i
		}
	}

	return -1
}

func (this *TwoThreeFourNode) getValueIndex(value int) int {
	index := 0
	for ; index < len(this.Values);index++ {
		if value < this.Values[index] {
			break
		}
	}

	return index
}
func (this *TwoThreeFourNode) getSubIndex(subNode *TwoThreeFourNode) int {
	for i, tSubNode := range this.SubNodes {
		if tSubNode == subNode {
			return i
		}
	}

	return -1
}
func (this *TwoThreeFourNode) getSubNode(index int) (subNode *TwoThreeFourNode) {
	if this.IsNil() {
		return
	}

	if index >= len(this.SubNodes) {
		return
	}

	return this.SubNodes[index]
}

func (this *TwoThreeFourNode) Height() int {
	if this == nil {
		return 0
	}

	if len(this.SubNodes) == 0 {
		return 1
	}

	return 1 + this.SubNodes[0].Height()
}
/*
Root
    - Child1
        - Grandchild1
        - Grandchild2
    - Child2
        - Grandchild3
        - Grandchild4
*/

func (this *TwoThreeFourNode) print(index, ceng int, prefix string) {
	Print := func(line string) {
		p1 := "-"
		if prefix == "" {
			p1 = ""
		}
		fmt.Printf("%s%s %s序%d层%d\n",prefix, p1, line, index, ceng)
	}

	if this.IsNil() {
		Print("[NIL TREE]")
		return
	}

	vcnt := len(this.Values)
	if vcnt == 0 || vcnt > two_three_four_max || (len(this.SubNodes) != 0 && vcnt != len(this.SubNodes) - 1) {
		panic(fmt.Sprintf("print find invalid node:%v", this))
	}

	Print(fmt.Sprintf("%v", this.Values))
	for i, subNode := range this.SubNodes {
		subNode.print(i, ceng + 1, prefix + print_prefix)
	}
}

func (this *TwoThreeFourNode) printTree() {

}

func (this *TwoThreeFourNode) addValue(value int) {
	if len(this.Values) == 0 {
		this.Values = []int{ value }
		return
	}

	index := this.getValueIndex(value)
	old_values := this.Values
	this.Values = make([]int, len(this.Values)+1)

	copy(this.Values, old_values[:index])
	copy(this.Values[index+1:], old_values[index:])
	this.Values[index] = value
}

func (this *TwoThreeFourNode) removeSubNode(index int) {
	this.SubNodes = append(this.SubNodes[:index], this.SubNodes[index+1:]...)
}

func (this *TwoThreeFourNode) addSubNode(index int, subNode *TwoThreeFourNode) {
	subNode.Parent = this

	if len(this.SubNodes) == 0 {
		this.SubNodes = []*TwoThreeFourNode{ subNode }
		return
	}

	oldNodes := this.SubNodes
	this.SubNodes = make([]*TwoThreeFourNode, len(this.SubNodes)+1)

	copy(this.SubNodes, oldNodes[:index])
	copy(this.SubNodes[index+1:], oldNodes[index:])
	this.SubNodes[index] = subNode
}

func (this *TwoThreeFourNode) find(value int) (fcnt int) {
	if this.IsNil() {
		return -1
	}

	curNode := this
	for curNode != nil {
		fcnt++
		var index int
		for ; index < curNode.Count();index++ {
			v := curNode.Values[index]
			if v == value {
				return
			} else if value < v {
				break
			}
		}

		curNode = curNode.getSubNode(index)
	}
	fcnt = -1
	return
}

//分裂的节点
func (this *TwoThreeFourNode) split() (newNode *TwoThreeFourNode){
	if this == nil {
		return
	}

	//必须是满节点
	vcnt := len(this.Values)
	if vcnt != two_three_four_max {
		return
	}

	//分裂的节点一定只有一个值
	newNode = &TwoThreeFourNode{
		Values: []int{ this.Values[1] },
	}

	leftNode := &TwoThreeFourNode{
		Values: []int{ this.Values[0] },
		Parent: newNode,
	}

	rightNode := &TwoThreeFourNode{
		Values: []int{ this.Values[2] },
		Parent: newNode,
	}

	var tmpNode *TwoThreeFourNode
	vcnt = len(this.SubNodes)
	for i := 0;i < vcnt;i++ {
		tmpNode = rightNode
		if i < vcnt / 2 {
			tmpNode = leftNode
		}

		this.SubNodes[i].Parent = tmpNode
		tmpNode.SubNodes = append(tmpNode.SubNodes, this.SubNodes[i])
	}

	newNode.SubNodes = append(newNode.SubNodes, leftNode, rightNode)
	return
}

//合并的节点
func (this *TwoThreeFourNode) merge(oldNode, newNode *TwoThreeFourNode) (root *TwoThreeFourNode, err error) {
	if this.IsNil() || newNode.IsNil() {
		root = newNode
		return
	}

	if newNode.Count() != 1 {
		err = fmt.Errorf("TwoThreeFourNode::merge newNode must be 2 point")
		return
	}

	if this.Count() >= two_three_four_max {
		err = fmt.Errorf("TwoThreeFourNode::merge rootNode must not be full point")
		return
	}

	subIndex := this.IndexSubNode(oldNode)
	if subIndex < 0 {
		err = fmt.Errorf("TwoThreeFourNode::merge oldNode must be subNode")
		return
	}

	//插值
	this.addValue(newNode.Values[0])

	//插子节点
	this.removeSubNode(subIndex)
	for i, subNode := range newNode.SubNodes {
		this.addSubNode(subIndex+i, subNode)
	}

	root = this
	return
}

/*
	1、从根节点往下找叶子节点插入
	2、遇到满（4）节点就split增加了height 往上合并merge降低层数，从而保证层数不变
	3、2的操作可以保证满4节点向上移动，下次插入就会再次执行split操作
*/
func (this *TwoThreeFourNode) insert(value int) (root *TwoThreeFourNode) {
	if this == nil {
		this = &TwoThreeFourNode{
			Values: []int{ value },
		}
		root = this
		return
	}

	if this.find(value) >= 0 {
		return
	}


	curNode := this
	root = this
	var err error
	for {
		if curNode.IsFull() {
			newNode := curNode.split()
			parentNode := curNode.Parent
			if parentNode != nil {
				parentNode, err = parentNode.merge(curNode, newNode)
				if err != nil {
					panic(fmt.Errorf("%v", err))
				}
			} else {
				root = newNode
			}
			curNode = newNode
		} else if curNode.IsLeaf() { //叶子节点直接加上
			curNode.addValue(value)
			break
		} else {
			curNode = curNode.getSubNode(curNode.getValueIndex(value))
		}
	}

	return
}

/*
	1、节点的路径信息 path.length === tree.height
*/
func (this *TwoThreeFourNode) path(value int) (path []int) {
	curNode := this
	var findIndex = -1
	for curNode != nil {
		if findIndex >= 0 {
			path = append(path, 0)
			curNode = curNode.getSubNode(0)
 			continue
		}

		var index int
		for ;index < len(curNode.Values);index++ {
			v := curNode.Values[index]
			if value < v {
				break
			} else if value == v {
				findIndex = index
				break
			}
		}

		path = append(path, index)
		curNode = curNode.getSubNode(index)
	}

	return
}

/*
	让所在2节点转换成3，4节点，方便后续移动值到叶子节点删除
	1、2节点，合并右边2节点
	2、2节点，合并左边2节点
	3、2节点，偷取右边3，4节点一个值
	4、2节点，偷取左边3，4节点一个值
*/
func (this *TwoThreeFourNode) transfer() *TwoThreeFourNode {
	//必须是2节点 且
	if this.IsNil() || this.Parent == nil || !this.isTwo() {
		return this
	}

	pNode := this.Parent
	var rightNode, leftNode *TwoThreeFourNode
	index := pNode.getSubIndex(this) //最左边
	pKeyIndex := index
	dtype := 0 //对应上面四种情况
	if index == 0 {
		rightNode = pNode.getSubNode(index+1)
		leftNode = this
		if len(rightNode.Values) == 1 {
			dtype = 1
		} else {
			dtype = 3
		}
	} else if index == len(pNode.SubNodes) - 1 { //最右边
		rightNode = this
		leftNode = pNode.getSubNode(index-1)
		pKeyIndex -= 1
		if len(leftNode.Values) == 1 {
			dtype = 2
		} else {
			dtype = 4
		}
	} else { //中间
		lNode := pNode.getSubNode(index-1)
		rNode := pNode.getSubNode(index+1)
		lNodeLen := len(lNode.Values)
		rNodeLen := len(rNode.Values)
		if lNodeLen <= rNodeLen {
			leftNode = lNode
			rightNode = this
			pKeyIndex -= 1
			if lNodeLen == 1 {
				dtype = 2
			} else {
				dtype = 4
			}
		} else {
			leftNode = this
			rightNode = rNode
			if rNodeLen == 1 {
				dtype = 1
			} else {
				dtype = 3
			}
		}
	}

	var newNode *TwoThreeFourNode
	if dtype == 1 || dtype == 2 {
		pNodeLen := len(pNode.Values)
		if pNodeLen == 1 { //父亲是2节点
			pNode.addValue(leftNode.Values[0])
			pNode.addValue(rightNode.Values[0])
			newNode = pNode

			//合并子节点
			pNode.SubNodes = nil
			var i int
			leftNode.Parent = nil
			for _, subNode := range leftNode.SubNodes {
				pNode.addSubNode(i, subNode)
				i++
			}

			leftNode.Parent = nil
			for _, subNode := range rightNode.SubNodes {
				pNode.addSubNode(i, subNode)
				i++
			}
		} else { //父节点偷取一个节点,合并左右节点
			pValue := pNode.Values[pKeyIndex]
			pNode.Values = append(pNode.Values[:pKeyIndex], pNode.Values[pKeyIndex+1:]...)
			pNode.removeSubNode(pKeyIndex) //删除两个节点
			pNode.removeSubNode(pKeyIndex)

			//合并成一个4节点
			newNode = &TwoThreeFourNode{}
			newNode.addValue(pValue)
			newNode.addValue(leftNode.Values[0])
			newNode.addValue(rightNode.Values[0])
			var i int
			for _, subNode := range leftNode.SubNodes {
				newNode.addSubNode(i, subNode)
				i++
			}

			for _, subNode := range rightNode.SubNodes {
				newNode.addSubNode(i, subNode)
				i++
			}

			//挂载新节点
			pNode.addSubNode(pKeyIndex, newNode)
		}
	} else if dtype == 3 {
		newNode = this
		pValue := pNode.Values[pKeyIndex]
		pNode.Values[pKeyIndex] = rightNode.Values[0]
		rightNode.Values = rightNode.Values[1:]

		leftNode.addValue(pValue)
		if len(rightNode.SubNodes) > 0 {
			leftNode.addSubNode(len(leftNode.SubNodes), rightNode.getSubNode(0))
			rightNode.SubNodes = rightNode.SubNodes[1:]
		}
	} else if dtype == 4 {
		newNode = this

		pValue := pNode.Values[pKeyIndex]
		pNode.Values[pKeyIndex] = leftNode.Values[len(leftNode.Values) - 1]
		leftNode.Values = leftNode.Values[:len(leftNode.Values) - 1]

		rightNode.addValue(pValue)
		if len(leftNode.SubNodes) > 0 {
			rightNode.addSubNode(0, leftNode.getSubNode(len(leftNode.SubNodes) - 1))
			leftNode.SubNodes = leftNode.SubNodes[:len(leftNode.SubNodes) - 1]
		}
	} else {
		panic("transfer invalid dtype")
	}

	return newNode
}

func (this *TwoThreeFourNode) transferValue(value int) {
	path := this.path(value)
	curNode := this

	for i := 0;i < len(path);i++{
		subIndex := path[i]
		nextNode := curNode.getSubNode(subIndex) //先获取，不然transfer会导致路径失效

		curNode.transfer()

		curNode = nextNode
		//fmt.Println("--------------------------------")
		//this.print(0,1, "")
	}
}
func (this *TwoThreeFourNode) remove(value int) bool {
	if this.find(value) < 0 {
		return false
	}

	//WRONG!!!
	//this.transferValue(value)
	//path := this.path(value)
	//curNode := this
	//for i := 0;i < len(path); i++{
	//	subIndex := path[i]
	//	nextNode := curNode.getSubNode(subIndex)
	//	var findIndex  = -1
	//	for vi, v := range curNode.Values {
	//		if value == v {
	//			findIndex = vi
	//			break
	//		}
	//	}
	//
	//	if findIndex >= 0 {
	//		if curNode.IsLeaf() { // 在叶子节点上直接删除
	//			curNode.Values = append(curNode.Values[:findIndex], curNode.Values[findIndex+1:]...)
	//			return true
	//		}
	//
	//		//往下一层不停的替换
	//		valueIndex := len(nextNode.Values) - 1
	//		curNode.Values[findIndex] = nextNode.Values[valueIndex]
	//		nextNode.Values = append(nextNode.Values[:valueIndex], nextNode.Values[valueIndex+1:]...)
	//		nextNode.addValue(value)
	//	}
	//
	//	curNode = nextNode
	//}


	//替换到叶子上
	curNode := this
	var dstNode *TwoThreeFourNode
	var dstIndex int
	for curNode != nil {
		if curNode.IsLeaf() && dstNode != nil {
			ncount := curNode.Count()
			dstNode.Values[dstIndex] = curNode.Values[ncount - 1]
			curNode.Values[ncount - 1] = value
			break
		} else {
			var findex = curNode.Count()
			if dstNode == nil {
				for i, v := range curNode.Values {
					if value == v { //找到了
						dstNode = curNode
						dstIndex = i
						findex = i
						break
					} else if value < v{
						findex = i
						break
					}
				}
			}

			curNode = curNode.getSubNode(findex)
		}
	}

	//删除
	curNode = this
	var b_node bool
	for curNode != nil {
		b_node := !b_node && dstNode != nil && curNode == dstNode
		if curNode.isTwo() {
			curNode = curNode.transfer() //转换
		}

		findIndex := curNode.Count()
		for i, v := range curNode.Values {
			if value == v { //找到了
				findIndex = i
				if curNode.IsLeaf() {
					curNode.Values = append(curNode.Values[:i], curNode.Values[i+1:]...)
					return true
				}
			} else if value < v{
				findIndex = i
				break
			}
		}

		if b_node {
			findIndex -= 1
		}
		curNode = curNode.getSubNode(findIndex)
	}

	return false
}


type TwoThreeFourTree struct {
	Root *TwoThreeFourNode
}

func (this *TwoThreeFourTree) Height() int {
	return this.Root.Height()
}

func (this *TwoThreeFourTree) Find(value int) (fcnt int) {
	return this.Root.find(value)
}

func (this *TwoThreeFourTree) Insert(value int) {
	this.Root = this.Root.insert(value)
}

func (this *TwoThreeFourTree) Print() {
	this.Root.print(0,1, "")
}

func (this *TwoThreeFourTree) Path(value int) []int {
	return this.Root.path(value)
}

func (this *TwoThreeFourTree) TransValue(value int) {
	this.Root.transferValue(value)
}

func (this *TwoThreeFourTree) Remove(value int) bool {
	ok := this.Root.remove(value)
	if len(this.Root.Values) == 0 {
		this.Root = nil
	}
	return ok
}