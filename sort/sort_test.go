package sort

import (
	"fmt"
	"testing"
)

type S struct {
	v int
}

func(s *S) String() string {
	return fmt.Sprintf("&{%d}", s.v)
}

var arr = []int{3,5,1,4,2,7,6,11}

//go test -v -run bubble .
func Test_bubble(t *testing.T) {
	s := &BubbleSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	arr_s := []*S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, &S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(*S)
		vb := b.(*S)
		return va.v < vb.v
	})
	fmt.Println(arr_s, err)
	t.Log("test bubble success")
}

//go test -v -run select .
func Test_select(t *testing.T) {
	s := &SelectSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	arr_s := []*S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, &S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(*S)
		vb := b.(*S)
		return va.v < vb.v
	})
	fmt.Println(arr_s, err)

	t.Log("test select success")
}

//go test -v -run insert .
func Test_insert(t *testing.T) {
	s := &InsertSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	arr_s := []*S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, &S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(*S)
		vb := b.(*S)
		return va.v < vb.v
	})
	fmt.Println(arr_s, err)

	t.Log("test insert success")
}


