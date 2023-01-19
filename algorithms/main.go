package main

import "log"

// selectionSort 选择排序
// 每次选择一个最小值放到最左边
func selectionSort(arr []int) {
	if arr == nil || len(arr) < 2 {
		// 已经有序，不用排序
		return
	}
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		if i != minIndex {
			// 当前不是最小才需要交换
			swap(&arr[i], &arr[minIndex])
		}
	}
}

// bubbleSort 冒泡排序
func bubbleSort(arr []int) {
	if arr == nil || len(arr) < 2 {
		return
	}
	n := len(arr)
	swapped := false
	for end := n - 1; end > 0; end-- {
		for i := 0; i < end; i++ {
			if arr[i] > arr[i+1] {
				swapped = true
				swap(&arr[i], &arr[i+1])
			}
		}
		if swapped == false {
			// 没有交换发生，说明已经有序
			return
		}
	}
}

// 不借用第三个变量实现交换
func swap(a, b *int) {
	//方式一
	//*a = *a + *b
	//*b = *a - *b // b = a + b -b
	//*a = *a - *b // a = a + b - a
	//方式二
	*a = *a ^ *b
	*b = *a ^ *b
	*a = *a ^ *b
}

func main() {
	arr := []int{10, 3, 0, -9, -1, 2, 100, -22, 54, 23, 7}
	log.Println(arr)

	//selectionSort(arr)
	bubbleSort(arr)

	log.Println(arr)
}
