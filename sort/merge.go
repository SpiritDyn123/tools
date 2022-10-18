package sort

import (
	"reflect"
)

/*
	归并排序:分治递归
*/

type MergeSort struct {
	baseSort
}

func (this *MergeSort) sort(arr reflect.Value, cf CompareFunc) {
	arr_l := arr.Len()
	if arr_l < 2 {
		return
	}

	//split
	sp_index := arr_l / 2
	this.sort(arr.Slice(0, sp_index), cf)
	this.sort(arr.Slice(sp_index, arr_l), cf)

	//merge
	arr_s := reflect.MakeSlice(arr.Type(), arr_l, arr_l)
	var i, j, k int
	for {
		condi := i < sp_index
		condj := j < (arr_l - sp_index)
		if condi && condj && cf(arr.Index(i).Interface(), arr.Index(j + sp_index).Interface()) || condi && !condj {
			arr_s.Index(k).Set(arr.Index(i))
			i++
			k++
		} else if condj {
			arr_s.Index(k).Set(arr.Index(j + sp_index))
			j++
			k++
		} else {
			break
		}
	}

	reflect.Copy(arr, arr_s)
}

func (this *MergeSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	this.sort(rv, cf)

	return nil
}

