package main

import "log"

// -------------------- 排序算法 ------------------------
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
// 两两比较，大的往右
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

// insertionSort 插入排序（最常用 最差O(N^2) 最好 O(N)）
// 将当前的数依次和左边比较，小的往左交换，直到左不再大于0或者不再有左
func insertionSort(arr []int) {
	if arr == nil || len(arr) < 2 {
		return
	}
	n := len(arr)
	for i := 1; i < n; i++ {
		for j := i - 1; j >= 0 && arr[j] > arr[j+1]; j-- {
			//1、左不再大于0
			//2、不再有左
			swap(&arr[j], &arr[j+1])
		}
	}
}

// -------------------- 查找算法 ------------------------

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

	// ------- 排序算法 --------
	//selectionSort(arr)
	//bubbleSort(arr)
	insertionSort(arr)

	// ------- 查找算法 --------

	log.Println(arr)
}
