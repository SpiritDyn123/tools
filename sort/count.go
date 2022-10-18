package sort

import (
	"fmt"
	"reflect"
)

/*
	计数排序:适合大量指定数字范围大量重复的排序
	桶排序: 优化的计数排序,结合简单的桶算法，分桶；计数排序就相当于一种特殊桶排序，桶数量是数组长度
	基数排序: 固定19个桶（正负数），把数字当成字符串，移位进行分桶
*/

type ICountValue interface {
	Value() int
}

type CountSort struct {
	baseSort
}

func (this *CountSort) GetMValue(values interface{}, cf CompareFunc, min bool) (ICountValue, error) {
	rv := reflect.ValueOf(values)
	r_mv :=  reflect.New(rv.Type().Elem()).Elem()
	if r_mv.Kind() == reflect.Ptr && r_mv.IsNil() && rv.Len() > 0 {
		r_mv.Set(reflect.ValueOf(rv.Index(0).Interface()))
	}
	m_value := r_mv.Interface()
	if _, ok := m_value.(ICountValue); !ok {
		return nil, fmt.Errorf("%s is not implete ICountValue", r_mv.Type())
	}

	for i := 0;i < rv.Len();i++ {
		tmp_v := rv.Index(i).Interface()
		if min {
			if cf(tmp_v, m_value) {
				m_value = tmp_v
			}
		} else {
			if cf(m_value, tmp_v) {
				m_value = tmp_v
			}
		}
	}

	return  m_value.(ICountValue), nil
}


func (this *CountSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	if rv.Len() <= 0 {
		return nil
	}

	min_value, err := this.GetMValue(values, cf, true)
	if err != nil {
		return err
	}

	max_value, err := this.GetMValue(values, cf, false)
	if err != nil {
		return err
	}

	slen := max_value.Value() - min_value.Value() + 1
	s_arr := reflect.MakeSlice(reflect.SliceOf(rv.Type()), slen, slen)
	for i := 0;i < rv.Len();i++ {
		index := rv.Index(i).Interface().(ICountValue).Value() - min_value.Value()
		i_arr := s_arr.Index(index)
		if i_arr.IsNil() {
			i_arr = reflect.New(rv.Type()).Elem()
		}
		i_arr = reflect.Append(i_arr, rv.Index(i))
		s_arr.Index(index).Set(i_arr)
	}

	//arr := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Cap())
	var ai int
	for i := 0;i < s_arr.Len();i++ {
		cur_arr := s_arr.Index(i)
		if cur_arr.IsNil() {
			continue
		}

		for j := 0;j < cur_arr.Len();j++ {
			//arr.Index(ai).Set(cur_arr.Index(j))
			rv.Index(ai).Set(cur_arr.Index(j))
			ai++
		}
	}

	//reflect.ValueOf(&values).Elem().Set(arr) //无用
	return nil
}

const (
	Default_bucket_count = 5
)

type BucketSort struct {
	CountSort

	Bucket_count int
}

func (this *BucketSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	if rv.Len() <= 0 {
		return nil
	}

	min_value, err := this.GetMValue(values, cf, true)
	if err != nil {
		return err
	}

	max_value, err := this.GetMValue(values, cf, false)
	if err != nil {
		return err
	}

	if this.Bucket_count <= 0 {
		this.Bucket_count = Default_bucket_count
	}

	slen := (max_value.Value() - min_value.Value()) / this.Bucket_count + 1 //分桶
	s_arr := reflect.MakeSlice(reflect.SliceOf(rv.Type()), slen, slen)
	for i := 0;i < rv.Len();i++ {
		cur_v := rv.Index(i).Interface().(ICountValue).Value()
		index := (cur_v - min_value.Value()) / this.Bucket_count
		i_arr := s_arr.Index(index)
		if i_arr.IsNil() {
			i_arr = reflect.New(rv.Type()).Elem()
		}
		i_arr = reflect.Append(i_arr, rv.Index(i))
		s_arr.Index(index).Set(i_arr)
	}

	//桶排序用插入排序
	i_s := InsertSort{}
	var ai int
	for i := 0;i < s_arr.Len();i++ {
		cur_arr := s_arr.Index(i)
		if cur_arr.IsNil() {
			continue
		}

		i_s.Sort(cur_arr.Interface(), cf)
		for j := 0;j < cur_arr.Len();j++ {
			rv.Index(ai).Set(cur_arr.Index(j))
			ai++
		}
	}

	return nil
}

//基数排序
type RadixSort struct {
	CountSort
}

func (this *RadixSort) Sort(values interface{}, cf CompareFunc) error  {
	err := this.baseSort.Sort(values, cf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(values)
	if rv.Len() <= 0 {
		return nil
	}

	f_v := rv.Index(0).Interface()
	if _, ok := f_v.(ICountValue); !ok {
		return fmt.Errorf("%s is not implete ICountValue", rv.Index(0).Type())
	}

	//0-8负数，9是0,10-18正数
	mod := 1
	pass := true
	for pass {
		pass = false
		s_arr := reflect.MakeSlice(reflect.SliceOf(rv.Type()), 19, 19)
		for i := 0;i < rv.Len();i++ {
			cur_v := rv.Index(i).Interface().(ICountValue).Value()
			bucket := (cur_v / mod) % 10 + 9
			if cur_v / (mod * 10) > 0 {
				pass = true
			}

			i_arr := s_arr.Index(bucket)
			if i_arr.IsNil() {
				i_arr = reflect.New(rv.Type()).Elem()
			}

			i_arr = reflect.Append(i_arr, rv.Index(i))
			s_arr.Index(bucket).Set(i_arr)
		}

		var ai int
		for i := 0;i < s_arr.Len();i++ {
			cur_arr := s_arr.Index(i)
			if cur_arr.IsNil() {
				continue
			}

			for j := 0;j < cur_arr.Len();j++ {
				rv.Index(ai).Set(cur_arr.Index(j))
				ai++
			}
		}

		mod *= 10
	}

	return nil
}

