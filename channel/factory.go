package channel

import (
    "fmt"
    "time"
)

var ch chan int

func pump() chan int {
    ch = make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            ch <- i
        }
    }()
    return ch
}

func suck(ch chan int) {
    for {
        receiver, ok := <-ch
        if ok {
            fmt.Println("success:", receiver)
        } else {
            break
        }
    }
    // 下面这种写法也可以，一直读到管道关闭
    //for v := range ch {
    //	fmt.Println("success:", v)
    //}
}

func Factory() {
    stream := pump()
    go suck(stream)
    time.Sleep(time.Millisecond)
}
