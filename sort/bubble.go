package sort

import (
	"reflect"
)

/*
	冒泡排序
*/

type BubbleSort struct {
	baseSort
}

func (this *BubbleSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	sorted := false
	for i := 0;i < rv.Len() - 1;i++ {
		for j := 1;j < rv.Len() - i;j++ {
			if cf(rv.Index(j).Interface(), rv.Index(j-1).Interface()) {
				tmp_v := rv.Index(j).Interface()
				rv.Index(j).Set(rv.Index(j-1))
				rv.Index(j-1).Set(reflect.ValueOf(tmp_v))
				sorted = true
			}
		}

		if !sorted { //说明已经排序结束
			break
		}
	}
	return nil
}

