package sort

import "reflect"

/*
	堆排序
*/

type HeapSort struct {
	baseSort
}

func (this *HeapSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	__sort := func(arr reflect.Value) { //log(n)
		arr_l := arr.Len()
		for i := arr_l / 2 - 1;i >= 0;i-- {
			//root->left->right i->2i+1->2i+2
			sun_index := 2*i+1
			if sun_index < arr_l { //左子树
				tmp_v := arr.Index(i).Interface()
				if cf(arr.Index(sun_index).Interface(), tmp_v) {
					arr.Index(i).Set(arr.Index(sun_index))
					arr.Index(sun_index).Set(reflect.ValueOf(tmp_v))
				}
			}

			sun_index = 2*i+2
			if sun_index < arr_l {
				tmp_v := arr.Index(i).Interface()
				if cf(arr.Index(sun_index).Interface(), tmp_v) {
					arr.Index(i).Set(arr.Index(sun_index))
					arr.Index(sun_index).Set(reflect.ValueOf(tmp_v))
				}
			}
		}
	}

	rv := reflect.ValueOf(values)
	arr_l := rv.Len()
	for i := 0;i < arr_l;i++ { //n
		__sort(rv.Slice(i, arr_l))
	}

	return nil
}
