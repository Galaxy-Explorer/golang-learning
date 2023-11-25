package channel

import (
    "fmt"
    "time"
)

func DeadLock4() {
    var ch chan int

    go func() {
        fmt.Println("look this 1")
        r := <-ch
        fmt.Println("look this 2", r)
    }()
    time.Sleep(time.Second)
    ch <- 1
    fmt.Println()
}
