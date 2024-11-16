package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	//"github.com/shirou/gopsutil/v4/mem"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	//fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(25 * time.Second)
	//fmt.Println("Hello World!")
}

func main() {
	// var printMemory = func() {
	// 	v, _ := mem.VirtualMemory()
	// 	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// }

	var printGoMemory = func() {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc = %v MiB\n", m.Alloc/(1024*1024))
	}
	defer ants.Release()

	runTimes := 10000

	// Use the common pool.
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}
	wg.Wait()
	//printMemory()
	printGoMemory()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks with common pool.\n")

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	//printMemory()
	printGoMemory()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks with a pool and an added func, result is %d\n", sum)
	// if sum != 499500 {
	// 	panic("the final result is wrong!!!")
	// }

	// Use the MultiPool and set the capacity of the 10 goroutine pools to unlimited.
	// If you use -1 as the pool size parameter, the size will be unlimited.
	// There are two load-balancing algorithms for pools: ants.RoundRobin and ants.LeastTasks.
	// mp, err := ants.NewMultiPool(10, -1, ants.RoundRobin)
	// if err != nil {
	// 	panic(err)
	// }
	// defer mp.ReleaseTimeout(5 * time.Second)
	// for i := 0; i < runTimes; i++ {
	// 	wg.Add(1)
	// 	err = mp.Submit(syncCalculateSum)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// wg.Wait()
	// //printMemory()
	// printGoMemory()
	// fmt.Printf("running goroutines: %d\n", mp.Running())
	// fmt.Printf("finish all tasks with multipool.\n")

	// Use the MultiPoolFunc and set the capacity of 10 goroutine pools to (runTimes/10).
	mpf, err := ants.NewMultiPoolWithFunc(10, runTimes/10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	}, ants.LeastTasks)
	if err != nil {
		panic(err)
	}
	defer mpf.ReleaseTimeout(5 * time.Second)
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = mpf.Invoke(int32(i))
	}
	wg.Wait()
	//printMemory()
	printGoMemory()
	fmt.Printf("running goroutines: %d\n", mpf.Running())
	fmt.Printf("finish all tasks with multipool and func, result is %d\n", sum)
	// if sum != 499500*2 {
	// 	panic("the final result is wrong!!!")
	// }
}
