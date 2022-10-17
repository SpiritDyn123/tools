package sort

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



	return nil
}
