package sort

import (
	"reflect"
)

/*
	快速排序:快速排序应该算是在冒泡排序基础上的递归分治法,适合顺序性较弱的随机数列
	方法是以第一个元素为基点,从最右边最左边j,i哨兵交换,相遇结束
*/

type QuickSort struct {
	baseSort
}

func (this *QuickSort) sort(arr reflect.Value, cf CompareFunc) {
	if arr.Len() <= 1 {
		return
	}

	var i  = 1
	var j = arr.Len() - 1
	tmp_v := arr.Index(0).Interface()
	for {
		for ;j != i;j-- {
			if cf(arr.Index(j).Interface(), tmp_v) {
				break
			}
		}

		for ;i != j;i++ {
			if cf(tmp_v, arr.Index(i).Interface()) {
				break
			}
		}

		if i != j {
			i_tmpv := arr.Index(i).Interface()
			arr.Index(i).Set(arr.Index(j))
			arr.Index(j).Set(reflect.ValueOf(i_tmpv))
		} else {
			if cf(arr.Index(i).Interface(), tmp_v) {
				arr.Index(0).Set(arr.Index(i))
				arr.Index(i).Set(reflect.ValueOf(tmp_v))
			}
			break
		}
	}

	this.sort(arr.Slice(0, i), cf)
	this.sort(arr.Slice(i, arr.Len()), cf)
}


//分治递归类型的冒泡算法
func (this *QuickSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	this.sort(rv, cf)

	return nil
}
