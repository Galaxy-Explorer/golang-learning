package waitGroup

import (
    "fmt"
    "sync"
)

func aggregation3() {
    tasksNum := 10

    dataCh := make(chan int)
    resp := make([]int, 0, tasksNum)

    go func() {
        var wg sync.WaitGroup
        for i := 0; i < tasksNum; i++ {
            wg.Add(1)
            go func(ch chan int, i int) {
                defer wg.Done()
                ch <- i
            }(dataCh, i)
        }
        wg.Wait()
        close(dataCh)
    }()

    for x := range dataCh {
        resp = append(resp, x)
    }

    fmt.Printf("resp: %+v", resp)
}
