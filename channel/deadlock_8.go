package channel

import "fmt"

func DeadLock8() {
    c := make(chan int, 5)
    recv := <-c
    fmt.Println(recv)
}
