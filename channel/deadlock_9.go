package channel

import "fmt"

func DeadLock9() {
    var ch = make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        close(ch)
    }()
    for {
        x, ok := <-ch
        if ok {
            fmt.Println("接收到的值：", x)
        } else {
            fmt.Println("end:", x, ok)
            break
        }
    }
}
