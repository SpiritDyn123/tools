package sort

import (
	"reflect"
)

/*
	插入排序
	希尔排序：分段插入排序
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


type ShellSort struct {
	baseSort
}

func (this *ShellSort) Sort(values interface{}, cf CompareFunc) error  {
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

