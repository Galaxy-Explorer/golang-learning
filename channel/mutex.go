package channel

import (
    "fmt"
)

type semaphore chan struct{}

func (s semaphore) lock() {
    s <- struct{}{}
}

func (s semaphore) unlock() {
    <-s
}

func Mutex() {
    test := 0
    sem := make(semaphore)
    for i := 0; i < 10; i++ {
        go func(i int) {
            sem.lock()
            test = i
        }(i)
        sem.unlock()
    }
    close(sem)
    fmt.Println(test)
}
