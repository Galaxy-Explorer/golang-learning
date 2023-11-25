package channel

import (
    "fmt"
    "time"
)

func DeadLock6() {
    var ch = make(chan int)
    go func() {
        defer fmt.Println("end")
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()

    go func() {
        for {
            x := <-ch
            fmt.Println("接收到的值：", x)
        }
    }()

    time.Sleep(5 * time.Second)
}
