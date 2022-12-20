package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var c int64
var d int64
var e int64
var l sync.Mutex
var f int64
var g int64

func demoNormalAdd(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		c++
	}
	wg.Done()
}

func demoMutexAdd(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		l.Lock()
		e++
		l.Unlock()
	}
	wg.Done()
}

func demoAtomicAdd(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&d, 1) // 原子加
	}
	wg.Done()
}

func demoNormalCas(wg *sync.WaitGroup) {
	old := f
	old += 1
	f = old
	wg.Done()
}

func demoAtomicCas(wg *sync.WaitGroup) {
	//CAS操作类似常见的乐观锁机制；该操作在进行交换前首先确保变量的值未被更改，即仍然保持参数 old 所记录的值，满足此前提下才进行交换操作。
	//！！当有大量的goroutine 对变量进行读写操作时，可能导致CAS操作无法成功，这时可以利用for循环多次尝试。
	for {
		old := atomic.LoadInt64(&g)
		if atomic.CompareAndSwapInt64(&g, old, old+1) {
			break
		}
	}
	wg.Done()
}

/**
是什么？
	原子操作，保证并发安全
有什么用？
	代码中的加锁操作因为涉及内核态的上下文切换会比较耗时、代价比较高。针对基本数据类型我们还可以使用原子操作来保证并发安全，因为原子操作是Go语言提供的方法它在用户态就可以完成，因此性能比加锁操作更好。
	原子操作由底层硬件支持，而锁则是由操作系统提供的API实现，若实现相同的功能，前者通常会更有效率
*/

func main() {
	c = 0
	wg := sync.WaitGroup{}

	//---- 基础运用

	// 普通add 输出 9680 | 8656 等不确定值  不是并发安全的
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go demoNormalAdd(&wg)
	}

	// 加锁版add 每次都输出 10000  是并发安全的，但是加锁性能开销大
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go demoMutexAdd(&wg)
	}

	// 原子add 每次都输出 10000  是并发安全，性能优于加锁版
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go demoAtomicAdd(&wg)
	}

	wg.Wait()
	fmt.Printf("demoNormalAdd c = %d\n", c) //\n 不加会输出一个 %
	fmt.Printf("demoAtomicAdd d = %d\n", d)
	fmt.Printf("demoMutexAdd e = %d\n", e)

	//----- 拓展 1 比较并交换操作 CAS
	// 普通 CAS 操作, 输出91 | 81 等不确定值
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go demoNormalCas(&wg)
	}

	// 原子 CAS， 每次输出都是100
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go demoAtomicCas(&wg)
	}

	wg.Wait()
	fmt.Printf("demoNormalCas f = %d\n", f)
	fmt.Printf("demoAtomicCas g = %d\n", g)

	//----- 拓展 2 载入操作
	// 载入操作都以Load为前缀（可避免读取到写入一半的数据）：
	//func LoadInt32(addr *int32) (val int32)
	//func LoadInt64(addr *int64) (val int64)
	//func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
	//func LoadUint32(addr *uint32) (val uint32)
	//func LoadUint64(addr *uint64) (val uint64)
	//func LoadUintptr(addr *uintptr) (val uintptr)

	//----- 拓展 3 载入操作
	//sync/atomic包中的类型Value（相当于一个容器），可用来**“原子地”**存储或加载任意类型的值。
	//func(v *Value) Load() (x interface{}): 读操作，从线程安全的v中读取上一步存放的内容
	//func(v *Value) Store(x interface{}): 写操作，将原始的变量x存放在atomic.Value类型中；

	//value存储的类型要求：
	//
	//不能存储nil（存nil会抛出panic）；
	//value中存储的第一个值，决定了其后续的值类型（以后只能存储此类型的值）；
	//尝试存储不同的类型，会抛出panic；
}
