package channel

func DeadLock7() {
    c := make(chan int, 5)
    for i := 0; i < 10; i++ {
        c <- i
    }
}
