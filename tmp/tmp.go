package main

import "reflect"

func (this *CountSort) GetMaxValue(values interface{}, cf CompareFunc) interface{} {
	rv := reflect.ValueOf(values)
	m_value := reflect.New(rv.Elem().Type()).Interface()
	for i := 0;i < rv.Len();i++ {
		tmp_v := rv.Index(i).Interface()
		if cf(m_value, tmp_v) {
			m_value = tmp_v
		}
	}

	return m_value
}

func main() {

}
