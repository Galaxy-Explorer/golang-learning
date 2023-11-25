package channel

import "fmt"

func DeadLock1() {
    c := make(chan int)
    c <- 1
    fmt.Println(<-c)
}
