package channel

import "fmt"

func DeadLock5() {
    var ch = make(chan int)
    go func() {
        // defer close(ch)
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
    for {
        x, ok := <-ch
        if ok {
            fmt.Println("接收到的值：", x)
        } else {
            break
        }
    }
}
