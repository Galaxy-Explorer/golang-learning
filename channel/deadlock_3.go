package channel

import (
    "fmt"
    "time"
)

func DeadLock3() {
    var ch chan int

    go func() {
        fmt.Println("look this 1")
        ch <- 10
        fmt.Println("look this 2")
    }()
    time.Sleep(time.Second)
    r := <-ch
    fmt.Println(r)
    fmt.Println()
}
