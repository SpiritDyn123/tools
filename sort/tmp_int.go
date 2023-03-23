package main

import (
	"fmt"
	"math/rand"
	"time"
)

func shuffle(data []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
}

//堆排序
func headSort(arr []int, min bool) {
	if len(arr) == 0 {
		return
	}

	beginPos := len(arr) / 2 - 1
	for i := beginPos;i >= 0;i-- {
		value := arr[i]
		leftIndex := 2 * i + 1
		if leftIndex < len(arr) {
			value2 := arr[leftIndex]
			if min {
				if value2 < value {
					arr[i], arr[leftIndex] = value2, value
				}
			} else {
				if value2 > value {
					arr[i], arr[leftIndex] = value2, value
				}
			}
		}

		value = arr[i]
		rightIndex := leftIndex + 1
		if rightIndex < len(arr) {
			value2 := arr[rightIndex]
			if min {
				if value2 < value {
					arr[i], arr[rightIndex] = value2, value
				}
			} else {
				if value2 > value {
					arr[i], arr[rightIndex] = value2, value
				}
			}
		}
	}

	headSort(arr[1:], min)
}

//快速排序
func quickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	sv := arr[0]
	var si = 0
	i := 1
	j := len(arr) - 1
	for {
		for ;i != j;j-- {
			if arr[j] > sv {
				arr[si] = arr[j]
				si = j
				break
			}
		}

		for ;i != j;i++ {
			if arr[i] < sv {
				arr[si] = arr[i]
				si = i
				break
			}
		}

		if i == j {
			if si != i { //表示并未替换过
				if arr[i] > sv {
					arr[i], arr[si] = sv, arr[i]
				}
			} else {
				arr[si] = sv
			}
			break
		} else {
			//todo 如果不用si，可以食用 i,j互相替换
			//arr[i], arr[j] = arr[j], arr[i]
		}
	}

	quickSort(arr[:i])
	quickSort(arr[i:])
}

//选择排序（直接排序）
func zhijieSelectSort(arr []int) {
	for i := 0; i < len(arr) - 1;i++ {
		for j := i + 1; j < len(arr);j++ {
			if arr[j] < arr[i] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}


//冒泡排序
func bubbleSort(data []int) {
	for i := 0;i < len(data);i++ {
		sorted := false
		for j := 1; j < len(data) - i;j++ {
			if data[j-1] > data[j] {
				data[j-1], data[j] = data[j], data[j-1]
				sorted = true
			}
		}

		if !sorted {
			break
		}
	}
}

//插入排序
func insertSort(data []int) {
	for i := 1;i < len(data);i++ {
		var ri int
		var sv = data[i]
		for j := i-1;j >= 0;j-- {
			if sv >= data[j] {
				ri = j+1
				break
			}
			data[j+1] = data[j]
		}

		if ri != i {
			data[ri] = sv
		}
	}
}

//希尔排序
func shellSort(data []int, seed int) { //宏观的insertSort
	if seed <= 0 {
		seed = 2
	}

	bLen := len(data)
	for gap := bLen / seed;gap > 0; gap /= seed{
		//插入排序
		for i := 0;i < gap;i++ { //组数
			//按照gap分组执行插入排序
			beginOffset := i //方便理解
			for j := beginOffset + gap; j < bLen - (gap - 1 - beginOffset) ;j += gap {
				sv := data[j]
				ri := beginOffset
				for k := j - gap; k >= beginOffset;k -= gap {
					if sv > data[k] {
						ri = k + gap
						break
					}
					data[k+gap] = data[k]
				}

				if ri != j {
					data[ri] = sv
				}
			}
		}
	}
}

//归并排序
func mergeSort(data []int) {
	bLen := len(data)
	if bLen <= 1 {
		return
	}

	si := bLen / 2
	d1 := data[:si]
	d2 := data[si:]
	mergeSort(d1)
	mergeSort(d2)

	tmpData := make([]int, 0, bLen)
	var i, j int
	for ;i < len(d1) && j < len(d2);{
		if d1[i] <= d2[j] {
			tmpData = append(tmpData, d1[i])
			i++
		} else {
			tmpData = append(tmpData, d2[j])
			j++
		}
	}

	if i >= len(d1) {
		tmpData = append(tmpData, d2[j:]...)
	} else {
		tmpData = append(tmpData, d1[i:]...)
	}

	tmpData = append(tmpData, d2[j:]...)
	copy(data, tmpData)
}

func main() {
	bLen := 50
	data := make([]int, 0, bLen)
	for i := 0; i < bLen;i++ {
		data = append(data, i)
	}

	shuffle(data)
	fmt.Println("sort before:", data)
	shellSort(data, 2)
	fmt.Println("sorted:", data)


}
