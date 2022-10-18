package sort

import (
	"reflect"
)

/*
	插入排序：对于越有序的效果越好 n2
	希尔排序：分段插入排序,基于上次的有序分组，插入排序更有效 nlogn->n2
*/

type InsertSort struct {
	baseSort
}

func (this *InsertSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	for i := 1;i < rv.Len();i++ {
		tmp_v := rv.Index(i).Interface()
		ri := 0
		for j := i-1;j >= 0;j-- {
			if cf(rv.Index(j).Interface(), tmp_v) {
				ri = j+1
				break
			}

			rv.Index(j+1).Set(rv.Index(j)) //移动
		}

		if ri != i {
			rv.Index(ri).Set(reflect.ValueOf(tmp_v))
		}
	}

	return nil
}


//最常见增量设定算法：第一个增量取待排序记录个数的一半，然后逐次减半，最后一个增量为 1;gap=length/2, gap /= 2
type ShellSort struct {
	baseSort
}

func (this *ShellSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	gap := rv.Len() / 2
	for gap != 0 {
		for i := 0;i < gap;i++ {
			//执行插入排序 深入仔细的看i,j,k
			for j := i+gap;j < rv.Len();j+=gap {
				tmp_v := rv.Index(j).Interface()
				rj := i
				for k := j - gap;k >= i;k -= gap { //往前走一格
					if cf(rv.Index(k).Interface(), tmp_v) {
						rj = k+gap
						break
					}

					rv.Index(k+gap).Set(rv.Index(k))
				}


				if rj != j {
					rv.Index(rj).Set(reflect.ValueOf(tmp_v))
				}
			}
		}

		gap /= 2
	}

	return nil
}

