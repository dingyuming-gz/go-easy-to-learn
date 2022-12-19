package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
)

var sum int32

func addFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n) //输出不是顺序的，不保证顺序
}

func sleepFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func panicFunc() {
	panic("我挂了")
}

func demoPool() {
	// 释放ants的默认协程池
	defer ants.Release()

	var wg sync.WaitGroup

	// 要执行的任务函数
	syncFunc := func() {
		defer wg.Done()
		sleepFunc()
	}

	runTimes := 100
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		// 提交任务到默认协程池
		_ = ants.Submit(syncFunc)
	}

	wg.Wait()
	fmt.Printf("demoPool running goroutines: %d\n", ants.Running()) //当前运行的goroutine的数量
	fmt.Printf("finish all tasks.\n")
}

func demoPoolWithFunc() {

	var wg sync.WaitGroup

	// 初始化协程池
	// 10为goroutine池的容量
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		defer wg.Done()
		addFunc(i)
	})

	// 释放协程池
	defer p.Release()

	// 提交任务
	runTimes := 100
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}

	wg.Wait()
	fmt.Printf("demoPoolWithFunc running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}

func demoReboot() {
	// 释放ants的默认协程池
	defer ants.Release()

	//由于在 demoPool 中 Release 了 默认协程池，这里 又使用了 默认协程池, 所以这里需要 Reboot 一下，
	ants.Reboot() //调用 Reboot() 后仍然可以使用已释放的池
	//调用已经 Release 的池子，会 死锁

	var wg sync.WaitGroup

	// 要执行的任务函数
	syncFunc := func() {
		defer wg.Done()
		sleepFunc()
	}

	runTimes := 10
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		// 提交任务到默认协程池
		_ = ants.Submit(syncFunc)
	}

	wg.Wait()
	fmt.Printf("demoReboot running goroutines: %d\n", ants.Running()) //当前运行的goroutine的数量
	fmt.Printf("finish all tasks.\n")
}

func demoPanic() {

	//默认池子被释放，也不影响 NewPool
	//ants.Release()

	// NewPool 初始化一个新的 pool，容量为 10
	// 无 WithPanicHandler 时，会输出堆栈 信息
	//pool, _ := ants.NewPool(10)

	// 指定 WithPanicHandler， 只输出格式化 信息
	pool, _ := ants.NewPool(10, ants.WithPanicHandler(func(i interface{}) {
		fmt.Printf("panic msg: %v\n", i)
	}))

	// 释放
	defer pool.Release()

	var wg sync.WaitGroup

	// 要执行的任务函数
	panicFunc := func() {
		defer wg.Done()
		panicFunc() //box with bomb
	}

	runTimes := 10
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		// 提交任务到协程池
		_ = pool.Submit(panicFunc)
	}

	wg.Wait()
	fmt.Printf("demoPanic running goroutines: %d\n", ants.Running()) //当前运行的goroutine的数量
	fmt.Printf("finish all tasks.\n")
}

/**
ants 是什么？
	ants 是一个高性能且低损耗的 goroutine 池

ants 有什么用？
    自动调度海量的 goroutines，复用 goroutines
	定期清理过期的 goroutines，进一步节省资源
	提供了大量有用的接口：任务提交、获取运行中的 goroutine 数量、动态调整 Pool 大小、释放 Pool、重启 Pool
	优雅处理 panic，防止程序崩溃
	资源复用，极大节省内存使用量；在大规模批量并发任务场景下比原生 goroutine 并发具有更高的性能
	非阻塞机制

https://github.com/panjf2000/ants
*/

/*
*
ants支持两种协程池，Pool和PoolWithFunc,差别在于，Pool每个任务是一个函数，
PoolWithFunc任务函数在初始化时确定，每个任务只是不同的入参调用该函数。
*/
func main() {

	//1、Pool
	//demoPool()

	//2、PoolWithFunc
	//demoPoolWithFunc()

	//3、使用已释放的池
	//需要取消 demoPool() 注释
	//demoReboot()

	//4、优雅处理 panic，防止程序崩溃
	demoPanic()

	time.Sleep(2 * time.Second)
	fmt.Println("I AM END") //虽然前面 panic 了，但这里 正常输出
}
