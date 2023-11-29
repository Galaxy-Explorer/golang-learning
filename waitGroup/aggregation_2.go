package waitGroup

import (
    "fmt"
    "sync"
)

func aggregation2() {
    tasksNum := 10

    dataCh := make(chan int)
    stopCh := make(chan struct{})
    resp := make([]int, 0, tasksNum)
    // 启动读goroutine
    // 有可能存在main groutine结束了，但是该读goroutine没有执行完成
    go func() {
        for data := range dataCh {
            resp = append(resp, data)
        }
        stopCh <- struct{}{}
    }()

    // 保证获取到所有数据后，通过 channel 传递到读协程手中
    var wg sync.WaitGroup
    for i := 0; i < tasksNum; i++ {
        wg.Add(1)
        go func(ch chan int, i int) {
            defer wg.Done()
            ch <- i
        }(dataCh, i)
    }
    // 确保所有取数据的协程都完成了工作，才关闭 ch
    wg.Wait()
    close(dataCh)

    <-stopCh

    fmt.Printf("resp: %+v", resp)
}
