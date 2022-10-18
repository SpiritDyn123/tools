package sort

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type S struct {
	v int
}

func(s *S) String() string {
	return fmt.Sprintf("&{%d}", s.v)
}

func (s S) Value() int {
	return s.v
}

var arr = []int{3,5,1,4,2,7,6,11,1,4}
func Shuffle(arr []int) {
	//重新打乱
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

//go test -v -run bubble .
func Test_bubble(t *testing.T) {
	s := &BubbleSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

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

	Shuffle(arr)

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

//go test -v -run quick .
func Test_quick(t *testing.T) {
	s := &QuickSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

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

	t.Log("test quick success")
}

//go test -v -run insert .
func Test_insert(t *testing.T) {
	s := &InsertSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

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

//go test -v -run shell .
func Test_shell(t *testing.T) {
	s := &ShellSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

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

	t.Log("test shell success")
}


//go test -v -run heap .
func Test_heap(t *testing.T) {
	s := &HeapSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

	arr_s := []*S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, &S{arr[i] })
	}
	rand.Seed(time.Now().UnixNano())
	arr_s = append(arr_s, &S{arr[rand.Intn(len(arr))]})
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(*S)
		vb := b.(*S)
		return va.v > vb.v
	})
	fmt.Println(arr_s, err)

	t.Log("test heap success")
}

//go test -v -run count .
func Test_count(t *testing.T) {
	s := &CountSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

	arr_s := []S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(S)
		vb := b.(S)
		return va.v < vb.v
	})
	fmt.Printf("%+v, %p, %v\n", arr_s, &arr_s, err)

	t.Log("test count success")
}


//go test -v -run bucket .
func Test_bucket(t *testing.T) {
	s := &BucketSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

	arr_s := []S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(S)
		vb := b.(S)
		return va.v < vb.v
	})
	fmt.Printf("%+v, %p, %v\n", arr_s, &arr_s, err)

	t.Log("test bucket success")
}

//go test -v -run radix .
func Test_radix(t *testing.T) {
	arr := []int{1,1, -1, -2, -10, 22, 6, 317, 8}
	s := &RadixSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

	arr_s := []S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(S)
		vb := b.(S)
		return va.v < vb.v
	})
	fmt.Printf("%+v, %p, %v\n", arr_s, &arr_s, err)

	t.Log("test radix success")
}

//go test -v -run merge .
func Test_merge(t *testing.T) {
	arr := []int{1,1, -1, -2, -10, 22, 6, 317, 8}
	s := &MergeSort{}
	err := s.Sort(arr, IntCompareFunc)
	fmt.Println(arr, err)

	Shuffle(arr)

	arr_s := []S{}
	for i := 0;i < len(arr);i++ {
		arr_s = append(arr_s, S{arr[i] })
	}
	err = s.Sort(arr_s, func(a, b interface{}) bool {
		va := a.(S)
		vb := b.(S)
		return va.v < vb.v
	})
	fmt.Printf("%+v, %p, %v\n", arr_s, &arr_s, err)

	t.Log("test merge success")
}








