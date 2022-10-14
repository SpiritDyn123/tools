package sort

import (
	"errors"
	"reflect"
)

type CompareFunc func(a, b interface{}) bool
func IntCompareFunc(a, b interface{}) bool {
	va := a.(int)
	vb := b.(int)
	return va < vb
}

func UintCompareFunc(a, b interface{}) bool {
	va := a.(uint)
	vb := b.(uint)
	return va < vb
}

func Int32CompareFunc(a, b interface{}) bool {
	va := a.(int32)
	vb := b.(int32)
	return va < vb
}

func Uint32CompareFunc(a, b interface{}) bool {
	va := a.(uint32)
	vb := b.(uint32)
	return va < vb
}

func Int64CompareFunc(a, b interface{}) bool {
	va := a.(int64)
	vb := b.(int64)
	return va < vb
}

func Uint64CompareFunc(a, b interface{}) bool {
	va := a.(uint64)
	vb := b.(uint64)
	return va < vb
}

func Int16CompareFunc(a, b interface{}) bool {
	va := a.(int16)
	vb := b.(int16)
	return va < vb
}

func Uint16CompareFunc(a, b interface{}) bool {
	va := a.(uint16)
	vb := b.(uint16)
	return va < vb
}

func Int8CompareFunc(a, b interface{}) bool {
	va := a.(int8)
	vb := b.(int8)
	return va < vb
}

func Uint8CompareFunc(a, b interface{}) bool {
	va := a.(uint8)
	vb := b.(uint8)
	return va < vb
}

func StringCompareFunc(a, b interface{}) bool {
	va := a.(string)
	vb := b.(string)
	return va < vb
}

type ISort interface{
	Sort(values interface{}, cf CompareFunc)
}

type baseSort struct {

}

func(this *baseSort) Sort(values interface{}, cf CompareFunc) error {
	rt := reflect.TypeOf(values)
	if rt.Kind() != reflect.Slice {
		return errors.New("values must be slice")
	}

	if rt.Elem().Kind() == reflect.Interface {
		return errors.New("values elem must be not interface")
	}

	return nil
}