package main

import "fmt"

func main() {
    var ch = make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < 10; i++ {
            ch <- i
        }

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
