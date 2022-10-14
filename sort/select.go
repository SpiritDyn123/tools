package sort

import "reflect"

/*
	选择排序
*/

type SelectSort struct {
	baseSort
}

func (this *SelectSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	for i := 0;i < rv.Len() - 1;i++ {
		ri := i
		for j := i+1;j < rv.Len();j++ {
			if cf(rv.Index(j).Interface(), rv.Index(ri).Interface()) {
				ri = j
			}
		}

		if i != ri {
			tmp_v := rv.Index(ri).Interface()
			rv.Index(ri).Set(rv.Index(i))
			rv.Index(i).Set(reflect.ValueOf(tmp_v))
		}
	}
	return nil
}

