package channel

import "fmt"

func DeadLock2() {
    c := make(chan int)
    fmt.Println(<-c)
    c <- 1
}
